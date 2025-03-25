package data

type Resume struct {
	UserTgId               string `json:"user_tg_id"`
	Skills                 string `json:"skills"`
	Region                 string `json:"region"`
	Salary                 string `json:"salary"`
	EmployeeResponsibility string `json:"employee_responsibility"`
	EmploymentType         string `json:"employment_type"`
	WorkFormat             string `json:"work_format"`
}
