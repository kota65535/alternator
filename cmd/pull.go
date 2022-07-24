package cmd

import (
	_ "embed"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

//go:embed pull.tmpl
var pullUsage string

func init() {
	c := &cobra.Command{
		Use:   "pull <database-url>",
		Short: "Show the current remote database schemas.",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			PullCmd(args[0])
		},
	}
	rootCmd.AddCommand(c)
	c.SetUsageTemplate(pullUsage)
}

func PullCmd(uri string) {
	dbUri, err := NewDatabaseUri(uri)
	cobra.CheckErr(err)

	alternator, err := NewAlternator(dbUri)
	cobra.CheckErr(err)
	defer alternator.Close()

	schemas, err := alternator.FetchSchemas()
	cobra.CheckErr(err)

	// Show remote database schemas
	ePrintf(strings.Repeat("―", width))
	if len(schemas) == 0 {
		bPrintln("No database.")
		os.Exit(0)
	}
	for _, s := range schemas {
		fmt.Println(s.Database.String())
		fmt.Println()
		for _, t := range s.Tables {
			fmt.Println(t.String())
			fmt.Println()
		}
	}
}
