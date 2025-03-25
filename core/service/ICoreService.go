package service

import (
	"context"
	"shuv1wolf/skillmatch/core/data"
)

type ICoreService interface {
	GetResumeById(ctx context.Context, userId string) (data.Resume, error)
	AddResume(ctx context.Context, userId string, textResume string) (data.Resume, error)
	FindJob(ctx context.Context, userId string) (string, error)
}
