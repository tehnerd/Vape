package measurements

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/tehnerd/vape/internal/api"
)

var stopCmd = &cobra.Command{
	Use:   "stop <id>",
	Short: "Stop a measurement",
	Long: `Stop a running measurement.

Examples:
  vape measurements stop 12345678`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid measurement ID: %s", args[0])
		}

		client := api.NewClient()
		err = client.StopMeasurement(id)
		if err != nil {
			return fmt.Errorf("failed to stop measurement: %w", err)
		}

		fmt.Printf("Measurement %d stopped\n", id)
		return nil
	},
}
