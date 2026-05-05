package llm

import (
	"context"
	"testing"
)

func TestMockClient_Ask(t *testing.T) {
	mock := &MockClient{Response: `{"root_cause": "test", "remediation": "fix it", "action_cmd": "kubectl fix"}`}

	result, err := mock.Ask(context.Background(), "test prompt")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != mock.Response {
		t.Errorf("expected %q, got %q", mock.Response, result)
	}
}
