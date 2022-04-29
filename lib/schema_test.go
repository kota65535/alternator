package lib

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStrip(t *testing.T) {
	str := "foo/*!12345 bar */baz/* comment */"

	stripped := stripConditionalComments(str)

	assert.Equal(t, "foo bar baz/* comment */", stripped)
}
