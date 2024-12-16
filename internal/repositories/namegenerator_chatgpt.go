package repositories

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	openai "github.com/sashabaranov/go-openai"
	v1 "k8s.io/api/apps/v1"
)

type nameGeneratorChatGPT struct {
	client *openai.Client
	model  string
}

func NewNameGeneratorChatGPT(token string, http *http.Client, model string) *nameGeneratorChatGPT {
	config := openai.DefaultConfig(token)
	config.HTTPClient = http

	c := openai.NewClientWithConfig(config)

	return &nameGeneratorChatGPT{c, model}
}

// Generate implements NameGenerator.
func (n *nameGeneratorChatGPT) Generate(ctx context.Context, deployment *v1.Deployment) (string, error) {
	bytes, err := json.Marshal(deployment)
	if err != nil {
		return "", err
	}

	resp, err := n.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model:  n.model,
		Stream: false,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: "You are apart of a system to system workflow. Give only short one line answers. Be unhinged",
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: fmt.Sprintf("Give me a funny kubernetes deployment name that is less than 253 characters for the following deployment: \n %s", string(bytes)),
			},
		},
	})
	if err != nil {
		return "", err
	}

	return resp.Choices[0].Message.Content, nil
}

var _ NameGenerator = &nameGeneratorChatGPT{}
