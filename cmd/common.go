package cmd

import (
	"database/sql"
	"fmt"
	"github.com/emirpasic/gods/sets/hashset"
	"github.com/kota65535/alternator/lib"
	"github.com/kota65535/alternator/parser"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"io/ioutil"
	"net/url"
	"os"
	"strings"
)

var SupportedDialects = hashset.New("mysql")

type DatabaseUrl struct {
	Dialect  string
	Host     string
	DbName   string
	User     string
	Password string
}

func (r DatabaseUrl) Dsn() string {
	return fmt.Sprintf("%s%s@(%s)/%s", r.User, optS(r.Password, ":%s"), r.Host, r.DbName)
}

func ParseDatabaseUrl(u string) DatabaseUrl {
	p, err := url.Parse(u)
	if err != nil {
		bPrintf("Error: failed to parse URL: %s.\n", u)
		if !strings.Contains(u, "://") {
			ePrintln("URL format: (dialect)://(username)[:(password)]@(hostname)[:(port)][/(database)]")
			ePrintln("Examples:")
			ePrintln("  mysql://root@localhost")
			ePrintln("  mysql://root@localhost:13306/db1")
			ePrintln("  mysql://root:mystrongpassword@mydb.dev.example.com/db1")
		}
		os.Exit(1)
	}
	dialect := p.Scheme
	if !SupportedDialects.Contains(dialect) {
		bPrintf("Error: unsupported dialect: %s\n", dialect)
		ePrintf("Supported dialects: %v\n", SupportedDialects)
		os.Exit(1)
	}

	host := p.Host
	dbName := ""
	if len(p.Path) > 0 {
		dbName = p.Path[1:]
	}
	user := p.User.Username()
	password, _ := p.User.Password()

	return DatabaseUrl{dialect, host, dbName, user, password}
}

func ConnectToDb(dbUrl DatabaseUrl) *sql.DB {
	logrus.Debug(dbUrl.Dsn())
	db, err := sql.Open(dbUrl.Dialect, dbUrl.Dsn())
	cobra.CheckErr(err)
	return db
}

func FetchGlobalConfig(db *sql.DB) *parser.GlobalConfig {
	rows1, err := db.Query("SHOW GLOBAL VARIABLES")
	defer rows1.Close()
	cobra.CheckErr(err)
	var name string
	var value string
	variables := map[string]string{}
	for rows1.Next() {
		err = rows1.Scan(&name, &value)
		variables[name] = value
	}

	rows2, err := db.Query("SHOW CHARACTER SET")
	defer rows2.Close()
	cobra.CheckErr(err)
	var charset string
	var description string
	var collation string
	var maxLen string
	charsetToCollation := map[string]string{}
	for rows2.Next() {
		err = rows2.Scan(&charset, &description, &collation, &maxLen)
		charsetToCollation[charset] = collation
	}

	return &parser.GlobalConfig{
		CharacterSetServer:   variables["character_set_server"],
		CharacterSetDatabase: variables["character_set_database"],
		CollationServer:      variables["collation_server"],
		CharsetToCollation:   charsetToCollation,
	}
}

func GetAlterations(path string, db *sql.DB, dbUrl DatabaseUrl, config *parser.GlobalConfig) lib.DatabaseAlterations {
	bPrint("Reading local schema file... ")
	toSchemas := ReadSchemas(path, config)
	bPrintln("done.")
	fromSchemas := FetchSchemas(db, dbUrl, config)

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

func ReadSchemas(filename string, config *parser.GlobalConfig) []lib.Schema {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		cobra.CheckErr(err)
	}
	return lib.NewSchemas(string(b), config)
}

func FetchSchemas(db *sql.DB, dbUrl DatabaseUrl, config *parser.GlobalConfig) []lib.Schema {
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
