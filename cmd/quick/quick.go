package quick

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/tehnerd/vape/internal/api"
	"github.com/tehnerd/vape/internal/config"
	"github.com/tehnerd/vape/internal/models"
	"github.com/tehnerd/vape/internal/output"
)

var (
	quickAF         int
	quickProbes     int
	quickProbeType  string
	quickProbeValue string
	quickWait       bool
	quickTimeout    int
)

var QuickCmd = &cobra.Command{
	Use:   "quick",
	Short: "Run quick one-off measurements",
	Long: `Run quick one-off measurements with sensible defaults.

These commands create a measurement, optionally wait for results, and display them.

Examples:
  vape quick ping google.com
  vape quick traceroute 8.8.8.8 --probes 5
  vape quick dns 8.8.8.8 example.com
  vape quick http example.com --wait`,
}

func init() {
	QuickCmd.PersistentFlags().IntVar(&quickAF, "af", config.GetDefaultAF(), "Address family (4 or 6)")
	QuickCmd.PersistentFlags().IntVar(&quickProbes, "probes", config.GetDefaultProbes(), "Number of probes")
	QuickCmd.PersistentFlags().StringVar(&quickProbeType, "probe-type", "area", "Probe selection type")
	QuickCmd.PersistentFlags().StringVar(&quickProbeValue, "probe-value", "WW", "Probe selection value")
	QuickCmd.PersistentFlags().BoolVar(&quickWait, "wait", false, "Wait for results")
	QuickCmd.PersistentFlags().IntVar(&quickTimeout, "timeout", 120, "Timeout in seconds when waiting for results")

	QuickCmd.AddCommand(pingCmd)
	QuickCmd.AddCommand(tracerouteCmd)
	QuickCmd.AddCommand(dnsCmd)
	QuickCmd.AddCommand(sslCmd)
	QuickCmd.AddCommand(httpCmd)
	QuickCmd.AddCommand(ntpCmd)
}

func createAndRun(msmType models.MeasurementType, target string, description string, extraDef func(*models.MeasurementDefinition)) error {
	definition := models.MeasurementDefinition{
		Type:        msmType,
		Target:      target,
		AF:          quickAF,
		Description: description,
	}

	if extraDef != nil {
		extraDef(&definition)
	}

	probeSpec := models.ProbeSpec{
		Requested: quickProbes,
		Type:      quickProbeType,
		Value:     quickProbeValue,
	}

	req := &models.MeasurementRequest{
		Definitions: []models.MeasurementDefinition{definition},
		Probes:      []models.ProbeSpec{probeSpec},
		IsOneoff:    true,
	}

	client := api.NewClient()
	resp, err := client.CreateMeasurement(req)
	if err != nil {
		return fmt.Errorf("failed to create measurement: %w", err)
	}

	if len(resp.Measurements) == 0 {
		return fmt.Errorf("no measurement created")
	}

	msmID := resp.Measurements[0]
	fmt.Printf("Created measurement %d\n", msmID)

	if !quickWait {
		fmt.Printf("Run 'vape measurements results %d' to view results\n", msmID)
		return nil
	}

	fmt.Println("Waiting for results...")
	return waitForResults(client, msmID)
}

func waitForResults(client *api.Client, msmID int) error {
	timeout := time.Duration(quickTimeout) * time.Second
	start := time.Now()
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if time.Since(start) > timeout {
				return fmt.Errorf("timeout waiting for results")
			}

			results, err := client.GetMeasurementResults(msmID, 0, 0, nil, 0)
			if err != nil {
				continue
			}

			if len(results) > 0 {
				fmt.Printf("\nReceived %d results:\n\n", len(results))
				formatter := output.GetFormatter()
				return formatter.Format(results, os.Stdout)
			}

			fmt.Print(".")
		}
	}
}
