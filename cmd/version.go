package cmd

import (
    "fmt"
    "github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
    Use:   "version",
    Short: "shows the current version + build information",
    Run: func(cmd *cobra.Command, args []string) {
        fmt.Printf("Version:    %s\n", Version)
        fmt.Printf("Git Hash:   %s\n", Commit)
        fmt.Printf("Build Time: %s\n", Date)
    },
}

func init() {
    RootCmd.AddCommand(versionCmd)
}