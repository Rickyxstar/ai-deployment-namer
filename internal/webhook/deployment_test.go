package webhook

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/rickyxstar/ai-deployment-namer/internal/repositories"
)

func TestMakeNameSafe(t *testing.T) {
	type testCase struct {
		value string
		want  string
	}

	testCases := []testCase{
		{
			value: "Redis-Rises-Again-From-The-Ashes-Of-Doom-And-Gloom-In-the-Cloud",
			want:  "redis-rises-again-from-the-ashes-of-doom-and-gloom-in-the-cloud",
		},
		{
			value: "REDIS-RULES-OKAY-FINE-IT'S-GONNA-RUN-WITH-3-CONTAINERS-AND-ITS-GONNA-BE-FINE-I-GUARANTEE-IT",
			want:  "redis-rules-okay-fine-it-s-gonna-run-with-3-containers-and-its",
		},
		{
			value: "Redis-Rampage-In-Production-May-The-Odds-Be-Ever-In-Favor-Of-Slowing-It-Down-Deployment-v1.0",
			want:  "redis-rampage-in-production-may-the-odds-be-ever-in-favor-of-sl.0",
		},
		{
			value: "Redis-Go-Bananas-3x-Party-Destruction",
			want:  "redis-go-bananas-3x-party-destruction",
		},
		{
			value: "Redis-rises-again-not-really-sorry-we're-still-testing",
			want:  "redis-rises-again-not-really-sorry-we-re-still-testing",
		},
	}

	for _, c := range testCases {
		got := ensureValidDNSSubdomain(c.value)
		assert.Equal(t, c.want, got)
	}
}

func TestDeploymentNamer(t *testing.T) {
	testsCases := []struct {
		name          string
		obj           runtime.Object
		generator     *repositories.MockNameGenerator
		expectedName  string
		expectedError error
	}{
		{
			name: "Valid Deployment",
			obj: &appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Name: "not-funny-name",
				},
			},
			generator:     &repositories.MockNameGenerator{GeneratedName: "funny-name"},
			expectedName:  "funny-name",
			expectedError: nil,
		},
		{
			name:          "Invalid Object Type",
			obj:           &appsv1.DaemonSet{},
			generator:     &repositories.MockNameGenerator{},
			expectedName:  "",
			expectedError: fmt.Errorf("expected a deployment but got a *v1.DaemonSet"),
		},
		{
			name:          "Generator Error",
			obj:           &appsv1.Deployment{},
			generator:     &repositories.MockNameGenerator{Err: fmt.Errorf("could not generate name")},
			expectedName:  "",
			expectedError: fmt.Errorf("could not generate name"),
		},
	}

	for _, c := range testsCases {
		t.Run(c.name, func(t *testing.T) {
			dn := NewDeploymentNamer(c.generator)
			err := dn.Default(context.TODO(), c.obj)

			if c.expectedError != nil {
				assert.Equal(t, c.expectedError, err)
			} else {
				assert.NoError(t, err)

				deployment, ok := c.obj.(*appsv1.Deployment)
				if !ok {
					t.Fatalf("Object is not a deployment")
				}

				assert.Equal(t, c.expectedName, deployment.Name)
			}
		})
	}
}
