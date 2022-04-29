package parser

import (
	"fmt"
	"strings"
)

const Indent = 4

type Statement interface {
	String() string
	StringWithFormat(int) string
}

type CreateDatabaseStatement struct {
	IfNotExists     bool
	DbName          string
	DatabaseOptions DatabaseOptions
}

func (r CreateDatabaseStatement) String() string {
	return r.StringWithFormat(Indent)
}

func (r CreateDatabaseStatement) StringWithFormat(indent int) string {
	databaseOptions := strings.Join(addIndent(r.DatabaseOptions.Strings(), indent), "\n")
	return fmt.Sprintf("CREATE DATABASE `%s`%s;",
		r.DbName,
		optS(databaseOptions, "\n%s"))
}

type CreateTableStatement struct {
	DbName            string
	Temporary         bool
	IfNotExists       bool
	TableName         string
	CreateDefinitions []interface{}
	TableOptions      TableOptions
}

func (r CreateTableStatement) String() string {
	return r.StringWithFormat(Indent)
}

func (r CreateTableStatement) StringWithFormat(indent int) string {
	var column []string
	var primary []string
	var unique []string
	var foreign []string
	var check []string
	var index []string
	var fulltext []string

	for _, d := range r.GetColumns() {
		column = append(column, d.String())
	}
	column = Align(column)
	for _, d := range r.GetPrimaryKeys() {
		primary = append(primary, d.String())
	}
	for _, d := range r.GetUniqueKeys() {
		unique = append(unique, d.String())
	}
	for _, d := range r.GetForeignKeys() {
		foreign = append(foreign, d.String())
	}
	for _, d := range r.GetCheckConstraints() {
		check = append(check, d.String())
	}
	for _, d := range r.GetIndexes() {
		index = append(index, d.String())
	}
	for _, d := range r.GetFullTextIndexes() {
		fulltext = append(fulltext, d.String())
	}

	var defStrs []string
	defStrs = append(defStrs, column...)
	defStrs = append(defStrs, primary...)
	defStrs = append(defStrs, unique...)
	defStrs = append(defStrs, foreign...)
	defStrs = append(defStrs, check...)
	defStrs = append(defStrs, index...)
	defStrs = append(defStrs, fulltext...)

	defStrs = addIndent(defStrs, indent)

	tableOptions := strings.Join(addIndent(r.TableOptions.Strings(), indent), "\n")

	return fmt.Sprintf("CREATE TABLE %s`%s`\n(\n%s\n)%s;",
		optS(r.DbName, "`%s`."),
		r.TableName,
		strings.Join(defStrs, ",\n"),
		optS(tableOptions, "\n%s"))
}

func (r CreateTableStatement) GetColumns() []*ColumnDefinition {
	var ret []*ColumnDefinition
	for _, cd := range r.CreateDefinitions {
		if d, ok := cd.(*ColumnDefinition); ok {
			ret = append(ret, d)
		}
	}
	return ret
}

func (r CreateTableStatement) GetIndexes() []*IndexDefinition {
	var ret []*IndexDefinition
	for _, cd := range r.CreateDefinitions {
		if d, ok := cd.(*IndexDefinition); ok {
			ret = append(ret, d)
		}
	}
	return ret
}

func (r CreateTableStatement) GetFullTextIndexes() []*FullTextIndexDefinition {
	var ret []*FullTextIndexDefinition
	for _, cd := range r.CreateDefinitions {
		if d, ok := cd.(*FullTextIndexDefinition); ok {
			ret = append(ret, d)
		}
	}
	return ret
}

func (r CreateTableStatement) GetPrimaryKeys() []*PrimaryKeyDefinition {
	var ret []*PrimaryKeyDefinition
	for _, cd := range r.CreateDefinitions {
		if d, ok := cd.(*PrimaryKeyDefinition); ok {
			ret = append(ret, d)
		}
	}
	return ret
}

func (r CreateTableStatement) GetUniqueKeys() []*UniqueKeyDefinition {
	var ret []*UniqueKeyDefinition
	for _, cd := range r.CreateDefinitions {
		if d, ok := cd.(*UniqueKeyDefinition); ok {
			ret = append(ret, d)
		}
	}
	return ret
}

func (r CreateTableStatement) GetForeignKeys() []*ForeignKeyDefinition {
	var ret []*ForeignKeyDefinition
	for _, cd := range r.CreateDefinitions {
		if d, ok := cd.(*ForeignKeyDefinition); ok {
			ret = append(ret, d)
		}
	}
	return ret
}

func (r CreateTableStatement) GetCheckConstraints() []*CheckConstraintDefinition {
	var ret []*CheckConstraintDefinition
	for _, cd := range r.CreateDefinitions {
		if d, ok := cd.(*CheckConstraintDefinition); ok {
			ret = append(ret, d)
		}
	}
	return ret
}

type UseStatement struct {
	DbName string
}

func (r UseStatement) String() string {
	return fmt.Sprintf("USE `%s`;", r.DbName)
}

func (r UseStatement) StringWithFormat(indent int) string {
	return r.String()
}
