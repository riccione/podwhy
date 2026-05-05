package llm

import (
	"context"
	"fmt"

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

func (c *Client) Ask(ctx context.Context, prompt string) (string, error) {
	req := &api.GenerateRequest{
		Model:  c.model,
		Prompt: prompt,
		Stream: new(bool),
	}

	resp, err := c.client.Generate(ctx, req)
	if err != nil {
		return "", fmt.Errorf("ollama generate failed: %w", err)
	}

	return resp.Response, nil
}
