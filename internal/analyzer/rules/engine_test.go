package rules

import (
	"testing"

	"github.com/podwhy/podwhy/internal/observer"
)

func TestCheckOOM(t *testing.T) {
	podCtx := &observer.PodContext{
		PodName:   "test-pod",
		Namespace: "default",
		Events:    []string{"Container killed due to OOMKilled"},
	}

	result, found := checkOOM(podCtx)
	if !found {
		t.Error("expected OOM to be detected")
	}
	if result.RootCause == "" {
		t.Error("expected root cause to be set")
	}
}

func TestCheckImagePull(t *testing.T) {
	podCtx := &observer.PodContext{
		PodName:   "test-pod",
		Namespace: "default",
		Events:    []string{"Failed to pull image: ImagePullBackOff"},
	}

	result, found := checkImagePull(podCtx)
	if !found {
		t.Error("expected ImagePullBackOff to be detected")
	}
	if result.ActionCmd == "" {
		t.Error("expected action command to be set")
	}
}

func TestCheckCrashLoop(t *testing.T) {
	podCtx := &observer.PodContext{
		PodName:   "test-pod",
		Namespace: "default",
		Events:    []string{"Back-off restarting failed container: CrashLoopBackOff"},
	}

	result, found := checkCrashLoop(podCtx)
	if !found {
		t.Error("expected CrashLoopBackOff to be detected")
	}
	_ = result
}

func TestCheckScheduling(t *testing.T) {
	podCtx := &observer.PodContext{
		PodName:   "test-pod",
		Namespace: "default",
		Events:    []string{"0/3 nodes are available: 3 Insufficient cpu: FailedScheduling"},
	}

	result, found := checkScheduling(podCtx)
	if !found {
		t.Error("expected FailedScheduling to be detected")
	}
	_ = result
}

func TestEngine_Check(t *testing.T) {
	engine := NewEngine()

	podCtx := &observer.PodContext{
		PodName:   "test-pod",
		Namespace: "default",
		Events:    []string{"Container killed due to OOMKilled"},
	}

	result, found := engine.Check(podCtx)
	if !found {
		t.Error("expected rule to match")
	}
	if result == nil {
		t.Error("expected result to be non-nil")
	}
}
