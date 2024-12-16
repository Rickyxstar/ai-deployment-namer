# AI Deployment Namer

This Helm Chart installs the AI Deployment Namer onto Kubernetes.

#### ChatGPT

Grab an [OpenAI api key](https://help.openai.com/en/articles/4936850-where-do-i-find-my-openai-api-key)

#### Deploy

```
helm install ai-deployment-namer \
 -n ai-deployment-namer --create-namespace \
 --set chatgpt.token=<your api key> \
./ai-deployment-namer
```

#### Values

| Parameter     | Description                                                         | Type   | Required                      | Default Value                               |
| ------------- | ------------------------------------------------------------------- | ------ | ----------------------------- | ------------------------------------------- |
| webhook.image | Docker image for the webhook service.                               | String | Yes                           | rickyxstar/ai-deployment-namer:0.1.0        |
| nameGenerator | The name generator to use. Currently supports chatgpt and ollama.   | String | Yes                           | chatgpt                                     |
| model         | The specific language model to use within the chosen nameGenerator. | String | Yes                           | gpt-4o-mini                                 |
| ollama.image  | Docker image for the Ollama service (if used).                      | String | No                            | ollama/ollama:latest                        |
| chatgpt.token | OpenAI API token for using ChatGPT as the name generator.           | String | Yes if nameGenerator: chatgpt | "" (Empty string - requires manual setting) |

Note: The chatgpt.token is crucial if you are using chatgpt as your nameGenerator. You must replace the empty string ("") with your actual OpenAI API token. Other name generators might require different configuration parameters. The ollama settings are only relevant if you intend to integrate with the Ollama LLM service.
