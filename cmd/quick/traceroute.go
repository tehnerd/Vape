package quick

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tehnerd/vape/internal/models"
)

var traceProtocol string
var traceMaxHops int

var tracerouteCmd = &cobra.Command{
	Use:     "traceroute <target>",
	Aliases: []string{"tr", "trace"},
	Short:   "Run a quick traceroute measurement",
	Long: `Run a quick traceroute measurement to a target.

Examples:
  vape quick traceroute google.com
  vape quick traceroute 8.8.8.8 --protocol UDP
  vape quick trace example.com --probes 10 --wait`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		target := args[0]
		description := fmt.Sprintf("VAPE quick traceroute to %s", target)

		return createAndRun(models.TypeTraceroute, target, description, func(def *models.MeasurementDefinition) {
			if traceProtocol != "" {
				def.Protocol = traceProtocol
			}
		})
	},
}

func init() {
	tracerouteCmd.Flags().StringVar(&traceProtocol, "protocol", "ICMP", "Protocol (ICMP, UDP, TCP)")
	tracerouteCmd.Flags().IntVar(&traceMaxHops, "max-hops", 32, "Maximum number of hops")
}
