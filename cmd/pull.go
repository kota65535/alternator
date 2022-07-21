package cmd

import (
	_ "embed"
	"fmt"
	"github.com/emirpasic/gods/sets/hashset"
	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

//go:embed pull.tmpl
var pullUsage string

var IgnoredDatabases = hashset.New("information_schema", "mysql", "performance_schema", "sys")

func init() {
	c := &cobra.Command{
		Use:   "pull <database-url>",
		Short: "Show the current database schema",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			pullCmd(args[0])
		},
	}
	rootCmd.AddCommand(c)
	c.SetUsageTemplate(pullUsage)
}

func pullCmd(url string) {
	dbUrl := ParseDatabaseUrl(url)
	bPrintf("Connecting to database... ")
	db := ConnectToDb(dbUrl)
	bPrintf("done.")
	defer db.Close()
	bPrintf("Fetching remote server global config... ")
	config := FetchGlobalConfig(db)
	schemas := FetchSchemas(db, dbUrl, config)

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
