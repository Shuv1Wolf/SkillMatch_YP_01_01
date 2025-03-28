package persistence

import (
	"context"
	"shuv1wolf/skillmatch/core/data"

	cpg "github.com/pip-services4/pip-services4-go/pip-services4-postgres-go/persistence"
)

type CorePostgresPersistence struct {
	cpg.IdentifiablePostgresPersistence[data.Resume, string]
}

func NewCorePostgresPersistence() *CorePostgresPersistence {
	c := &CorePostgresPersistence{}
	c.IdentifiablePostgresPersistence = *cpg.InheritIdentifiablePostgresPersistence[data.Resume, string](c, "resumes")
	return c
}

func (c *CorePostgresPersistence) DefineSchema() {
	c.ClearSchema()
	c.IdentifiablePostgresPersistence.DefineSchema()
	c.EnsureSchema("CREATE TABLE " + c.QuotedTableName() + ` (
        "id" TEXT PRIMARY KEY,
        "skills" TEXT,
        "region" TEXT,
        "salary" TEXT,
        "employee_responsibility" TEXT,
        "employment_type" TEXT,
        "work_format" TEXT
    )`)
}

func (c *CorePostgresPersistence) Create(ctx context.Context, resume data.Resume) (data.Resume, error) {
	if item, err := c.GetOneById(ctx, resume.Id); err == nil && item.Id != "" {
		return c.IdentifiablePostgresPersistence.Update(ctx, resume)
	} else {
		return c.IdentifiablePostgresPersistence.Create(ctx, resume)
	}
}
