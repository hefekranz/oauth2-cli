package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"oauth2-cli/internal"
	"os"
	"sort"
)

func init() {
	SessionCmd.AddCommand(SessionListCmd)
}

var SessionListCmd = &cobra.Command{
	Use:   "list",
	Short: "list sessions",
	Run: func(cmd *cobra.Command, args []string) {
		files, err := ReadDir(internal.GetSessionDir())
		if err != nil {
			fatalOnErr(err)
		}

		for _, v := range files {
			fmt.Println(v.Name())
		}
	},
}

func ReadDir(dirname string) ([]os.FileInfo, error) {
	f, err := os.Open(dirname)
	if err != nil {
		return nil, err
	}
	list, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		return nil, err
	}
	sort.Slice(list, func(i, j int) bool { return list[i].Name() < list[j].Name() })
	return list, nil
}