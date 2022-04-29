package cmd

import (
	"github.com/kota65535/alternator/lib"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"strings"
)

func init() {
	c := &cobra.Command{
		Use:   "diff <schema-file> <database-url>",
		Short: "Diff schema change of database",
		Args:  cobra.RangeArgs(2, 2),
		Run: func(cmd *cobra.Command, args []string) {
			diffCmd(args[0], args[1])
		},
	}
	rootCmd.AddCommand(c)
}

func diffCmd(local string, remote string) *lib.DatabaseAlterations {
	alt := getAlterations(local, remote)

	ePrintln(strings.Repeat("―", width))
	bPrintln("Schema Diff:")
	bPrintln()
	for _, s := range alt.Diff() {
		printlnDiff(s)
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
