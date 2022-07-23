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
		Short: "Show database schema changes required by the schema file",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			PlanCmd(args[0], args[1])
		},
	}
	rootCmd.AddCommand(c)
	c.SetUsageTemplate(planUsage)
}

func PlanCmd(local string, remote string) *lib.DatabaseAlterations {
	dbUrl := ParseDatabaseUrl(remote)

	bPrint("Connecting to database... ")
	Db = ConnectToDb(dbUrl)
	defer Db.Close()
	bPrintln("done.")

	bPrint("Fetching remote server global config... ")
	config := FetchGlobalConfig()
	bPrintln("done.")

	alt := GetAlterations(local, dbUrl, config)

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
	return &alt
}
