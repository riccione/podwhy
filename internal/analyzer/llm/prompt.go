package llm

import (
	"fmt"

	"github.com/podwhy/podwhy/internal/observer"
)

func BuildPrompt(podCtx *observer.PodContext) string {
	return fmt.Sprintf(`Analyze this Kubernetes pod issue and provide a diagnosis.

Pod: %s
Namespace: %s

Pod YAML:
%s

Warning Events:
%s

Container Logs (last 10 lines):
%s

Provide a JSON response with:
- root_cause: Brief description of the issue
- remediation: How to fix it
- action_cmd: kubectl command to apply the fix

Response:`,
		podCtx.PodName,
		podCtx.Namespace,
		podCtx.CleanedYAML,
		formatEvents(podCtx.Events),
		podCtx.Logs,
	)
}

func formatEvents(events []string) string {
	if len(events) == 0 {
		return "No warning events"
	}
	result := ""
	for _, e := range events {
		result += "- " + e + "\n"
	}
	return result
}
