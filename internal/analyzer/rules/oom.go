package rules

import (
	"fmt"

	"github.com/podwhy/podwhy/internal/observer"
)

func checkOOM(podCtx *observer.PodContext) (*RuleResult, bool) {
	for _, event := range podCtx.Events {
		if contains(event, "OOMKilled") || contains(event, "OutOfMemory") {
			return &RuleResult{
				RootCause:   "Container was terminated due to OOMKilled",
				Remediation: "Increase memory limits/requests for the container",
				ActionCmd:   fmt.Sprintf("kubectl set resources pod %s -n %s --limits=memory=<new-limit>", podCtx.PodName, podCtx.Namespace),
			}, true
		}
	}
	return nil, false
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && contains(s[1:], substr))
}
