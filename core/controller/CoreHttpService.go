package controller

import (
	"context"

	cref "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
	cservices "github.com/pip-services4/pip-services4-go/pip-services4-http-go/controllers"
)

type CoreHttpController struct {
	cservices.CommandableHttpController
}

func NewCoreHttpController() *CoreHttpController {
	c := &CoreHttpController{}
	c.CommandableHttpController = *cservices.InheritCommandableHttpController(c, "v1/core")
	c.DependencyResolver.Put(context.Background(), "service", cref.NewDescriptor("core", "service", "*", "*", "1.0"))
	return c
}
