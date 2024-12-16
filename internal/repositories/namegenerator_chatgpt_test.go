package repositories

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/rickyxstar/ai-deployment-namer/internal/common"
	openai "github.com/sashabaranov/go-openai"
	"github.com/stretchr/testify/assert"
	appsv1 "k8s.io/api/apps/v1"
)

func TestChatGPT(t *testing.T) {
	t.Run("Handles a response from ChatGPT", func(t *testing.T) {
		want := "Redis-rises-again-not-really-sorry-we're-still-testing"
		model := "4o"

		g := NewNameGeneratorChatGPT("token", &http.Client{
			Transport: &common.MockTransport{
				ResponseFunc: func(req *http.Request) (*http.Response, error) {
					switch req.URL.String() {
					case "https://api.openai.com/v1/chat/completions":
						res, err := json.Marshal(openai.ChatCompletionResponse{
							Choices: []openai.ChatCompletionChoice{
								{
									Message: openai.ChatCompletionMessage{
										Content: want,
									},
								},
							},
						})
						if err != nil {
							return nil, errors.New("could not marshal json")
						}

						return &http.Response{
							StatusCode: http.StatusAccepted,
							Body:       io.NopCloser(bytes.NewReader(res)),
							Header:     make(http.Header),
						}, nil

					default:
						return &http.Response{
							StatusCode: http.StatusNotFound,
							Body:       io.NopCloser(bytes.NewBufferString(`{"error": "not found"}`)),
							Header:     make(http.Header),
						}, nil
					}
				},
			},
		}, model)

		got, err := g.Generate(context.TODO(), &appsv1.Deployment{})

		assert.NoError(t, err)
		assert.Equal(t, want, got)
	})
}
