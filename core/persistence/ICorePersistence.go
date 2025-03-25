package persistence

import (
	"context"
	"shuv1wolf/skillmatch/core/data"
)

type ICorePersistence interface {
	GetOneById(ctx context.Context, id string) (data.Resume, error)
	Create(ctx context.Context, resume data.Resume) (data.Resume, error)
}
