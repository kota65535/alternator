package cmd

import (
	_ "embed"
	"fmt"
	"github.com/kota65535/alternator/lib"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

//go:embed apply.tmpl
var applyUsage string

type ApplyParams struct {
	AutoApprove bool
}

func init() {
	var params ApplyParams

	c := &cobra.Command{
		Use:   "apply <schema-file> <database-url>",
		Short: "Update the remote database schema according to the local schema file.",
		Long:  "Update the remote database schema according to the local schema file.",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			ApplyCmd(args[0], args[1], params)
		},
	}
	c.Flags().BoolVar(&params.AutoApprove, "auto-approve", false, "Approve automatically")
	rootCmd.AddCommand(c)
	c.SetUsageTemplate(applyUsage)
}

func ApplyCmd(path string, uri string, params ApplyParams) *lib.DatabaseAlterations {
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
	bPrintln()

	// Apply
	ePrintln(strings.Repeat("―", width))
	if !params.AutoApprove {
		if !confirm("Do you want to apply?") {
			os.Exit(0)
		}
		ePrintln()
	}

	for _, s := range alt.Statements() {
		ePrintf("Executing: %s\n", s)
		_, err := alternator.Db.Exec(s)
		cobra.CheckErr(err)
	}
	bPrintln("\nFinished!")

	return alt
}
