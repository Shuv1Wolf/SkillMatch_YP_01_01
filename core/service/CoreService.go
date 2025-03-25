package service

import (
	"context"
	"fmt"
	"shuv1wolf/skillmatch/core/data"
	persist "shuv1wolf/skillmatch/core/persistence"

	cconf "github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
	cref "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
	ccmd "github.com/pip-services4/pip-services4-go/pip-services4-rpc-go/commands"
)

type CoreService struct {
	persistence persist.ICorePersistence
	commandSet  *CoreCommandSet
}

func NewCoreService() *CoreService {
	c := &CoreService{}
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
	fmt.Println(userId)
	fmt.Println(textResume)

	return data.Resume{}, nil
}

func (c *CoreService) FindJob(ctx context.Context, userId string) (string, error) {

	return "JOBS", nil
}
