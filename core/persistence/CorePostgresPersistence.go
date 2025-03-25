package persistence

import (
	"shuv1wolf/skillmatch/core/data"

	cpg "github.com/pip-services4/pip-services4-go/pip-services4-postgres-go/persistence"
)

type CorePostgresPersistence struct {
	cpg.IdentifiableJsonPostgresPersistence[data.Resume, string]
}

func NewCorePostgresPersistence() *CorePostgresPersistence {
	c := &CorePostgresPersistence{}
	c.IdentifiableJsonPostgresPersistence = *cpg.InheritIdentifiableJsonPostgresPersistence[data.Resume, string](c, "resumes")
	return c
}

func (c *CorePostgresPersistence) DefineSchema() {
	c.ClearSchema()
	c.IdentifiableJsonPostgresPersistence.DefineSchema()
	c.EnsureTable("", "")
}
