package analyzer

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/podwhy/podwhy/internal/analyzer/llm"
	"github.com/podwhy/podwhy/internal/analyzer/rules"
	"github.com/podwhy/podwhy/internal/observer"
)

type Engine struct {
	ruleEngine *rules.Engine
	llmEngine  *llm.Client
}

func NewAnalyzer(model string) (*Engine, error) {
	llmClient, err := llm.NewClient(model)
	if err != nil {
		return nil, fmt.Errorf("failed to create LLM client: %w", err)
	}

	return &Engine{
		ruleEngine: rules.NewEngine(),
		llmEngine:  llmClient,
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

	llmDiagnosis, err := e.askLLM(ctx, podCtx)
	if err != nil {
		return nil, fmt.Errorf("LLM failed: %w", err)
	}

	llmDiagnosis.Source = "LLMEngine"
	return llmDiagnosis, nil
}

func (e *Engine) askLLM(ctx context.Context, podCtx *observer.PodContext) (*Diagnosis, error) {
	prompt := llm.BuildPrompt(podCtx)

	response, err := e.llmEngine.Ask(ctx, prompt)
	if err != nil {
		return nil, err
	}

	var result struct {
		RootCause   string `json:"root_cause"`
		Remediation string `json:"remediation"`
		ActionCmd   string `json:"action_cmd"`
	}

	if err := json.Unmarshal([]byte(response), &result); err != nil {
		return &Diagnosis{
			RootCause:   "LLM analysis completed",
			Remediation: response,
			ActionCmd:   "",
		}, nil
	}

	return &Diagnosis{
		RootCause:   result.RootCause,
		Remediation: result.Remediation,
		ActionCmd:   result.ActionCmd,
	}, nil
}
