package probes

import (
	"github.com/spf13/cobra"
)

var ProbesCmd = &cobra.Command{
	Use:     "probes",
	Aliases: []string{"probe", "p"},
	Short:   "Manage probes",
	Long:    "List and get information about RIPE Atlas probes",
}

func init() {
	ProbesCmd.AddCommand(listCmd)
	ProbesCmd.AddCommand(getCmd)
}
