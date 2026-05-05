package llm

import (
	"context"
	"fmt"
	"time"

	"github.com/ollama/ollama/api"
)

type Client struct {
	client *api.Client
	model  string
}

func NewClient(model string) (*Client, error) {
	client, err := api.ClientFromEnvironment()
	if err != nil {
		return nil, fmt.Errorf("failed to create ollama client: %w", err)
	}

	return &Client{
		client: client,
		model:  model,
	}, nil
}

func (c *Client) CheckConnection(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := c.client.Version(ctx)
	if err != nil {
		return fmt.Errorf("ollama server not reachable: %w\n\nTo start Ollama:\n  ollama serve\n\nTo pull model %s:\n  ollama pull %s", err, c.model, c.model)
	}
	return nil
}

func (c *Client) Ask(ctx context.Context, prompt string, systemPrompt string) (string, error) {
	var response string

	req := &api.GenerateRequest{
		Model:  c.model,
		Prompt: prompt,
		System: systemPrompt,
		Stream: new(bool),
	}

	err := c.client.Generate(ctx, req, func(resp api.GenerateResponse) error {
		response += resp.Response
		return nil
	})

	if err != nil {
		return "", fmt.Errorf("ollama generate failed: %w", err)
	}

	return response, nil
}
