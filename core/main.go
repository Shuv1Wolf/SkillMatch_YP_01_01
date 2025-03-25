package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
)

type Vacancy struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Area struct {
		Name string `json:"name"`
	} `json:"area"`
	Employer struct {
		Name string `json:"name"`
	} `json:"employer"`
	AlternateURL string `json:"alternate_url"`
}

type Response struct {
	Items []Vacancy `json:"items"`
}

func main() {
	baseURL := "https://api.hh.ru/vacancies"

	// Параметры поиска
	params := url.Values{}
	params.Add("text", "golang") // Поисковый запрос
	params.Add("area", "1")      // Москва, ID региона
	params.Add("per_page", "10") // Кол-во вакансий на странице
	params.Add("page", "0")      // Номер страницы

	fullURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	resp, err := http.Get(fullURL)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	var result Response
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Fatal(err)
	}

	// Вывод результатов
	for _, vacancy := range result.Items {
		description := getVacancyDescription(vacancy.ID)
		fmt.Printf("Вакансия: %s\nКомпания: %s\nГород: %s\nСсылка: %s\nОписание:\n%s\n\n",
			vacancy.Name, vacancy.Employer.Name, vacancy.Area.Name, vacancy.AlternateURL, description)
	}

}
func getVacancyDescription(id string) string {
	url := fmt.Sprintf("https://api.hh.ru/vacancies/%s", id)

	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Ошибка при получении описания: %v", err)
		return ""
	}
	defer resp.Body.Close()

	var data struct {
		Description string `json:"description"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		log.Printf("Ошибка декодирования описания: %v", err)
		return ""
	}

	return data.Description
}
