package repositories

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"slices"

	ollama "github.com/ollama/ollama/api"
	appsv1 "k8s.io/api/apps/v1"
)

var stream = false

type nameGeneratorOllama struct {
	client *ollama.Client
	model  string
}

func NewNameGeneratorOllama(base string, http *http.Client, model string) *nameGeneratorOllama {
	u, err := url.Parse(base)
	if err != nil {
		panic(err)
	}

	c := ollama.NewClient(u, http)

	return &nameGeneratorOllama{c, model}
}

func (n *nameGeneratorOllama) Generate(ctx context.Context, deployment *appsv1.Deployment) (string, error) {
	bytes, err := json.Marshal(deployment)
	if err != nil {
		return "", err
	}

	models, err := n.client.List(ctx)
	if err != nil {

	}

	found := slices.ContainsFunc(models.Models, func(m ollama.ListModelResponse) bool {
		return m.Name == n.model
	})

	if !found {
		if err = n.client.Pull(
			ctx,
			&ollama.PullRequest{
				Model:  n.model,
				Stream: &stream,
			},
			func(p ollama.ProgressResponse) error { return nil },
		); err != nil {
			return "", err
		}
	}

	var genRequest = ollama.GenerateRequest{
		Stream: &stream,
		Model:  n.model,
		System: "You are apart of a system to system workflow. Give only short one line answers. Be a little unhinged",
	}

	var resp *string = nil
	genRequest.Prompt = fmt.Sprintf("Give me a funny kubernetes deployment name that is less than 253 characters for the following deployment: \n %s", string(bytes))

	err = n.client.Generate(
		ctx,
		&genRequest,
		func(r ollama.GenerateResponse) error {
			resp = &r.Response
			return nil
		},
	)

	if err != nil {
		return "", err
	}

	if resp == nil {
		return "", errors.New("could not generate a pod name")
	}

	return *resp, nil
}

var _ NameGenerator = &nameGeneratorOllama{}
