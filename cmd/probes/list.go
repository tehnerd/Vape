package probes

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/tehnerd/vape/internal/api"
	"github.com/tehnerd/vape/internal/models"
	"github.com/tehnerd/vape/internal/output"
)

var (
	listCountry  string
	listASN      int
	listStatus   string
	listIsAnchor bool
	listLimit    int
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List probes",
	Long: `List RIPE Atlas probes.

Examples:
  vape probes list
  vape probes list --country US
  vape probes list --asn 15169
  vape probes list --status Connected --limit 50`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client := api.NewClient()

		params := &models.ProbeListParams{
			CountryCode: listCountry,
			ASN:         listASN,
			Limit:       listLimit,
		}

		if listStatus != "" {
			params.Status = models.ProbeStatusName(listStatus)
		}

		if cmd.Flags().Changed("anchor") {
			params.IsAnchor = &listIsAnchor
		}

		resp, err := client.ListProbes(params)
		if err != nil {
			return fmt.Errorf("failed to list probes: %w", err)
		}

		formatter := output.GetFormatter()
		if _, ok := formatter.(*output.JSONFormatter); ok {
			return formatter.Format(resp, os.Stdout)
		}

		if len(resp.Results) == 0 {
			fmt.Println("No probes found")
			return nil
		}

		fmt.Printf("Found %d probes (showing %d)\n\n", resp.Count, len(resp.Results))
		return formatter.Format(resp.Results, os.Stdout)
	},
}

func init() {
	listCmd.Flags().StringVar(&listCountry, "country", "", "Filter by country code (e.g., US, DE, JP)")
	listCmd.Flags().IntVar(&listASN, "asn", 0, "Filter by ASN")
	listCmd.Flags().StringVar(&listStatus, "status", "", "Filter by status (Connected, Disconnected, Abandoned)")
	listCmd.Flags().BoolVar(&listIsAnchor, "anchor", false, "Filter anchors only")
	listCmd.Flags().IntVar(&listLimit, "limit", 25, "Maximum number of results")
}
