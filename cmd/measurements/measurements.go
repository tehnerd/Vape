package measurements

import (
	"github.com/spf13/cobra"
)

var MeasurementsCmd = &cobra.Command{
	Use:     "measurements",
	Aliases: []string{"msm", "m"},
	Short:   "Manage measurements",
	Long:    "Create, list, get, and manage RIPE Atlas measurements",
}

func init() {
	MeasurementsCmd.AddCommand(createCmd)
	MeasurementsCmd.AddCommand(listCmd)
	MeasurementsCmd.AddCommand(getCmd)
	MeasurementsCmd.AddCommand(resultsCmd)
	MeasurementsCmd.AddCommand(stopCmd)
}
