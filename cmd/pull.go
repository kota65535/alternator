package cmd

import (
	"database/sql"
	_ "embed"
	"fmt"
	"github.com/emirpasic/gods/sets/hashset"
	_ "github.com/go-sql-driver/mysql"
	"github.com/kota65535/alternator/lib"
	"github.com/kota65535/alternator/parser"
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
	dbUrl := parseDatabaseUrl(url)
	bPrintf("Connecting to database... ")
	db := connectToDb(dbUrl)
	bPrintf("done.")
	defer db.Close()
	bPrintf("Fetching remote server global config... ")
	config := parser.FetchGlobalConfig(db)
	schemas := fetchSchemas(db, dbUrl, config)

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

func connectToDb(dbUrl DatabaseUrl) *sql.DB {
	dPrintln(dbUrl.Dsn())
	db, err := sql.Open(dbUrl.Dialect, dbUrl.Dsn())
	cobra.CheckErr(err)
	return db
}

func fetchSchemas(db *sql.DB, dbUrl DatabaseUrl, config *parser.GlobalConfig) []lib.Schema {
	var schemas []lib.Schema
	if dbUrl.DbName != "" {
		bPrintf("Fetching schemas of database '%s'...\n", dbUrl.DbName)
		schemas = []lib.Schema{fetchFromDatabase(db, dbUrl.DbName, config)}
	} else {
		bPrintf("Fetching user-defined all database schemas...\n")
		schemas = fetchFromDatabases(db, config)
	}
	return schemas
}

func fetchFromDatabases(db *sql.DB, config *parser.GlobalConfig) []lib.Schema {

	databases := listUserDefinedDatabases(db)

	var schemas []lib.Schema
	for _, d := range databases {
		stmts := fetchFromDatabase(db, d, config)
		schemas = append(schemas, stmts)
	}

	return schemas
}

func fetchFromDatabase(db *sql.DB, dbName string, config *parser.GlobalConfig) lib.Schema {
	var strs []string

	strs = append(strs, getCreateDatabase(db, dbName))
	strs = append(strs, fmt.Sprintf("USE `%s`", dbName))

	tables := listTables(db, dbName)

	for _, t := range tables {
		strs = append(strs, getCreateTable(db, dbName, t))
	}

	return lib.NewSchemas(strings.Join(strs, ";\n"), config)[0]
}

func getCreateDatabase(db *sql.DB, name string) string {
	rows, err := db.Query(fmt.Sprintf("SHOW CREATE DATABASE `%s`", name))
	defer rows.Close()
	cobra.CheckErr(err)
	var dbName string
	var statement string
	for rows.Next() {
		_ = rows.Scan(&dbName, &statement)
	}
	return statement
}

func getCreateTable(db *sql.DB, dbName string, tableName string) string {
	rows, err := db.Query(fmt.Sprintf("SHOW CREATE TABLE `%s`.`%s`", dbName, tableName))
	defer rows.Close()
	cobra.CheckErr(err)
	var statement string
	for rows.Next() {
		_ = rows.Scan(&tableName, &statement)
	}
	return statement
}

func listDatabases(db *sql.DB) []string {
	rows, err := db.Query("SHOW DATABASES")
	defer rows.Close()
	cobra.CheckErr(err)

	var databases []string
	var database string
	for rows.Next() {
		_ = rows.Scan(&database)
		databases = append(databases, database)
	}
	return databases
}

func listUserDefinedDatabases(db *sql.DB) []string {
	databases := listDatabases(db)
	ret := []string{}
	for _, d := range databases {
		if !IgnoredDatabases.Contains(d) {
			ret = append(ret, d)
		}
	}
	return ret
}

func listTables(db *sql.DB, dbName string) []string {
	rows, err := db.Query(fmt.Sprintf("SHOW TABLES FROM `%s`", dbName))
	defer rows.Close()
	cobra.CheckErr(err)

	var tables []string
	var table string
	for rows.Next() {
		_ = rows.Scan(&table)
		tables = append(tables, table)
	}
	return tables
}
