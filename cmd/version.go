package cmd

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/cobra"
)

var (
	Version = "unset"
)

func init() {
	c := &cobra.Command{
		Use:   "version",
		Short: "Show version",
		Run: func(cmd *cobra.Command, args []string) {
			versionCmd()
		},
	}

	rootCmd.AddCommand(c)
}

func versionCmd() {
	fmt.Println(Version)
}
