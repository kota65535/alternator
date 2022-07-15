package lib

import (
	"fmt"
	"github.com/emirpasic/gods/sets/linkedhashset"
	"github.com/kota65535/alternator/parser"
	"sort"
)

type IndexAlterations struct {
	Added       []*AddedIndex
	Modified    []*ModifiedIndex
	Dropped     []*DroppedIndex
	Renamed     []*RenamedIndex
	Retained    []*RetainedIndex
	alterations []Alteration
}

func NewIndexAlterations(
	from []*parser.IndexDefinition,
	to []*parser.IndexDefinition,
	columnOrder map[string]int) IndexAlterations {

	fromMap := map[string]*parser.IndexDefinition{}
	fromSet := linkedhashset.New()
	for _, t := range from {
		keys := t.StringKeyPartList()
		fromMap[keys] = t
		fromSet.Add(keys)
	}
	toMap := map[string]*parser.IndexDefinition{}
	toSet := linkedhashset.New()
	for _, t := range to {
		keys := t.StringKeyPartList()
		toMap[keys] = t
		toSet.Add(keys)
	}

	indexOrder := getIndexOrder(from, to, columnOrder)

	var added []*AddedIndex
	var dropped []*DroppedIndex
	var modified []*ModifiedIndex
	var renamed []*RenamedIndex
	var retained []*RetainedIndex

	for _, v := range difference(fromSet, toSet).Values() {
		s := v.(string)
		dropped = append(dropped, &DroppedIndex{
			This:       fromMap[s],
			Sequential: Sequential{indexOrder[s]},
		})
	}
	for _, v := range difference(toSet, fromSet).Values() {
		s := v.(string)
		added = append(added, &AddedIndex{
			This:       toMap[s],
			Sequential: Sequential{indexOrder[s]},
		})
	}
	for _, v := range intersection(fromSet, toSet).Values() {
		s := v.(string)
		t1 := fromMap[s]
		t2 := toMap[s]
		if t1.EqualsExceptIndexName(*t2) {
			if t2.IndexName != "" && t1.IndexName != t2.IndexName {
				renamed = append(renamed, &RenamedIndex{
					From:       t1,
					To:         t2,
					Sequential: Sequential{indexOrder[s]},
				})
			} else {
				retained = append(retained, &RetainedIndex{
					This:       t2,
					Sequential: Sequential{indexOrder[s]},
				})
			}
		} else {
			modified = append(modified, &ModifiedIndex{
				From:       t1,
				To:         t2,
				Sequential: Sequential{indexOrder[s]},
			})
		}
	}

	return IndexAlterations{
		Added:    added,
		Modified: modified,
		Renamed:  renamed,
		Dropped:  dropped,
		Retained: retained,
	}
}

func (r IndexAlterations) Statements() []string {
	ret := []string{}
	for _, b := range r.Alterations() {
		ret = append(ret, b.Statements()...)
	}
	ret = parser.Align(ret)
	return ret
}

func (r IndexAlterations) Diff() []string {
	ret := []string{}
	for _, b := range r.Alterations() {
		ret = append(ret, b.Diff()...)
	}
	ret = parser.Align(ret)
	return ret
}

func (r *IndexAlterations) Alterations() []Alteration {
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

func (r *IndexAlterations) Equivalent() bool {
	return len(r.Dropped) == 0 && len(r.Renamed) == 0 && len(r.Modified) == 0 && len(r.Added) == 0
}

// HandleColumnDrop ensures that dropping index is run before dropping column
func (r *IndexAlterations) HandleColumnDrop(drop Alteration, columnName string) {
	for _, d := range r.Dropped {
		if keyPartContains(d.This.KeyPartList, columnName) {
			drop.AddDependsOn(d)
		}
	}
}

// HandleColumnRename handles key part change caused by column rename.
func (r *IndexAlterations) HandleColumnRename(fromName string, toName string) {
	removed := map[Alteration]bool{}
	// Changing key part is first considered as drop & add a foreign key, so we will find the pair
	for _, d := range r.Dropped {
		for _, a := range r.Added {
			if keyPartContains(d.This.KeyPartList, fromName) && keyPartContains(a.This.KeyPartList, toName) {
				// update key parts
				d.This.KeyPartList = keyPartReplace(d.This.KeyPartList, fromName, toName)
				if d.This.EqualsExceptIndexName(*a.This) {
					r.Retained = append(r.Retained, &RetainedIndex{
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
	r.Dropped = RemoveIf(r.Dropped, func(m *DroppedIndex) bool { return removed[m] })
	r.Added = RemoveIf(r.Added, func(m *AddedIndex) bool { return removed[m] })
}

type AddedIndex struct {
	This *parser.IndexDefinition
	Sequential
	Dependent
	Prefixable
}

func (r AddedIndex) Statements() []string {
	return []string{fmt.Sprintf("ADD %s", r.This.String())}
}

func (r AddedIndex) Diff() []string {
	return []string{fmt.Sprintf("+ %s", r.This.String())}
}

func (r AddedIndex) Id() string {
	return keyPartId(r.This.KeyPartList)
}

type ModifiedIndex struct {
	From *parser.IndexDefinition
	To   *parser.IndexDefinition
	Sequential
	Dependent
	Prefixable
}

func (r ModifiedIndex) Statements() []string {
	indexName := r.From.IndexName
	if indexName == "" {
		indexName = fmt.Sprintf("<unknown index name of '%s'>", r.From.StringKeyPartList())
	}
	return []string{fmt.Sprintf("ALTER INDEX `%s`%s", indexName,
		optS(r.To.IndexOptions.Diff(r.From.IndexOptions).String(), " %s"))}
}

func (r ModifiedIndex) Diff() []string {
	return []string{fmt.Sprintf("~ %s\t-> %s", r.From.String(), r.To.String())}
}

func (r ModifiedIndex) Id() string {
	return keyPartId(r.To.KeyPartList)
}

type DroppedIndex struct {
	This *parser.IndexDefinition
	Sequential
	Dependent
	Prefixable
}

func (r DroppedIndex) Statements() []string {
	indexName := r.This.IndexName
	if indexName == "" {
		indexName = fmt.Sprintf("<unknown index name of '%s'>", r.This.StringKeyPartList())
	}
	return []string{fmt.Sprintf("DROP INDEX `%s`", indexName)}
}

func (r DroppedIndex) Diff() []string {
	return []string{fmt.Sprintf("- %s", r.This.String())}
}

func (r DroppedIndex) Id() string {
	return keyPartId(r.This.KeyPartList)
}

type RenamedIndex struct {
	From *parser.IndexDefinition
	To   *parser.IndexDefinition
	Sequential
	Dependent
	Prefixable
}

func (r RenamedIndex) Statements() []string {
	return []string{fmt.Sprintf("RENAME INDEX `%s` TO `%s`", r.From.IndexName, r.To.IndexName)}
}

func (r RenamedIndex) Diff() []string {
	return []string{fmt.Sprintf("~ %s\t-> %s", r.From.String(), r.To.String())}
}

func (r RenamedIndex) Id() string {
	return keyPartId(r.From.KeyPartList)
}

type RetainedIndex struct {
	This *parser.IndexDefinition
	Sequential
	Dependent
	Prefixable
}

func (r RetainedIndex) Statements() []string {
	return []string{}
}

func (r RetainedIndex) Diff() []string {
	return []string{fmt.Sprintf("  %s", r.This.String())}
}

func (r RetainedIndex) Id() string {
	return keyPartId(r.This.KeyPartList)
}

func getIndexOrder(from []*parser.IndexDefinition, to []*parser.IndexDefinition, columnOrder map[string]int) map[string]int {
	all := []*parser.IndexDefinition{}
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

func keyPartContains(keyParts []parser.KeyPart, columnName string) bool {
	return ContainsIf(keyParts, func(r parser.KeyPart) bool { return r.Column == columnName })
}

func keyPartId(keyParts []parser.KeyPart) string {
	return parser.JoinT(keyParts, "\000", "")
}

func keyPartReplace(keyParts []parser.KeyPart, fromName string, toName string) []parser.KeyPart {
	ret := []parser.KeyPart{}
	for _, k := range keyParts {
		if k.Column == fromName {
			k.Column = toName
			ret = append(ret, k)
		} else {
			ret = append(ret, k)
		}
	}
	return ret
}
