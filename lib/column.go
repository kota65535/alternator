package lib

import (
	"fmt"
	"github.com/kota65535/alternator/parser"
	"reflect"
)

type ColumnAlterations struct {
	Added       []*AddedColumn
	Modified    []*ModifiedColumn
	Dropped     []*DroppedColumn
	Renamed     []*RenamedColumn
	Moved       []*MovedColumn
	Retained    []*RetainedColumn
	ColumnOrder map[string]int
	alterations []Alteration
}

func NewColumnAlterations(from []*parser.ColumnDefinition, to []*parser.ColumnDefinition) ColumnAlterations {

	fromColumnNames := map[string][]*parser.ColumnDefinition{}
	toColumnNames := map[string][]*parser.ColumnDefinition{}

	for _, t := range from {
		fromColumnNames[t.ColumnName] = append(fromColumnNames[t.ColumnName], t)
	}
	for _, t := range to {
		toColumnNames[t.ColumnName] = append(fromColumnNames[t.ColumnName], t)
	}

	var addedOrMoved []*AddedColumn
	var droppedOrMoved []*DroppedColumn
	var modified []*ModifiedColumn
	var moved []*MovedColumn
	var renamed []*RenamedColumn
	var retained []*RetainedColumn

	p1 := 0
	p2 := 0
	seq := 0
	for p1 < len(from) || p2 < len(to) {
		// The rest of c2 columns should be added
		if p1 == len(from) {
			t2 := to[p2]
			var after *parser.ColumnDefinition
			if p2 > 0 {
				after = to[p2-1]
			}
			addedOrMoved = append(addedOrMoved, &AddedColumn{
				This:       t2,
				After:      after,
				Sequential: Sequential{seq},
			})
			p2 += 1
			seq += 1
			continue
		}
		// The rest of c1 columns should be dropped
		if p2 == len(to) {
			t1 := from[p1]
			droppedOrMoved = append(droppedOrMoved, &DroppedColumn{
				This:       t1,
				Sequential: Sequential{seq},
			})
			p1 += 1
			seq += 1
			continue
		}

		c1 := from[p1]
		c2 := to[p2]

		// Skip if equals
		if reflect.DeepEqual(c1, c2) {
			retained = append(retained, &RetainedColumn{
				From:       c1,
				To:         c2,
				Sequential: Sequential{seq},
			})
			p1 += 1
			p2 += 1
			seq += 1
			continue
		}

		_, fromHasC2 := fromColumnNames[c2.ColumnName]
		_, toHasC1 := toColumnNames[c1.ColumnName]

		// Renamed if column definitions are the same and both table do not have the opponent's column name
		if !fromHasC2 && !toHasC1 && columnDefsEqual(*c1, *c2) {
			renamed = append(renamed, &RenamedColumn{
				From:       c1,
				To:         c2,
				Sequential: Sequential{seq},
			})
			p1 += 1
			p2 += 1
			seq += 1
			continue
		}

		// add column if it exists only in 'to'
		if !fromHasC2 {
			var after *parser.ColumnDefinition
			if p2 > 0 {
				after = to[p2-1]
			}
			addedOrMoved = append(addedOrMoved, &AddedColumn{
				This:       c2,
				After:      after,
				Sequential: Sequential{seq},
			})
			p2 += 1
			seq += 1
			continue
		}
		// remove column if it exists only in 'from'
		if !toHasC1 {
			droppedOrMoved = append(droppedOrMoved, &DroppedColumn{
				This:       c1,
				Sequential: Sequential{seq},
			})
			p1 += 1
			seq += 1
			continue
		}

		if c1.ColumnName == c2.ColumnName {
			// column is modified
			modified = append(modified, &ModifiedColumn{
				From:       c1,
				To:         c2,
				Sequential: Sequential{seq},
			})
			p1 += 1
			p2 += 1
			seq += 1
		} else {
			droppedOrMoved = append(droppedOrMoved, &DroppedColumn{
				This:       c1,
				Sequential: Sequential{seq},
			})
			p1 += 1
			seq += 1
		}
	}

	movedColumnNames := map[string]bool{}
	for _, a := range addedOrMoved {
		for _, d := range droppedOrMoved {
			if a.This.ColumnName == d.This.ColumnName {
				moved = append(moved, &MovedColumn{
					From:       d.This,
					To:         a.This,
					After:      a.After,
					Sequential: Sequential{a.SeqNum()},
				})
				movedColumnNames[a.This.ColumnName] = true
			}
		}
	}

	var added []*AddedColumn
	var dropped []*DroppedColumn
	for _, a := range addedOrMoved {
		if movedColumnNames[a.This.ColumnName] {
			continue
		}
		added = append(added, a)
	}
	for _, d := range droppedOrMoved {
		if movedColumnNames[d.This.ColumnName] {
			continue
		}
		dropped = append(dropped, d)
	}

	columnOrder := getColumnOrder(from, to)
	return ColumnAlterations{
		Added:       added,
		Modified:    modified,
		Renamed:     renamed,
		Moved:       moved,
		Dropped:     dropped,
		Retained:    retained,
		ColumnOrder: columnOrder,
	}
}

func (r ColumnAlterations) Statements() []string {
	ret := []string{}
	for _, b := range r.Alterations() {
		ret = append(ret, b.Statements()...)
	}
	return parser.Align(ret)
}

func (r ColumnAlterations) Diff() []string {
	ret := []string{}
	for _, b := range r.Alterations() {
		ret = append(ret, b.Diff()...)
	}
	ret = parser.Align(ret)
	return ret
}

func (r *ColumnAlterations) Alterations() []Alteration {
	if len(r.alterations) > 0 {
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
	for _, a := range r.Moved {
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

func (r *ColumnAlterations) Equivalent() bool {
	return len(r.Dropped) == 0 && len(r.Renamed) == 0 && len(r.Modified) == 0 && len(r.Added) == 0 && len(r.Moved) == 0
}

type AddedColumn struct {
	This  *parser.ColumnDefinition
	After *parser.ColumnDefinition
	Sequential
	Dependent
	Prefixable
}

func (r AddedColumn) Statements() []string {
	columnPos := "FIRST"
	if r.After != nil {
		columnPos = fmt.Sprintf("AFTER `%s`", r.After.ColumnName)
	}
	return []string{fmt.Sprintf("ADD COLUMN\t%s", r.This.StringWithPos(columnPos))}
}

func (r AddedColumn) Diff() []string {
	return []string{fmt.Sprintf("+ %s", r.This.String())}
}

func (r AddedColumn) Id() string {
	return r.This.ColumnName
}

func (r AddedColumn) IsValid() bool {
	return true
}

type ModifiedColumn struct {
	From *parser.ColumnDefinition
	To   *parser.ColumnDefinition
	Sequential
	Dependent
	Prefixable
}

func (r ModifiedColumn) Statements() []string {
	return []string{fmt.Sprintf("MODIFY COLUMN\t%s", r.To.String())}
}

func (r ModifiedColumn) Diff() []string {
	return []string{fmt.Sprintf("~ %s\t-> %s", r.From.String(), r.To.String())}
}

func (r ModifiedColumn) Id() string {
	return r.To.ColumnName
}

func (r ModifiedColumn) IsValid() bool {
	return !columnDefsEqual(*r.From, *r.To)
}

type MovedColumn struct {
	From  *parser.ColumnDefinition
	To    *parser.ColumnDefinition
	After *parser.ColumnDefinition
	Sequential
	Dependent
	Prefixable
}

func (r MovedColumn) Statements() []string {
	columnPos := "FIRST"
	if r.After.ColumnName != "" {
		columnPos = fmt.Sprintf("AFTER `%s`", r.After.ColumnName)
	}
	return []string{fmt.Sprintf("MODIFY COLUMN\t%s", r.To.StringWithPos(columnPos))}
}

func (r MovedColumn) Diff() []string {
	return []string{fmt.Sprintf("@ %s", r.To.String())}
}

func (r MovedColumn) Id() string {
	return r.To.ColumnName
}

func (r MovedColumn) IsValid() bool {
	return true
}

type DroppedColumn struct {
	This *parser.ColumnDefinition
	Sequential
	Dependent
	Prefixable
}

func (r DroppedColumn) Statements() []string {
	return []string{fmt.Sprintf("DROP COLUMN\t`%s`", r.This.ColumnName)}
}

func (r DroppedColumn) Diff() []string {
	return []string{fmt.Sprintf("- %s", r.This.String())}
}

func (r DroppedColumn) Id() string {
	return r.This.ColumnName
}

func (r DroppedColumn) IsValid() bool {
	return true
}

type RenamedColumn struct {
	From *parser.ColumnDefinition
	To   *parser.ColumnDefinition
	Sequential
	Dependent
	Prefixable
}

func (r RenamedColumn) Statements() []string {
	return []string{fmt.Sprintf("CHANGE COLUMN\t`%s` %s", r.From.ColumnName, r.To.String())}
}

func (r RenamedColumn) Diff() []string {
	return []string{fmt.Sprintf("~ %s\t-> %s", r.From.String(), r.To.String())}
}

func (r RenamedColumn) Id() string {
	return r.To.ColumnName
}

func (r RenamedColumn) IsValid() bool {
	return r.From.ColumnName != r.To.ColumnName
}

type RetainedColumn struct {
	From *parser.ColumnDefinition
	To   *parser.ColumnDefinition
	Sequential
	Dependent
	Prefixable
}

func (r RetainedColumn) Statements() []string {
	return []string{}
}

func (r RetainedColumn) Diff() []string {
	return []string{fmt.Sprintf("  %s", r.To.String())}
}

func (r RetainedColumn) Id() string {
	return r.To.ColumnName
}

func (r RetainedColumn) IsValid() bool {
	return true
}

func getColumnOrder(from []*parser.ColumnDefinition, to []*parser.ColumnDefinition) map[string]int {
	ret := map[string]int{}
	p1 := 0
	p2 := 0
	seq := 0
	for p1 < len(from) || p2 < len(to) {
		if p1 >= len(from) {
			ret[to[p2].ColumnName] = seq
			p2 += 1
			seq += 1
			continue
		}
		if p2 >= len(to) {
			ret[from[p1].ColumnName] = seq
			p1 += 1
			seq += 1
			continue
		}
		ret[from[p2].ColumnName] = seq
		ret[to[p1].ColumnName] = seq + 1
		p1 += 1
		p2 += 1
		seq += 2
	}
	return ret
}

// Compare 2 column definitions, ignoring column name.
func columnDefsEqual(c1 parser.ColumnDefinition, c2 parser.ColumnDefinition) bool {
	c1.ColumnName = ""
	c2.ColumnName = ""
	return reflect.DeepEqual(c1, c2)
}
