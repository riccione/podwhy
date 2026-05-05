package observer

import (
	"context"
	"fmt"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func FetchPod(ctx context.Context, client kubernetes.Interface, name, namespace string) (*v1.Pod, error) {
	pod, err := client.CoreV1().Pods(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, fmt.Errorf("pod %q not found in namespace %q", name, namespace)
		}
		return nil, err
	}
	return pod, nil
}

func FetchEvents(ctx context.Context, client kubernetes.Interface, pod *v1.Pod) ([]string, error) {
	opts := metav1.ListOptions{
		FieldSelector: fmt.Sprintf("involvedObject.name=%s", pod.Name),
	}

	events, err := client.CoreV1().Events(pod.Namespace).List(ctx, opts)
	if err != nil {
		return nil, err
	}

	var warnings []string
	for _, event := range events.Items {
		if event.Type == v1.EventTypeWarning {
			warnings = append(warnings, event.Message)
		}
	}

	return warnings, nil
}
