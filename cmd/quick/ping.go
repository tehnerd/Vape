package quick

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tehnerd/vape/internal/models"
)

var pingPackets int
var pingSize int

var pingCmd = &cobra.Command{
	Use:   "ping <target>",
	Short: "Run a quick ping measurement",
	Long: `Run a quick ping measurement to a target.

Examples:
  vape quick ping google.com
  vape quick ping 8.8.8.8 --probes 20
  vape quick ping example.com --af 6 --wait`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		target := args[0]
		description := fmt.Sprintf("VAPE quick ping to %s", target)

		return createAndRun(models.TypePing, target, description, func(def *models.MeasurementDefinition) {
			if pingPackets > 0 {
				def.Packets = pingPackets
			}
			if pingSize > 0 {
				def.Size = pingSize
			}
		})
	},
}

func init() {
	pingCmd.Flags().IntVar(&pingPackets, "packets", 3, "Number of packets")
	pingCmd.Flags().IntVar(&pingSize, "size", 48, "Packet size in bytes")
}
