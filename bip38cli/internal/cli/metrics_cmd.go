package cli

import (
	"encoding/json"
	"fmt"

	"github.com/carlosrabelo/bip38cli/bip38cli/internal/metrics"
	"github.com/spf13/cobra"
)

var metricsCmd = &cobra.Command{
	Use:   "metrics",
	Short: "Show operation metrics",
	Long: `Display collected metrics for encrypt, decrypt, and intermediate operations.

Includes operation counts, average durations, and success rates.

Examples:
  bip38cli metrics
  bip38cli metrics --output-format json`,
	RunE: runMetrics,
}

func init() {
	rootCmd.AddCommand(metricsCmd)
}

func runMetrics(cmd *cobra.Command, _ []string) error {
	snap := metrics.GetSnapshot()

	switch outputFormat(cmd) {
	case "json":
		out, err := json.MarshalIndent(&snap, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal metrics: %v", err)
		}
		fmt.Println(string(out))
	default:
		fmt.Printf("Uptime:               %s\n", snap.Uptime.Round(1000000))
		fmt.Println()
		fmt.Printf("Encrypt  - count: %d  errors: %d  avg: %s  success: %.2f%%\n",
			snap.EncryptCount, snap.EncryptErrors, snap.AverageEncryptTime.Round(1000000),
			metrics.GetMetrics().EncryptSuccessRate())
		fmt.Printf("Decrypt  - count: %d  errors: %d  avg: %s  success: %.2f%%\n",
			snap.DecryptCount, snap.DecryptErrors, snap.AverageDecryptTime.Round(1000000),
			metrics.GetMetrics().DecryptSuccessRate())
		fmt.Printf("Intermediate - count: %d  errors: %d  avg: %s  success: %.2f%%\n",
			snap.IntermediateCount, snap.IntermediateErrors, snap.AverageIntermediateTime.Round(1000000),
			metrics.GetMetrics().IntermediateSuccessRate())
	}

	return nil
}
