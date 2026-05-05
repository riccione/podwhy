package observer

import (
	"testing"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestSanitizePod_RemovesManagedFields(t *testing.T) {
	pod := &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-pod",
			Namespace: "default",
			ManagedFields: []metav1.ManagedFieldsEntry{
				{Manager: "kubectl"},
			},
		},
		Status: v1.PodStatus{
			Phase: v1.PodRunning,
			Conditions: []v1.PodCondition{
				{Type: v1.PodReady, Status: v1.ConditionTrue},
				{Type: v1.PodScheduled, Status: v1.ConditionFalse},
			},
			ContainerStatuses: []v1.ContainerStatus{
				{Name: "app", Ready: true},
			},
		},
	}

	cleaned, err := SanitizePod(pod)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cleaned.ManagedFields != nil {
		t.Error("expected managedFields to be nil")
	}

	if len(cleaned.Status.Conditions) != 1 {
		t.Errorf("expected 1 condition, got %d", len(cleaned.Status.Conditions))
	}

	yamlStr, err := PodToYAML(cleaned)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if contains(yamlStr, "managedFields") {
		t.Error("output YAML should not contain managedFields")
	}
}

func TestIsHealthy(t *testing.T) {
	if !IsHealthy([]string{}) {
		t.Error("expected healthy with no warnings")
	}
	if IsHealthy([]string{"disk pressure"}) {
		t.Error("expected not healthy with warnings")
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && contains(s[1:], substr))
}
