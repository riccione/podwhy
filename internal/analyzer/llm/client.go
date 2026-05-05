package llm

import (
	"context"
	"fmt"
)

type Client struct {
	model string
}

func NewClient(model string) (*Client, error) {
	return &Client{
		model: model,
	}, nil
}

func (c *Client) Ask(ctx context.Context, prompt string) (string, error) {
	return fmt.Sprintf("LLM analysis for prompt: %s", truncate(prompt, 100)), nil
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}
