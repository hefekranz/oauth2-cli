package cmd

import (
	"github.com/spf13/cobra"
	"oauth2-cli/internal"
	"os"
)

func init() {
	SessionCmd.AddCommand(SessionDeleteCmd)
}

var SessionDeleteCmd = &cobra.Command{
	Use:   "delete [sessionId]",
	Short: "deletes a session",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		deleteSession(args[0])
	},
}

func deleteSession(id string) {
	err := os.Remove(internal.GetSessionFilePath(id))
	if err != nil {
		fatalOnErr(err)
	}
}
