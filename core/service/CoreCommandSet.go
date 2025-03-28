package service

import (
	"context"
	"shuv1wolf/skillmatch/core/data"

	cconv "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/convert"
	exec "github.com/pip-services4/pip-services4-go/pip-services4-components-go/exec"
	cvalid "github.com/pip-services4/pip-services4-go/pip-services4-data-go/validate"
	ccmd "github.com/pip-services4/pip-services4-go/pip-services4-rpc-go/commands"
)

type CoreCommandSet struct {
	ccmd.CommandSet
	controller    ICoreService
	coreConvertor cconv.IJSONEngine[data.Resume]
}

func NewCoreCommandSet(controller ICoreService) *CoreCommandSet {
	c := &CoreCommandSet{
		CommandSet:    *ccmd.NewCommandSet(),
		controller:    controller,
		coreConvertor: cconv.NewDefaultCustomTypeJsonConvertor[data.Resume](),
	}

	c.AddCommand(c.makeAddResumeCommand())
	c.AddCommand(c.makeGetResumeByIdCommand())
	c.AddCommand(c.makeFindJobCommand())

	return c
}

func (c *CoreCommandSet) makeAddResumeCommand() ccmd.ICommand {
	return ccmd.NewCommand(
		"add_resume",
		cvalid.NewObjectSchema().
			WithRequiredProperty("user_id", cconv.String).
			WithRequiredProperty("text_resume", cconv.String),
		func(ctx context.Context, args *exec.Parameters) (result any, err error) {
			return c.controller.AddResume(ctx, args.GetAsString("user_id"), args.GetAsString("text_resume"))
		})
}

func (c *CoreCommandSet) makeGetResumeByIdCommand() ccmd.ICommand {
	return ccmd.NewCommand(
		"get_resume",
		cvalid.NewObjectSchema().
			WithRequiredProperty("user_id", cconv.String),
		func(ctx context.Context, args *exec.Parameters) (result any, err error) {
			return c.controller.GetResumeById(ctx, args.GetAsString("user_id"))
		})
}

func (c *CoreCommandSet) makeFindJobCommand() ccmd.ICommand {
	return ccmd.NewCommand(
		"find_job",
		cvalid.NewObjectSchema().
			WithRequiredProperty("user_id", cconv.String),
		func(ctx context.Context, args *exec.Parameters) (result any, err error) {
			return c.controller.FindJob(ctx, args.GetAsString("user_id"))
		})
}
