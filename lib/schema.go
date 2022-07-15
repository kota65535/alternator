package lib

import (
	"github.com/emirpasic/gods/sets/hashset"
	"github.com/emirpasic/gods/stacks/arraystack"
	"github.com/kota65535/alternator/parser"
	"github.com/spf13/cobra"
	"regexp"
	"sort"
	"strings"
)

type Schema struct {
	Database *parser.CreateDatabaseStatement
	Tables   []*parser.CreateTableStatement
}

var TypeDefaultFieldLen = map[string]string{
	"bit":       "1",
	"tinyint":   "4",
	"smallint":  "6",
	"mediumint": "9",
	"int":       "11",
	"bigint":    "20",
	"decimal":   "10,0",
}

func (r Schema) String() string {
	statements := []string{}

	statements = append(statements, r.Database.StringWithFormat(4))
	for _, t := range r.Tables {
		statements = append(statements, t.StringWithFormat(4))
	}
	return strings.Join(statements, "\n")
}

func NewSchemas(str string) []Schema {

	stripped := stripConditionalComments(str)

	p := parser.NewParser(strings.NewReader(stripped))
	statements, err := p.Parse()
	cobra.CheckErr(err)

	schema := normalizeStatements(statements)

	// Sort table name by their dependencies (foreign keys)
	//for _, s := range schema {
	//	s.Tables = NewDag(s.Tables).Sort()
	//}

	return schema
}

func stripOutermostParentheses(str string) string {
	stack := arraystack.New()
	pairIdx := map[int]int{}
	for i := 0; i < len(str); i++ {
		switch str[i] {
		case '(':
			stack.Push(i)
		case ')':
			j, ok := stack.Pop()
			if !ok {
				return str
			}
			pairIdx[j.(int)] = i
		}
	}
	nStrip := 0
	for i := 0; i < len(str); i++ {
		if pairIdx[i] != len(str)-(i+1) {
			break
		}
		nStrip += 1
	}

	return str[nStrip : len(str)-nStrip]
}

func stripConditionalComments(str string) string {
	startPattern := regexp.MustCompile("/\\*!\\d*")
	endPattern := regexp.MustCompile("\\*/")

	started := false
	pos := 0
	ret := ""
	for pos < len(str) {
		if started {
			loc := endPattern.FindStringIndex(str[pos:])
			if loc == nil {
				ret += str[pos:]
				break
			}
			ret += str[pos : pos+loc[0]]
			pos += loc[1]
			started = false
		} else {
			loc := startPattern.FindStringIndex(str[pos:])
			if loc == nil {
				ret += str[pos:]
				break
			}
			ret += str[pos : pos+loc[0]]
			pos += loc[1]
			started = true
		}
	}
	return ret
}

func normalizeDataType(t interface{}) interface{} {
	if it, ok := t.(parser.IntegerType); ok {
		// Unset If field length is default
		if it.FieldLen == TypeDefaultFieldLen[it.Name] {
			it.FieldLen = ""
		}
		// decimal(20,0) -> decimal(20)
		// decimal(10,0) -> decimal
		// decimal(10) -> decimal
		if it.Name == "decimal" {
			fl := strings.Split(it.FieldLen, ",")
			if len(fl) == 2 && fl[1] == "0" {
				it.FieldLen = fl[0]
			}
		}
		return it
	}
	return t
}

func normalizeStatements(statements []parser.Statement) []Schema {
	currentDbName := ""
	defaultCharsets := map[string]string{}
	defaultCollates := map[string]string{}
	schemas := map[string]*Schema{}
	for i, _ := range statements {
		s := statements[i]
		if us, ok := s.(parser.UseStatement); ok {
			currentDbName = us.DbName
		}
		if cds, ok := s.(parser.CreateDatabaseStatement); ok {
			schemas[cds.DbName] = &Schema{
				Database: &cds,
				Tables:   []*parser.CreateTableStatement{},
			}
			defaultCharsets[cds.DbName] = cds.DatabaseOptions.DefaultCharset
			defaultCollates[cds.DbName] = cds.DatabaseOptions.DefaultCollate
		}
		if cts, ok := s.(parser.CreateTableStatement); ok {
			// Current DB name set by USE statement
			if cts.DbName == "" {
				cts.DbName = currentDbName
			}

			defaultCharset := defaultCharsets[currentDbName]
			defaultCollate := defaultCollates[currentDbName]

			// Unset default charset/collation if they match database's default
			if cts.TableOptions.DefaultCharset == defaultCharset {
				cts.TableOptions.DefaultCharset = ""
			} else {
				defaultCharset = cts.TableOptions.DefaultCharset
			}
			if cts.TableOptions.DefaultCollate == defaultCollate {
				cts.TableOptions.DefaultCollate = ""
			} else {
				defaultCollate = cts.TableOptions.DefaultCollate
			}
			// Unset if engine is InnoDB, which is default
			if cts.TableOptions.Engine == "InnoDB" {
				cts.TableOptions.Engine = ""
			}

			var createDefinitions []interface{}

			columns := cts.GetColumns()
			primaryKeys := cts.GetPrimaryKeys()
			uniqueKeys := cts.GetUniqueKeys()
			indexes := cts.GetIndexes()
			fullTexts := cts.GetFullTextIndexes()
			foreignKeys := cts.GetForeignKeys()
			checks := cts.GetCheckConstraints()

			// ========== Column modifications ==========

			for _, v := range columns {

				v.DataType = normalizeDataType(v.DataType)

				// Unset charset/collation if equal to the table's default
				if dt, ok := v.DataType.(parser.StringType); ok {
					if dt.Charset == defaultCharset {
						dt.Charset = ""
					}
					if dt.Collation == defaultCollate {
						dt.Collation = ""
					}
					v.DataType = dt
				}

				// Unset if nullability is NULL, which is default
				if v.ColumnOptions.Nullability == "NULL" {
					v.ColumnOptions.Nullability = ""
				}

				// Unset if default is NULL, which is default
				if v.ColumnOptions.Default == "NULL" {
					v.ColumnOptions.Default = ""
				}

				// Separate to primary key definition
				if v.ColumnOptions.Primary {
					primaryKeys = append(primaryKeys, &parser.PrimaryKeyDefinition{
						KeyPartList: []parser.KeyPart{{Column: v.ColumnName}},
					})
					v.ColumnOptions.Primary = false
				}
				// Separate to unique key definition
				if v.ColumnOptions.Unique {
					uniqueKeys = append(uniqueKeys, &parser.UniqueKeyDefinition{
						KeyPartList: []parser.KeyPart{{Column: v.ColumnName}},
					})
					v.ColumnOptions.Unique = false
				}
				// Separate to foreign key definition
				if v.ColumnOptions.ReferenceDefinition.TableName != "" {
					foreignKeys = append(foreignKeys, &parser.ForeignKeyDefinition{
						KeyPartList: []parser.KeyPart{{Column: v.ColumnName}},
						ReferenceDefinition: parser.ReferenceDefinition{
							TableName:        v.ColumnOptions.ReferenceDefinition.TableName,
							KeyPartList:      v.ColumnOptions.ReferenceDefinition.KeyPartList,
							ReferenceOptions: v.ColumnOptions.ReferenceDefinition.ReferenceOptions,
						},
					})
					v.ColumnOptions.ReferenceDefinition = parser.ReferenceDefinition{}
				}
				// Separate to check constraint definition
				if v.ColumnOptions.CheckConstraintDefinition.Check != "" {
					checks = append(checks, &parser.CheckConstraintDefinition{
						Check:                  v.ColumnOptions.CheckConstraintDefinition.Check,
						CheckConstraintOptions: v.ColumnOptions.CheckConstraintDefinition.CheckConstraintOptions,
					})
					v.ColumnOptions.CheckConstraintDefinition = parser.CheckConstraintDefinition{}
				}
			}

			// Unset if index order is ASC, which is default
			for _, p := range primaryKeys {
				for i, _ := range p.KeyPartList {
					k := &p.KeyPartList[i]
					if k.Order == "ASC" {
						k.Order = ""
					}
				}
			}
			// Unset if index order is ASC, which is default
			for _, p := range indexes {
				for i, _ := range p.KeyPartList {
					k := &p.KeyPartList[i]
					if k.Order == "ASC" {
						k.Order = ""
					}
				}
			}
			// Unset if index order is ASC, which is default
			for _, p := range fullTexts {
				for i, _ := range p.KeyPartList {
					k := &p.KeyPartList[i]
					if k.Order == "ASC" {
						k.Order = ""
					}
				}
			}
			// Unset if index order is ASC, which is default
			for _, p := range foreignKeys {
				for i, _ := range p.KeyPartList {
					k := &p.KeyPartList[i]
					if k.Order == "ASC" {
						k.Order = ""
					}
				}
			}

			// Add NOT NULL for primary key column
			for _, p := range primaryKeys {
				for _, c := range columns {
					if keyPartContains(p.KeyPartList, c.ColumnName) {
						c.ColumnOptions.Nullability = "NOT NULL"
					}
				}
			}

			// Determine key name of unique keys.
			// | constraint name | index name | key name        |
			// |-----------------|------------|-----------------|
			// | n               | n          | generated       |
			// | n               | y          | index name      |
			// | y               | n          | constraint name |
			// | y               | y          | index name      |
			for _, v := range uniqueKeys {
				if v.ConstraintName != "" && v.IndexName == "" {
					v.IndexName = v.ConstraintName
				}
				// Unset constraint name because it is not shown in 'SHOW CREATE TABLE' statement output
				v.ConstraintName = ""
			}

			// Determine key name of foreign keys.
			// | constraint name | index name | key name        |
			// |-----------------|------------|-----------------|
			// | n               | n          | generated       |
			// | n               | y          | index name      |
			// | y               | n          | constraint name |
			// | y               | y          | constraint name |
			for _, v := range foreignKeys {
				if v.ConstraintName != "" {
					v.IndexName = v.ConstraintName
				}
			}

			// Remove index definitions created by foreign keys
			fkIndexes := hashset.New()
			for _, f := range foreignKeys {
				for _, idx := range indexes {
					if arraysEqual(f.KeyPartList, idx.KeyPartList) {
						f.IndexName = idx.IndexName
						fkIndexes.Add(idx)
					}
				}
			}
			indexes = RemoveIf(indexes, func(e *parser.IndexDefinition) bool {
				return fkIndexes.Contains(e)
			})

			// Unset if reference option is RESTRICT or NO ACTION, which is default
			for _, v := range foreignKeys {
				if Contains([]string{"RESTRICT", "NO ACTION"}, v.ReferenceDefinition.ReferenceOptions.OnDelete) {
					v.ReferenceDefinition.ReferenceOptions.OnDelete = ""
				}
				if Contains([]string{"RESTRICT", "NO ACTION"}, v.ReferenceDefinition.ReferenceOptions.OnUpdate) {
					v.ReferenceDefinition.ReferenceOptions.OnUpdate = ""
				}
			}

			for _, c := range checks {
				c.Check = stripOutermostParentheses(c.Check)
			}

			for _, c := range columns {
				createDefinitions = append(createDefinitions, c)
			}
			for _, c := range primaryKeys {
				createDefinitions = append(createDefinitions, c)
			}
			for _, c := range uniqueKeys {
				createDefinitions = append(createDefinitions, c)
			}
			for _, c := range foreignKeys {
				createDefinitions = append(createDefinitions, c)
			}
			for _, c := range checks {
				createDefinitions = append(createDefinitions, c)
			}
			for _, c := range indexes {
				createDefinitions = append(createDefinitions, c)
			}
			for _, c := range fullTexts {
				createDefinitions = append(createDefinitions, c)
			}

			cts.CreateDefinitions = createDefinitions

			schemas[cts.DbName].Tables = append(schemas[cts.DbName].Tables, &cts)
		}
	}

	// Sort database names alphabetically
	dbNames := []string{}
	for k, _ := range schemas {
		dbNames = append(dbNames, k)
	}
	sort.Strings(dbNames)

	var ret []Schema
	for _, k := range dbNames {
		ret = append(ret, *schemas[k])
	}

	return ret
}
