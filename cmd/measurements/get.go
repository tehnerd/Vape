package measurements

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
	Short: "Get measurement details",
	Long: `Get details of a specific measurement.

Examples:
  vape measurements get 12345678`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid measurement ID: %s", args[0])
		}

		client := api.NewClient()
		measurement, err := client.GetMeasurement(id)
		if err != nil {
			return fmt.Errorf("failed to get measurement: %w", err)
		}

		formatter := output.GetFormatter()
		return formatter.Format(measurement, os.Stdout)
	},
}
