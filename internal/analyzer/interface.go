package analyzer

import (
	"context"

	"github.com/podwhy/podwhy/internal/observer"
)

type Diagnosis struct {
	Source      string
	RootCause   string
	Remediation string
	ActionCmd   string
}

type Analyzer interface {
	Diagnose(ctx context.Context, podCtx *observer.PodContext) (*Diagnosis, error)
}

type RuleEngine interface {
	Check(podCtx *observer.PodContext) (*Diagnosis, bool)
}

type LLMEngine interface {
	Ask(ctx context.Context, podCtx *observer.PodContext) (*Diagnosis, error)
}
