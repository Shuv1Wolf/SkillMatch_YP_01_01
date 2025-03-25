package build

import (
	"shuv1wolf/skillmatch/core/controller"
	"shuv1wolf/skillmatch/core/persistence"
	"shuv1wolf/skillmatch/core/service"

	cbuild "github.com/pip-services4/pip-services4-go/pip-services4-components-go/build"
	cref "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
)

type CoreServiceFactory struct {
	cbuild.Factory
}

func NewCoreServiceFactory() *CoreServiceFactory {
	c := &CoreServiceFactory{
		Factory: *cbuild.NewFactory(),
	}

	postgresPersistenceDescriptor := cref.NewDescriptor("core", "persistence", "postgres", "*", "1.0")
	serviceDescriptor := cref.NewDescriptor("core", "service", "default", "*", "1.0")
	httpcontrollerV1Descriptor := cref.NewDescriptor("core", "controller", "http", "*", "1.0")

	c.RegisterType(postgresPersistenceDescriptor, persistence.NewCorePostgresPersistence)
	c.RegisterType(serviceDescriptor, service.NewCoreService)
	c.RegisterType(httpcontrollerV1Descriptor, controller.NewCoreHttpController)

	return c
}
