package lib

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"strings"
	"testing"
)

func TestGetAlteredTables(t *testing.T) {
	alt := getAlteredDatabases(t, "test/table/from.sql", "test/table/to.sql")
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

	b1, err := os.ReadFile("test/table/alter.sql")
	require.NoError(t, err)
	assert.Equal(t, string(b1), strings.Join(statements, "\n"))

	b2, err := os.ReadFile("test/table/diff.txt")
	require.NoError(t, err)
	assert.Equal(t, string(b2), strings.Join(diff, "\n"))

	b3, err := os.ReadFile("test/table/diff_from.txt")
	require.NoError(t, err)
	assert.Equal(t, string(b3), strings.Join(diffFrom, "\n"))

	b4, err := os.ReadFile("test/table/diff_to.txt")
	require.NoError(t, err)
	assert.Equal(t, string(b4), strings.Join(diffTo, "\n"))
}
