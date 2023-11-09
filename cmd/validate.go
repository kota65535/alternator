package cmd

import (
	_ "embed"
	"fmt"
	"github.com/emirpasic/gods/sets/hashset"
	"github.com/kota65535/alternator/lib"
	"github.com/kota65535/alternator/parser"
	"github.com/spf13/cobra"
	"io/ioutil"
)

//go:embed validate.tmpl
var validateUsage string

func init() {
	c := &cobra.Command{
		Use:   "validate <schema-file>",
		Short: "Validate local schema file.",
		Long:  "Validate local schema file.",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			ValidateCmd(args[0])
		},
	}
	rootCmd.AddCommand(c)
	c.SetUsageTemplate(validateUsage)
}

func ValidateCmd(path string) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		cobra.CheckErr(fmt.Errorf("failed to read shema file: %s : %w", path, err))
	}
	_, err = lib.NewSchemas(string(b), &parser.GlobalConfig{}, hashset.New())
	cobra.CheckErr(err)
}
