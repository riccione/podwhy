package rules

import (
	"fmt"

	"github.com/podwhy/podwhy/internal/analyzer"
	"github.com/podwhy/podwhy/internal/observer"
)

func checkImagePull(podCtx *observer.PodContext) (*analyzer.Diagnosis, bool) {
	for _, event := range podCtx.Events {
		if contains(event, "ImagePullBackOff") || contains(event, "ErrImagePull") {
			return &analyzer.Diagnosis{
				RootCause:   "Image pull failed - ImagePullBackOff or ErrImagePull detected",
				Remediation: "Check image name, tag, and registry credentials",
				ActionCmd:   fmt.Sprintf("kubectl describe pod %s -n %s", podCtx.PodName, podCtx.Namespace),
			}, true
		}
	}
	return nil, false
}
