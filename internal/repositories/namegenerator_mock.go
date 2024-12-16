package repositories

import (
	"context"

	appsv1 "k8s.io/api/apps/v1"
)

type MockNameGenerator struct {
	GeneratedName string
	Err           error
}

func (m *MockNameGenerator) Generate(ctx context.Context, deployment *appsv1.Deployment) (string, error) {
	return m.GeneratedName, m.Err
}

var _ NameGenerator = &MockNameGenerator{}
