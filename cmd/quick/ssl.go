package quick

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tehnerd/vape/internal/models"
)

var (
	sslPort     int
	sslHostname string
)

var sslCmd = &cobra.Command{
	Use:     "ssl <target>",
	Aliases: []string{"tls", "cert"},
	Short:   "Run a quick SSL certificate measurement",
	Long: `Run a quick SSL certificate measurement to check TLS certificates.

Examples:
  vape quick ssl example.com
  vape quick ssl example.com --port 8443
  vape quick ssl 1.2.3.4 --hostname example.com
  vape quick cert google.com --wait`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		target := args[0]
		description := fmt.Sprintf("VAPE quick SSL check for %s", target)

		return createAndRun(models.TypeSSL, target, description, func(def *models.MeasurementDefinition) {
			if sslPort > 0 {
				def.Port = sslPort
			}
			def.Hostname = sslHostname
		})
	},
}

func init() {
	sslCmd.Flags().IntVar(&sslPort, "port", 443, "Target port")
	sslCmd.Flags().StringVar(&sslHostname, "hostname", "", "SNI server name sent in TLS handshake (target is the connect IP/host)")
}
