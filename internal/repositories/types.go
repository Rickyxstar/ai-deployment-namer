package repositories

import (
	"context"

	appsv1 "k8s.io/api/apps/v1"
)

type NameGenerator interface {
	Generate(ctx context.Context, deployment *appsv1.Deployment) (string, error)
}
