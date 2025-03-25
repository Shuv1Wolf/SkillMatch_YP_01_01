package containers

import (
	factory "shuv1wolf/skillmatch/core/build"

	cproc "github.com/pip-services4/pip-services4-go/pip-services4-container-go/container"
	rbuild "github.com/pip-services4/pip-services4-go/pip-services4-http-go/build"
	cpg "github.com/pip-services4/pip-services4-go/pip-services4-postgres-go/build"
)

type CoreProcess struct {
	cproc.ProcessContainer
}

func NewCoreProcess() *CoreProcess {
	c := &CoreProcess{
		ProcessContainer: *cproc.NewProcessContainer("core", "Core Skill Match microservice"),
	}

	c.AddFactory(factory.NewCoreServiceFactory())
	c.AddFactory(rbuild.NewDefaultHttpFactory())
	c.AddFactory(cpg.NewDefaultPostgresFactory())

	return c
}
