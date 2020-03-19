package cmd

import (
	"github.com/spf13/cobra"
	"oauth2-cli/internal"
)

func init() {
	SessionCmd.AddCommand(SessionGetCmd)
}

var SessionGetCmd = &cobra.Command{
	Use:   "get [sessionId]",
	Short: "displays a session",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		getSession(args[0])
	},
}

func getSession(id string) {
	session, err := internal.LoadSession(id)
	if err != nil {
		fatalOnErr(err)
	}

	if internal.SessionIsExpired(session) {
		internal.PrintWarning("Expired")
	}
	printJson(session)
}

