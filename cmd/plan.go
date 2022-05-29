package cmd

import (
	_ "embed"
	"fmt"
	"github.com/kota65535/alternator/lib"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"io/ioutil"
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
			planCmd(args[0], args[1])
		},
	}
	rootCmd.AddCommand(c)
	c.SetUsageTemplate(planUsage)
}

func planCmd(local string, remote string) *lib.DatabaseAlterations {

	alt := getAlterations(local, remote)

	ePrintln(strings.Repeat("―", width))
	bPrintln("Schema diff:")
	bPrintln()
	for _, s := range alt.Diff() {
		printlnDiff(s)
	}
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

func getAlterations(to string, from string) lib.DatabaseAlterations {
	toSchemas := readSchemas(to)
	fromSchemas := fetchSchemas(from)

	logrus.Debug("Showing local file schema")
	for _, s := range toSchemas {
		logrus.Debug(s.String())
	}
	logrus.Debug("Showing remote database schema")
	for _, s := range fromSchemas {
		logrus.Debug(s.String())
	}

	return lib.NewDatabaseAlterations(fromSchemas, toSchemas)
}

func readSchemas(filename string) []lib.Schema {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		cobra.CheckErr(err)
	}
	return lib.NewSchemas(string(b))
}
