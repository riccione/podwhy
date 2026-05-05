package explain

import (
	"context"
	"fmt"
	"os"

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

		ctx2 := observer.NewPodContext(pod, warnings, "")

		if observer.IsHealthy(warnings) {
			fmt.Println("Pod is healthy - no warnings detected")
			return nil
		}

		if debug {
			fmt.Println("--- Cleaned Pod YAML ---")
			fmt.Println(ctx2.CleanedYAML)
			fmt.Println("--- Warning Events ---")
			for _, w := range ctx2.Events {
				fmt.Println(w)
			}
		}

		return nil
	},
}

func init() {
	rootCmd.Flags().StringP("namespace", "n", "default", "Kubernetes namespace")
	rootCmd.Flags().BoolP("debug", "d", false, "Print cleaned YAML and events")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
