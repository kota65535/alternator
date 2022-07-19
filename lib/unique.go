package lib

import (
	"fmt"
	"github.com/emirpasic/gods/sets/linkedhashset"
	"github.com/kota65535/alternator/parser"
	"reflect"
	"sort"
)

type UniqueKeyAlterations struct {
	Added       []*AddedUniqueKey
	Modified    []*ModifiedUniqueKey
	Dropped     []*DroppedUniqueKey
	Renamed     []*RenamedUniqueKey
	Retained    []*RetainedUniqueKey
	alterations []Alteration
}

func NewUniqueAlterations(
	from []*parser.UniqueKeyDefinition,
	to []*parser.UniqueKeyDefinition,
	columnOrder map[string]int) UniqueKeyAlterations {

	fromMap := map[string]*parser.UniqueKeyDefinition{}
	fromSet := linkedhashset.New()
	for _, t := range from {
		keys := t.StringKeyPartList()
		fromMap[keys] = t
		fromSet.Add(keys)
	}
	toMap := map[string]*parser.UniqueKeyDefinition{}
	toSet := linkedhashset.New()
	for _, t := range to {
		keys := t.StringKeyPartList()
		toMap[keys] = t
		toSet.Add(keys)
	}

	uniqueKeyOrder := getUniqueKeyOrder(from, to, columnOrder)

	var added []*AddedUniqueKey
	var dropped []*DroppedUniqueKey
	var modified []*ModifiedUniqueKey
	var retained []*RetainedUniqueKey
	var renamed []*RenamedUniqueKey

	for _, v := range difference(fromSet, toSet).Values() {
		s := v.(string)
		dropped = append(dropped, &DroppedUniqueKey{
			This:       fromMap[s],
			Sequential: Sequential{uniqueKeyOrder[s]},
		})
	}
	for _, v := range difference(toSet, fromSet).Values() {
		s := v.(string)
		added = append(added, &AddedUniqueKey{
			This:       toMap[s],
			Sequential: Sequential{uniqueKeyOrder[s]},
		})
	}
	for _, v := range intersection(fromSet, toSet).Values() {
		s := v.(string)
		t1 := fromMap[s]
		t2 := toMap[s]
		if uniqueDefsEqual(*t1, *t2) {
			retained = append(retained, &RetainedUniqueKey{
				This:       t2,
				Sequential: Sequential{uniqueKeyOrder[s]},
			})
		} else if t2.ConstraintName != "" && t1.ConstraintName != t2.ConstraintName {
			dropped = append(dropped, &DroppedUniqueKey{
				This:       fromMap[s],
				Sequential: Sequential{uniqueKeyOrder[s]},
			})
			added = append(added, &AddedUniqueKey{
				This:       toMap[s],
				Sequential: Sequential{uniqueKeyOrder[s]},
			})
		} else if t2.IndexName != "" && t1.IndexName != t2.IndexName {
			renamed = append(renamed, &RenamedUniqueKey{
				From:       t1,
				To:         t2,
				Sequential: Sequential{uniqueKeyOrder[s]},
			})
		} else {
			modified = append(modified, &ModifiedUniqueKey{
				From:       t1,
				To:         t2,
				Sequential: Sequential{uniqueKeyOrder[s]},
			})
		}
	}

	return UniqueKeyAlterations{
		Added:    added,
		Modified: modified,
		Dropped:  dropped,
		Renamed:  renamed,
		Retained: retained,
	}
}

func (r UniqueKeyAlterations) Statements() []string {
	ret := []string{}
	for _, b := range r.Alterations() {
		ret = append(ret, b.Statements()...)
	}
	ret = parser.Align(ret)
	return ret
}

func (r UniqueKeyAlterations) Diff() []string {
	ret := []string{}
	for _, b := range r.Alterations() {
		ret = append(ret, b.Diff()...)
	}
	ret = parser.Align(ret)
	return ret
}

func (r *UniqueKeyAlterations) Alterations() []Alteration {
	if r.alterations != nil {
		return r.alterations
	}
	alterations := []Alteration{}
	for _, a := range r.Dropped {
		for _, b := range r.Added {
			if a.This.IndexName == b.This.IndexName {
				b.AddDependsOn(b)
			}
		}
		alterations = append(alterations, a)
	}
	for _, a := range r.Renamed {
		for _, b := range r.Added {
			if a.From.IndexName == b.This.IndexName {
				b.AddDependsOn(b)
			}
		}
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

func (r *UniqueKeyAlterations) Equivalent() bool {
	return len(r.Dropped) == 0 && len(r.Renamed) == 0 && len(r.Modified) == 0 && len(r.Added) == 0
}

// HandleColumnDrop ensures that dropping a unique key is run before dropping a column
func (r *UniqueKeyAlterations) HandleColumnDrop(drop Alteration, columnName string) {
	for _, d := range r.Dropped {
		if keyPartContains(d.This.KeyPartList, columnName) {
			drop.AddDependsOn(d)
		}
	}
}

type AddedUniqueKey struct {
	This *parser.UniqueKeyDefinition
	Sequential
	Dependent
	Prefixable
}

func (r AddedUniqueKey) Statements() []string {
	return []string{fmt.Sprintf("ADD %s", r.This.String())}
}

func (r AddedUniqueKey) Diff() []string {
	return []string{fmt.Sprintf("+ %s", r.This.String())}
}

func (r AddedUniqueKey) Id() string {
	return keyPartId(r.This.KeyPartList)
}

type ModifiedUniqueKey struct {
	From *parser.UniqueKeyDefinition
	To   *parser.UniqueKeyDefinition
	Sequential
	Dependent
	Prefixable
}

func (r ModifiedUniqueKey) Statements() []string {
	indexName := r.From.IndexName
	if indexName == "" {
		indexName = fmt.Sprintf("<unknown index name of '%s'>", r.From.StringKeyPartList())
	}
	return []string{fmt.Sprintf("ALTER INDEX `%s`%s", indexName,
		optS(r.To.IndexOptions.Diff(r.From.IndexOptions).String(), " %s"))}
}

func (r ModifiedUniqueKey) Diff() []string {
	return []string{fmt.Sprintf("~ %s\t-> %s", r.From.String(), r.To.String())}
}

func (r ModifiedUniqueKey) Id() string {
	return keyPartId(r.To.KeyPartList)
}

type DroppedUniqueKey struct {
	This *parser.UniqueKeyDefinition
	Sequential
	Dependent
	Prefixable
}

func (r DroppedUniqueKey) Statements() []string {
	indexName := r.This.IndexName
	if indexName == "" {
		indexName = fmt.Sprintf("<unknown index name of '%s'>", r.This.StringKeyPartList())
	}
	return []string{fmt.Sprintf("DROP INDEX `%s`", indexName)}
}

func (r DroppedUniqueKey) Diff() []string {
	return []string{fmt.Sprintf("- %s", r.This.String())}
}

func (r DroppedUniqueKey) Id() string {
	return keyPartId(r.This.KeyPartList)
}

type RenamedUniqueKey struct {
	From *parser.UniqueKeyDefinition
	To   *parser.UniqueKeyDefinition
	Sequential
	Dependent
	Prefixable
}

func (r RenamedUniqueKey) Statements() []string {
	return []string{fmt.Sprintf("RENAME INDEX `%s` TO `%s`", r.From.IndexName, r.To.IndexName)}
}

func (r RenamedUniqueKey) Diff() []string {
	return []string{fmt.Sprintf("~ %s\t-> %s", r.From.String(), r.To.String())}
}

func (r RenamedUniqueKey) Id() string {
	return keyPartId(r.From.KeyPartList)
}

type RetainedUniqueKey struct {
	This *parser.UniqueKeyDefinition
	Sequential
	Dependent
	Prefixable
}

func (r RetainedUniqueKey) Statements() []string {
	return []string{}
}

func (r RetainedUniqueKey) Diff() []string {
	return []string{fmt.Sprintf("  %s", r.This.String())}
}

func (r RetainedUniqueKey) Id() string {
	return keyPartId(r.This.KeyPartList)
}

func getUniqueKeyOrder(from []*parser.UniqueKeyDefinition, to []*parser.UniqueKeyDefinition, columnOrder map[string]int) map[string]int {
	all := []*parser.UniqueKeyDefinition{}
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

// Compare 2 foreign key definitions.
// index name and constraint name of 'to' is ignored on comparison if they are empty
func uniqueDefsEqual(from parser.UniqueKeyDefinition, to parser.UniqueKeyDefinition) bool {
	fi := from.IndexName
	ti := to.IndexName
	fc := from.ConstraintName
	tc := to.ConstraintName
	from.IndexName = ""
	to.IndexName = ""
	from.ConstraintName = ""
	to.ConstraintName = ""
	return reflect.DeepEqual(from, to) && (ti == "" || fi == ti) && (tc == "" || fc == tc)
}
