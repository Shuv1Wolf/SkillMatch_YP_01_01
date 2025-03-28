package service

import (
	"context"
	"encoding/json"
	"fmt"
	"shuv1wolf/skillmatch/core/clients"
	"shuv1wolf/skillmatch/core/data"
	"shuv1wolf/skillmatch/core/helpers"
	persist "shuv1wolf/skillmatch/core/persistence"

	cconf "github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
	cref "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
	ccmd "github.com/pip-services4/pip-services4-go/pip-services4-rpc-go/commands"
	"github.com/sashabaranov/go-openai"
)

type CoreService struct {
	persistence persist.ICorePersistence
	commandSet  *CoreCommandSet
	hhClient    clients.HHClient
	llmClient   clients.OpenAIClient
}

func NewCoreService() *CoreService {
	c := &CoreService{
		hhClient:  *clients.NewHHClient(),
		llmClient: *clients.NewOpenAIClient(),
	}
	return c
}

func (c *CoreService) Configure(ctx context.Context, config *cconf.ConfigParams) {
	// Read configuration parameters here...
}

func (c *CoreService) SetReferences(ctx context.Context, references cref.IReferences) {
	locator := cref.NewDescriptor("core", "persistence", "*", "*", "1.0")
	p, err := references.GetOneRequired(locator)
	if p != nil && err == nil {
		if _pers, ok := p.(persist.ICorePersistence); ok {
			c.persistence = _pers
			return
		}
	}
	panic(cref.NewReferenceError(ctx, locator))
}

func (c *CoreService) GetCommandSet() *ccmd.CommandSet {
	if c.commandSet == nil {
		c.commandSet = NewCoreCommandSet(c)
	}
	return &c.commandSet.CommandSet
}

func (c *CoreService) GetResumeById(ctx context.Context, id string) (data.Resume, error) {
	return c.persistence.GetOneById(ctx, id)
}

func (c *CoreService) AddResume(ctx context.Context, userId string, textResume string) (data.Resume, error) {
	prompt := []openai.ChatCompletionMessage{
		{
			Role: openai.ChatMessageRoleSystem,
			Content: `As a resume parser, convert the following resume text into the specified JSON structure. 
			Extract and map the relevant details from [resume_text] to the fields below. 
			Ensure all fields are filled; use "N/A" if unavailable. Output only valid JSON, nothing else.
			Use only the allowed formats for each field as described below.

			JSON structure:
			{
			"skills": "[comma_separated_skills]",
			"region": "[location]",
			"salary": "[expected_salary]",
			"employee_responsibility": "[key_responsibilities]",
			"employment_type": "[one of: full, part, project, volunteer, probation]",
			"work_format": "[one of: remote, flyInFlyOut, shift, flexible, fullDay, partDay]"
			}
			
			Text resume:
			` + textResume + `
			Important formatting rules:
			"skills" - should be a comma-separated list of skills (e.g., "Golang, Docker, Kubernetes"). It shouldn't be more than 3-4!
			"salary" - should be a numeric string (e.g., "120000") without currency symbols.
			"employment_type" - must strictly match one of: full, part, project, volunteer, or probation.
			"work_format" - must strictly match one of: remote, flyInFlyOut, shift, flexible, fullDay, or partDay.
			If a value is missing or ambiguous, return "N/A".
			`,
		},
	}

	respData, err := c.llmClient.Chat(ctx, prompt)
	if err != nil {
		return data.Resume{}, err
	}
	jsonStr, err := helpers.ExtractJSON(respData)
	if err != nil {
		fmt.Println("Ошибка извлечения JSON:", err)
		return data.Resume{}, err
	}

	var resume data.Resume
	resume.Id = userId

	err = json.Unmarshal([]byte(jsonStr), &resume)
	if err != nil {
		return data.Resume{}, err
	}

	result, err := c.persistence.Create(ctx, resume)
	if err != nil {
		return data.Resume{}, err
	}

	return result, nil
}

func (c *CoreService) FindJob(ctx context.Context, userId string) ([]*clients.Vacancy, error) {
	resume, err := c.persistence.GetOneById(ctx, userId)
	if err != nil {
		return nil, err
	}

	vacancies, err := c.hhClient.FindVacanciesByResume(resume)
	if err != nil {
		return nil, err
	}

	if len(vacancies) == 0 {
		return nil, nil
	}

	limit := 3
	if len(vacancies) < 3 {
		limit = len(vacancies)
	}

	resp := make([]*clients.Vacancy, 0, limit)

	for i := 0; i < limit; i++ {
		text, err := c.hhClient.GetVacancyText(vacancies[i].Id)
		if err != nil {
			return nil, err
		}

		prompt := []openai.ChatCompletionMessage{
			{
				Role: openai.ChatMessageRoleSystem,
				Content: `As a job matching assistant, evaluate how well the given job vacancy matches the candidate's resume.

				Use the structured resume data and the plain-text vacancy description below.
				Score the match as a percentage from 0 to 100, where:
				90-100% = excellent match
				70-89% = good match
				50-69% = moderate match
				below 50% = weak or poor match

				Consider the following factors:
				Skills alignment (importance: high)
				Responsibilities fit
				Salary expectations
				Work format compatibility (e.g., remote/fullDay)
				Employment type (e.g., full-time/project)
				Just give me the numbers as an answer!

				Resume:` + resume.ResumeToJSONString() + `
				Vacancy:` + text,
			},
		}

		respData, err := c.llmClient.Chat(ctx, prompt)
		if err != nil {
			return nil, err
		}

		vacancies[i].Score = respData
		resp = append(resp, &vacancies[i])
	}

	return resp, nil
}
