package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(SessionCmd)
}

var SessionCmd = &cobra.Command{
	Use:   "session",
	Short: "manage current sessions",
}