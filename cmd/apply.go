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
		Short: "Update the database schema according to the schema file",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			ApplyCmd(args[0], args[1], params)
		},
	}
	c.Flags().BoolVar(&params.AutoApprove, "auto-approve", false, "Approve automatically")
	rootCmd.AddCommand(c)
	c.SetUsageTemplate(applyUsage)
}

func ApplyCmd(local string, remote string, params ApplyParams) *lib.DatabaseAlterations {
	dbUrl := ParseDatabaseUrl(remote)

	bPrint("Connecting to database... ")
	Db = ConnectToDb(dbUrl)
	bPrintln("done.")
	defer Db.Close()

	bPrint("Fetching remote server global config... ")
	config := FetchGlobalConfig()
	bPrintln("done.")

	alt := GetAlterations(local, dbUrl, config)

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
		_, err := Db.Exec(s)
		cobra.CheckErr(err)
	}
	bPrintln("\nFinished!")
	return &alt
}
