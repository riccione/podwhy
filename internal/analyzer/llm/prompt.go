package llm

import (
	"fmt"

	"github.com/podwhy/podwhy/internal/observer"
)

const systemPrompt = `You are a Kubernetes SRE. Analyze the provided YAML and Events. You must output your response in the following format:
PROBLEM: <Short description>
REASON: <Detailed root cause>
SOLUTION: <Step-by-step fix>
COMMAND: <The exact kubectl patch or edit command>`

func GetSystemPrompt() string {
	return systemPrompt
}

func BuildPrompt(podCtx *observer.PodContext) string {
	return fmt.Sprintf(`Analyze this Kubernetes pod issue and provide a diagnosis.

Pod: %s
Namespace: %s

Pod YAML:
%s

Warning Events:
%s

Container Logs (last 10 lines):
%s`,
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
