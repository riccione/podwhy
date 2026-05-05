package rules

import (
	"fmt"

	"github.com/podwhy/podwhy/internal/analyzer"
	"github.com/podwhy/podwhy/internal/observer"
)

func checkCrashLoop(podCtx *observer.PodContext) (*analyzer.Diagnosis, bool) {
	for _, event := range podCtx.Events {
		if contains(event, "CrashLoopBackOff") {
			return &analyzer.Diagnosis{
				RootCause:   "Container is in CrashLoopBackOff state",
				Remediation: "Check container logs for application errors and ensure the container starts correctly",
				ActionCmd:   fmt.Sprintf("kubectl logs %s -n %s --previous", podCtx.PodName, podCtx.Namespace),
			}, true
		}
	}
	return nil, false
}
