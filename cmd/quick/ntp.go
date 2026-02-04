package quick

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tehnerd/vape/internal/models"
)

var ntpCmd = &cobra.Command{
	Use:   "ntp <target>",
	Short: "Run a quick NTP measurement",
	Long: `Run a quick NTP measurement to check time servers.

Examples:
  vape quick ntp pool.ntp.org
  vape quick ntp time.google.com --probes 20
  vape quick ntp time.cloudflare.com --wait`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		target := args[0]
		description := fmt.Sprintf("VAPE quick NTP check for %s", target)

		return createAndRun(models.TypeNTP, target, description, nil)
	},
}
