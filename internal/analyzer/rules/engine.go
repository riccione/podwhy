package rules

import (
	"github.com/podwhy/podwhy/internal/analyzer"
	"github.com/podwhy/podwhy/internal/observer"
)

type Engine struct {
	rules []func(*observer.PodContext) (*analyzer.Diagnosis, bool)
}

func NewEngine() *Engine {
	return &Engine{
		rules: []func(*observer.PodContext) (*analyzer.Diagnosis, bool){
			checkOOM,
			checkImagePull,
			checkCrashLoop,
			checkScheduling,
		},
	}
}

func (e *Engine) Check(podCtx *observer.PodContext) (*analyzer.Diagnosis, bool) {
	for _, rule := range e.rules {
		if diagnosis, found := rule(podCtx); found {
			return diagnosis, true
		}
	}
	return nil, false
}
