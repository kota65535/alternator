package cmd

import (
	"bufio"
	"database/sql"
	_ "embed"
	"fmt"
	"github.com/kota65535/alternator/lib"
	"github.com/spf13/cobra"
	"log"
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
			applyCmd(args[0], args[1], params)
		},
	}
	c.Flags().BoolVar(&params.AutoApprove, "auto-approve", false, "Approve automatically")
	rootCmd.AddCommand(c)
	c.SetUsageTemplate(applyUsage)
}

func applyCmd(local string, remote string, params ApplyParams) *lib.DatabaseAlterations {
	alt := getAlterations(local, remote)

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

	ePrintln(strings.Repeat("―", width))
	if !params.AutoApprove {
		if !confirm("Do you want to apply?") {
			os.Exit(0)
		}
		ePrintln()
	}

	dbUrl := parseDatabaseUrl(remote)
	db, err := sql.Open(dbUrl.Dialect, dbUrl.Dsn())
	defer db.Close()
	cobra.CheckErr(err)

	for _, s := range alt.Statements() {
		ePrintf("Executing: %s\n", s)
		_, err := db.Exec(s)
		cobra.CheckErr(err)
	}
	bPrintln("\nFinished!")
	return &alt
}

func confirm(s string) bool {
	reader := bufio.NewReader(os.Stdin)

	for {
		bPrintf("%s [y/n]: ", s)

		response, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		response = strings.ToLower(strings.TrimSpace(response))

		if response == "y" || response == "yes" {
			return true
		} else if response == "n" || response == "no" {
			return false
		}
	}
}
