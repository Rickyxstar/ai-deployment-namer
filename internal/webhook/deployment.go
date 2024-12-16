package webhook

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/rickyxstar/ai-deployment-namer/internal/repositories"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

type DeploymentNamer struct {
	generator repositories.NameGenerator
}

func NewDeploymentNamer(generator repositories.NameGenerator) *DeploymentNamer {
	return &DeploymentNamer{generator}
}

func (a *DeploymentNamer) Default(ctx context.Context, obj runtime.Object) error {
	log := logf.FromContext(ctx)
	deployment, ok := obj.(*appsv1.Deployment)
	if !ok {
		return fmt.Errorf("expected a deployment but got a %T", obj)
	}

	log.Info("Generating a better name", "deployment", deployment.Name, "namespace", deployment.Namespace)

	name, err := a.generator.Generate(ctx, deployment)
	if err != nil {
		log.Error(err, "could not generate deployment name")
		return err
	}

	deployment.Name = ensureValidDNSSubdomain(name)

	log.Info("Set a better name", "deployment", deployment.Name, "namespace", deployment.Namespace)

	return nil
}

// EnsureValidDNSSubdomain transforms the input string into a valid DNS subdomain name per RFC 1123.
func ensureValidDNSSubdomain(name string) string {
	// Convert to lowercase
	name = strings.ToLower(name)

	// Replace invalid characters with a hyphen
	re := regexp.MustCompile(`[^a-z0-9.-]`)
	name = re.ReplaceAllString(name, "-")

	// Split into labels and ensure each is valid
	labels := strings.Split(name, ".")
	for i, label := range labels {
		if len(label) > 63 {
			label = label[:63] // Truncate to 63 characters
		}
		label = strings.Trim(label, "-") // Remove leading and trailing hyphens
		if len(label) == 0 {
			label = "a" // Replace empty labels with a valid default
		}
		labels[i] = label
	}

	// Join the labels back together
	name = strings.Join(labels, ".")

	// Ensure the total length is no more than 253 characters
	if len(name) > 253 {
		name = name[:253]
	}

	// Ensure the final string starts and ends with an alphanumeric character
	name = strings.Trim(name, "-.")

	if len(name) == 0 {
		return "a" // Default to "a" if the string ends up empty
	}

	return name
}
