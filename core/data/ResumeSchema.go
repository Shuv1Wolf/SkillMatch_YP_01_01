package data

import (
	cconv "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/convert"
	cvalid "github.com/pip-services4/pip-services4-go/pip-services4-data-go/validate"
)

type ResumeSchema struct {
	cvalid.ObjectSchema
}

func NewResumeSchema() *ResumeSchema {
	c := ResumeSchema{}
	c.ObjectSchema = *cvalid.NewObjectSchema()

	c.WithOptionalProperty("id", cconv.String)
	c.WithRequiredProperty("skills", cconv.String)
	c.WithOptionalProperty("region", cconv.String)
	c.WithRequiredProperty("salary", cconv.String)
	c.WithOptionalProperty("employee_responsibility", cconv.String)
	c.WithOptionalProperty("employment_type", cconv.String)
	c.WithOptionalProperty("work_format", cconv.String)
	return &c
}
