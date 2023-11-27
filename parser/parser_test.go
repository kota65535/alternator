package parser

import (
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	logrus.SetLevel(logrus.DebugLevel)
	code := m.Run()
	os.Exit(code)
}

func TestCreateDb(t *testing.T) {

	f, err := os.Open("test/db/input.sql")
	require.NoError(t, err)

	p := NewParser(f)
	r, err := p.Parse()
	require.NoError(t, err)

	assert.Equal(t, []Statement{
		CreateDatabaseStatement{
			DbName:          "db1",
			DatabaseOptions: &DatabaseOptions{},
		},
		CreateDatabaseStatement{
			IfNotExists: true,
			DbName:      "db2",
			DatabaseOptions: &DatabaseOptions{
				DefaultCharset:    "utf8mb4",
				DefaultCollate:    "utf8mb4_bin",
				DefaultEncryption: "'Y'",
			},
		},
		CreateDatabaseStatement{
			IfNotExists: true,
			DbName:      "db3",
			DatabaseOptions: &DatabaseOptions{
				DefaultCharset:    "utf8mb4",
				DefaultCollate:    "utf8mb4_unicode_ci",
				DefaultEncryption: "'N'",
			},
		},
	}, r)

	b1, err := os.ReadFile("test/db/output1.sql")
	b2, err := os.ReadFile("test/db/output2.sql")
	b3, err := os.ReadFile("test/db/output3.sql")

	assert.Equal(t, string(b1), r[0].String())
	assert.Equal(t, string(b2), r[1].String())
	assert.Equal(t, string(b3), r[2].String())
}

func TestUseDb(t *testing.T) {

	f, err := os.Open("test/use/input.sql")
	p := NewParser(f)
	r, err := p.Parse()

	require.NoError(t, err)

	assert.Equal(t, []Statement{
		UseStatement{
			DbName: "db1",
		},
		UseStatement{
			DbName: "db2",
		},
	}, r)

	b1, err := os.ReadFile("test/use/output1.sql")
	b2, err := os.ReadFile("test/use/output2.sql")

	assert.Equal(t, string(b1), r[0].String())
	assert.Equal(t, string(b2), r[1].String())

}

func TestCreateTableWithOptions(t *testing.T) {
	f, err := os.Open("test/table/options/input.sql")
	p := NewParser(f)
	r, err := p.Parse()

	require.NoError(t, err)

	assert.Equal(t, []Statement{
		CreateTableStatement{
			TableName: "t1",
			CreateDefinitions: []interface{}{
				&ColumnDefinition{
					ColumnName: "int1",
					DataType: IntegerType{
						Name: "int",
					},
				},
			},
			TableOptions: TableOptions{
				AutoExtendedSize:         "1",
				AutoIncrement:            "1",
				AvgRowLength:             "1",
				DefaultCharset:           "utf8mb4",
				Checksum:                 "1",
				DefaultCollate:           "utf8mb4_bin",
				Comment:                  "'foo'",
				Compression:              "'ZLIB'",
				Connection:               "'connect_string'",
				DataDirectory:            "'path1'",
				IndexDirectory:           "'path2'",
				DelayKeyWrite:            "1",
				Encryption:               "'Y'",
				Engine:                   "INNODB",
				EngineAttribute:          "'attr1'",
				InsertMethod:             "FIRST",
				KeyBlockSize:             "1",
				MaxRows:                  "1",
				MinRows:                  "1",
				PackKeys:                 "1",
				Password:                 "'password'",
				RowFormat:                "DYNAMIC",
				SecondaryEngineAttribute: "'attr2'",
				StatsAutoRecalc:          "1",
				StatsPersistent:          "1",
				StatsSamplePages:         "1",
				TableSpace:               "tbl_space",
				TableSpaceStorage:        "DISK",
				Union:                    []string{"t2", "t3"},
			},
		},
	}, r)

	b1, err := os.ReadFile("test/table/options/output.sql")

	assert.Equal(t, string(b1), r[0].String())
}

func TestCreateTableWithExpressions(t *testing.T) {
	f, err := os.Open("test/table/expression/input.sql")
	p := NewParser(f)
	r, err := p.Parse()

	require.NoError(t, err)

	b1, err := os.ReadFile("test/table/expression/output.sql")

	assert.Equal(t, string(b1), r[0].String())
}

func TestCreateTableWithNumericTypes(t *testing.T) {

	f, err := os.Open("test/table/numeric/input.sql")
	p := NewParser(f)
	r, err := p.Parse()

	require.NoError(t, err)

	assert.Equal(t,
		CreateTableStatement{
			TableName: "t1",
			CreateDefinitions: []interface{}{
				&ColumnDefinition{
					ColumnName: "bit1",
					DataType: IntegerType{
						Name: "bit",
					},
				},
				&ColumnDefinition{
					ColumnName: "bit2",
					DataType: IntegerType{
						Name:     "bit",
						FieldLen: "1",
					},
					ColumnOptions: ColumnOptions{
						Nullability: "NOT NULL",
						Default:     "0b001",
					},
				},
				&ColumnDefinition{
					ColumnName: "tinyint1",
					DataType: IntegerType{
						Name: "tinyint",
					},
				},
				&ColumnDefinition{
					ColumnName: "tinyint2",
					DataType: IntegerType{
						Name:     "tinyint",
						FieldLen: "1",
						Unsigned: true,
						Zerofill: true,
					},
					ColumnOptions: ColumnOptions{
						Nullability: "NOT NULL",
						Default:     "1",
					},
				},
				&ColumnDefinition{
					ColumnName: "bool1",
					DataType: IntegerType{
						Name: "bool",
					},
				},
				&ColumnDefinition{
					ColumnName: "bool2",
					DataType: IntegerType{
						Name: "bool",
					},
					ColumnOptions: ColumnOptions{
						Nullability: "NOT NULL",
						Default:     "TRUE",
					},
				},
				&ColumnDefinition{
					ColumnName: "smallint1",
					DataType: IntegerType{
						Name: "smallint",
					},
				},
				&ColumnDefinition{
					ColumnName: "smallint2",
					DataType: IntegerType{
						Name:     "smallint",
						FieldLen: "1",
						Unsigned: true,
						Zerofill: true,
					},
					ColumnOptions: ColumnOptions{
						Nullability: "NOT NULL",
						Default:     "0x123",
					},
				},
				&ColumnDefinition{
					ColumnName: "mediumint1",
					DataType: IntegerType{
						Name: "mediumint",
					},
				},
				&ColumnDefinition{
					ColumnName: "mediumint2",
					DataType: IntegerType{
						Name:     "mediumint",
						FieldLen: "1",
						Unsigned: true,
						Zerofill: true,
					},
					ColumnOptions: ColumnOptions{
						Nullability: "NOT NULL",
						Default:     "0x0123",
					},
				},
				&ColumnDefinition{
					ColumnName: "int1",
					DataType: IntegerType{
						Name: "int",
					},
				},
				&ColumnDefinition{
					ColumnName: "int2",
					DataType: IntegerType{
						Name:     "int",
						FieldLen: "1",
						Unsigned: true,
						Zerofill: true,
					},
					ColumnOptions: ColumnOptions{
						Nullability: "NOT NULL",
						Default:     "1",
					},
				},
				&ColumnDefinition{
					ColumnName: "bigint1",
					DataType: IntegerType{
						Name: "bigint",
					},
				},
				&ColumnDefinition{
					ColumnName: "bigint2",
					DataType: IntegerType{
						Name:     "bigint",
						FieldLen: "1",
						Unsigned: true,
						Zerofill: true,
					},
					ColumnOptions: ColumnOptions{
						Nullability: "NOT NULL",
						Default:     "1",
					},
				},
				&ColumnDefinition{
					ColumnName: "decimal1",
					DataType: FixedPointType{
						Name: "decimal",
					},
				},
				&ColumnDefinition{
					ColumnName: "decimal2",
					DataType: FixedPointType{
						Name:     "decimal",
						FieldLen: "2",
					},
				},
				&ColumnDefinition{
					ColumnName: "decimal3",
					DataType: FixedPointType{
						Name:       "decimal",
						FieldLen:   "2",
						FieldScale: "1",
						Unsigned:   true,
						Zerofill:   true,
					},
					ColumnOptions: ColumnOptions{
						Nullability: "NOT NULL",
						Default:     "1",
					},
				},
				&ColumnDefinition{
					ColumnName: "float1",
					DataType: FloatingPointType{
						Name: "float",
					},
				},
				&ColumnDefinition{
					ColumnName: "float2",
					DataType: FloatingPointType{
						Name:       "float",
						FieldLen:   "2",
						FieldScale: "1",
						Unsigned:   true,
						Zerofill:   true,
					},
					ColumnOptions: ColumnOptions{
						Nullability: "NOT NULL",
						Default:     "(rand() * rand())",
					},
				},
				&ColumnDefinition{
					ColumnName: "double1",
					DataType: FloatingPointType{
						Name: "double",
					},
				},
				&ColumnDefinition{
					ColumnName: "double2",
					DataType: FloatingPointType{
						Name:       "double",
						FieldLen:   "2",
						FieldScale: "1",
						Unsigned:   true,
						Zerofill:   true,
					},
					ColumnOptions: ColumnOptions{
						Nullability: "NOT NULL",
						Default:     "1.1",
					},
				},
				&ColumnDefinition{
					ColumnName: "double3",
					DataType: FloatingPointType{
						Name: "double",
					},
					ColumnOptions: ColumnOptions{
						GeneratedAs: "(sqrt(`double1` * `double2`))",
					},
				},
			},
		},
		r[0])

	b1, err := os.ReadFile("test/table/numeric/output.sql")

	assert.Equal(t, string(b1), r[0].String())
}

func TestCreateTableWithStringTypes(t *testing.T) {

	f, err := os.Open("test/table/string/input.sql")
	p := NewParser(f)
	r, err := p.Parse()

	require.NoError(t, err)

	assert.Equal(t,
		CreateTableStatement{
			TableName: "t1",
			CreateDefinitions: []interface{}{
				&ColumnDefinition{
					ColumnName: "char1",
					DataType: StringType{
						Name: "char",
					},
				},
				&ColumnDefinition{
					ColumnName: "char2",
					DataType: StringType{
						Name:      "char",
						FieldLen:  "1",
						Charset:   "utf8mb4",
						Collation: "utf8mb4_bin",
					},
					ColumnOptions: ColumnOptions{
						Nullability: "NOT NULL",
						Default:     "'a'",
					},
				},
				&ColumnDefinition{
					ColumnName: "varchar1",
					DataType: StringType{
						Name:     "varchar",
						FieldLen: "1",
					},
				},
				&ColumnDefinition{
					ColumnName: "varchar2",
					DataType: StringType{
						Name:      "varchar",
						FieldLen:  "2",
						Charset:   "utf8mb4",
						Collation: "utf8mb4_bin",
					},
					ColumnOptions: ColumnOptions{
						Nullability: "NOT NULL",
						Default:     "'a'",
					},
				},
				&ColumnDefinition{
					ColumnName: "binary1",
					DataType: StringType{
						Name: "binary",
					},
				},
				&ColumnDefinition{
					ColumnName: "binary2",
					DataType: StringType{
						Name:     "binary",
						FieldLen: "1",
					},
					ColumnOptions: ColumnOptions{
						Nullability: "NOT NULL",
						Default:     "(uuid_to_bin(uuid()))",
					},
				},
				&ColumnDefinition{
					ColumnName: "varbinary1",
					DataType: StringType{
						Name:     "varbinary",
						FieldLen: "1",
					},
				},
				&ColumnDefinition{
					ColumnName: "varbinary2",
					DataType: StringType{
						Name:     "varbinary",
						FieldLen: "2",
					},
					ColumnOptions: ColumnOptions{
						Nullability: "NOT NULL",
						Default:     "'a'",
					},
				},
				&ColumnDefinition{
					ColumnName: "tinyblob1",
					DataType: StringType{
						Name: "tinyblob",
					},
				},
				&ColumnDefinition{
					ColumnName: "tinyblob2",
					DataType: StringType{
						Name: "tinyblob",
					},
					ColumnOptions: ColumnOptions{
						Nullability: "NOT NULL",
					},
				},
				&ColumnDefinition{
					ColumnName: "tinytext1",
					DataType: StringType{
						Name: "tinytext",
					},
					ColumnOptions: ColumnOptions{},
				},
				&ColumnDefinition{
					ColumnName: "tinytext2",
					DataType: StringType{
						Name:      "tinytext",
						Charset:   "utf8mb4",
						Collation: "utf8mb4_bin",
					},
					ColumnOptions: ColumnOptions{
						Nullability: "NOT NULL",
					},
				},
				&ColumnDefinition{
					ColumnName: "blob1",
					DataType: StringType{
						Name: "blob",
					},
				},
				&ColumnDefinition{
					ColumnName: "blob2",
					DataType: StringType{
						Name:     "blob",
						FieldLen: "1",
					},
					ColumnOptions: ColumnOptions{
						Nullability: "NOT NULL",
					},
				},
				&ColumnDefinition{
					ColumnName: "text1",
					DataType: StringType{
						Name: "text",
					},
				},
				&ColumnDefinition{
					ColumnName: "text2",
					DataType: StringType{
						Name:      "text",
						FieldLen:  "1",
						Charset:   "utf8mb4",
						Collation: "utf8mb4_bin",
					},
					ColumnOptions: ColumnOptions{
						Nullability: "NOT NULL",
					},
				},
				&ColumnDefinition{
					ColumnName: "mediumblob1",
					DataType: StringType{
						Name: "mediumblob",
					},
				},
				&ColumnDefinition{
					ColumnName: "mediumblob2",
					DataType: StringType{
						Name: "mediumblob",
					},
					ColumnOptions: ColumnOptions{
						Nullability: "NOT NULL",
					},
				},
				&ColumnDefinition{
					ColumnName: "mediumtext1",
					DataType: StringType{
						Name: "mediumtext",
					},
				},
				&ColumnDefinition{
					ColumnName: "mediumtext2",
					DataType: StringType{
						Name:      "mediumtext",
						Charset:   "utf8mb4",
						Collation: "utf8mb4_bin",
					},
					ColumnOptions: ColumnOptions{
						Nullability: "NOT NULL",
					},
				},
				&ColumnDefinition{
					ColumnName: "longblob1",
					DataType: StringType{
						Name: "longblob",
					},
				},
				&ColumnDefinition{
					ColumnName: "longblob2",
					DataType: StringType{
						Name: "longblob",
					},
					ColumnOptions: ColumnOptions{
						Nullability: "NOT NULL",
					},
				},
				&ColumnDefinition{
					ColumnName: "longtext1",
					DataType: StringType{
						Name: "longtext",
					},
				},
				&ColumnDefinition{
					ColumnName: "longtext2",
					DataType: StringType{
						Name:      "longtext",
						Charset:   "utf8mb4",
						Collation: "utf8mb4_bin",
					},
					ColumnOptions: ColumnOptions{
						Nullability: "NOT NULL",
					},
				},
				&ColumnDefinition{
					ColumnName: "enum1",
					DataType: StringListType{
						Name:   "enum",
						Values: []string{"'a'"},
					},
				},
				&ColumnDefinition{
					ColumnName: "enum2",
					DataType: StringListType{
						Name:      "enum",
						Values:    []string{"'a'", "'b'"},
						Charset:   "utf8mb4",
						Collation: "utf8mb4_bin",
					},
					ColumnOptions: ColumnOptions{
						Nullability: "NOT NULL",
						Default:     "'a'",
					},
				},
				&ColumnDefinition{
					ColumnName: "set1",
					DataType: StringListType{
						Name:   "set",
						Values: []string{"'a'"},
					},
				},
				&ColumnDefinition{
					ColumnName: "set2",
					DataType: StringListType{
						Name:      "set",
						Values:    []string{"'a'", "'b'"},
						Charset:   "utf8mb4",
						Collation: "utf8mb4_bin",
					},
					ColumnOptions: ColumnOptions{
						Nullability: "NOT NULL",
						Default:     "'a'",
					},
				},
			}},
		r[0])

	b1, err := os.ReadFile("test/table/string/output.sql")

	assert.Equal(t, string(b1), r[0].String())
}

func TestCreateTableWithDateAndTimeTypes(t *testing.T) {

	f, err := os.Open("test/table/date/input.sql")
	p := NewParser(f)
	r, err := p.Parse()

	require.NoError(t, err)

	assert.Equal(t,
		CreateTableStatement{
			TableName: "t1",
			CreateDefinitions: []interface{}{
				&ColumnDefinition{
					ColumnName: "date1",
					DataType: DateAndTimeType{
						Name: "date",
					},
				},
				&ColumnDefinition{
					ColumnName: "date2",
					DataType: DateAndTimeType{
						Name: "date",
					},
					ColumnOptions: ColumnOptions{
						Nullability: "NOT NULL",
						Default:     "(CURRENT_DATE + INTERVAL 1 YEAR)",
					},
				},
				&ColumnDefinition{
					ColumnName: "datetime1",
					DataType: DateAndTimeType{
						Name: "datetime",
					},
				},
				&ColumnDefinition{
					ColumnName: "datetime2",
					DataType: DateAndTimeType{
						Name:     "datetime",
						FieldLen: "1",
					},
					ColumnOptions: ColumnOptions{
						Nullability: "NOT NULL",
					},
				},
				&ColumnDefinition{
					ColumnName: "timestamp1",
					DataType: DateAndTimeType{
						Name: "timestamp",
					},
					ColumnOptions: ColumnOptions{},
				},
				&ColumnDefinition{
					ColumnName: "timestamp2",
					DataType: DateAndTimeType{
						Name:     "timestamp",
						FieldLen: "1",
					},
					ColumnOptions: ColumnOptions{
						Nullability: "NOT NULL",
						Default:     "CURRENT_TIMESTAMP",
						OnUpdate:    "CURRENT_TIMESTAMP",
					},
				},
				&ColumnDefinition{
					ColumnName: "time1",
					DataType: DateAndTimeType{
						Name: "time",
					},
					ColumnOptions: ColumnOptions{},
				},
				&ColumnDefinition{
					ColumnName: "time2",
					DataType: DateAndTimeType{
						Name:     "time",
						FieldLen: "1",
					},
					ColumnOptions: ColumnOptions{
						Nullability: "NOT NULL",
					},
				},
				&ColumnDefinition{
					ColumnName: "year1",
					DataType: DateAndTimeType{
						Name: "year",
					},
					ColumnOptions: ColumnOptions{},
				},
				&ColumnDefinition{
					ColumnName: "year2",
					DataType: DateAndTimeType{
						Name:     "year",
						FieldLen: "4",
					},
					ColumnOptions: ColumnOptions{
						Nullability: "NOT NULL",
					},
				},
			},
		},
		r[0])

	b1, err := os.ReadFile("test/table/date/output.sql")

	assert.Equal(t, string(b1), r[0].String())

}

func TestCreateTableWithindexes(t *testing.T) {

	f, err := os.Open("test/table/index/input.sql")
	p := NewParser(f)
	r, err := p.Parse()

	require.NoError(t, err)

	assert.Equal(t,
		CreateTableStatement{
			TableName: "t1",
			CreateDefinitions: []interface{}{
				&ColumnDefinition{
					ColumnName: "int1",
					DataType: IntegerType{
						Name: "int",
					},
				},
				&ColumnDefinition{
					ColumnName: "int2",
					DataType: IntegerType{
						Name: "int",
					},
				},
				&ColumnDefinition{
					ColumnName: "int3",
					DataType: IntegerType{
						Name: "int",
					},
				},
				&ColumnDefinition{
					ColumnName: "varchar1",
					DataType: StringType{
						Name:     "varchar",
						FieldLen: "10",
					},
				},
				&ColumnDefinition{
					ColumnName: "varchar2",
					DataType: StringType{
						Name:     "varchar",
						FieldLen: "10",
					},
				},
				&ColumnDefinition{
					ColumnName: "varchar3",
					DataType: StringType{
						Name:     "varchar",
						FieldLen: "10",
					},
				},
				&IndexDefinition{
					KeyPartList: []KeyPart{
						{Column: "int1"},
					},
				},
				&IndexDefinition{
					IndexName: "idx1",
					KeyPartList: []KeyPart{
						{Column: "int2", Order: "ASC"},
						{Column: "int3", Order: "DESC"},
					},
				},
				&IndexDefinition{
					IndexName: "idx2",
					KeyPartList: []KeyPart{
						{Column: "varchar1", Length: "5"},
					},
					IndexOptions: IndexOptions{
						IndexType:    "BTREE",
						KeyBlockSize: "1",
						Visibility:   "VISIBLE",
						Comment:      "'foo'",
					},
				},
				&FullTextIndexDefinition{
					KeyPartList: []KeyPart{
						{Column: "varchar2"},
					},
					IndexOptions: IndexOptions{},
				},
				&FullTextIndexDefinition{
					IndexName: "idx3",
					KeyPartList: []KeyPart{
						{Column: "varchar3"},
					},
					IndexOptions: IndexOptions{
						KeyBlockSize: "1",
						Parser:       "ngram",
						Visibility:   "VISIBLE",
						Comment:      "'foo'",
					},
				},
			},
		},
		r[0])

	b1, err := os.ReadFile("test/table/index/output.sql")

	assert.Equal(t, string(b1), r[0].String())
}

func TestCreateTableWithConstraints(t *testing.T) {

	f, err := os.Open("test/table/constraint/input.sql")
	p := NewParser(f)
	r, err := p.Parse()

	require.NoError(t, err)

	assert.Equal(t,
		CreateTableStatement{
			TableName: "t2",
			CreateDefinitions: []interface{}{
				&ColumnDefinition{
					ColumnName: "int1",
					DataType: IntegerType{
						Name: "int",
					},
				},
				&ColumnDefinition{
					ColumnName: "int2",
					DataType: IntegerType{
						Name: "int",
					},
				},
				&PrimaryKeyDefinition{
					KeyPartList: []KeyPart{
						{Column: "int1", Order: "DESC"},
					},
				},
				&UniqueKeyDefinition{
					KeyPartList: []KeyPart{
						{Column: "int2", Order: "DESC"},
					},
				},
				&ForeignKeyDefinition{
					KeyPartList: []KeyPart{
						{Column: "int1"},
					},
					ReferenceDefinition: ReferenceDefinition{
						TableName: "t1",
						KeyPartList: []KeyPart{
							{Column: "int2"},
						},
					},
				},
			},
		},
		r[1])

	assert.Equal(t,
		CreateTableStatement{
			TableName: "t3",
			CreateDefinitions: []interface{}{
				&ColumnDefinition{
					ColumnName: "int1",
					DataType: IntegerType{
						Name: "int",
					},
				},
				&ColumnDefinition{
					ColumnName: "int2",
					DataType: IntegerType{
						Name: "int",
					},
				},
				&PrimaryKeyDefinition{
					ConstraintName: "u1",
					KeyPartList: []KeyPart{
						{Column: "int1"},
						{Column: "int2"},
					},
					IndexOptions: IndexOptions{
						IndexType:    "BTREE",
						KeyBlockSize: "1",
						Visibility:   "VISIBLE",
						Comment:      "'foo'",
					},
				},
				&UniqueKeyDefinition{
					ConstraintName: "u2",
					KeyPartList: []KeyPart{
						{Column: "int1", Order: "ASC"},
						{Column: "int2"},
					},
					IndexOptions: IndexOptions{
						IndexType:    "BTREE",
						KeyBlockSize: "1",
						Visibility:   "VISIBLE",
						Comment:      "'foo'",
					},
				},
				&ForeignKeyDefinition{
					ConstraintName: "u3",
					IndexName:      "i3",
					KeyPartList: []KeyPart{
						{Column: "int1", Order: "ASC"},
						{Column: "int2"},
					},
					ReferenceDefinition: ReferenceDefinition{
						TableName: "t1",
						KeyPartList: []KeyPart{
							{Column: "int1"},
							{Column: "int2"},
						},
						ReferenceOptions: ReferenceOptions{
							Match:    "FULL",
							OnUpdate: "CASCADE",
							OnDelete: "RESTRICT",
						},
					},
				},
			},
		},
		r[2])

	b1, err := os.ReadFile("test/table/constraint/output1.sql")
	b2, err := os.ReadFile("test/table/constraint/output2.sql")

	assert.Equal(t, string(b1), r[1].String())
	assert.Equal(t, string(b2), r[2].String())
}

func TestCreateTableWithOtherTypes(t *testing.T) {

	f, err := os.Open("test/table/others/input.sql")
	p := NewParser(f)
	r, err := p.Parse()

	require.NoError(t, err)

	assert.Equal(t,
		CreateTableStatement{
			TableName: "t1",
			CreateDefinitions: []interface{}{
				&ColumnDefinition{
					ColumnName: "json1",
					DataType: JsonType{
						Name: "json",
					},
				},
				&ColumnDefinition{
					ColumnName: "geometry1",
					DataType: SpatialType{
						Name: "geometry",
					},
				},
				&ColumnDefinition{
					ColumnName: "point1",
					DataType: SpatialType{
						Name: "point",
					},
				},
				&ColumnDefinition{
					ColumnName: "linestring1",
					DataType: SpatialType{
						Name: "linestring",
					},
				},
				&ColumnDefinition{
					ColumnName: "polygon1",
					DataType: SpatialType{
						Name: "polygon",
					},
				},
				&ColumnDefinition{
					ColumnName: "multipoint1",
					DataType: SpatialType{
						Name: "multipoint",
					},
				},
				&ColumnDefinition{
					ColumnName: "multilinestring1",
					DataType: SpatialType{
						Name: "multilinestring",
					},
				},
				&ColumnDefinition{
					ColumnName: "multipolygon1",
					DataType: SpatialType{
						Name: "multipolygon",
					},
				},
				&ColumnDefinition{
					ColumnName: "geometrycollection1",
					DataType: SpatialType{
						Name: "geometrycollection",
					},
				},
			},
		},
		r[0])

	b1, err := os.ReadFile("test/table/others/output.sql")

	assert.Equal(t, string(b1), r[0].String())
}

func TestCreateTableWithCheckConstraints(t *testing.T) {

	f, err := os.Open("test/table/check/input.sql")
	p := NewParser(f)
	_, _ = p.Parse()

	require.NoError(t, err)

}

func TestCreateTableWithPartitions(t *testing.T) {

	f, err := os.Open("test/table/partition/input.sql")
	p := NewParser(f)
	r, err := p.Parse()

	require.NoError(t, err)

	assert.Equal(t,
		CreateTableStatement{
			TableName: "t1",
			CreateDefinitions: []interface{}{
				&ColumnDefinition{
					ColumnName: "int1",
					DataType: IntegerType{
						Name: "int",
					},
				},
				&PrimaryKeyDefinition{
					KeyPartList: []KeyPart{{Column: "int1"}},
				},
			},
			Partitions: PartitionConfig{
				PartitionBy: PartitionBy{
					Type:       "HASH",
					Expression: "`int1`",
				},
				PartitionDefinitions: []PartitionDefinition{},
			},
		}, r[0])

	b1, err := os.ReadFile("test/table/partition/output1.sql")
	assert.Equal(t, string(b1), r[0].String())

	assert.Equal(t,
		CreateTableStatement{
			TableName: "t2",
			CreateDefinitions: []interface{}{
				&ColumnDefinition{
					ColumnName: "int1",
					DataType: IntegerType{
						Name: "int",
					},
				},
				&ColumnDefinition{
					ColumnName: "date1",
					DataType: DateAndTimeType{
						Name: "date",
					},
				},
			},
			Partitions: PartitionConfig{
				PartitionBy: PartitionBy{
					Type:       "RANGE",
					Expression: "year(`date1`)",
				},
				Partitions: "3",
				SubpartitionBy: PartitionBy{
					Type:       "HASH",
					Expression: "to_days(`date1`)",
				},
				Subpartitions: "2",
				PartitionDefinitions: []PartitionDefinition{
					{
						Name:            "p0",
						Operator:        "LESS THAN",
						ValueExpression: "1990",
						Subpartitions:   []SubpartitionDefinition{},
					},
					{
						Name:            "p1",
						Operator:        "LESS THAN",
						ValueExpression: "2000",
						Subpartitions:   []SubpartitionDefinition{},
					},
					{
						Name:            "p2",
						Operator:        "LESS THAN",
						ValueExpression: "MAXVALUE",
						Subpartitions:   []SubpartitionDefinition{},
					},
				},
			},
		}, r[1])

	b2, err := os.ReadFile("test/table/partition/output2.sql")
	assert.Equal(t, string(b2), r[1].String())

	assert.Equal(t,
		CreateTableStatement{
			TableName: "t3",
			CreateDefinitions: []interface{}{
				&ColumnDefinition{
					ColumnName: "double1",
					DataType: FloatingPointType{
						Name: "double",
					},
				},
				&ColumnDefinition{
					ColumnName: "double2",
					DataType: FloatingPointType{
						Name: "double",
					},
				},
			},
			Partitions: PartitionConfig{
				PartitionBy: PartitionBy{
					Type:    "RANGE",
					Columns: []string{`double1`, `double2`},
				},
				Partitions: "3",
				PartitionDefinitions: []PartitionDefinition{
					{
						Name:            "p0",
						Operator:        "LESS THAN",
						ValueExpression: "1990.1",
						Subpartitions:   []SubpartitionDefinition{},
					},
					{
						Name:            "p1",
						Operator:        "LESS THAN",
						ValueExpression: "2000.1",
						Subpartitions:   []SubpartitionDefinition{},
					},
					{
						Name:            "p2",
						Operator:        "LESS THAN",
						ValueExpression: "MAXVALUE",
						Subpartitions:   []SubpartitionDefinition{},
					},
				},
			},
		}, r[2])

	b3, err := os.ReadFile("test/table/partition/output3.sql")
	assert.Equal(t, string(b3), r[2].String())
}
