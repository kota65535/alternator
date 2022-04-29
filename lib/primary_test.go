package lib

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"strings"
	"testing"
)

func TestPrimaryKeyColumnRename(t *testing.T) {
	alt := getAlteredDatabases(t, "test/table/primary/1/from.sql", "test/table/primary/1/to.sql")
	statements := alt.Statements()
	diff := alt.Diff()

	for _, s := range diff {
		fmt.Println(s)
	}

	b1, err := ioutil.ReadFile("test/table/primary/1/alter.sql")
	require.NoError(t, err)
	assert.Equal(t, string(b1), strings.Join(statements, "\n"))

	b2, err := ioutil.ReadFile("test/table/primary/1/diff.txt")
	require.NoError(t, err)
	assert.Equal(t, string(b2), strings.Join(diff, "\n"))
}

func TestPrimaryKeyColumnModification(t *testing.T) {
	alt := getAlteredDatabases(t, "test/table/primary/2/from.sql", "test/table/primary/2/to.sql")
	statements := alt.Statements()
	diff := alt.Diff()

	for _, s := range diff {
		fmt.Println(s)
	}

	b1, err := ioutil.ReadFile("test/table/primary/2/alter.sql")
	require.NoError(t, err)
	assert.Equal(t, string(b1), strings.Join(statements, "\n"))

	b2, err := ioutil.ReadFile("test/table/primary/2/diff.txt")
	require.NoError(t, err)
	assert.Equal(t, string(b2), strings.Join(diff, "\n"))
}
