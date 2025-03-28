import os
import logging
import docx
import PyPDF2
from pathlib import Path
from aiogram import Bot, Dispatcher, types, F
from aiogram.filters import Command
from aiogram.types import BotCommand, Message, ContentType
import asyncio
import httpx
from dotenv import load_dotenv

logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

load_dotenv()

API_TOKEN = os.getenv('TG_TOKEN')
RESUMES_DIR = "user_resumes"
os.makedirs(RESUMES_DIR, exist_ok=True)

bot = Bot(token=API_TOKEN)
dp = Dispatcher()


def extract_text_from_pdf(file_path: str) -> str:
    with open(file_path, 'rb') as file:
        reader = PyPDF2.PdfReader(file)
        text = "\n".join([page.extract_text() for page in reader.pages])
    return text


def extract_text_from_docx(file_path: str) -> str:
    doc = docx.Document(file_path)
    return "\n".join([para.text for para in doc.paragraphs])


async def process_resume(user_id: str, resume_text: str):
    api_url = "http://127.0.0.1:8080/v1/core/add_resume"
    payload = {
        "user_id": f"{user_id}",
        "text_resume": resume_text
    }

    headers = {
        "Content-Type": "application/json",
        "Accept": "*/*",                      
        "User-Agent": "SkillMatchBot/1.0",                 
        "Host": "127.0.0.1:8080"                           
    }

    try:
        async with httpx.AsyncClient(timeout=100.0) as client:
            response = await client.post(api_url, json=payload, headers=headers)
            response.raise_for_status()
            
            data_size_kb = len(resume_text.encode('utf-8')) / 1024
            return (
                "✅ Ваше резюме успешно сохранено!\n"
                f"Размер: {data_size_kb:.1f} КБ\n"
                "Теперь вы можете использовать /find_jobs для поиска вакансий"
            )
    except httpx.RequestError as e:
        return f"❌ Ошибка при сохранении резюме: {str(e)}"


async def set_commands():
    commands = [
        BotCommand(command="/start", description="Начать работу с ботом"),
        BotCommand(command="/help", description="Помощь и список команд"),
        BotCommand(command="/resume", description="Добавить/изменить резюме"),
        BotCommand(command="/find_jobs", description="Найти вакансии")
    ]
    await bot.set_my_commands(commands)


@dp.message(Command("start"))
async def send_welcome(message: types.Message):
    welcome_text = (
        "Приветствуем в <b>SkillMatch BOT</b>!\n\n"
        "Этот бот поможет вам быстро найти подходящие вакансии.\n\n"
        "Отправьте /resume чтобы загрузить резюме\n"
        "Используйте /help для списка команд"
    )
    await message.answer(welcome_text, parse_mode='HTML')


@dp.message(Command("help"))
async def send_help(message: types.Message):
    help_text = (
        "<b>Доступные команды:</b>\n\n"
        "/resume - Загрузить или изменить резюме\n"
        "/find_jobs - Найти вакансии по вашему резюме\n\n"
        "Поддерживаемые форматы: PDF, DOCX, TXT или текст"
    )
    await message.answer(help_text, parse_mode='HTML')


from aiogram import types
from aiogram.filters import Command
from aiogram.types import Message
import httpx

@dp.message(Command("find_jobs"))
async def request_jobs(message: Message):
    api_url = "http://127.0.0.1:8080/v1/core/find_job"
    payload = {
        "user_id": f"{message.from_user.id}",
    }

    headers = {
        "Content-Type": "application/json",
        "Accept": "*/*",
        "User-Agent": "SkillMatchBot/1.0",
        "Host": "127.0.0.1:8080"
    }

    try:
        async with httpx.AsyncClient(timeout=100.0) as client:
            response = await client.post(api_url, json=payload, headers=headers)
            response.raise_for_status()

            vacancies = response.json()

            # Преобразуем score к числу
            def parse_score(v):
                score_str = str(v.get("Score", "0")).replace("%", "")
                try:
                    return int(score_str)
                except ValueError:
                    return 0

            # Сортировка по убыванию score
            vacancies.sort(key=parse_score, reverse=True)

            if not vacancies:
                return await message.answer("❌ Вакансии не найдены по вашему резюме.")

            text = "🔎 *Вот вакансии, подходящие вам по резюме:*\n\n"

            for v in vacancies:
                score = str(v.get("Score", "N/A"))
                text += (
                    f"📌 *{v['name']}*\n"
                    f"🔗 [Открыть вакансию]({v['link']})\n"
                    f"🎯 Совпадение: *{score}%*\n\n"
                )

            await message.answer(text, parse_mode="Markdown")

    except httpx.RequestError as e:
        await message.answer(f"❌ Ошибка при поиске вакансий: {str(e)}")




@dp.message(Command("resume"))
async def request_resume(message: Message):
    await message.answer(
        "Пожалуйста, отправьте ваше резюме:\n"
        "- Как файл (PDF, DOCX или TXT)\n"
        "- Или текстовым сообщением"
    )


@dp.message(F.content_type.in_({ContentType.DOCUMENT, ContentType.TEXT}))
async def handle_resume(message: Message):
    user_id = message.from_user.id
    text = ""

    if message.content_type == ContentType.TEXT:
        text = message.text
        await message.answer("Резюме получено в виде текста. Обрабатываю...")

    elif message.content_type == ContentType.DOCUMENT:
        file = await bot.get_file(message.document.file_id)
        file_ext = Path(message.document.file_name).suffix.lower()

        if file_ext not in ('.pdf', '.docx', '.txt'):
            await message.answer("Пожалуйста, отправьте файл в формате PDF, DOCX или TXT")
            return

        file_path = Path(RESUMES_DIR) / f"temp_{user_id}{file_ext}"
        await bot.download_file(file.file_path, destination=file_path)

        await message.answer(f"Файл {message.document.file_name} получен. Обрабатываю...")

        try:
            if file_ext == '.pdf':
                text = extract_text_from_pdf(file_path)
            elif file_ext == '.docx':
                text = extract_text_from_docx(file_path)
            elif file_ext == '.txt':
                with open(file_path, 'r', encoding='utf-8') as f:
                    text = f.read()

            os.remove(file_path)

        except Exception as e:
            await message.answer("Ошибка при обработке файла. Попробуйте другой формат.")
            logger.error(f"Error processing file: {e}")
            return

    if text:
        resp = await process_resume(user_id, text)
        await message.answer(resp)
    else:
        await message.answer("Не удалось извлечь текст из резюме. Попробуйте другой формат.")


async def main():
    await set_commands()
    await dp.start_polling(bot)


if __name__ == '__main__':
    asyncio.run(main())