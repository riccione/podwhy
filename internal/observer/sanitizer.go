package observer

import (
	"bytes"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
)

func SanitizePod(pod *v1.Pod) (*v1.Pod, error) {
	pod = pod.DeepCopy()

	pod.ManagedFields = nil
	pod.OwnerReferences = nil

	cleanConditions := []v1.PodCondition{}
	for _, c := range pod.Status.Conditions {
		if c.Status != v1.ConditionTrue {
			cleanConditions = append(cleanConditions, c)
		}
	}
	pod.Status.Conditions = cleanConditions

	return pod, nil
}

func PodToYAML(pod *v1.Pod) (string, error) {
	obj, err := runtime.DefaultUnstructuredConverter.ToUnstructured(pod)
	if err != nil {
		return "", err
	}

	delete(obj, "managedFields")

	var buf bytes.Buffer
	enc := yaml.NewYAMLOrJSONEncoder(&buf, yaml.SimpleMetaFactory{})
	if err := enc.Encode(obj); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func IsHealthy(warnings []string) bool {
	return len(warnings) == 0
}

func NewPodContext(pod *v1.Pod, warnings []string, logs string) *PodContext {
	cleanPod, _ := SanitizePod(pod)
	yamlStr, _ := PodToYAML(cleanPod)

	return &PodContext{
		PodName:     pod.Name,
		Namespace:   pod.Namespace,
		CleanedYAML: yamlStr,
		Events:      warnings,
		Logs:        logs,
	}
}
