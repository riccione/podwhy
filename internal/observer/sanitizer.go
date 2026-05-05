package observer

import (
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/json"
	"sigs.k8s.io/yaml"
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

	data, err := json.Marshal(obj)
	if err != nil {
		return "", err
	}

	yamlBytes, err := yaml.JSONToYAML(data)
	if err != nil {
		return "", err
	}

	return string(yamlBytes), nil
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
