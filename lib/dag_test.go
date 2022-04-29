package lib

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type TestAlteration struct {
	id string
	Sequential
	Dependent
	Prefixable
}

func (r TestAlteration) Statements() []string {
	return []string{}
}

func (r TestAlteration) Diff() []string {
	return []string{}
}

func (r TestAlteration) Id() string {
	return r.id
}

func TestDAG(t *testing.T) {
	nodes := []Alteration{
		&TestAlteration{id: "a"},
		&TestAlteration{id: "b"},
		&TestAlteration{id: "c"},
		&TestAlteration{id: "d"},
		&TestAlteration{id: "e"},
	}

	dag := NewDag(nodes)
	sorted := dag.Sort()

	ids := []string{}
	for _, s := range sorted {
		ids = append(ids, s.Id())
	}
	assert.Equal(t, []string{"a", "b", "c", "d", "e"}, ids)
}

func TestDAGWithSequence(t *testing.T) {
	nodes := []Alteration{
		&TestAlteration{id: "a", Sequential: Sequential{3}},
		&TestAlteration{id: "b", Sequential: Sequential{5}},
		&TestAlteration{id: "c", Sequential: Sequential{2}},
		&TestAlteration{id: "d", Sequential: Sequential{4}},
		&TestAlteration{id: "e", Sequential: Sequential{1}},
	}

	dag := NewDag(nodes)
	sorted := dag.Sort()

	ids := []string{}
	for _, s := range sorted {
		ids = append(ids, s.Id())
	}
	assert.Equal(t, []string{"e", "c", "a", "d", "b"}, ids)
}

func TestDAGWithDependencies(t *testing.T) {
	nodes := []Alteration{
		&TestAlteration{id: "a", Sequential: Sequential{3}},
		&TestAlteration{id: "b", Sequential: Sequential{5}},
		&TestAlteration{id: "c", Sequential: Sequential{2}},
		&TestAlteration{id: "d", Sequential: Sequential{4}},
		&TestAlteration{id: "e", Sequential: Sequential{1}},
	}
	// 5 -> 3
	nodes[0].AddDependsOn(nodes[1])
	// 4 -> 2
	nodes[2].AddDependsOn(nodes[3])

	dag := NewDag(nodes)
	sorted := dag.Sort()

	ids := []string{}
	for _, s := range sorted {
		ids = append(ids, s.Id())
	}
	assert.Equal(t, []string{"e", "d", "c", "b", "a"}, ids)
}
