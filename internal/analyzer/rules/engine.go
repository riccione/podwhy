package rules

import (
	"github.com/podwhy/podwhy/internal/observer"
)

type RuleResult struct {
	RootCause   string
	Remediation string
	ActionCmd   string
}

type Engine struct {
	rules []func(*observer.PodContext) (*RuleResult, bool)
}

func NewEngine() *Engine {
	return &Engine{
		rules: []func(*observer.PodContext) (*RuleResult, bool){
			checkOOM,
			checkImagePull,
			checkCrashLoop,
			checkScheduling,
		},
	}
}

func (e *Engine) Check(podCtx *observer.PodContext) (*RuleResult, bool) {
	for _, rule := range e.rules {
		if result, found := rule(podCtx); found {
			return result, true
		}
	}
	return nil, false
}
