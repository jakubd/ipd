package cmd

import (
	"github.com/jakubd/ipd"
	"github.com/spf13/cobra"
)

var resolveDomains = false

var rootCmd = &cobra.Command{
	Use:   "ip <any ip>",
	Short: "get info on a single ip",
	Long:  "get information on a single ip",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		LookupOne(args[0])
	},
}

func LookupOne(ip string) {
	ipd.OutputLookup(ip, false, resolveDomains)
}

func Execute() error {
	cobra.MinimumNArgs(1)
	return rootCmd.Execute()
}

func init() {
	rootCmd.Flags().BoolVarP(&resolveDomains, "resolve", "r", false, "resolve domains and urls")
	rootCmd.AddCommand(pipeCmd)
}
