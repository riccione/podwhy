package analyzer

import (
	"context"
	"fmt"

	"github.com/podwhy/podwhy/internal/analyzer/llm"
	"github.com/podwhy/podwhy/internal/analyzer/rules"
	"github.com/podwhy/podwhy/internal/observer"
)

type Engine struct {
	ruleEngine *rules.Engine
	llmEngine  *llm.LLMProvider
}

func NewAnalyzer(model string) (*Engine, error) {
	llmProvider, err := llm.NewProvider(model)
	if err != nil {
		return nil, fmt.Errorf("failed to create LLM provider: %w", err)
	}

	return &Engine{
		ruleEngine: rules.NewEngine(),
		llmEngine:  llmProvider,
	}, nil
}

func (e *Engine) Diagnose(ctx context.Context, podCtx *observer.PodContext) (*Diagnosis, error) {
	if result, found := e.ruleEngine.Check(podCtx); found {
		return &Diagnosis{
			Source:      "RuleEngine",
			RootCause:   result.RootCause,
			Remediation: result.Remediation,
			ActionCmd:   result.ActionCmd,
		}, nil
	}

	llmResult, err := e.llmEngine.Ask(ctx, podCtx)
	if err != nil {
		return nil, fmt.Errorf("LLM failed: %w", err)
	}

	return &Diagnosis{
		Source:      "LLMEngine",
		RootCause:   llmResult.RootCause,
		Remediation: llmResult.Remediation,
		ActionCmd:   llmResult.ActionCmd,
	}, nil
}
