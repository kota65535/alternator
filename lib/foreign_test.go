package lib

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"strings"
	"testing"
)

func TestForeignKeys(t *testing.T) {
	alt := getAlteredDatabases(t, "test/table/foreign/1/from.sql", "test/table/foreign/1/to.sql")
	statements := alt.Statements()
	fmt.Println("========== Statements ==========")
	for _, s := range alt.Statements() {
		fmt.Println(s)
	}
	fmt.Println("=====================")
	diff := alt.Diff()
	diffFrom := alt.FromString()
	diffTo := alt.ToString()
	for _, s := range diff {
		fmt.Println(s)
	}
	fmt.Println("==========")
	for _, s := range diffFrom {
		fmt.Println(s)
	}
	fmt.Println("==========")
	for _, s := range diffTo {
		fmt.Println(s)
	}

	b1, err := ioutil.ReadFile("test/table/foreign/1/alter.sql")
	require.NoError(t, err)
	assert.Equal(t, string(b1), strings.Join(statements, "\n"))

	b2, err := ioutil.ReadFile("test/table/foreign/1/diff.txt")
	require.NoError(t, err)
	assert.Equal(t, string(b2), strings.Join(diff, "\n"))

	b3, err := ioutil.ReadFile("test/table/foreign/1/diff_from.txt")
	require.NoError(t, err)
	assert.Equal(t, string(b3), strings.Join(diffFrom, "\n"))

	b4, err := ioutil.ReadFile("test/table/foreign/1/diff_to.txt")
	require.NoError(t, err)
	assert.Equal(t, string(b4), strings.Join(diffTo, "\n"))
}

func TestForeignKeysWithTableRename(t *testing.T) {
	alt := getAlteredDatabases(t, "test/table/foreign/2/from.sql", "test/table/foreign/2/to.sql")
	statements := alt.Statements()
	diff := alt.Diff()
	diffFrom := alt.FromString()
	diffTo := alt.ToString()
	for _, s := range diff {
		fmt.Println(s)
	}
	fmt.Println("==========")
	for _, s := range diffFrom {
		fmt.Println(s)
	}
	fmt.Println("==========")
	for _, s := range diffTo {
		fmt.Println(s)
	}

	b1, err := ioutil.ReadFile("test/table/foreign/2/alter.sql")
	require.NoError(t, err)
	assert.Equal(t, string(b1), strings.Join(statements, "\n"))

	b2, err := ioutil.ReadFile("test/table/foreign/2/diff.txt")
	require.NoError(t, err)
	assert.Equal(t, string(b2), strings.Join(diff, "\n"))

	b3, err := ioutil.ReadFile("test/table/foreign/2/diff_from.txt")
	require.NoError(t, err)
	assert.Equal(t, string(b3), strings.Join(diffFrom, "\n"))

	b4, err := ioutil.ReadFile("test/table/foreign/2/diff_to.txt")
	require.NoError(t, err)
	assert.Equal(t, string(b4), strings.Join(diffTo, "\n"))
}

func TestForeignKeysWithColumnRename(t *testing.T) {
	alt := getAlteredDatabases(t, "test/table/foreign/3/from.sql", "test/table/foreign/3/to.sql")
	statements := alt.Statements()
	diff := alt.Diff()
	diffFrom := alt.FromString()
	diffTo := alt.ToString()
	for _, s := range diff {
		fmt.Println(s)
	}
	fmt.Println("==========")
	for _, s := range diffFrom {
		fmt.Println(s)
	}
	fmt.Println("==========")
	for _, s := range diffTo {
		fmt.Println(s)
	}

	b1, err := ioutil.ReadFile("test/table/foreign/3/alter.sql")
	require.NoError(t, err)
	assert.Equal(t, string(b1), strings.Join(statements, "\n"))

	b2, err := ioutil.ReadFile("test/table/foreign/3/diff.txt")
	require.NoError(t, err)
	assert.Equal(t, string(b2), strings.Join(diff, "\n"))

	b3, err := ioutil.ReadFile("test/table/foreign/3/diff_from.txt")
	require.NoError(t, err)
	assert.Equal(t, string(b3), strings.Join(diffFrom, "\n"))

	b4, err := ioutil.ReadFile("test/table/foreign/3/diff_to.txt")
	require.NoError(t, err)
	assert.Equal(t, string(b4), strings.Join(diffTo, "\n"))
}

func TestForeignKeysWithColumnModification(t *testing.T) {
	alt := getAlteredDatabases(t, "test/table/foreign/4/from.sql", "test/table/foreign/4/to.sql")
	statements := alt.Statements()
	diff := alt.Diff()
	diffFrom := alt.FromString()
	diffTo := alt.ToString()
	for _, s := range diff {
		fmt.Println(s)
	}
	fmt.Println("==========")
	for _, s := range diffFrom {
		fmt.Println(s)
	}
	fmt.Println("==========")
	for _, s := range diffTo {
		fmt.Println(s)
	}

	b1, err := ioutil.ReadFile("test/table/foreign/4/alter.sql")
	require.NoError(t, err)
	assert.Equal(t, string(b1), strings.Join(statements, "\n"))

	b2, err := ioutil.ReadFile("test/table/foreign/4/diff.txt")
	require.NoError(t, err)
	assert.Equal(t, string(b2), strings.Join(diff, "\n"))

	b3, err := ioutil.ReadFile("test/table/foreign/4/diff_from.txt")
	require.NoError(t, err)
	assert.Equal(t, string(b3), strings.Join(diffFrom, "\n"))

	b4, err := ioutil.ReadFile("test/table/foreign/4/diff_to.txt")
	require.NoError(t, err)
	assert.Equal(t, string(b4), strings.Join(diffTo, "\n"))
}
