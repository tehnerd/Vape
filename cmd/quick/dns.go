package quick

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tehnerd/vape/internal/models"
)

var (
	dnsQueryType      string
	dnsQueryClass     string
	dnsUseResolver    bool
)

var dnsCmd = &cobra.Command{
	Use:   "dns <resolver> <query>",
	Short: "Run a quick DNS measurement",
	Long: `Run a quick DNS measurement.

The first argument is the DNS resolver to query (use 'probe' to use probe's resolver).
The second argument is the domain name to query.

Examples:
  vape quick dns 8.8.8.8 example.com
  vape quick dns 1.1.1.1 google.com --query-type AAAA
  vape quick dns probe example.com --query-type MX`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		resolver := args[0]
		query := args[1]

		useProbeResolver := resolver == "probe"
		target := resolver
		if useProbeResolver {
			target = ""
		}

		description := fmt.Sprintf("VAPE quick DNS query for %s", query)

		return createAndRun(models.TypeDNS, target, description, func(def *models.MeasurementDefinition) {
			def.QueryArgument = query
			def.QueryType = dnsQueryType
			def.QueryClass = dnsQueryClass
			def.UseProbeResolver = useProbeResolver
		})
	},
}

func init() {
	dnsCmd.Flags().StringVar(&dnsQueryType, "query-type", "A", "DNS query type (A, AAAA, MX, TXT, etc.)")
	dnsCmd.Flags().StringVar(&dnsQueryClass, "query-class", "IN", "DNS query class")
	dnsCmd.Flags().BoolVar(&dnsUseResolver, "use-resolver", false, "Use probe's resolver (deprecated, use 'probe' as resolver)")
}
