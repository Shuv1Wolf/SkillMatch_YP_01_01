package data

import "encoding/json"

type Resume struct {
	Id                     string `json:"id"`
	Skills                 string `json:"skills"`
	Region                 string `json:"region"`
	Salary                 string `json:"salary"`
	EmployeeResponsibility string `json:"employee_responsibility"`
	EmploymentType         string `json:"employment_type"`
	WorkFormat             string `json:"work_format"`
}

func (c *Resume) ResumeToJSONString() string {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return ""
	}
	return string(data)
}
