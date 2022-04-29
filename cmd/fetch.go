package cmd

import (
	"database/sql"
	"fmt"
	"github.com/emirpasic/gods/sets/hashset"
	_ "github.com/go-sql-driver/mysql"
	"github.com/kota65535/alternator/lib"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

var IgnoredDatabases = hashset.New("information_schema", "mysql", "performance_schema", "sys")

func init() {
	c := &cobra.Command{
		Use:   "fetch <database-url>",
		Short: "Fetch schemas from database",
		Args:  cobra.RangeArgs(1, 1),
		Run: func(cmd *cobra.Command, args []string) {
			fetchCmd(args[0])
		},
	}

	rootCmd.AddCommand(c)
}

func fetchCmd(url string) {
	schemas := fetchSchemas(url)

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

func fetchSchemas(url string) []lib.Schema {
	dbUrl := parseDatabaseUrl(url)
	dPrintln(dbUrl.Dsn())
	db, err := sql.Open(dbUrl.Dialect, dbUrl.Dsn())
	defer db.Close()
	cobra.CheckErr(err)

	var schemas []lib.Schema
	if dbUrl.DbName != "" {
		bPrintf("Fetching schemas of database '%s'...\n", dbUrl.DbName)
		schemas = []lib.Schema{fetchFromDatabase(db, dbUrl.DbName)}
	} else {
		bPrintf("Fetching user-defined all database schemas...\n")
		schemas = fetchFromDatabases(db)
	}
	return schemas
}

func fetchFromDatabases(db *sql.DB) []lib.Schema {

	databases := listUserDefinedDatabases(db)

	var schemas []lib.Schema
	for _, d := range databases {
		stmts := fetchFromDatabase(db, d)
		schemas = append(schemas, stmts)
	}

	return schemas
}

func fetchFromDatabase(db *sql.DB, dbName string) lib.Schema {
	var strs []string

	strs = append(strs, getCreateDatabase(db, dbName))
	strs = append(strs, fmt.Sprintf("USE `%s`", dbName))

	tables := listTables(db, dbName)

	for _, t := range tables {
		strs = append(strs, getCreateTable(db, dbName, t))
	}

	return lib.NewSchemas(strings.Join(strs, ";\n"))[0]
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
