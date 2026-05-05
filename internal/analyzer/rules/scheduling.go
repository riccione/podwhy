package rules

import (
	"fmt"

	"github.com/podwhy/podwhy/internal/analyzer"
	"github.com/podwhy/podwhy/internal/observer"
)

func checkScheduling(podCtx *observer.PodContext) (*analyzer.Diagnosis, bool) {
	for _, event := range podCtx.Events {
		if contains(event, "FailedScheduling") {
			return &analyzer.Diagnosis{
				RootCause:   "Pod scheduling failed - insufficient resources or node constraints",
				Remediation: "Check node resources, taints, tolerations, or affinity rules",
				ActionCmd:   fmt.Sprintf("kubectl describe pod %s -n %s", podCtx.PodName, podCtx.Namespace),
			}, true
		}
	}
	return nil, false
}
