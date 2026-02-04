package credits

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/tehnerd/vape/internal/api"
	"github.com/tehnerd/vape/internal/output"
)

var CreditsCmd = &cobra.Command{
	Use:   "credits",
	Short: "Show credit balance",
	Long: `Show your RIPE Atlas credit balance and estimates.

Examples:
  vape credits`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client := api.NewClient()
		credits, err := client.GetCredits()
		if err != nil {
			return fmt.Errorf("failed to get credits: %w", err)
		}

		formatter := output.GetFormatter()
		return formatter.Format(credits, os.Stdout)
	},
}
