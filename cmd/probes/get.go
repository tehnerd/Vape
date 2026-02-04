package probes

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/tehnerd/vape/internal/api"
	"github.com/tehnerd/vape/internal/output"
)

var getCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "Get probe details",
	Long: `Get details of a specific probe.

Examples:
  vape probes get 12345`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid probe ID: %s", args[0])
		}

		client := api.NewClient()
		probe, err := client.GetProbe(id)
		if err != nil {
			return fmt.Errorf("failed to get probe: %w", err)
		}

		formatter := output.GetFormatter()
		return formatter.Format(probe, os.Stdout)
	},
}
