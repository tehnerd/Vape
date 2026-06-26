package measurements

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tehnerd/vape/internal/api"
	"github.com/tehnerd/vape/internal/output"
)

var (
	resultsStart      int64
	resultsStop       int64
	resultsProbeIDs   string
	resultsLimit      int
	resultsCertFormat string
)

var resultsCmd = &cobra.Command{
	Use:   "results <id>",
	Short: "Get measurement results",
	Long: `Get results of a specific measurement.

Examples:
  vape measurements results 12345678
  vape measurements results 12345678 --limit 10
  vape measurements results 12345678 --probe-ids 1234,5678`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid measurement ID: %s", args[0])
		}

		var probeIDs []int
		if resultsProbeIDs != "" {
			for _, pidStr := range strings.Split(resultsProbeIDs, ",") {
				pid, err := strconv.Atoi(strings.TrimSpace(pidStr))
				if err != nil {
					return fmt.Errorf("invalid probe ID: %s", pidStr)
				}
				probeIDs = append(probeIDs, pid)
			}
		}

		client := api.NewClient()
		results, err := client.GetMeasurementResults(id, resultsStart, resultsStop, probeIDs, resultsLimit)
		if err != nil {
			return fmt.Errorf("failed to get results: %w", err)
		}

		if len(results) == 0 {
			fmt.Println("No results available yet")
			return nil
		}

		formatter := output.GetFormatter()
		if tf, ok := formatter.(*output.TableFormatter); ok {
			tf.SSLCertMode = output.ParseSSLCertMode(resultsCertFormat)
		}
		return formatter.Format(results, os.Stdout)
	},
}

func init() {
	resultsCmd.Flags().Int64Var(&resultsStart, "start", 0, "Start timestamp (Unix epoch)")
	resultsCmd.Flags().Int64Var(&resultsStop, "stop", 0, "Stop timestamp (Unix epoch)")
	resultsCmd.Flags().StringVar(&resultsProbeIDs, "probe-ids", "", "Comma-separated list of probe IDs")
	resultsCmd.Flags().IntVar(&resultsLimit, "limit", 0, "Maximum number of results")
	resultsCmd.Flags().StringVar(&resultsCertFormat, "cert-format", "names", "SSL cert output: 'names' (SAN + cert names) or 'full' (base64 + decoded)")
}
