package llm

import "context"

type MockClient struct {
	Response string
}

func (m *MockClient) Ask(ctx context.Context, prompt string) (string, error) {
	return m.Response, nil
}
