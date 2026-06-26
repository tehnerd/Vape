package measurements

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/tehnerd/vape/internal/api"
	"github.com/tehnerd/vape/internal/config"
	"github.com/tehnerd/vape/internal/models"
	"github.com/tehnerd/vape/internal/output"
)

var (
	createType        string
	createTarget      string
	createAF          int
	createProbes      int
	createDescription string
	createOneoff      bool
	createProbeType   string
	createProbeValue  string
	createProtocol    string
	// DNS specific
	createQueryClass    string
	createQueryType     string
	createQueryArgument string
	createUseResolver   bool
	// HTTP specific
	createMethod string
	createPath   string
	// SSL/HTTP specific
	createPort     int
	createHostname string
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new measurement",
	Long: `Create a new RIPE Atlas measurement.

Examples:
  vape measurements create --type ping --target google.com
  vape measurements create --type traceroute --target 8.8.8.8 --probes 20 --af 4
  vape measurements create --type dns --target 8.8.8.8 --query-argument example.com --query-type A
  vape measurements create --type http --target example.com --port 443
  vape measurements create --type sslcert --target 1.2.3.4 --hostname example.com --port 443`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if createType == "" {
			return fmt.Errorf("--type is required")
		}
		if createTarget == "" && createType != "dns" {
			return fmt.Errorf("--target is required")
		}

		definition := models.MeasurementDefinition{
			Type:        models.MeasurementType(createType),
			Target:      createTarget,
			AF:          createAF,
			Description: createDescription,
		}

		// Set type-specific options
		switch createType {
		case "ping", "traceroute":
			if createProtocol != "" {
				definition.Protocol = createProtocol
			}
		case "dns":
			definition.QueryClass = createQueryClass
			definition.QueryType = createQueryType
			definition.QueryArgument = createQueryArgument
			definition.UseProbeResolver = createUseResolver
		case "http":
			definition.Method = createMethod
			definition.Path = createPath
			if createPort > 0 {
				definition.Port = createPort
			}
		case "sslcert":
			if createPort > 0 {
				definition.Port = createPort
			}
			definition.Hostname = createHostname
		}

		probeSpec := models.ProbeSpec{
			Requested: createProbes,
			Type:      createProbeType,
			Value:     createProbeValue,
		}

		req := &models.MeasurementRequest{
			Definitions: []models.MeasurementDefinition{definition},
			Probes:      []models.ProbeSpec{probeSpec},
			IsOneoff:    createOneoff,
		}

		client := api.NewClient()
		resp, err := client.CreateMeasurement(req)
		if err != nil {
			return fmt.Errorf("failed to create measurement: %w", err)
		}

		formatter := output.GetFormatter()
		if _, ok := formatter.(*output.JSONFormatter); ok {
			formatter.Format(resp, os.Stdout)
		} else {
			fmt.Printf("Created measurement(s): %v\n", resp.Measurements)
			for _, id := range resp.Measurements {
				fmt.Printf("View results: vape measurements results %d\n", id)
			}
		}

		return nil
	},
}

func init() {
	createCmd.Flags().StringVarP(&createType, "type", "t", "", "Measurement type (ping, traceroute, dns, sslcert, ntp, http)")
	createCmd.Flags().StringVar(&createTarget, "target", "", "Target hostname or IP")
	createCmd.Flags().IntVar(&createAF, "af", config.GetDefaultAF(), "Address family (4 or 6)")
	createCmd.Flags().IntVar(&createProbes, "probes", config.GetDefaultProbes(), "Number of probes")
	createCmd.Flags().StringVar(&createDescription, "description", "", "Measurement description")
	createCmd.Flags().BoolVar(&createOneoff, "oneoff", true, "Run as one-off measurement")
	createCmd.Flags().StringVar(&createProbeType, "probe-type", "area", "Probe selection type (area, country, asn, prefix, msm)")
	createCmd.Flags().StringVar(&createProbeValue, "probe-value", "WW", "Probe selection value")
	createCmd.Flags().StringVar(&createProtocol, "protocol", "", "Protocol (ICMP, UDP, TCP)")
	// DNS flags
	createCmd.Flags().StringVar(&createQueryClass, "query-class", "IN", "DNS query class")
	createCmd.Flags().StringVar(&createQueryType, "query-type", "A", "DNS query type (A, AAAA, MX, etc.)")
	createCmd.Flags().StringVar(&createQueryArgument, "query-argument", "", "DNS query argument (domain name)")
	createCmd.Flags().BoolVar(&createUseResolver, "use-probe-resolver", false, "Use probe's resolver")
	// HTTP flags
	createCmd.Flags().StringVar(&createMethod, "method", "GET", "HTTP method")
	createCmd.Flags().StringVar(&createPath, "path", "/", "HTTP path")
	// Port flag
	createCmd.Flags().IntVar(&createPort, "port", 0, "Target port (for SSL/HTTP)")
	// SSL flags
	createCmd.Flags().StringVar(&createHostname, "hostname", "", "SNI server name for sslcert (sent in TLS handshake; --target is the connect IP/host)")

	createCmd.MarkFlagRequired("type")
}
