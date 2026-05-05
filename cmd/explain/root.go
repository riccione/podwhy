package explain

import (
	"context"
	"fmt"
	"os"

	"github.com/podwhy/podwhy/internal/analyzer"
	"github.com/podwhy/podwhy/internal/observer"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "explain [pod-name]",
	Short: "Explain why a Kubernetes pod is having issues",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		podName := args[0]
		namespace, _ := cmd.Flags().GetString("namespace")
		debug, _ := cmd.Flags().GetBool("debug")
		model, _ := cmd.Flags().GetString("model")

		client, err := observer.NewClient()
		if err != nil {
			return fmt.Errorf("failed to create k8s client: %w", err)
		}

		ctx := context.Background()
		pod, err := observer.FetchPod(ctx, client, podName, namespace)
		if err != nil {
			return err
		}

		warnings, err := observer.FetchEvents(ctx, client, pod)
		if err != nil {
			return fmt.Errorf("failed to fetch events: %w", err)
		}

		podCtx := observer.NewPodContext(pod, warnings, "")

		if observer.IsHealthy(warnings) {
			fmt.Println("Pod is healthy - no warnings detected")
			return nil
		}

		engine, err := analyzer.NewAnalyzer(model)
		if err != nil {
			return fmt.Errorf("failed to create analyzer: %w", err)
		}

		diagnosis, err := engine.Diagnose(ctx, podCtx)
		if err != nil {
			return fmt.Errorf("diagnosis failed: %w", err)
		}

		fmt.Printf("Source: %s\n", diagnosis.Source)
		fmt.Printf("Root Cause: %s\n", diagnosis.RootCause)
		fmt.Printf("Remediation: %s\n", diagnosis.Remediation)
		if diagnosis.ActionCmd != "" {
			fmt.Printf("Action: %s\n", diagnosis.ActionCmd)
		}

		if debug {
			fmt.Println("\n--- Cleaned Pod YAML ---")
			fmt.Println(podCtx.CleanedYAML)
			fmt.Println("--- Warning Events ---")
			for _, w := range podCtx.Events {
				fmt.Println(w)
			}
		}

		return nil
	},
}

func init() {
	rootCmd.Flags().StringP("namespace", "n", "default", "Kubernetes namespace")
	rootCmd.Flags().BoolP("debug", "d", false, "Print cleaned YAML and events")
	rootCmd.Flags().StringP("model", "m", "phi3", "LLM model name")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
