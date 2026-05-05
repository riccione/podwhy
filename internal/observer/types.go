package observer

type PodContext struct {
	PodName     string
	Namespace   string
	CleanedYAML string
	Events      []string
	Logs        string
}
