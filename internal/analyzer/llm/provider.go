package llm

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/podwhy/podwhy/internal/observer"
)

type LLMProvider struct {
	client *Client
	model  string
}

type LLMResponse struct {
	RootCause   string
	Remediation string
	ActionCmd   string
}

func NewProvider(model string) (*LLMProvider, error) {
	client, err := NewClient(model)
	if err != nil {
		return nil, err
	}

	return &LLMProvider{
		client: client,
		model:  model,
	}, nil
}

func (p *LLMProvider) Ask(ctx context.Context, podCtx *observer.PodContext) (*LLMResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	if err := p.client.CheckConnection(ctx); err != nil {
		return nil, err
	}

	prompt := BuildPrompt(podCtx)
	systemPrompt := GetSystemPrompt()

	response, err := p.client.Ask(ctx, prompt, systemPrompt)
	if err != nil {
		return nil, fmt.Errorf("LLM request failed: %w", err)
	}

	return ParseResponse(response), nil
}

func ParseResponse(response string) *LLMResponse {
	result := &LLMResponse{}

	lines := strings.Split(response, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		switch {
		case strings.HasPrefix(line, "PROBLEM:"):
			result.RootCause = strings.TrimSpace(strings.TrimPrefix(line, "PROBLEM:"))
		case strings.HasPrefix(line, "REASON:"):
			result.Remediation = strings.TrimSpace(strings.TrimPrefix(line, "REASON:"))
		case strings.HasPrefix(line, "SOLUTION:"):
			result.Remediation += "\n" + strings.TrimSpace(strings.TrimPrefix(line, "SOLUTION:"))
		case strings.HasPrefix(line, "COMMAND:"):
			result.ActionCmd = strings.TrimSpace(strings.TrimPrefix(line, "COMMAND:"))
		}
	}

	return result
}
