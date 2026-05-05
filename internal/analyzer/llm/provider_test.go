package llm

import (
	"context"
	"testing"

	"github.com/podwhy/podwhy/internal/observer"
)

func TestMockProvider_Ask(t *testing.T) {
	expected := &LLMResponse{
		RootCause:   "test problem",
		Remediation: "fix it",
		ActionCmd:   "kubectl fix",
	}

	mock := &MockProvider{Response: expected}
	result, err := mock.Ask(context.Background(), &observer.PodContext{})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.RootCause != expected.RootCause {
		t.Errorf("expected %q, got %q", expected.RootCause, result.RootCause)
	}
}

func TestParseResponse(t *testing.T) {
	input := `PROBLEM: Pod OOMKilled
REASON: Container exceeded memory limit
SOLUTION: Increase memory limit
COMMAND: kubectl set resources pod test -n default --limits=memory=512Mi`

	result := ParseResponse(input)

	if result.RootCause != "Pod OOMKilled" {
		t.Errorf("unexpected root cause: %s", result.RootCause)
	}
	if result.ActionCmd != "kubectl set resources pod test -n default --limits=memory=512Mi" {
		t.Errorf("unexpected action cmd: %s", result.ActionCmd)
	}
}
