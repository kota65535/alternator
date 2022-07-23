package parser

import (
	"fmt"
	"github.com/emirpasic/gods/maps/linkedhashmap"
	"reflect"
	"strings"
)

type GlobalConfig struct {
	CharacterSetServer   string
	CharacterSetDatabase string
	CollationServer      string
	CharsetToCollation   map[string]string
}

type DatabaseOptions struct {
	DefaultCharset    string
	DefaultCollate    string
	DefaultEncryption string

	GlobalConfig *GlobalConfig
}

func (r DatabaseOptions) ActualDefaultCharset() string {
	if r.DefaultCharset != "" {
		return r.DefaultCharset
	} else {
		return r.GlobalConfig.CharacterSetServer
	}
}

func (r DatabaseOptions) ActualDefaultCollate() string {
	if r.DefaultCollate != "" {
		return r.DefaultCollate
	} else if r.DefaultCharset != "" {
		return r.GlobalConfig.CharsetToCollation[r.DefaultCharset]
	} else {
		return r.GlobalConfig.CollationServer
	}
}
func (r DatabaseOptions) String() string {
	return strings.Join(r.Strings(), " ")
}

func (r DatabaseOptions) Map() *linkedhashmap.Map {
	ret := linkedhashmap.New()
	if r.DefaultCharset != "" {
		ret.Put("DEFAULT CHARACTER SET", r.DefaultCharset)
	}
	if r.DefaultCollate != "" {
		ret.Put("DEFAULT COLLATE", r.DefaultCollate)
	}
	if r.DefaultEncryption != "" {
		ret.Put("DEFAULT ENCRYPTION", r.DefaultEncryption)
	}
	return ret
}

func (r DatabaseOptions) Strings() []string {
	ret := []string{}
	m := r.Map()
	for _, k := range m.Keys() {
		v, ok := m.Get(k)
		if ok {
			ret = append(ret, fmt.Sprintf("%s = %s", k, v))
		}
	}
	return ret
}

type ColumnDefinition struct {
	ColumnName    string
	DataType      interface{}
	ColumnOptions ColumnOptions
}

func (r ColumnDefinition) Equals(a ColumnDefinition) bool {
	r2, ok1 := r.DataType.(StringType)
	a2, ok2 := a.DataType.(StringType)
	if ok1 && ok2 {
		e := r2.ActualCharset() == a2.ActualCharset() && r2.ActualCollate() == a2.ActualCollate()
		r2.DefaultCharset = ""
		r2.DefaultCollation = ""
		a2.DefaultCharset = ""
		a2.DefaultCollation = ""
		r.DataType = r2
		a.DataType = a2
		return reflect.DeepEqual(r, a) && e
	}
	return reflect.DeepEqual(r, a)
}

func (r ColumnDefinition) EqualsExceptColumnName(c ColumnDefinition) bool {
	r.ColumnName = ""
	c.ColumnName = ""
	return r.Equals(c)
}

func (r ColumnDefinition) String() string {
	return fmt.Sprintf("`%s`\t%s\t%s",
		r.ColumnName,
		r.DataType,
		r.ColumnOptions.String())
}

func (r ColumnDefinition) StringWithPos(pos string) string {
	columnOptionWithPos := fmt.Sprintf("%s%s", optS(r.ColumnOptions.String(), "%s "), pos)
	return fmt.Sprintf("`%s`\t%s\t%s",
		r.ColumnName,
		r.DataType,
		columnOptionWithPos)
}

type ColumnOptions struct {
	Nullability               string
	Default                   string
	GeneratedAs               string
	GeneratedColumnType       string
	Visibility                string
	AutoIncrement             bool
	Unique                    bool
	Primary                   bool
	Comment                   string
	ReferenceDefinition       ReferenceDefinition
	CheckConstraintDefinition CheckConstraintDefinition
	OnUpdate                  string
	Srid                      string
}

func (r ColumnOptions) String() string {
	return strings.Join(r.Strings(), " ")
}

func (r ColumnOptions) Strings() []string {
	var strs []string
	if r.Nullability != "" {
		strs = append(strs, r.Nullability)
	}
	if r.Default != "" {
		strs = append(strs, fmt.Sprintf("DEFAULT %s", r.Default))
	}
	if r.GeneratedAs != "" {
		strs = append(strs, fmt.Sprintf("GENERATED ALWAYS AS %s", r.GeneratedAs))
	}
	if r.GeneratedColumnType != "" {
		strs = append(strs, r.GeneratedColumnType)
	}
	if r.Visibility != "" {
		strs = append(strs, r.Visibility)
	}
	if r.AutoIncrement {
		strs = append(strs, "AUTO_INCREMENT")
	}
	if r.Unique {
		strs = append(strs, "UNIQUE KEY")
	}
	if r.Primary {
		strs = append(strs, "PRIMARY KEY")
	}
	if r.Comment != "" {
		strs = append(strs, fmt.Sprintf("COMMENT %s", r.Comment))
	}
	if r.ReferenceDefinition.TableName != "" {
		strs = append(strs, r.ReferenceDefinition.String())
	}
	if r.CheckConstraintDefinition.Check != "" {
		strs = append(strs, r.CheckConstraintDefinition.String())
	}
	if r.OnUpdate != "" {
		strs = append(strs, fmt.Sprintf("ON UPDATE %s", r.OnUpdate))
	}

	return strs
}

func (r ColumnOptions) Diff(o ColumnOptions) ColumnOptions {
	return structDifference(r, o)
}

type IntegerType struct {
	Name     string
	FieldLen string
	Unsigned bool
	Zerofill bool
}

func (t IntegerType) String() string {
	return fmt.Sprintf("%s%s%s%s",
		t.Name,
		optS(t.FieldLen, "(%s)"),
		optB(t.Unsigned, " unsigned"),
		optB(t.Zerofill, " zerofill"))
}

type FixedPointType struct {
	Name       string
	FieldLen   string
	FieldScale string
	Unsigned   bool
	Zerofill   bool
}

func (t FixedPointType) String() string {
	lenAndPlace := ""
	if t.FieldLen != "" {
		lenAndPlace += fmt.Sprintf("(%s", t.FieldLen)
		if t.FieldScale != "" {
			lenAndPlace += fmt.Sprintf(", %s", t.FieldScale)
		}
		lenAndPlace += ")"
	}
	return fmt.Sprintf("%s%s%s%s",
		t.Name,
		lenAndPlace,
		optB(t.Unsigned, " unsigned"),
		optB(t.Zerofill, " zerofill"))
}

type FloatingPointType struct {
	Name       string
	FieldLen   string
	FieldScale string
	Unsigned   bool
	Zerofill   bool
}

func (t FloatingPointType) String() string {
	lenAndPlace := ""
	if t.FieldLen != "" {
		lenAndPlace += fmt.Sprintf("(%s", t.FieldLen)
		if t.FieldScale != "" {
			lenAndPlace += fmt.Sprintf(", %s", t.FieldScale)
		}
		lenAndPlace += ")"
	}
	return fmt.Sprintf("%s%s%s%s",
		t.Name,
		lenAndPlace,
		optB(t.Unsigned, " unsigned"),
		optB(t.Zerofill, " zerofill"))
}

type DateAndTimeType struct {
	Name     string
	FieldLen string
}

func (t DateAndTimeType) String() string {
	return fmt.Sprintf("%s%s",
		t.Name,
		optS(t.FieldLen, "(%s)"))
}

type StringType struct {
	Name             string
	FieldLen         string
	Charset          string
	Collation        string
	DefaultCharset   string
	DefaultCollation string
}

func (r StringType) ActualCharset() string {
	if r.Charset != "" {
		return r.Charset
	} else {
		return r.DefaultCharset
	}
}

func (r StringType) ActualCollate() string {
	if r.Collation != "" {
		return r.Collation
	} else {
		return r.DefaultCollation
	}
}

func (t StringType) String() string {
	return fmt.Sprintf("%s%s%s%s",
		t.Name,
		optS(t.FieldLen, "(%s)"),
		optB(t.Charset != "" && t.Charset != t.DefaultCharset, fmt.Sprintf(" CHARACTER SET %s", t.Charset)),
		optB(t.Collation != "" && t.Collation != t.DefaultCollation, fmt.Sprintf(" COLLATE %s", t.Collation)))
}

type StringListType struct {
	Name      string
	Values    []string
	Charset   string
	Collation string
}

func (t StringListType) String() string {
	values := ""
	if len(t.Values) > 0 {
		values = fmt.Sprintf("(%s)", strings.Join(t.Values, ", "))
	}
	return fmt.Sprintf("%s%s%s%s",
		t.Name,
		values,
		optS(t.Charset, " CHARACTER SET %s"),
		optS(t.Collation, " COLLATE %s"))
}

type JsonType struct {
	Name string
}

func (t JsonType) String() string {
	return t.Name
}

type SpatialType struct {
	Name string
}

func (t SpatialType) String() string {
	return t.Name
}

type ReferenceDefinition struct {
	TableName        string
	KeyPartList      []KeyPart
	ReferenceOptions ReferenceOptions
}

func (r ReferenceDefinition) String() string {
	return fmt.Sprintf("REFERENCES `%s` (%s)%s",
		r.TableName,
		JoinT(r.KeyPartList, ", ", ""),
		optS(r.ReferenceOptions.String(), " %s"))
}

type ReferenceOptions struct {
	Match    string
	OnDelete string
	OnUpdate string
}

func (r ReferenceOptions) String() string {
	return strings.Join(r.Strings(), " ")
}

func (r ReferenceOptions) Strings() []string {
	var strs []string
	if r.Match != "" {
		strs = append(strs, fmt.Sprintf("MATCH %s", r.Match))
	}
	if r.OnDelete != "" {
		strs = append(strs, fmt.Sprintf("ON DELETE %s", r.OnDelete))
	}
	if r.OnUpdate != "" {
		strs = append(strs, fmt.Sprintf("ON UPDATE %s", r.OnUpdate))
	}
	return strs
}

type CheckConstraintDefinition struct {
	ConstraintName         string
	Check                  string
	CheckConstraintOptions CheckConstraintOptions
}

func (r CheckConstraintDefinition) String() string {
	return fmt.Sprintf("%sCHECK (%s)%s",
		optS(r.ConstraintName, "CONSTRAINT `%s` "),
		stripSequentialSpaces(r.Check),
		optS(r.CheckConstraintOptions.String(), " %s"))
}

type CheckConstraintOptions struct {
	Enforcement string
}

func (r CheckConstraintOptions) String() string {
	return strings.Join(r.Strings(), " ")
}

func (r CheckConstraintOptions) Strings() []string {
	var strs []string
	if r.Enforcement != "" {
		strs = append(strs, r.Enforcement)
	}
	return strs
}

func (r CheckConstraintOptions) Diff(o CheckConstraintOptions) CheckConstraintOptions {
	return structDifference(r, o)
}

type KeyPart struct {
	Column     string
	Length     string
	Order      string
	Expression string
}

func (r KeyPart) String() string {
	columnOrExpr := ""
	if r.Column != "" {
		columnOrExpr = fmt.Sprintf("`%s`", r.Column)
	} else {
		columnOrExpr = r.Expression
	}
	return fmt.Sprintf("%s%s%s",
		columnOrExpr,
		optS(r.Length, "(%s)"),
		optS(r.Order, " %s"))
}

type IndexDefinition struct {
	IndexName    string
	KeyPartList  []KeyPart
	IndexOptions IndexOptions
}

func (r IndexDefinition) EqualsExceptIndexName(a IndexDefinition) bool {
	r.IndexName = ""
	a.IndexName = ""
	return r.Equals(a)
}
func (r IndexDefinition) Equals(a IndexDefinition) bool {
	if len(r.KeyPartList) != len(a.KeyPartList) {
		return false
	}
	return r.IndexName == a.IndexName && reflect.DeepEqual(r.IndexOptions, a.IndexOptions)
}

func (r IndexDefinition) String() string {
	return fmt.Sprintf("INDEX%s (%s)%s",
		optS(r.IndexName, " `%s`"),
		JoinT(r.KeyPartList, ", ", ""),
		optS(r.IndexOptions.String(), " %s"))
}

func (r IndexDefinition) StringKeyPartList() string {
	return fmt.Sprintf("INDEX (%s)", JoinT(r.KeyPartList, ", ", ""))
}

type FullTextIndexDefinition struct {
	IndexName    string
	KeyPartList  []KeyPart
	IndexOptions IndexOptions
}

func (r FullTextIndexDefinition) String() string {
	return fmt.Sprintf("FULLTEXT INDEX%s (%s)%s",
		optS(r.IndexName, " `%s`"),
		JoinT(r.KeyPartList, ", ", ""),
		optS(r.IndexOptions.String(), " %s"))
}

func (r FullTextIndexDefinition) StringKeyPartList() string {
	return fmt.Sprintf("FULLTEXT INDEX (%s)", JoinT(r.KeyPartList, ", ", ""))
}

type PrimaryKeyDefinition struct {
	ConstraintName string
	KeyPartList    []KeyPart
	IndexOptions   IndexOptions
}

func (r PrimaryKeyDefinition) String() string {
	return fmt.Sprintf("%sPRIMARY KEY (%s)%s",
		optS(r.ConstraintName, "CONSTRAINT `%s` "),
		JoinT(r.KeyPartList, ", ", ""),
		optS(r.IndexOptions.String(), " %s"))
}

func (r PrimaryKeyDefinition) StringKeyPartList() string {
	return fmt.Sprintf("PRIMARY KEY (%s)", JoinT(r.KeyPartList, ", ", ""))
}

type UniqueKeyDefinition struct {
	ConstraintName string
	IndexName      string
	KeyPartList    []KeyPart
	IndexOptions   IndexOptions
}

func (r UniqueKeyDefinition) String() string {
	return fmt.Sprintf("%sUNIQUE KEY%s (%s)%s",
		optS(r.ConstraintName, "CONSTRAINT `%s` "),
		optS(r.IndexName, " `%s`"),
		JoinT(r.KeyPartList, ", ", ""),
		optS(r.IndexOptions.String(), " %s"))
}

func (r UniqueKeyDefinition) StringKeyPartList() string {
	return fmt.Sprintf("UNIQUE KEY (%s)", JoinT(r.KeyPartList, ", ", ""))
}

type IndexOptions struct {
	IndexType    string
	KeyBlockSize string
	Parser       string
	Comment      string
	Visibility   string
}

func (r IndexOptions) String() string {
	return strings.Join(r.Strings(), " ")
}

func (r IndexOptions) Strings() []string {
	var strs []string
	if r.IndexType != "" {
		strs = append(strs, fmt.Sprintf("USING %s", r.IndexType))
	}
	if r.KeyBlockSize != "" {
		strs = append(strs, fmt.Sprintf("KEY_BLOCK_SIZE %s", r.KeyBlockSize))
	}
	if r.Parser != "" {
		strs = append(strs, fmt.Sprintf("WITH PARSER %s", r.Parser))
	}
	if r.Comment != "" {
		strs = append(strs, fmt.Sprintf("COMMENT %s", r.Comment))
	}
	if r.Visibility != "" {
		strs = append(strs, r.Visibility)
	}

	return strs
}

func (r IndexOptions) Diff(o IndexOptions) IndexOptions {
	return structDifference(r, o)
}

type ForeignKeyDefinition struct {
	ConstraintName      string
	IndexName           string
	KeyPartList         []KeyPart
	ReferenceDefinition ReferenceDefinition
}

func (r ForeignKeyDefinition) String() string {
	return fmt.Sprintf("%sFOREIGN KEY%s (%s)%s",
		optS(r.ConstraintName, "CONSTRAINT `%s` "),
		optS(r.IndexName, " `%s`"),
		JoinT(r.KeyPartList, ", ", ""),
		optS(r.ReferenceDefinition.String(), " %s"))
}

func (r ForeignKeyDefinition) StringKeyPartList() string {
	return fmt.Sprintf("FOREIGN KEY (%s)", JoinT(r.KeyPartList, ", ", ""))
}

type TableOptions struct {
	AutoExtendedSize         string
	AutoIncrement            string
	AvgRowLength             string
	DefaultCharset           string
	Checksum                 string
	DefaultCollate           string
	Comment                  string
	Compression              string
	Connection               string
	DataDirectory            string
	IndexDirectory           string
	DelayKeyWrite            string
	Encryption               string
	Engine                   string
	EngineAttribute          string
	InsertMethod             string
	KeyBlockSize             string
	MaxRows                  string
	MinRows                  string
	PackKeys                 string
	Password                 string
	RowFormat                string
	SecondaryEngineAttribute string
	StatsAutoRecalc          string
	StatsPersistent          string
	StatsSamplePages         string
	TableSpace               string
	TableSpaceStorage        string
	Union                    []string
	DatabaseOptions          *DatabaseOptions
}

func (r TableOptions) ActualDefaultCharset() string {
	if r.DefaultCharset != "" {
		return r.DefaultCharset
	} else {
		return r.DatabaseOptions.ActualDefaultCharset()
	}
}

func (r TableOptions) ActualDefaultCollate() string {
	if r.DefaultCollate != "" {
		return r.DefaultCollate
	} else {
		return r.DatabaseOptions.ActualDefaultCollate()
	}
}

func (r TableOptions) Map() *linkedhashmap.Map {
	ret := linkedhashmap.New()
	if r.AutoExtendedSize != "" {
		ret.Put("AUTOEXTENDED_SIZE", r.AutoExtendedSize)
	}
	if r.AutoIncrement != "" {
		ret.Put("AUTO_INCREMENT", r.AutoIncrement)
	}
	if r.AvgRowLength != "" {
		ret.Put("AVG_ROW_LENGTH", r.AvgRowLength)
	}
	if r.DefaultCharset != "" && (r.DatabaseOptions == nil || r.DefaultCharset != r.DatabaseOptions.DefaultCharset) {
		ret.Put("DEFAULT CHARACTER SET", r.DefaultCharset)
	}
	if r.Checksum != "" {
		ret.Put("CHECKSUM", r.Checksum)
	}
	if r.DefaultCollate != "" && (r.DatabaseOptions == nil || r.DefaultCollate != r.DatabaseOptions.DefaultCollate) {
		ret.Put("DEFAULT COLLATE", r.DefaultCollate)
	}
	if r.Comment != "" {
		ret.Put("COMMENT", r.Comment)
	}
	if r.Compression != "" {
		ret.Put("COMPRESSION", r.Compression)
	}
	if r.Connection != "" {
		ret.Put("CONNECTION", r.Connection)
	}
	if r.DataDirectory != "" {
		ret.Put("DATA DIRECTORY", r.DataDirectory)
	}
	if r.IndexDirectory != "" {
		ret.Put("INDEX DIRECTORY", r.IndexDirectory)
	}
	if r.DelayKeyWrite != "" {
		ret.Put("DELAY_KEY_WRITE", r.DelayKeyWrite)
	}
	if r.Encryption != "" {
		ret.Put("ENCRYPTION", r.Encryption)
	}
	if r.Engine != "" {
		ret.Put("ENGINE", r.Engine)
	}
	if r.EngineAttribute != "" {
		ret.Put("ENGINE_ATTRIBUTE", r.EngineAttribute)
	}
	if r.InsertMethod != "" {
		ret.Put("INSERT_METHOD", r.InsertMethod)
	}
	if r.KeyBlockSize != "" {
		ret.Put("KEY_BLOCK_SIZE", r.KeyBlockSize)
	}
	if r.MaxRows != "" {
		ret.Put("MAX_ROWS", r.MaxRows)
	}
	if r.MinRows != "" {
		ret.Put("MIN_ROWS", r.MinRows)
	}
	if r.PackKeys != "" {
		ret.Put("PACK_KEYS", r.PackKeys)
	}
	if r.Password != "" {
		ret.Put("PASSWORD", r.Password)
	}
	if r.RowFormat != "" {
		ret.Put("ROW_FORMAT", r.RowFormat)
	}
	if r.SecondaryEngineAttribute != "" {
		ret.Put("SECONDARY_ENGINE_ATTRIBUTE", r.SecondaryEngineAttribute)
	}
	if r.StatsAutoRecalc != "" {
		ret.Put("STATS_AUTO_RECALC", r.StatsAutoRecalc)
	}
	if r.StatsPersistent != "" {
		ret.Put("STATS_PERSISTENT", r.StatsPersistent)
	}
	if r.StatsSamplePages != "" {
		ret.Put("STATS_SAMPLE_PAGES", r.StatsSamplePages)
	}
	if r.TableSpace != "" {
		ret.Put("TABLESPACE", fmt.Sprintf("%s %s", r.TableSpace, optS(r.TableSpaceStorage, "STORAGE %s")))
	}
	if len(r.Union) != 0 {
		ret.Put("UNION", fmt.Sprintf("(%s)", JoinS(r.Union, ", ", "`")))
	}
	return ret
}

func (r TableOptions) Strings() []string {
	ret := []string{}
	m := r.Map()
	for _, k := range m.Keys() {
		v, ok := m.Get(k)
		if ok {
			ret = append(ret, fmt.Sprintf("%s = %s", k, v))
		}
	}
	return ret
}

func (r TableOptions) String() string {
	return strings.Join(r.Strings(), " ")
}

type PartitionConfig struct {
	PartitionBy          PartitionBy
	Partitions           string
	SubpartitionBy       PartitionBy
	Subpartitions        string
	PartitionDefinitions []PartitionDefinition
}

func (r PartitionConfig) Strings() []string {
	ret := []string{}
	if r.PartitionBy.Type == "" {
		return ret
	}
	ret = append(ret, fmt.Sprintf("PARTITION BY %s", r.PartitionBy.String()))
	if r.Partitions != "" {
		ret = append(ret, fmt.Sprintf("PARTITIONS %s", r.Partitions))
	}
	if r.SubpartitionBy.Type != "" {
		ret = append(ret, fmt.Sprintf("SUBPARTITION BY %s", r.SubpartitionBy.String()))
	}
	if r.Subpartitions != "" {
		ret = append(ret, fmt.Sprintf("SUBPARTITIONS %s", r.Subpartitions))
	}
	if len(r.PartitionDefinitions) > 0 {
		ret = append(ret, "(")
		for i, p := range r.PartitionDefinitions {
			if i == len(r.PartitionDefinitions)-1 {
				ret = append(ret, strings.Repeat(" ", 4)+p.String())
			} else {
				ret = append(ret, strings.Repeat(" ", 4)+p.String()+",")
			}
		}
		ret = append(ret, ")")
	}
	return ret
}

func (r PartitionConfig) String() string {
	return strings.Join(r.Strings(), "\n")
}

type PartitionDefinition struct {
	Name             string
	Operator         string
	ValueExpression  string
	PartitionOptions PartitionOptions
	Subpartitions    []SubpartitionDefinition
}

func (r PartitionDefinition) Strings() []string {
	ret := []string{}
	ret = append(ret, fmt.Sprintf("PARTITION %s VALUES %s (%s)", r.Name, r.Operator, r.ValueExpression))
	ret = append(ret, r.PartitionOptions.Strings()...)
	for _, p := range r.Subpartitions {
		ret = append(ret, p.Strings()...)
	}
	return ret
}

func (r PartitionDefinition) String() string {
	return strings.Join(r.Strings(), "\n")
}

type SubpartitionDefinition struct {
	Name             string
	PartitionOptions PartitionOptions
}

func (r SubpartitionDefinition) Strings() []string {
	ret := []string{}
	ret = append(ret, fmt.Sprintf("SUBPARTITION %s", r.Name))
	ret = append(ret, r.PartitionOptions.Strings()...)
	return ret
}

func (r SubpartitionDefinition) String() string {
	return strings.Join(r.Strings(), "\n")
}

type PartitionBy struct {
	Type       string
	Expression string
	Columns    []string
}

func (r PartitionBy) String() string {
	var values string
	if r.Expression != "" {
		values = fmt.Sprintf("(%s)", r.Expression)
	} else {
		values = "COLUMNS (" + JoinS(r.Columns, ", ", "`") + ")"
	}
	return fmt.Sprintf("%s %s", r.Type, values)
}

type PartitionOptions struct {
	Engine         string
	Comment        string
	DataDirectory  string
	IndexDirectory string
	MaxRows        string
	MinRows        string
	TableSpace     string
}

func (r PartitionOptions) Strings() []string {
	ret := []string{}
	m := r.Map()
	for _, k := range m.Keys() {
		v, ok := m.Get(k)
		if ok {
			ret = append(ret, fmt.Sprintf("%s = %s", k, v))
		}
	}
	return ret
}

func (r PartitionOptions) Map() *linkedhashmap.Map {
	ret := linkedhashmap.New()
	if r.Comment != "" {
		ret.Put("COMMENT", fmt.Sprintf("'%s'", r.Comment))
	}
	if r.DataDirectory != "" {
		ret.Put("DATA DIRECTORY", fmt.Sprintf("'%s'", r.DataDirectory))
	}
	if r.IndexDirectory != "" {
		ret.Put("INDEX DIRECTORY", fmt.Sprintf("'%s'", r.IndexDirectory))
	}
	if r.Engine != "" {
		ret.Put("ENGINE", r.Engine)
	}
	if r.MaxRows != "" {
		ret.Put("MAX_ROWS", r.MaxRows)
	}
	if r.MinRows != "" {
		ret.Put("MIN_ROWS", r.MinRows)
	}
	if r.TableSpace != "" {
		ret.Put("TABLESPACE", r.TableSpace)
	}
	return ret
}
