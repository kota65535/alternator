package lib

import (
	"fmt"
	"github.com/emirpasic/gods/sets/linkedhashset"
	"github.com/kota65535/alternator/parser"
	"reflect"
	"sort"
)

type FullTextIndexAlterations struct {
	Added       []*AddedFullTextIndex
	Modified    []*ModifiedFullTextIndex
	Dropped     []*DroppedFullTextIndex
	Renamed     []*RenamedFullTextIndex
	Retained    []*RetainedFullTextIndex
	alterations []Alteration
}

func NewFullTextIndexAlteration(
	from []*parser.FullTextIndexDefinition,
	to []*parser.FullTextIndexDefinition,
	columnOrder map[string]int) FullTextIndexAlterations {

	fromMap := map[string]*parser.FullTextIndexDefinition{}
	fromSet := linkedhashset.New()
	for _, t := range from {
		keys := t.StringKeyPartList()
		fromMap[keys] = t
		fromSet.Add(keys)
	}
	toMap := map[string]*parser.FullTextIndexDefinition{}
	toSet := linkedhashset.New()
	for _, t := range to {
		keys := t.StringKeyPartList()
		toMap[keys] = t
		toSet.Add(keys)
	}

	fullTextIndexOrder := getFullTextIndexOrder(from, to, columnOrder)

	var added []*AddedFullTextIndex
	var dropped []*DroppedFullTextIndex
	var modified []*ModifiedFullTextIndex
	var retained []*RetainedFullTextIndex
	var renamed []*RenamedFullTextIndex

	for _, v := range difference(fromSet, toSet).Values() {
		s := v.(string)
		dropped = append(dropped, &DroppedFullTextIndex{
			This:       fromMap[s],
			Sequential: Sequential{fullTextIndexOrder[s]},
		})
	}
	for _, v := range difference(toSet, fromSet).Values() {
		s := v.(string)
		added = append(added, &AddedFullTextIndex{
			This:       toMap[s],
			Sequential: Sequential{fullTextIndexOrder[s]},
		})
	}
	for _, v := range intersection(fromSet, toSet).Values() {
		s := v.(string)
		t1 := fromMap[s]
		t2 := toMap[s]
		if fullTextIndexDefsEqual(*t1, *t2) {
			if t2.IndexName != "" && t1.IndexName != t2.IndexName {
				renamed = append(renamed, &RenamedFullTextIndex{
					From:       t1,
					To:         t2,
					Sequential: Sequential{fullTextIndexOrder[s]},
				})
			} else {
				retained = append(retained, &RetainedFullTextIndex{
					This:       t2,
					Sequential: Sequential{fullTextIndexOrder[s]},
				})
			}
		} else {
			modified = append(modified, &ModifiedFullTextIndex{
				From:       t1,
				To:         t2,
				Sequential: Sequential{fullTextIndexOrder[s]},
			})
		}
	}

	return FullTextIndexAlterations{
		Added:    added,
		Modified: modified,
		Renamed:  renamed,
		Dropped:  dropped,
		Retained: retained,
	}
}

func (r FullTextIndexAlterations) Statements() []string {
	ret := []string{}
	for _, b := range r.Alterations() {
		ret = append(ret, b.Statements()...)
	}
	ret = parser.Align(ret)
	return ret
}

func (r FullTextIndexAlterations) Diff() []string {
	ret := []string{}
	for _, b := range r.Alterations() {
		ret = append(ret, b.Diff()...)
	}
	ret = parser.Align(ret)
	return ret
}

func (r FullTextIndexAlterations) FromString() []string {
	ret := []string{}
	for _, a := range r.Alterations() {
		ret = append(ret, a.FromString()...)
	}
	ret = parser.Align(ret)
	return ret
}

func (r FullTextIndexAlterations) ToString() []string {
	ret := []string{}
	for _, a := range r.Alterations() {
		ret = append(ret, a.ToString()...)
	}
	ret = parser.Align(ret)
	return ret
}

func (r *FullTextIndexAlterations) Alterations() []Alteration {
	if r.alterations != nil {
		return r.alterations
	}
	alterations := []Alteration{}
	for _, a := range r.Dropped {
		alterations = append(alterations, a)
	}
	for _, a := range r.Renamed {
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

func (r *FullTextIndexAlterations) Equivalent() bool {
	return len(r.Dropped) == 0 && len(r.Renamed) == 0 && len(r.Modified) == 0 && len(r.Added) == 0
}

// HandleColumnDrop ensures that dropping index is run before dropping column
func (r *FullTextIndexAlterations) HandleColumnDrop(drop Alteration, columnName string) {
	for _, d := range r.Dropped {
		if keyPartContains(d.This.KeyPartList, columnName) {
			drop.AddDependsOn(d)
		}
	}
}

type AddedFullTextIndex struct {
	This *parser.FullTextIndexDefinition
	Sequential
	Dependent
	Prefixable
}

func (r AddedFullTextIndex) Statements() []string {
	return []string{fmt.Sprintf("ADD %s", r.This.String())}
}

func (r AddedFullTextIndex) Diff() []string {
	return []string{fmt.Sprintf("+ %s", r.This.String())}
}

func (r AddedFullTextIndex) FromString() []string {
	return []string{}
}

func (r AddedFullTextIndex) ToString() []string {
	return []string{r.This.String()}
}

func (r AddedFullTextIndex) Id() string {
	return keyPartId(r.This.KeyPartList)
}

type ModifiedFullTextIndex struct {
	From *parser.FullTextIndexDefinition
	To   *parser.FullTextIndexDefinition
	Sequential
	Dependent
	Prefixable
}

func (r ModifiedFullTextIndex) Statements() []string {
	indexName := r.From.IndexName
	if indexName == "" {
		indexName = fmt.Sprintf("<unknown index name of '%s'>", r.From.StringKeyPartList())
	}
	return []string{fmt.Sprintf("ALTER INDEX `%s`%s", indexName,
		optS(r.To.IndexOptions.Diff(r.From.IndexOptions).String(), " %s"))}
}

func (r ModifiedFullTextIndex) Diff() []string {
	return []string{fmt.Sprintf("~ %s\t-> %s", r.From.String(), r.To.String())}
}

func (r ModifiedFullTextIndex) FromString() []string {
	return []string{r.From.String()}
}

func (r ModifiedFullTextIndex) ToString() []string {
	return []string{r.To.String()}
}

func (r ModifiedFullTextIndex) Id() string {
	return keyPartId(r.From.KeyPartList)
}

type DroppedFullTextIndex struct {
	This *parser.FullTextIndexDefinition
	Sequential
	Dependent
	Prefixable
}

func (r DroppedFullTextIndex) Statements() []string {
	indexName := r.This.IndexName
	if indexName == "" {
		indexName = fmt.Sprintf("<unknown index name of '%s'>", r.This.StringKeyPartList())
	}
	return []string{fmt.Sprintf("DROP INDEX `%s`", indexName)}
}

func (r DroppedFullTextIndex) Diff() []string {
	return []string{fmt.Sprintf("- %s", r.This.String())}
}

func (r DroppedFullTextIndex) FromString() []string {
	return []string{r.This.String()}
}

func (r DroppedFullTextIndex) ToString() []string {
	return []string{}
}

func (r DroppedFullTextIndex) Id() string {
	return keyPartId(r.This.KeyPartList)
}

type RenamedFullTextIndex struct {
	From *parser.FullTextIndexDefinition
	To   *parser.FullTextIndexDefinition
	Sequential
	Dependent
	Prefixable
}

func (r RenamedFullTextIndex) Statements() []string {
	return []string{fmt.Sprintf("RENAME INDEX `%s` TO `%s`", r.From.IndexName, r.To.IndexName)}
}

func (r RenamedFullTextIndex) Diff() []string {
	return []string{fmt.Sprintf("~ %s\t-> %s", r.From.String(), r.To.String())}
}

func (r RenamedFullTextIndex) FromString() []string {
	return []string{r.From.String()}
}

func (r RenamedFullTextIndex) ToString() []string {
	return []string{r.To.String()}
}

func (r RenamedFullTextIndex) Id() string {
	return keyPartId(r.From.KeyPartList)
}

type RetainedFullTextIndex struct {
	This *parser.FullTextIndexDefinition
	Sequential
	Dependent
	Prefixable
}

func (r RetainedFullTextIndex) Statements() []string {
	return []string{}
}

func (r RetainedFullTextIndex) Diff() []string {
	return []string{fmt.Sprintf("  %s", r.This.String())}
}

func (r RetainedFullTextIndex) FromString() []string {
	return []string{r.This.String()}
}

func (r RetainedFullTextIndex) ToString() []string {
	return []string{r.This.String()}
}

func (r RetainedFullTextIndex) Id() string {
	return keyPartId(r.This.KeyPartList)
}

func getFullTextIndexOrder(from []*parser.FullTextIndexDefinition, to []*parser.FullTextIndexDefinition, columnOrder map[string]int) map[string]int {
	all := []*parser.FullTextIndexDefinition{}
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

func fullTextIndexDefsEqual(c1 parser.FullTextIndexDefinition, c2 parser.FullTextIndexDefinition) bool {
	c1.IndexName = ""
	c2.IndexName = ""
	return reflect.DeepEqual(c1, c2)
}
