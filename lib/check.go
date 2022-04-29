package lib

import (
	"fmt"
	"github.com/emirpasic/gods/sets/linkedhashset"
	"github.com/kota65535/alternator/parser"
	"reflect"
)

type CheckConstraintsAlterations struct {
	Added       []*AddedCheckConstraint
	Modified    []*ModifiedCheckConstraint
	Dropped     []*DroppedCheckConstraint
	Retained    []*RetainedCheckConstraint
	alterations []Alteration
}

func NewCheckConstraintAlterations(
	from []*parser.CheckConstraintDefinition,
	to []*parser.CheckConstraintDefinition) CheckConstraintsAlterations {

	fromMap := map[string]*parser.CheckConstraintDefinition{}
	fromSet := linkedhashset.New()
	for _, t := range from {
		fromMap[t.Check] = t
		fromSet.Add(t.Check)
	}
	toMap := map[string]*parser.CheckConstraintDefinition{}
	toSet := linkedhashset.New()
	for _, t := range to {
		toMap[t.Check] = t
		toSet.Add(t.Check)
	}

	checkConstraintOrder := getCheckConstraintOrder(from, to)

	var added []*AddedCheckConstraint
	var dropped []*DroppedCheckConstraint
	var modified []*ModifiedCheckConstraint
	var retained []*RetainedCheckConstraint

	for _, v := range difference(fromSet, toSet).Values() {
		s := v.(string)
		dropped = append(dropped, &DroppedCheckConstraint{
			This:       fromMap[s],
			Sequential: Sequential{checkConstraintOrder[s]},
		})
	}
	for _, v := range difference(toSet, fromSet).Values() {
		s := v.(string)
		added = append(added, &AddedCheckConstraint{
			This:       toMap[s],
			Sequential: Sequential{checkConstraintOrder[s]},
		})
	}
	for _, v := range intersection(fromSet, toSet).Values() {
		s := v.(string)
		t1 := fromMap[s]
		t2 := toMap[s]
		if checkConstraintDefsEqual(*t1, *t2) {
			retained = append(retained, &RetainedCheckConstraint{
				This:       t1,
				Sequential: Sequential{checkConstraintOrder[s]},
			})
		} else if t2.ConstraintName != "" && t1.ConstraintName != t2.ConstraintName {
			dropped = append(dropped, &DroppedCheckConstraint{
				This:       fromMap[s],
				Sequential: Sequential{checkConstraintOrder[s]},
			})
			added = append(added, &AddedCheckConstraint{
				This:       toMap[s],
				Sequential: Sequential{checkConstraintOrder[s]},
			})
		} else {
			modified = append(modified, &ModifiedCheckConstraint{
				From:       t1,
				To:         t2,
				Sequential: Sequential{checkConstraintOrder[s]},
			})
		}
	}

	return CheckConstraintsAlterations{
		Added:    added,
		Modified: modified,
		Dropped:  dropped,
		Retained: retained,
	}
}

func (r CheckConstraintsAlterations) Statements() []string {
	ret := []string{}
	for _, b := range r.Alterations() {
		ret = append(ret, b.Statements()...)
	}
	ret = parser.Align(ret)
	return ret
}

func (r CheckConstraintsAlterations) Diff() []string {
	ret := []string{}
	for _, b := range r.Alterations() {
		ret = append(ret, b.Diff()...)
	}
	ret = parser.Align(ret)
	return ret
}

func (r *CheckConstraintsAlterations) Alterations() []Alteration {
	if r.alterations != nil {
		return r.alterations
	}
	alterations := []Alteration{}
	for _, a := range r.Dropped {
		alterations = append(alterations, a)
	}
	for _, a := range r.Modified {
		alterations = append(alterations, a)
	}
	for _, a := range r.Added {
		alterations = append(alterations, a)
	}
	for _, a := range r.Retained {
		alterations = append(alterations, a)
	}
	r.alterations = NewDag(alterations).Sort()
	return r.alterations
}

func (r *CheckConstraintsAlterations) Equivalent() bool {
	return len(r.Dropped) == 0 && len(r.Modified) == 0 && len(r.Added) == 0
}

type AddedCheckConstraint struct {
	This *parser.CheckConstraintDefinition
	Sequential
	Dependent
	Prefixable
}

func (r AddedCheckConstraint) Statements() []string {
	return []string{fmt.Sprintf("ADD %s", r.This.String())}
}

func (r AddedCheckConstraint) Diff() []string {
	return []string{fmt.Sprintf("+ %s", r.This.String())}
}

func (r AddedCheckConstraint) Id() string {
	return r.This.Check
}

type ModifiedCheckConstraint struct {
	From *parser.CheckConstraintDefinition
	To   *parser.CheckConstraintDefinition
	Sequential
	Dependent
	Prefixable
}

func (r ModifiedCheckConstraint) Statements() []string {
	constraintName := r.From.ConstraintName
	if constraintName == "" {
		constraintName = fmt.Sprintf("<unknown constraint name of '%s'>", r.From.Check)
	}
	return []string{fmt.Sprintf("ALTER CHECK `%s`%s", constraintName,
		optS(r.To.CheckConstraintOptions.Diff(r.From.CheckConstraintOptions).String(), " %s"))}
}

func (r ModifiedCheckConstraint) Diff() []string {
	return []string{fmt.Sprintf("~ %s\t-> %s", r.From.String(), r.To.String())}
}

func (r ModifiedCheckConstraint) Id() string {
	return r.To.Check
}

type DroppedCheckConstraint struct {
	This *parser.CheckConstraintDefinition
	Sequential
	Dependent
	Prefixable
}

func (r DroppedCheckConstraint) Statements() []string {
	constraintName := r.This.ConstraintName
	if constraintName == "" {
		constraintName = fmt.Sprintf("<unknown constraint name of (%s)>", r.This.Check)
	}
	return []string{fmt.Sprintf("DROP CHECK `%s`", constraintName)}
}

func (r DroppedCheckConstraint) Diff() []string {
	return []string{fmt.Sprintf("- %s", r.This.String())}
}

func (r DroppedCheckConstraint) Id() string {
	return r.This.Check
}

type RetainedCheckConstraint struct {
	This *parser.CheckConstraintDefinition
	Sequential
	Dependent
	Prefixable
}

func (r RetainedCheckConstraint) Statements() []string {
	return []string{}
}

func (r RetainedCheckConstraint) Diff() []string {
	return []string{fmt.Sprintf("  %s", r.This.String())}
}

func (r RetainedCheckConstraint) Id() string {
	return r.This.Check
}

func getCheckConstraintOrder(from []*parser.CheckConstraintDefinition, to []*parser.CheckConstraintDefinition) map[string]int {
	ret := map[string]int{}
	p1 := 0
	p2 := 0
	seq := 0
	for p1 < len(from) || p2 < len(to) {
		if p1 >= len(from) {
			ret[to[p2].Check] = seq
			p2 += 1
			seq += 1
			continue
		}
		if p2 >= len(to) {
			ret[from[p1].Check] = seq
			p1 += 1
			seq += 1
			continue
		}
		ret[from[p2].Check] = seq
		ret[to[p1].Check] = seq + 1
		p1 += 1
		p2 += 1
		seq += 2
	}
	return ret
}

// Compare 2 check constraint definitions.
// constraint name of 'to' is ignored on comparison if they are empty
func checkConstraintDefsEqual(from parser.CheckConstraintDefinition, to parser.CheckConstraintDefinition) bool {
	fc := from.ConstraintName
	tc := to.ConstraintName
	from.ConstraintName = ""
	to.ConstraintName = ""
	return reflect.DeepEqual(from, to) && (tc == "" || fc == tc)
}
