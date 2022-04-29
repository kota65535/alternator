package lib

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"strings"
	"testing"
)

func TestUniqueKeys(t *testing.T) {
	alt := getAlteredDatabases(t, "test/table/unique/from.sql", "test/table/unique/to.sql")
	statements := alt.Statements()
	diff := alt.Diff()
	for _, s := range diff {
		fmt.Println(s)
	}

	b1, err := ioutil.ReadFile("test/table/unique/alter.sql")
	require.NoError(t, err)
	assert.Equal(t, string(b1), strings.Join(statements, "\n"))

	b2, err := ioutil.ReadFile("test/table/unique/diff.txt")
	require.NoError(t, err)
	assert.Equal(t, string(b2), strings.Join(diff, "\n"))
}
