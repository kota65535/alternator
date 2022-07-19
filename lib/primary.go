package lib

import (
	"fmt"
	"github.com/emirpasic/gods/sets/linkedhashset"
	"github.com/kota65535/alternator/parser"
	"reflect"
	"sort"
)

type PrimaryKeyAlterations struct {
	Added       []*AddedPrimaryKey
	Modified    []*ModifiedPrimaryKey
	Dropped     []*DroppedPrimaryKey
	Retained    []*RetainedPrimaryKey
	alterations []Alteration
}

func NewPrimaryKeyAlterations(
	from []*parser.PrimaryKeyDefinition,
	to []*parser.PrimaryKeyDefinition,
	columnOrder map[string]int) PrimaryKeyAlterations {

	fromMap := map[string]*parser.PrimaryKeyDefinition{}
	fromSet := linkedhashset.New()
	for _, t := range from {
		keys := t.StringKeyPartList()
		fromMap[keys] = t
		fromSet.Add(keys)
	}
	toMap := map[string]*parser.PrimaryKeyDefinition{}
	toSet := linkedhashset.New()
	for _, t := range to {
		keys := t.StringKeyPartList()
		toMap[keys] = t
		toSet.Add(keys)
	}

	primaryKeyOrder := getPrimaryKeyOrder(from, to, columnOrder)

	var added []*AddedPrimaryKey
	var dropped []*DroppedPrimaryKey
	var modified []*ModifiedPrimaryKey
	var retained []*RetainedPrimaryKey

	for _, v := range difference(fromSet, toSet).Values() {
		s := v.(string)
		dropped = append(dropped, &DroppedPrimaryKey{
			This:       fromMap[s],
			Sequential: Sequential{primaryKeyOrder[s]},
		})
	}
	for _, v := range difference(toSet, fromSet).Values() {
		s := v.(string)
		added = append(added, &AddedPrimaryKey{
			This:       toMap[s],
			Sequential: Sequential{primaryKeyOrder[s]},
		})
	}
	for _, v := range intersection(fromSet, toSet).Values() {
		s := v.(string)
		t1 := fromMap[s]
		t2 := toMap[s]
		if primaryKeyDefsEqual(*t1, *t2) {
			retained = append(retained, &RetainedPrimaryKey{
				This:       t2,
				Sequential: Sequential{primaryKeyOrder[s]},
			})
		} else {
			modified = append(modified, &ModifiedPrimaryKey{
				From:       t1,
				To:         t2,
				Sequential: Sequential{primaryKeyOrder[s]},
			})
		}
	}

	return PrimaryKeyAlterations{
		Added:    added,
		Modified: modified,
		Dropped:  dropped,
		Retained: retained,
	}
}

func (r PrimaryKeyAlterations) Statements() []string {
	ret := []string{}
	for _, b := range r.Alterations() {
		ret = append(ret, b.Statements()...)
	}
	ret = parser.Align(ret)
	return ret
}

func (r PrimaryKeyAlterations) Diff() []string {
	ret := []string{}
	for _, b := range r.Alterations() {
		ret = append(ret, b.Diff()...)
	}
	ret = parser.Align(ret)
	return ret
}

func (r *PrimaryKeyAlterations) Alterations() []Alteration {
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

func (r *PrimaryKeyAlterations) Equivalent() bool {
	return len(r.Dropped) == 0 && len(r.Modified) == 0 && len(r.Added) == 0
}

// HandleColumnRename handles key part change caused by column rename.
func (r *PrimaryKeyAlterations) HandleColumnRename(fromName string, toName string) {
	removed := map[Alteration]bool{}
	// Changing key part is first considered as drop & add a foreign key, so we will find the pair
	for _, d := range r.Dropped {
		for _, a := range r.Added {
			if keyPartContains(d.This.KeyPartList, fromName) && keyPartContains(a.This.KeyPartList, toName) {
				// update key parts
				d.This.KeyPartList = keyPartReplace(d.This.KeyPartList, fromName, toName)
				if primaryKeyDefsEqual(*d.This, *a.This) {
					r.Retained = append(r.Retained, &RetainedPrimaryKey{
						This:       a.This,
						Sequential: a.Sequential,
						Dependent:  a.Dependent,
					})
					removed[d] = true
					removed[a] = true
				}
			}
		}
	}
	r.Dropped = RemoveIf(r.Dropped, func(m *DroppedPrimaryKey) bool { return removed[m] })
	r.Added = RemoveIf(r.Added, func(m *AddedPrimaryKey) bool { return removed[m] })
}

func (r *PrimaryKeyAlterations) HandleColumnModify(alt Alteration, columnName string) {
	// Changing key part is first considered as drop & add a foreign key, so we will find the pair
	for _, d := range r.Dropped {
		if keyPartContains(d.This.KeyPartList, columnName) {
			alt.AddDependsOn(d)
		}
	}

	for _, a := range r.Added {
		if keyPartContains(a.This.KeyPartList, columnName) {
			a.AddDependsOn(alt)
		}
	}
}

// HandleColumnRename handles key part change caused by column rename.
func (r *PrimaryKeyAlterations) HandlePrimaryKeyDrop(alt Alteration, fromName string, toName string) {
	removed := map[Alteration]bool{}
	// Changing key part is first considered as drop & add a foreign key, so we will find the pair
	for _, d := range r.Dropped {
		for _, a := range r.Added {
			if keyPartContains(d.This.KeyPartList, fromName) && keyPartContains(a.This.KeyPartList, toName) {
				// update key parts
				d.This.KeyPartList = keyPartReplace(d.This.KeyPartList, fromName, toName)
				if primaryKeyDefsEqual(*d.This, *a.This) {
					r.Retained = append(r.Retained, &RetainedPrimaryKey{
						This:       a.This,
						Sequential: a.Sequential,
						Dependent:  a.Dependent,
					})
					removed[d] = true
					removed[a] = true
				} else {
					// still needs to drop & add, and must be altered after column rename
					a.AddDependsOn(alt)
					d.AddDependsOn(alt)
				}
			}
		}
	}
	r.Dropped = RemoveIf(r.Dropped, func(m *DroppedPrimaryKey) bool { return removed[m] })
	r.Added = RemoveIf(r.Added, func(m *AddedPrimaryKey) bool { return removed[m] })
}

type AddedPrimaryKey struct {
	This *parser.PrimaryKeyDefinition
	Sequential
	Dependent
	Prefixable
}

func (r AddedPrimaryKey) Statements() []string {
	return []string{fmt.Sprintf("ADD %s", r.This.String())}
}

func (r AddedPrimaryKey) Diff() []string {
	return []string{fmt.Sprintf("+ %s", r.This.String())}
}

func (r AddedPrimaryKey) Id() string {
	return keyPartId(r.This.KeyPartList)
}

type ModifiedPrimaryKey struct {
	From *parser.PrimaryKeyDefinition
	To   *parser.PrimaryKeyDefinition
	Sequential
	Dependent
	Prefixable
}

func (r ModifiedPrimaryKey) Statements() []string {
	constraintName := r.From.ConstraintName
	if constraintName == "" {
		constraintName = fmt.Sprintf("<unknown constraint name of '%s'>", r.From.StringKeyPartList())
	}
	return []string{
		"DROP PRIMARY KEY",
		fmt.Sprintf("ADD %s", r.To.String()),
	}
}

func (r ModifiedPrimaryKey) Diff() []string {
	return []string{fmt.Sprintf("~ %s\t-> %s", r.From.String(), r.To.String())}
}

func (r ModifiedPrimaryKey) Id() string {
	return keyPartId(r.To.KeyPartList)
}

type DroppedPrimaryKey struct {
	This *parser.PrimaryKeyDefinition
	Sequential
	Dependent
	Prefixable
}

func (r DroppedPrimaryKey) Statements() []string {
	return []string{"DROP PRIMARY KEY"}
}

func (r DroppedPrimaryKey) Diff() []string {
	return []string{fmt.Sprintf("- %s", r.This.String())}
}

func (r DroppedPrimaryKey) Id() string {
	return keyPartId(r.This.KeyPartList)

}

type RetainedPrimaryKey struct {
	This *parser.PrimaryKeyDefinition
	Sequential
	Dependent
	Prefixable
}

func (r RetainedPrimaryKey) Statements() []string {
	return []string{}
}

func (r RetainedPrimaryKey) Diff() []string {
	return []string{fmt.Sprintf("  %s", r.This.String())}
}

func (r RetainedPrimaryKey) Id() string {
	return keyPartId(r.This.KeyPartList)

}

func getPrimaryKeyOrder(from []*parser.PrimaryKeyDefinition, to []*parser.PrimaryKeyDefinition, columnOrder map[string]int) map[string]int {
	all := []*parser.PrimaryKeyDefinition{}
	all = append(all, from...)
	all = append(all, to...)
	sort.SliceStable(all, func(i, j int) bool {
		if len(all[i].KeyPartList) != len(all[j].KeyPartList) {
			return len(all[i].KeyPartList) < len(all[j].KeyPartList)
		}
		length := len(all[i].KeyPartList)
		for a := 0; a < length; a++ {
			if all[i].KeyPartList[a] != all[j].KeyPartList[a] {
				return columnOrder[all[i].KeyPartList[a].Column] < columnOrder[all[j].KeyPartList[a].Column]
			}
		}
		return true
	})
	ret := map[string]int{}
	for i, a := range all {
		ret[a.StringKeyPartList()] = i
	}
	return ret
}

// Compare 2 primary key definitions.
// constraint name of 'to' is ignored on comparison if they are empty
func primaryKeyDefsEqual(from parser.PrimaryKeyDefinition, to parser.PrimaryKeyDefinition) bool {
	fc := from.ConstraintName
	tc := to.ConstraintName
	from.ConstraintName = ""
	to.ConstraintName = ""
	return reflect.DeepEqual(from, to) && (tc == "" || fc == tc)
}
