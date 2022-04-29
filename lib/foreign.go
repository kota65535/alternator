package lib

import (
	"fmt"
	"github.com/emirpasic/gods/sets/linkedhashset"
	"github.com/kota65535/alternator/parser"
	"reflect"
	"sort"
	"strings"
)

type ForeignKeyAlterations struct {
	Added       []*AddedForeignKey
	Dropped     []*DroppedForeignKey
	Retained    []*RetainedForeignKey
	alterations []Alteration
}

func NewForeignKeyAlterations(
	from []*parser.ForeignKeyDefinition,
	to []*parser.ForeignKeyDefinition,
	columnOrder map[string]int) ForeignKeyAlterations {

	fromMap := map[string]*parser.ForeignKeyDefinition{}
	fromSet := linkedhashset.New()
	for _, t := range from {
		keys := t.StringKeyPartList()
		fromMap[keys] = t
		fromSet.Add(keys)
	}
	toMap := map[string]*parser.ForeignKeyDefinition{}
	toSet := linkedhashset.New()
	for _, t := range to {
		keys := t.StringKeyPartList()
		toMap[keys] = t
		toSet.Add(keys)
	}

	foreignKeyOrder := getForeignKeyOrder(from, to, columnOrder)

	var added []*AddedForeignKey
	var dropped []*DroppedForeignKey
	var retained []*RetainedForeignKey

	for _, v := range difference(fromSet, toSet).Values() {
		s := v.(string)
		dropped = append(dropped, &DroppedForeignKey{
			This:       fromMap[s],
			Sequential: Sequential{foreignKeyOrder[s]},
		})
	}
	for _, v := range difference(toSet, fromSet).Values() {
		s := v.(string)
		added = append(added, &AddedForeignKey{
			This:       toMap[s],
			Sequential: Sequential{foreignKeyOrder[s]},
		})
	}
	for _, v := range intersection(fromSet, toSet).Values() {
		s := v.(string)
		t1 := fromMap[s]
		t2 := toMap[s]
		if foreignKeyDefsEqual(*t1, *t2) {
			retained = append(retained, &RetainedForeignKey{
				This:       t1,
				Sequential: Sequential{foreignKeyOrder[s]},
			})
		} else {
			dropped = append(dropped, &DroppedForeignKey{
				This:       fromMap[s],
				Sequential: Sequential{foreignKeyOrder[s]},
			})
			added = append(added, &AddedForeignKey{
				This:       toMap[s],
				Sequential: Sequential{foreignKeyOrder[s]},
			})
		}
	}

	return ForeignKeyAlterations{
		Added:    added,
		Dropped:  dropped,
		Retained: retained,
	}
}

func (r ForeignKeyAlterations) Statements() []string {
	ret := []string{}
	for _, b := range r.Alterations() {
		ret = append(ret, b.Statements()...)
	}
	ret = parser.Align(ret)
	return ret
}

func (r ForeignKeyAlterations) Diff() []string {
	ret := []string{}
	for _, b := range r.Alterations() {
		ret = append(ret, b.Diff()...)
	}
	ret = parser.Align(ret)
	return ret
}

func (r *ForeignKeyAlterations) Alterations() []Alteration {
	if r.alterations != nil {
		return r.alterations
	}
	alterations := []Alteration{}
	for _, a := range r.Dropped {
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

func (r *ForeignKeyAlterations) Equivalent() bool {
	return len(r.Dropped) == 0 && len(r.Added) == 0
}

// HandleTableDrop ensures that foreign key drop is preceded to the table drop.
func (r *ForeignKeyAlterations) HandleTableDrop(alt Alteration, tableName string) {
	for _, d := range r.Dropped {
		if d.This.ReferenceDefinition.TableName == tableName {
			alt.AddDependsOn(d)
		}
	}
}

// HandleTableRename fixes falsy detection of foreign key modification caused by table renaming.
func (r *ForeignKeyAlterations) HandleTableRename(alt Alteration, fromName string, toName string) {
	removed := map[Alteration]bool{}
	// update table name of dropped foreign keys
	for _, d := range r.Dropped {
		if d.This.ReferenceDefinition.TableName == fromName {
			d.This.ReferenceDefinition.TableName = toName
		}
	}
	// update table name of dropped foreign keys
	for _, a := range r.Added {
		if a.This.ReferenceDefinition.TableName == toName {
			a.AddDependsOn(alt)
		}
	}
	// If dropped/added foreign key defs are equal, treat them as retained
	for _, d := range r.Dropped {
		for _, a := range r.Added {
			if d.This.ReferenceDefinition.TableName == a.This.ReferenceDefinition.TableName {
				if foreignKeyDefsEqual(*d.This, *a.This) {
					r.Retained = append(r.Retained, &RetainedForeignKey{
						This:       a.This,
						Sequential: a.Sequential,
						Dependent:  a.Dependent,
					})
					removed[d] = true
					removed[a] = true
				} else {
					// still needs to be modified, and must be altered after table rename
					d.AddDependsOn(alt)
					a.AddDependsOn(alt)
				}
			}
		}
	}
	r.Dropped = RemoveIf(r.Dropped, func(m *DroppedForeignKey) bool { return removed[m] })
	r.Added = RemoveIf(r.Added, func(m *AddedForeignKey) bool { return removed[m] })
}

// HandleColumnModify ensures the foreign keys are recreated when their key columns are modified.
func (r *ForeignKeyAlterations) HandleColumnModify(modify Alteration, columnName string) {
	for _, d := range r.Dropped {
		if Contains(d.This.KeyPartList, columnName) {
			modify.AddDependsOn(d)
		}
	}

	for _, a := range r.Added {
		if Contains(a.This.KeyPartList, columnName) {
			a.AddDependsOn(modify)
		}
	}
}

// HandleColumnDrop ensures that foreign key deletion is run before column deletion
func (r *ForeignKeyAlterations) HandleColumnDrop(drop Alteration, columnName string) {
	for _, f := range r.Dropped {
		if Contains(f.This.KeyPartList, columnName) {
			drop.AddDependsOn(f)
		}
	}
}

func (r *ForeignKeyAlterations) HandleRefColumnDrop(drop Alteration, tableName string, columnName string) {
	for _, f := range r.Dropped {
		if f.This.ReferenceDefinition.TableName == tableName && Contains(f.This.ReferenceDefinition.KeyPartList, columnName) {
			drop.AddDependsOn(f)
		}
	}
}

// HandleRefColumnRename handles key part change caused by column rename.
func (r *ForeignKeyAlterations) HandleRefColumnRename(tableName string, fromName string, toName string) {
	removed := map[Alteration]bool{}
	// Changing key part is first considered as drop & add a foreign key, so we will find the pair
	for _, d := range r.Dropped {
		for _, a := range r.Added {
			if d.This.ReferenceDefinition.TableName == tableName &&
				Contains(d.This.ReferenceDefinition.KeyPartList, fromName) &&
				Contains(a.This.ReferenceDefinition.KeyPartList, toName) {
				// update key parts
				d.This.ReferenceDefinition.KeyPartList = Replace(d.This.ReferenceDefinition.KeyPartList, fromName, toName)
				if foreignKeyDefsEqual(*d.This, *a.This) {
					r.Retained = append(r.Retained, &RetainedForeignKey{
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
	r.Dropped = RemoveIf(r.Dropped, func(m *DroppedForeignKey) bool { return removed[m] })
	r.Added = RemoveIf(r.Added, func(m *AddedForeignKey) bool { return removed[m] })
}

// HandleRefColumnModify handles referencing column modification that requires recreating the foreign key.
func (r *ForeignKeyAlterations) HandleRefColumnModify(modify Alteration, tableName string, columnName string) {
	removed := map[Alteration]bool{}
	// Retained foreign keys must be dropped before the column modification and added again after that
	for _, f := range r.Retained {
		if f.This.ReferenceDefinition.TableName == tableName && Contains(f.This.ReferenceDefinition.KeyPartList, columnName) {
			// Drop the FK before column modification
			droppedFK := &DroppedForeignKey{
				This:       f.This,
				Sequential: f.Sequential,
				Dependent:  f.Dependent,
			}
			modify.AddDependsOn(droppedFK)
			r.Dropped = append(r.Dropped, droppedFK)

			// Add the FK after column modification
			addedFK := &AddedForeignKey{
				This:       f.This,
				Sequential: f.Sequential,
				Dependent:  f.Dependent,
			}
			addedFK.AddDependsOn(modify)
			r.Added = append(r.Added, addedFK)

			// Remove the original retained FK
			removed[f] = true
		}
	}

	// Dropped foreign key must be dropped before the column modification
	for _, f := range r.Dropped {
		if f.This.ReferenceDefinition.TableName == tableName && Contains(f.This.ReferenceDefinition.KeyPartList, columnName) {
			modify.AddDependsOn(f)
		}
	}
	r.Retained = RemoveIf(r.Retained, func(m *RetainedForeignKey) bool { return removed[m] })
}

type AddedForeignKey struct {
	This *parser.ForeignKeyDefinition
	Sequential
	Dependent
	Prefixable
}

func (r AddedForeignKey) Statements() []string {
	return []string{fmt.Sprintf("ADD %s", r.This.String())}
}

func (r AddedForeignKey) Diff() []string {
	return []string{fmt.Sprintf("+ %s", r.This.String())}
}

func (r AddedForeignKey) Id() string {
	return strings.Join(r.This.KeyPartList, "\000")
}

type DroppedForeignKey struct {
	This *parser.ForeignKeyDefinition
	Sequential
	Dependent
	Prefixable
}

func (r DroppedForeignKey) Statements() []string {
	constraintName := r.This.ConstraintName
	if constraintName == "" {
		constraintName = fmt.Sprintf("<unknown constraint name of '%s'>", r.This.StringKeyPartList())
	}
	indexName := r.This.IndexName
	if indexName == "" {
		indexName = fmt.Sprintf("<unknown index name of '%s'>", r.This.StringKeyPartList())
	}
	return []string{
		fmt.Sprintf("DROP FOREIGN KEY `%s`", constraintName),
		fmt.Sprintf("DROP INDEX `%s`", indexName),
	}
}

func (r DroppedForeignKey) Diff() []string {
	return []string{fmt.Sprintf("- %s", r.This.String())}
}

func (r DroppedForeignKey) Id() string {
	return strings.Join(r.This.KeyPartList, "\000")
}

type RetainedForeignKey struct {
	This *parser.ForeignKeyDefinition
	Sequential
	Dependent
	Prefixable
}

func (r RetainedForeignKey) Statements() []string {
	return []string{}
}

func (r RetainedForeignKey) Diff() []string {
	return []string{fmt.Sprintf("  %s", r.This.String())}
}

func (r RetainedForeignKey) Id() string {
	return strings.Join(r.This.KeyPartList, "\000")
}

func getForeignKeyOrder(from []*parser.ForeignKeyDefinition, to []*parser.ForeignKeyDefinition, columnOrder map[string]int) map[string]int {
	all := []*parser.ForeignKeyDefinition{}
	all = append(all, from...)
	all = append(all, to...)
	sort.SliceStable(all, func(i, j int) bool {
		if len(all[i].KeyPartList) != len(all[j].KeyPartList) {
			return len(all[i].KeyPartList) < len(all[j].KeyPartList)
		}
		length := len(all[i].KeyPartList)
		for a := 0; a < length; a++ {
			if all[i].KeyPartList[a] != all[j].KeyPartList[a] {
				return columnOrder[all[i].KeyPartList[a]] < columnOrder[all[j].KeyPartList[a]]
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
func foreignKeyDefsEqual(from parser.ForeignKeyDefinition, to parser.ForeignKeyDefinition) bool {
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
