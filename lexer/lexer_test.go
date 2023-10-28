package lexer

import (
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"strings"
	"testing"
)

var Skipped = []TokenType{
	NewRegexpTokenType(-1, `\s`, 0),
	NewRegexpTokenType(-1, `#.*\n`, 0),
	NewRegexpTokenType(-1, `/\*[^!].*?\*/`, 0),
	NewRegexpTokenType(-1, `/\*!\d{5}`, 0),
	NewRegexpTokenType(-1, `\*/`, 0),
}

func TestMain(m *testing.M) {
	logrus.SetLevel(logrus.DebugLevel)
	code := m.Run()
	os.Exit(code)
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
			NewRegexpTokenType(2, "[a-zA-Z0-9]+", 2),
			NewRegexpTokenType(3, "MATCH (FULL|PARTIAL)", 1),
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

func TestSkippedTokenType(t *testing.T) {

	schema := "CREATE /* comment */ table `t1` /*!40100 (`id`) */"

	l := NewLexer(strings.NewReader(schema),
		[]TokenType{
			NewSimpleTokenType(1, "CREATE", true, 1),
			NewSimpleTokenType(2, "TABLE", true, 1),
			NewSimpleTokenType(3, "(", true, 1),
			NewSimpleTokenType(4, ")", true, 1),
			NewRegexpTokenType(5, "`[a-zA-Z0-9]+`", 2),
		},
		Skipped,
	)

	t1, err := l.Scan()
	t2, err := l.Scan()
	t3, err := l.Scan()
	t4, err := l.Scan()
	t5, err := l.Scan()
	t6, err := l.Scan()

	require.NoError(t, err)
	assert.Equal(t, 1, int(t1.Type.GetID()))
	assert.Equal(t, "CREATE", t1.Literal)
	assert.Equal(t, 2, int(t2.Type.GetID()))
	assert.Equal(t, "TABLE", t2.Literal)
	assert.Equal(t, 5, int(t3.Type.GetID()))
	assert.Equal(t, "`t1`", t3.Literal)
	assert.Equal(t, 3, int(t4.Type.GetID()))
	assert.Equal(t, "(", t4.Literal)
	assert.Equal(t, 5, int(t5.Type.GetID()))
	assert.Equal(t, "`id`", t5.Literal)
	assert.Equal(t, 4, int(t6.Type.GetID()))
	assert.Equal(t, ")", t6.Literal)
}
