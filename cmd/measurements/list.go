package measurements

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/tehnerd/vape/internal/api"
	"github.com/tehnerd/vape/internal/output"
)

var (
	listStatus string
	listType   string
	listTarget string
	listMine   bool
	listLimit  int
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List measurements",
	Long: `List RIPE Atlas measurements.

Examples:
  vape measurements list
  vape measurements list --mine
  vape measurements list --status Ongoing --type ping
  vape measurements list --target google.com`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client := api.NewClient()

		params := &api.MeasurementListParams{
			Status: listStatus,
			Type:   listType,
			Target: listTarget,
			Mine:   listMine,
			Limit:  listLimit,
		}

		resp, err := client.ListMeasurements(params)
		if err != nil {
			return fmt.Errorf("failed to list measurements: %w", err)
		}

		formatter := output.GetFormatter()
		if _, ok := formatter.(*output.JSONFormatter); ok {
			return formatter.Format(resp, os.Stdout)
		}

		if len(resp.Results) == 0 {
			fmt.Println("No measurements found")
			return nil
		}

		fmt.Printf("Found %d measurements (showing %d)\n\n", resp.Count, len(resp.Results))
		return formatter.Format(resp.Results, os.Stdout)
	},
}

func init() {
	listCmd.Flags().StringVar(&listStatus, "status", "", "Filter by status (Specified, Scheduled, Ongoing, Stopped)")
	listCmd.Flags().StringVar(&listType, "type", "", "Filter by type (ping, traceroute, dns, sslcert, ntp, http)")
	listCmd.Flags().StringVar(&listTarget, "target", "", "Filter by target")
	listCmd.Flags().BoolVar(&listMine, "mine", false, "Show only my measurements")
	listCmd.Flags().IntVar(&listLimit, "limit", 25, "Maximum number of results")
}
