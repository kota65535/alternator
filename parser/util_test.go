package parser

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAlign(t *testing.T) {
	lines := []string{
		"`bit`\tBIT(1)\tNOT NULL DEFAULT 1,\n",
		"`int`\tINT,\n",
		"`tinyint`\tTINYINT(1) UNSIGNED ZEROFILL\tNOT NULL DEFAULT 1,\n",
		"`bool`\tBOOLEAN\tNOT NULL DEFAULT 1,\n",
	}

	aligned := Align(lines)

	fmt.Println(aligned)

}

func TestStructDifference(t *testing.T) {

	t1 := TableOptions{
		AutoExtendedSize: "2",
		AvgRowLength:     "3",
	}
	t2 := TableOptions{
		AutoExtendedSize: "1",
		AutoIncrement:    "1",
	}
	diff := structDifference(t1, t2)

	assert.Equal(t, "2", diff.AutoExtendedSize)
	assert.Equal(t, "", diff.AutoIncrement)
	assert.Equal(t, "3", diff.AvgRowLength)
}
