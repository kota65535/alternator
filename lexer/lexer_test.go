package lexer

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

var Skipped = []TokenType{
	NewRegexpTokenType(-1, "\\s", false, 0),
	NewRegexpTokenType(-1, "#.*\n", false, 0),
	NewRegexpTokenType(-1, "--.*\n", false, 0),
	NewRegexpTokenType(-1, "/\\*.*\\*/\n", false, 0),
}

func TestSimpleTokenType(t *testing.T) {

	schema := "CREATE table `t1` (`id` inthogee)"

	l := NewLexer(strings.NewReader(schema),
		[]TokenType{
			NewSimpleTokenType(1, "CREATE", true, 1),
			NewSimpleTokenType(2, "TABLE", true, 1),
		},
		Skipped,
	)

	t1, err := l.Scan()
	t2, err := l.Scan()

	require.NoError(t, err)
	assert.Equal(t, 1, int(t1.Type.GetID()))
	assert.Equal(t, "CREATE", t1.Literal)
	assert.Equal(t, 2, int(t2.Type.GetID()))
	assert.Equal(t, "TABLE", t2.Literal)
}

func TestRegexTokenType(t *testing.T) {

	schema := "REFERENCES tbl MATCH FULL"

	l := NewLexer(strings.NewReader(schema),
		[]TokenType{
			NewSimpleTokenType(1, "REFERENCES", true, 1),
			NewRegexpTokenType(2, "[a-zA-Z0-9]+", true, 2),
			NewRegexpTokenType(3, "MATCH (FULL|PARTIAL)", true, 1),
		},
		Skipped,
	)

	t1, err := l.Scan()
	t2, err := l.Scan()
	t3, err := l.Scan()

	require.NoError(t, err)
	assert.Equal(t, 1, int(t1.Type.GetID()))
	assert.Equal(t, "REFERENCES", t1.Literal)
	assert.Equal(t, 2, int(t2.Type.GetID()))
	assert.Equal(t, "tbl", t2.Literal)
	assert.Equal(t, 3, int(t3.Type.GetID()))
	assert.Equal(t, "MATCH FULL", t3.Literal)
	assert.Equal(t, "FULL", t3.Submatches[0])
}
