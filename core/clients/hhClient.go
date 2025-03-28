package clients

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"shuv1wolf/skillmatch/core/data"
	"strings"
)

type Vacancy struct {
	Name  string `json:"name"`
	Link  string `json:"link"`
	Id    string
	Score string
}

type HHClient struct {
	baseURL string
	client  *http.Client
}

func NewHHClient() *HHClient {
	return &HHClient{
		baseURL: "https://api.hh.ru",
		client:  http.DefaultClient,
	}
}

type hhVacancyItem struct {
	Name string `json:"name"`
	Alt  string `json:"alternate_url"`
}

type hhSearchResponse struct {
	Items []hhVacancyItem `json:"items"`
}

// Moscow, SPB, Rostov-on-Don
var defaultAreaIDs = []string{"1", "2", "76"}

func (hh *HHClient) FindVacanciesByResume(resume data.Resume) ([]Vacancy, error) {
	endpoint := fmt.Sprintf("%s/vacancies", hh.baseURL)

	params := url.Values{}
	if resume.Skills != "" {
		params.Add("text", resume.Skills)
	}

	for _, id := range defaultAreaIDs {
		params.Add("area", id)
	}

	reqURL := fmt.Sprintf("%s?%s", endpoint, params.Encode())
	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := hh.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var searchResp hhSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&searchResp); err != nil {
		return nil, err
	}

	var vacancies []Vacancy
	for _, item := range searchResp.Items {
		vacancies = append(vacancies, Vacancy{
			Name: item.Name,
			Link: item.Alt,
			Id:   ExtractVacancyIDFromURL(item.Alt),
		})
	}

	return vacancies, nil
}

func (hh *HHClient) GetVacancyText(vacancyID string) (string, error) {
	endpoint := fmt.Sprintf("%s/vacancies/%s", hh.baseURL, vacancyID)

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return "", err
	}

	resp, err := hh.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("failed to get vacancy: status %d", resp.StatusCode)
	}

	var data struct {
		Description string `json:"description"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", err
	}

	return data.Description, nil
}

func ExtractVacancyIDFromURL(link string) string {
	parts := strings.Split(link, "/")
	return parts[len(parts)-1]
}
