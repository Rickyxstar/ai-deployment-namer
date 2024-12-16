package repositories

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"testing"

	ollama "github.com/ollama/ollama/api"
	"github.com/rickyxstar/ai-deployment-namer/internal/common"
	"github.com/stretchr/testify/assert"
	appsv1 "k8s.io/api/apps/v1"
)

func TestOllama(t *testing.T) {
	t.Run("Handles a response from Ollama when target model is missing", func(t *testing.T) {
		want := "Redis-rises-again-not-really-sorry-we're-still-testing"

		g := NewNameGeneratorOllama("http://127.0.0.1:11434", &http.Client{
			Transport: &common.MockTransport{
				ResponseFunc: func(req *http.Request) (*http.Response, error) {
					switch req.URL.String() {
					case "http://127.0.0.1:11434/api/generate":
						res, err := json.Marshal(ollama.GenerateResponse{
							Response: want,
							Done:     true,
						})
						if err != nil {
							return nil, errors.New("could not marshal json")
						}

						return &http.Response{
							StatusCode: http.StatusAccepted,
							Body:       io.NopCloser(bytes.NewReader(res)),
							Header:     make(http.Header),
						}, nil

					case "http://127.0.0.1:11434/api/tags":
						res, err := json.Marshal(ollama.ListResponse{})
						if err != nil {
							return nil, errors.New("could not marshal json")
						}

						return &http.Response{
							StatusCode: http.StatusAccepted,
							Body:       io.NopCloser(bytes.NewReader(res)),
							Header:     make(http.Header),
						}, nil

					case "http://127.0.0.1:11434/api/pull":
						res, err := json.Marshal(ollama.ProgressResponse{})
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
		}, "model")

		got, err := g.Generate(context.TODO(), &appsv1.Deployment{})

		assert.NoError(t, err)
		assert.Equal(t, want, got)
	})

	t.Run("Handles a response from Ollama when target model is found", func(t *testing.T) {
		want := "Redis-rises-again-not-really-sorry-we're-still-testing"
		model := "model:latest"

		g := NewNameGeneratorOllama("http://127.0.0.1:11434", &http.Client{
			Transport: &common.MockTransport{
				ResponseFunc: func(req *http.Request) (*http.Response, error) {
					switch req.URL.String() {
					case "http://127.0.0.1:11434/api/generate":
						res, err := json.Marshal(ollama.GenerateResponse{
							Response: want,
							Done:     true,
						})
						if err != nil {
							return nil, errors.New("could not marshal json")
						}

						return &http.Response{
							StatusCode: http.StatusAccepted,
							Body:       io.NopCloser(bytes.NewReader(res)),
							Header:     make(http.Header),
						}, nil

					case "http://127.0.0.1:11434/api/tags":
						res, err := json.Marshal(ollama.ListResponse{
							Models: []ollama.ListModelResponse{
								{
									Name: model,
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
