package lib

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"strings"
	"testing"
)

func TestGetAlteredTableOptions(t *testing.T) {
	alt := getAlteredDatabases(t, "test/table/option/from.sql", "test/table/option/to.sql")
	statements := alt.Statements()
	diff := alt.Diff()
	for _, s := range diff {
		fmt.Println(s)
	}

	b1, err := ioutil.ReadFile("test/table/option/alter.sql")
	require.NoError(t, err)
	assert.Equal(t, string(b1), strings.Join(statements, "\n"))

	b2, err := ioutil.ReadFile("test/table/option/diff.txt")
	require.NoError(t, err)
	assert.Equal(t, string(b2), strings.Join(diff, "\n"))
}
