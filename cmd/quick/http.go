package quick

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tehnerd/vape/internal/models"
)

var (
	httpMethod string
	httpPath   string
	httpPort   int
)

var httpCmd = &cobra.Command{
	Use:   "http <target>",
	Short: "Run a quick HTTP measurement",
	Long: `Run a quick HTTP measurement to fetch a URL.

Examples:
  vape quick http example.com
  vape quick http example.com --path /api/health
  vape quick http example.com --method HEAD --port 8080`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		target := args[0]
		description := fmt.Sprintf("VAPE quick HTTP request to %s", target)

		return createAndRun(models.TypeHTTP, target, description, func(def *models.MeasurementDefinition) {
			def.Method = httpMethod
			def.Path = httpPath
			if httpPort > 0 {
				def.Port = httpPort
			}
		})
	},
}

func init() {
	httpCmd.Flags().StringVar(&httpMethod, "method", "GET", "HTTP method")
	httpCmd.Flags().StringVar(&httpPath, "path", "/", "HTTP path")
	httpCmd.Flags().IntVar(&httpPort, "port", 80, "Target port")
}
