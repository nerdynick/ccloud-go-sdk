package cmd

import (
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List information from the API",
}

func init() {
	rootCmd.AddCommand(listCmd)
}
