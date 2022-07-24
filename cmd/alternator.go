package cmd

import (
	"database/sql"
	"fmt"
	"github.com/emirpasic/gods/sets/hashset"
	"github.com/go-sql-driver/mysql"
	"github.com/kota65535/alternator/lib"
	"github.com/kota65535/alternator/parser"
	"io/ioutil"
	"net/url"
	"strings"
)

var (
	SupportedDialects = hashset.New("mysql")
	IgnoredDatabases  = hashset.New("information_schema", "mysql", "performance_schema", "sys")
)

type DatabaseUri struct {
	Dialect  string
	Host     string
	DbName   string
	User     string
	Password string
}

func NewDatabaseUri(uri string) (*DatabaseUri, error) {
	p, err := url.Parse(uri)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URI: %s : %w", uri, err)
	}

	dialect := p.Scheme
	if !SupportedDialects.Contains(dialect) {
		return nil, fmt.Errorf("unsupported dialect: %s : %w", uri, err)
	}

	host := p.Host
	dbName := ""
	if len(p.Path) > 0 {
		dbName = p.Path[1:]
	}
	if dbName == "" && !managesAllDatabases {
		return nil, fmt.Errorf("database name is required")
	}

	user := p.User.Username()
	password, _ := p.User.Password()

	return &DatabaseUri{
		dialect,
		host,
		dbName,
		user,
		password,
	}, nil
}

func (r DatabaseUri) Dsn() string {
	return fmt.Sprintf("%s%s@(%s)/%s", r.User, optS(r.Password, ":%s"), r.Host, r.DbName)
}

func (r DatabaseUri) DsnWoDbName() string {
	return fmt.Sprintf("%s%s@(%s)/", r.User, optS(r.Password, ":%s"), r.Host)
}

type Alternator struct {
	DbUri        *DatabaseUri
	Db           *sql.DB
	GlobalConfig *parser.GlobalConfig
}

func NewAlternator(dbUri *DatabaseUri) (*Alternator, error) {
	// do not use database name because it may not exist in the remote server
	db, err := sql.Open(dbUri.Dialect, dbUri.DsnWoDbName())
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection. DSN = %s : %w", dbUri.DsnWoDbName(), err)
	}

	globalConfig, err := fetchGlobalConfig(db)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch global config : %w", err)
	}

	return &Alternator{
		dbUri,
		db,
		globalConfig,
	}, nil
}

func (r *Alternator) ReadSchemas(schema string) ([]*lib.Schema, error) {
	schemas, err := lib.NewSchemas(schema, r.GlobalConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create shema : %w", err)
	}
	return schemas, nil
}

func (r *Alternator) ReadSchemasFromFile(path string) ([]*lib.Schema, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read shema file: %s : %w", path, err)
	}
	return r.ReadSchemas(string(b))
}

func (r *Alternator) FetchSchemas() ([]*lib.Schema, error) {
	if r.DbUri.DbName != "" {
		schema, err := r.fetchFromDatabase(r.DbUri.DbName)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch remote database schema: %w", err)
		}
		// Schema is empty
		if schema == nil {
			return []*lib.Schema{}, nil
		}
		return []*lib.Schema{schema}, nil
	} else {
		schemas, err := r.fetchFromDatabases()
		if err != nil {
			return nil, fmt.Errorf("failed to fetch remote database schemas: %w", err)
		}
		return schemas, nil
	}
}

func (r *Alternator) GetAlterations(schema string) (*lib.DatabaseAlterations, []*lib.Schema, []*lib.Schema, error) {
	localSchemas, err := r.ReadSchemas(schema)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to read local shema file: %w", err)
	}
	remoteSchemas, err := r.FetchSchemas()
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to fetch remote schema: %w", err)
	}

	return lib.NewDatabaseAlterations(remoteSchemas, localSchemas), remoteSchemas, localSchemas, nil
}

func (r *Alternator) GetAlterationsFromFile(path string) (*lib.DatabaseAlterations, []*lib.Schema, []*lib.Schema, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to read shema file: %s : %w", path, err)
	}
	return r.GetAlterations(string(b))
}

func (r *Alternator) Close() error {
	err := r.Db.Close()
	if err != nil {
		return fmt.Errorf("error: failed to close DB connection")
	}
	return nil
}

func (r *Alternator) fetchFromDatabases() ([]*lib.Schema, error) {
	databases, err := r.listUserDefinedDatabases()
	if err != nil {
		return nil, fmt.Errorf("failed to list user defined databases : %w", err)
	}

	var schemas []*lib.Schema
	for _, d := range databases {
		schema, err := r.fetchFromDatabase(d)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch remote database schema : %w", err)
		}
		// Schema is empty
		if schema != nil {
			schemas = append(schemas, schema)
		}
	}

	return schemas, nil
}

func (r *Alternator) fetchFromDatabase(dbName string) (*lib.Schema, error) {
	var strs []string

	databaseSchema, err := r.getCreateDatabase(dbName)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch remote database creation statement: %w", err)
	}
	// Schema is empty
	if databaseSchema == "" {
		return nil, nil
	}

	strs = append(strs, databaseSchema)
	strs = append(strs, fmt.Sprintf("USE `%s`", dbName))

	tables, err := r.listTables(dbName)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch remote table names: %w", err)
	}

	for _, t := range tables {
		tableSchema, err := r.getCreateTable(dbName, t)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch remote table creation statement: %w", err)
		}
		strs = append(strs, tableSchema)
	}

	schemas, err := lib.NewSchemas(strings.Join(strs, ";\n"), r.GlobalConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create shema : %w", err)
	}

	return schemas[0], nil
}

func (r *Alternator) getCreateDatabase(name string) (string, error) {
	rows, err := r.Db.Query(fmt.Sprintf("SHOW CREATE DATABASE `%s`", name))
	if err != nil {
		if v, ok := err.(*mysql.MySQLError); ok && v.Number == 1049 {
			return "", nil
		} else {
			return "", fmt.Errorf("failed to query \"SHOW CREATE DATABASE\": %w", err)
		}
	}
	defer rows.Close()
	var dbName string
	var statement string
	for rows.Next() {
		_ = rows.Scan(&dbName, &statement)
	}
	return statement, nil
}

func (r *Alternator) getCreateTable(dbName string, tableName string) (string, error) {
	rows, err := r.Db.Query(fmt.Sprintf("SHOW CREATE TABLE `%s`.`%s`", dbName, tableName))
	if err != nil {
		return "", fmt.Errorf("failed to query \"SHOW CREATE TABLE\": %w", err)
	}
	defer rows.Close()
	var statement string
	for rows.Next() {
		_ = rows.Scan(&tableName, &statement)
	}
	return statement, nil
}

func (r *Alternator) listDatabases() ([]string, error) {
	rows, err := r.Db.Query("SHOW DATABASES")
	if err != nil {
		return nil, fmt.Errorf("failed to query \"SHOW DATABASES\": %w", err)
	}
	defer rows.Close()
	var databases []string
	var database string
	for rows.Next() {
		_ = rows.Scan(&database)
		databases = append(databases, database)
	}
	return databases, nil
}

func (r *Alternator) listUserDefinedDatabases() ([]string, error) {
	databases, err := r.listDatabases()
	if err != nil {
		return nil, fmt.Errorf("failed to list databases: %w", err)
	}
	filtered := []string{}
	for _, d := range databases {
		if !IgnoredDatabases.Contains(d) {
			filtered = append(filtered, d)
		}
	}
	return filtered, nil
}

func (r *Alternator) listTables(dbName string) ([]string, error) {
	rows, err := r.Db.Query(fmt.Sprintf("SHOW TABLES FROM `%s`", dbName))
	if err != nil {
		return nil, fmt.Errorf("failed to query \"SHOW TABLES FROM `%s`\": %w", dbName, err)
	}
	defer rows.Close()

	var tables []string
	var table string
	for rows.Next() {
		_ = rows.Scan(&table)
		tables = append(tables, table)
	}
	return tables, nil
}

func fetchGlobalConfig(db *sql.DB) (*parser.GlobalConfig, error) {
	rows1, err := db.Query("SHOW GLOBAL VARIABLES")
	if err != nil {
		return nil, fmt.Errorf("failed to query \"SHOW GLOBAL VARIABLES\": %w", err)
	}
	defer rows1.Close()
	var name string
	var value string
	variables := map[string]string{}
	for rows1.Next() {
		err = rows1.Scan(&name, &value)
		variables[name] = value
	}

	rows2, err := db.Query("SHOW CHARACTER SET")
	if err != nil {
		return nil, fmt.Errorf("failed to query \"SHOW CHARACTER SET\": %w", err)
	}
	defer rows2.Close()
	var charset string
	var description string
	var collation string
	var maxLen string
	charsetToCollation := map[string]string{}
	for rows2.Next() {
		err = rows2.Scan(&charset, &description, &collation, &maxLen)
		charsetToCollation[charset] = collation
	}

	if val, ok := variables["default_table_encryption"]; ok {
		if val == "ON" {
			variables["default_table_encryption"] = "'Y'"
		} else {
			variables["default_table_encryption"] = "'N'"
		}
	} else {
		variables["default_table_encryption"] = "'N'"
	}

	return &parser.GlobalConfig{
		CharacterSetServer:   variables["character_set_server"],
		CharacterSetDatabase: variables["character_set_database"],
		CollationServer:      variables["collation_server"],
		CharsetToCollation:   charsetToCollation,
		Encryption:           variables["default_table_encryption"],
	}, nil
}
