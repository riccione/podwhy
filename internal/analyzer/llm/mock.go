package llm

import (
	"context"

	"github.com/podwhy/podwhy/internal/observer"
)

type MockProvider struct {
	Response *LLMResponse
	Err      error
}

func (m *MockProvider) Ask(ctx context.Context, podCtx *observer.PodContext) (*LLMResponse, error) {
	return m.Response, m.Err
}
