package cmd

import (
	_ "embed"
	"fmt"
	"github.com/kota65535/alternator/lib"
	"github.com/spf13/cobra"
	"strings"
)

//go:embed plan.tmpl
var planUsage string

func init() {
	c := &cobra.Command{
		Use:   "plan <schema-file> <database-url>",
		Short: "Show remote database schema changes required by the local schema file.",
		Long:  "Show remote database schema changes required by the local schema file.",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			PlanCmd(args[0], args[1])
		},
	}
	rootCmd.AddCommand(c)
	c.SetUsageTemplate(planUsage)
}

func PlanCmd(path string, uri string) *lib.DatabaseAlterations {
	dbUri, err := NewDatabaseUri(uri)
	cobra.CheckErr(err)

	alternator, err := NewAlternator(dbUri)
	cobra.CheckErr(err)
	defer alternator.Close()

	alt, _, _, err := alternator.GetAlterationsFromFile(path)
	cobra.CheckErr(err)

	// Show diff
	ePrintln(strings.Repeat("―", width))
	bPrintln("Schema diff:")
	bPrintln()
	for _, s := range alt.Diff() {
		printlnDiff(s)
	}

	// Show statements to execute
	ePrintln(strings.Repeat("―", width))
	statements := alt.Statements()
	if len(statements) == 0 {
		bPrintln("Your database schema is up-to-date! No change required.")
		return nil
	}
	bPrintln("Statements to execute:")
	bPrintln()
	for _, s := range alt.Statements() {
		fmt.Println(s)
	}

	return alt
}
