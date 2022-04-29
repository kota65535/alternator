package lib

import (
	"fmt"
	"github.com/emirpasic/gods/sets/linkedhashset"
	"github.com/kota65535/alternator/parser"
	"strings"
)

type TableAlterations struct {
	Added              []*AddedTable
	Modified           []*ModifiedTable
	Dropped            []*DroppedTable
	Renamed            []*RenamedTable
	Retained           []*RetainedTable
	alterations        []Alteration
	elementAlterations []Alteration
}

func NewTableAlterations(from []*parser.CreateTableStatement, to []*parser.CreateTableStatement) TableAlterations {

	fromMap := map[string]*parser.CreateTableStatement{}
	fromSet := linkedhashset.New()
	for _, t := range from {
		fromMap[t.TableName] = t
		fromSet.Add(t.TableName)
	}
	toMap := map[string]*parser.CreateTableStatement{}
	toSet := linkedhashset.New()
	for _, t := range to {
		toMap[t.TableName] = t
		toSet.Add(t.TableName)
	}

	tableOrder := getTableOrder(from, to)
	dependencies := getTableDependencies(from)

	var addedOrRenamed []*AddedTable
	var droppedOrRenamed []*DroppedTable
	var renamed []*RenamedTable
	var modified []*ModifiedTable
	var retained []*RetainedTable

	// Tables not in 'to' DB have been deleted or renamed
	for _, v := range difference(fromSet, toSet).Values() {
		s := v.(string)
		droppedOrRenamed = append(droppedOrRenamed, &DroppedTable{
			This:       fromMap[s],
			Sequential: Sequential{tableOrder[s]},
		})
	}
	// Tables not in 'from' DB have added or renamed
	for _, v := range difference(toSet, fromSet).Values() {
		s := v.(string)
		addedOrRenamed = append(addedOrRenamed, &AddedTable{
			This:       toMap[s],
			Sequential: Sequential{tableOrder[s]},
		})
	}

	renamedFrom := map[string]bool{}
	renamedTo := map[string]bool{}
	// If two tables are equal except their names, they have been renamed
	for _, a := range addedOrRenamed {
		for _, d := range droppedOrRenamed {
			t1 := fromMap[d.Id()]
			t2 := toMap[a.Id()]

			elements := NewTableElementAlterations(t1, t2)

			if elements.Equivalent() {
				// create ForeignKeyAlterations instance in case of column modification foreign key is referencing
				renamed = append(renamed, &RenamedTable{
					From:        d.This,
					To:          a.This,
					ForeignKeys: elements.ForeignKeys,
					Sequential:  Sequential{tableOrder[a.Id()]},
				})
				renamedFrom[d.This.TableName] = true
				renamedTo[a.This.TableName] = true
			}
		}
	}

	// Remove renamed tables from the lists
	added := RemoveIf(addedOrRenamed, func(t *AddedTable) bool {
		return renamedTo[t.This.TableName]
	})
	dropped := RemoveIf(droppedOrRenamed, func(t *DroppedTable) bool {
		return renamedFrom[t.This.TableName]
	})

	for _, v := range intersection(fromSet, toSet).Values() {
		s := v.(string)
		t1 := fromMap[s]
		t2 := toMap[s]

		elements := NewTableElementAlterations(t1, t2)

		if elements.Equivalent() {
			// create ForeignKeyAlterations instance in case of column modification foreign key is referencing
			retained = append(retained, &RetainedTable{
				This:        t1,
				ForeignKeys: elements.ForeignKeys,
				Sequential:  Sequential{tableOrder[s]},
			})
		} else {
			modified = append(modified, &ModifiedTable{
				From:             t1,
				To:               t2,
				Columns:          elements.Columns,
				PrimaryKeys:      elements.PrimaryKeys,
				UniqueKeys:       elements.UniqueKeys,
				ForeignKeys:      elements.ForeignKeys,
				Indexes:          elements.Indexes,
				FullTextIndexes:  elements.FullTextIndexes,
				CheckConstraints: elements.CheckConstraints,
				TableOptions:     elements.TableOptions,
				Sequential:       Sequential{tableOrder[s]},
			})
		}
	}

	for _, mt := range modified {
		// handle table renaming referred by foreign keys
		for _, rt := range renamed {
			mt.ForeignKeys.HandleTableRename(rt, rt.From.TableName, rt.To.TableName)
		}
		// handle foreign key column drop
		for _, dc := range mt.Columns.Dropped {
			mt.ForeignKeys.HandleColumnDrop(dc, dc.This.ColumnName)
		}
		// handle FK & PK key column modification
		for _, mc := range mt.Columns.Modified {
			mt.PrimaryKeys.HandleColumnModify(mc, mc.From.ColumnName)
			mt.ForeignKeys.HandleColumnModify(mc, mc.From.ColumnName)
		}
	}

	for _, t1 := range modified {
		// handle column drop referred by foreign keys
		for _, c1 := range t1.Columns.Dropped {
			for _, t2 := range modified {
				t2.ForeignKeys.HandleRefColumnDrop(c1, t1.To.TableName, c1.This.ColumnName)
			}
		}
		for _, c1 := range t1.Columns.Renamed {
			for _, t2 := range modified {
				t2.ForeignKeys.HandleRefColumnRename(t1.To.TableName, c1.From.ColumnName, c1.To.ColumnName)
			}
		}
		// handle column modification referred by foreign keys
		for _, c1 := range t1.Columns.Modified {
			for _, t2 := range modified {
				t2.ForeignKeys.HandleRefColumnModify(c1, t1.To.TableName, c1.To.ColumnName)
			}
			for _, t2 := range retained {
				t2.ForeignKeys.HandleRefColumnModify(c1, t1.To.TableName, c1.To.ColumnName)
			}
			for _, t2 := range renamed {
				t2.ForeignKeys.HandleRefColumnModify(c1, t1.To.TableName, c1.To.ColumnName)
			}
		}
	}

	alterations := map[string]Alteration{}
	for _, t := range added {
		alterations[t.Id()] = t
	}
	for _, t := range modified {
		alterations[t.Id()] = t
	}
	for _, t := range renamed {
		alterations[t.Id()] = t
	}
	for _, t := range dropped {
		alterations[t.Id()] = t
	}
	for _, t := range retained {
		alterations[t.Id()] = t
	}
	for n, a := range alterations {
		for _, dep := range dependencies[n] {
			a.AddDependsOn(alterations[dep])
		}
	}

	return TableAlterations{
		Added:    added,
		Modified: modified,
		Renamed:  renamed,
		Dropped:  dropped,
		Retained: retained,
	}
}

type TableElementAlterations struct {
	Columns          *ColumnAlterations
	PrimaryKeys      *PrimaryKeyAlterations
	UniqueKeys       *UniqueKeyAlterations
	ForeignKeys      *ForeignKeyAlterations
	Indexes          *IndexAlterations
	FullTextIndexes  *FullTextIndexAlterations
	CheckConstraints *CheckConstraintsAlterations
	TableOptions     *TableOptionAlterations
}

func NewTableElementAlterations(t1 *parser.CreateTableStatement, t2 *parser.CreateTableStatement) TableElementAlterations {
	columns := NewColumnAlterations(t1.GetColumns(), t2.GetColumns())
	for _, c := range columns.Renamed {
		for _, p := range t1.GetPrimaryKeys() {
			p.KeyPartList = Replace(p.KeyPartList, c.From.ColumnName, c.To.ColumnName)
		}
		for _, p := range t1.GetUniqueKeys() {
			p.KeyPartList = Replace(p.KeyPartList, c.From.ColumnName, c.To.ColumnName)
		}
		for _, p := range t1.GetIndexes() {
			p.KeyPartList = Replace(p.KeyPartList, c.From.ColumnName, c.To.ColumnName)
		}
		for _, p := range t1.GetFullTextIndexes() {
			p.KeyPartList = Replace(p.KeyPartList, c.From.ColumnName, c.To.ColumnName)
		}
		for _, p := range t1.GetForeignKeys() {
			p.KeyPartList = Replace(p.KeyPartList, c.From.ColumnName, c.To.ColumnName)
		}
	}

	primaryKeys := NewPrimaryKeyAlterations(t1.GetPrimaryKeys(), t2.GetPrimaryKeys(), columns.ColumnOrder)
	uniqueKeys := NewUniqueAlterations(t1.GetUniqueKeys(), t2.GetUniqueKeys(), columns.ColumnOrder)
	indexes := NewIndexAlterations(t1.GetIndexes(), t2.GetIndexes(), columns.ColumnOrder)
	fullTextIndexes := NewFullTextIndexAlteration(t1.GetFullTextIndexes(), t2.GetFullTextIndexes(), columns.ColumnOrder)
	foreignKeys := NewForeignKeyAlterations(t1.GetForeignKeys(), t2.GetForeignKeys(), columns.ColumnOrder)
	checkConstraints := NewCheckConstraintAlterations(t1.GetCheckConstraints(), t2.GetCheckConstraints())
	tableOptions := NewTableOptionAlterations(&t1.TableOptions, &t2.TableOptions)

	return TableElementAlterations{
		Columns:          &columns,
		PrimaryKeys:      &primaryKeys,
		UniqueKeys:       &uniqueKeys,
		Indexes:          &indexes,
		FullTextIndexes:  &fullTextIndexes,
		ForeignKeys:      &foreignKeys,
		CheckConstraints: &checkConstraints,
		TableOptions:     &tableOptions,
	}
}

func (r TableElementAlterations) Equivalent() bool {
	return r.Columns.Equivalent() &&
		r.PrimaryKeys.Equivalent() &&
		r.UniqueKeys.Equivalent() &&
		r.Indexes.Equivalent() &&
		r.FullTextIndexes.Equivalent() &&
		r.ForeignKeys.Equivalent() &&
		r.CheckConstraints.Equivalent() &&
		r.TableOptions.Equivalent()
}

func (r TableAlterations) Statements() []string {
	ret := []string{}
	for _, a := range r.TableElementAlterations() {
		statements := []string{}
		for _, s := range a.Statements() {
			statements = append(statements, a.Prefix()+s)
		}
		statements = parser.Align(statements)
		for _, s := range statements {
			// TODO: more smartly
			if !strings.HasSuffix(s, ";") {
				s += ";"
			}
			ret = append(ret, s)
		}
	}
	return ret
}

func (r TableAlterations) Diff() []string {
	ret := []string{}
	for _, a := range r.Alterations() {
		ret = append(ret, a.Diff()...)
	}
	return ret
}

func (r *TableAlterations) TableElementAlterations() []Alteration {
	if r.elementAlterations != nil {
		return r.elementAlterations
	}
	alterations := []Alteration{}
	for _, a := range r.Added {
		for _, b := range a.Alterations() {
			b.SetSeqNum(a.SeqNum())
			alterations = append(alterations, b)
		}
	}
	for _, a := range r.Modified {
		for _, b := range a.Alterations() {
			b.SetSeqNum(a.SeqNum())
			alterations = append(alterations, b)
		}
	}
	for _, a := range r.Dropped {
		for _, b := range a.Alterations() {
			b.SetSeqNum(a.SeqNum())
			alterations = append(alterations, b)
		}
	}
	for _, a := range r.Renamed {
		for _, b := range a.Alterations() {
			b.SetSeqNum(a.SeqNum())
			alterations = append(alterations, b)
		}
	}
	for _, a := range r.Retained {
		for _, b := range a.Alterations() {
			b.SetSeqNum(a.SeqNum())
			alterations = append(alterations, b)
		}
	}

	r.elementAlterations = NewDag(alterations).Sort()
	return r.elementAlterations
}

func (r *TableAlterations) Alterations() []Alteration {
	if r.alterations != nil {
		return r.alterations
	}
	alterations := []Alteration{}
	for _, a := range r.Added {
		alterations = append(alterations, a)
	}
	for _, a := range r.Modified {
		alterations = append(alterations, a)
	}
	for _, a := range r.Dropped {
		alterations = append(alterations, a)
	}
	for _, a := range r.Renamed {
		alterations = append(alterations, a)
	}
	for _, a := range r.Retained {
		alterations = append(alterations, a)
	}

	r.alterations = NewDag(alterations).Sort()
	return r.alterations
}

type AddedTable struct {
	This *parser.CreateTableStatement
	Sequential
	Dependent
	Prefixable
}

func (r *AddedTable) Alterations() []Alteration {
	return []Alteration{r}
}

func (r AddedTable) Statements() []string {
	return []string{r.This.String()}
}

func (r AddedTable) Diff() []string {
	return []string{prefix(r.This.String(), "+ ")}
}

func (r AddedTable) Id() string {
	return r.This.TableName
}

type ModifiedTable struct {
	From             *parser.CreateTableStatement
	To               *parser.CreateTableStatement
	Columns          *ColumnAlterations
	PrimaryKeys      *PrimaryKeyAlterations
	UniqueKeys       *UniqueKeyAlterations
	Indexes          *IndexAlterations
	FullTextIndexes  *FullTextIndexAlterations
	ForeignKeys      *ForeignKeyAlterations
	CheckConstraints *CheckConstraintsAlterations
	TableOptions     *TableOptionAlterations
	Sequential
	Dependent
	Prefixable
}

func (r ModifiedTable) Alterations() []Alteration {
	alterations := []Alteration{}
	alterations = append(alterations, r.Columns.Alterations()...)
	alterations = append(alterations, r.PrimaryKeys.Alterations()...)
	alterations = append(alterations, r.UniqueKeys.Alterations()...)
	alterations = append(alterations, r.Indexes.Alterations()...)
	alterations = append(alterations, r.FullTextIndexes.Alterations()...)
	alterations = append(alterations, r.ForeignKeys.Alterations()...)
	alterations = append(alterations, r.CheckConstraints.Alterations()...)
	alterations = append(alterations, r.TableOptions.Alterations()...)

	for i, a := range alterations {
		a.SetPrefix(fmt.Sprintf("ALTER TABLE `%s`.`%s` ", r.To.DbName, r.To.TableName))
		a.SetSeqNum(i)
	}

	sorted := NewDag(alterations).Sort()
	return sorted
}

func (r ModifiedTable) Statements() []string {
	r.preprocessStatements()

	alterations := []Alteration{}
	alterations = append(alterations, r.Columns.Alterations()...)
	alterations = append(alterations, r.ForeignKeys.Alterations()...)

	for i, a := range alterations {
		a.SetSeqNum(i)
	}

	sorted := NewDag(alterations).Sort()

	statements := []string{}
	for _, a := range sorted {
		statements = append(statements, a.Statements()...)
	}

	return createStatements(r.From.DbName, r.From.TableName, statements)
}

func (r ModifiedTable) Diff() []string {
	defStrs := []string{}
	for _, s := range r.Columns.Diff() {
		defStrs = append(defStrs, fmt.Sprintf("%c%s%s", s[0], strings.Repeat(" ", 4+1), s[2:]))
	}
	for _, s := range r.PrimaryKeys.Diff() {
		defStrs = append(defStrs, fmt.Sprintf("%c%s%s", s[0], strings.Repeat(" ", 4+1), s[2:]))
	}
	for _, s := range r.UniqueKeys.Diff() {
		defStrs = append(defStrs, fmt.Sprintf("%c%s%s", s[0], strings.Repeat(" ", 4+1), s[2:]))
	}
	for _, s := range r.Indexes.Diff() {
		defStrs = append(defStrs, fmt.Sprintf("%c%s%s", s[0], strings.Repeat(" ", 4+1), s[2:]))
	}
	for _, s := range r.FullTextIndexes.Diff() {
		defStrs = append(defStrs, fmt.Sprintf("%c%s%s", s[0], strings.Repeat(" ", 4+1), s[2:]))
	}
	for _, s := range r.ForeignKeys.Diff() {
		defStrs = append(defStrs, fmt.Sprintf("%c%s%s", s[0], strings.Repeat(" ", 4+1), s[2:]))
	}
	for _, s := range r.CheckConstraints.Diff() {
		defStrs = append(defStrs, fmt.Sprintf("%c%s%s", s[0], strings.Repeat(" ", 4+1), s[2:]))
	}

	tsStrs := []string{}
	for _, s := range r.TableOptions.Diff() {
		tsStrs = append(tsStrs, fmt.Sprintf("%c%s%s", s[0], strings.Repeat(" ", 4+1), s[2:]))
	}
	tsStrs = parser.Align(tsStrs)
	tableOptions := strings.Join(tsStrs, "\n")

	return []string{
		fmt.Sprintf("  CREATE TABLE %s`%s`\n  (\n%s\n  )%s;",
			optS(r.To.DbName, "`%s`."),
			r.To.TableName,
			strings.Join(defStrs, ",\n"),
			optS(tableOptions, "\n%s")),
	}
}

func (r ModifiedTable) Id() string {
	return r.To.TableName
}

func (r *ModifiedTable) preprocessStatements() {
	for _, pk := range r.PrimaryKeys.Added {
		for i, _ := range r.Columns.Modified {
			c := r.Columns.Modified[i]
			if index(pk.This.KeyPartList, c.From.ColumnName) >= 0 {
				c.From.ColumnOptions.Nullability = "NOT NULL"
			}
		}
	}
}

type DroppedTable struct {
	This *parser.CreateTableStatement
	Sequential
	Dependent
	Prefixable
}

func (r *DroppedTable) Alterations() []Alteration {
	return []Alteration{r}
}

func (r DroppedTable) Statements() []string {
	return []string{fmt.Sprintf("DROP TABLE `%s`.`%s`", r.This.DbName, r.This.TableName)}
}

func (r DroppedTable) Diff() []string {
	return []string{prefix(r.This.String(), "- ")}
}

func (r DroppedTable) Id() string {
	return r.This.TableName
}

type RenamedTable struct {
	From        *parser.CreateTableStatement
	To          *parser.CreateTableStatement
	ForeignKeys *ForeignKeyAlterations
	Sequential
	Dependent
	Prefixable
}

func (r *RenamedTable) Alterations() []Alteration {
	alterations := []Alteration{}
	alterations = append(alterations, r)
	alterations = append(alterations, r.ForeignKeys.Alterations()...)
	return alterations
}

func (r RenamedTable) Statements() []string {
	ret := []string{}
	ret = append(ret, fmt.Sprintf("ALTER TABLE `%s`.`%s` RENAME TO `%s`.`%s`",
		r.From.DbName,
		r.From.TableName,
		r.From.DbName,
		r.To.TableName))
	ret = append(ret, r.ForeignKeys.Statements()...)
	return ret
}

func (r RenamedTable) Diff() []string {
	from := prefix(r.From.String(), "  ")
	to := prefix(r.To.String(), "  ")
	fromHead, _, _ := strings.Cut(from, "\n")
	toHead, _, _ := strings.Cut(to, "\n")
	ret := strings.Replace(from, fromHead, fmt.Sprintf("~ %s -> %s", fromHead[2:], toHead[2:]), 1)
	return []string{ret}
}

func (r RenamedTable) Id() string {
	return r.To.TableName
}

type RetainedTable struct {
	This        *parser.CreateTableStatement
	ForeignKeys *ForeignKeyAlterations
	Sequential
	Dependent
	Prefixable
}

func (r *RetainedTable) Alterations() []Alteration {
	alterations := []Alteration{}
	alterations = append(alterations, r)
	alterations = append(alterations, r.ForeignKeys.Alterations()...)
	return alterations
}

func (r RetainedTable) Statements() []string {
	return []string{}
}

func (r RetainedTable) Diff() []string {
	return []string{prefix(r.This.String(), "  ")}
}

func (r RetainedTable) Id() string {
	return r.This.TableName
}

func getTableOrder(from []*parser.CreateTableStatement, to []*parser.CreateTableStatement) map[string]int {
	ret := map[string]int{}
	p1 := 0
	p2 := 0
	seq := 0
	for p1 < len(from) || p2 < len(to) {
		if p1 >= len(from) {
			ret[to[p2].TableName] = seq
			p2 += 1
			seq += 1
			continue
		}
		if p2 >= len(to) {
			if _, ok := ret[from[p1].TableName]; !ok {
				ret[from[p1].TableName] = seq
			}
			p1 += 1
			seq += 1
			continue
		}
		ret[to[p2].TableName] = seq
		if _, ok := ret[from[p1].TableName]; !ok {
			ret[from[p1].TableName] = seq + 1
		}
		p1 += 1
		p2 += 1
		seq += 2
	}
	return ret
}

func getTableDependencies(statements []*parser.CreateTableStatement) map[string][]string {
	ret := map[string][]string{}
	for _, s := range statements {
		for _, d := range s.CreateDefinitions {
			if fk, ok := d.(*parser.ForeignKeyDefinition); ok {
				if !Contains(ret[s.TableName], fk.ReferenceDefinition.TableName) {
					ret[s.TableName] = append(ret[s.TableName], fk.ReferenceDefinition.TableName)
				}
			}
		}
	}
	return ret
}

func createStatements(db string, table string, statements []string) []string {
	ret := []string{}
	for _, s := range statements {
		ret = append(ret, fmt.Sprintf("ALTER TABLE `%s`.`%s` %s", db, table, s))
	}
	return ret
}
