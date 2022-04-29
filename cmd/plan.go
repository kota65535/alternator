package cmd

import (
	"fmt"
	"github.com/kota65535/alternator/lib"
	"github.com/spf13/cobra"
	"io/ioutil"
	"strings"
)

func init() {
	c := &cobra.Command{
		Use:   "plan <schema-file> <database-url>",
		Short: "Plan schema change of database",
		Args:  cobra.RangeArgs(2, 2),
		Run: func(cmd *cobra.Command, args []string) {
			planCmd(args[0], args[1])
		},
	}
	rootCmd.AddCommand(c)
}

func planCmd(local string, remote string) lib.DatabaseAlterations {
	alt := getAlterations(local, remote)

	ePrintln(strings.Repeat("―", width))
	bPrintln("Statements To Execute:")
	bPrintln()
	for _, s := range alt.Statements() {
		fmt.Println(s)
	}
	return alt
}

func readSchemas(filename string) []lib.Schema {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		cobra.CheckErr(err)
	}
	return lib.NewSchemas(string(b))
}
