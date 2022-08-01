package lib

import (
	"fmt"
	"github.com/emirpasic/gods/sets/linkedhashset"
	"github.com/kota65535/alternator/parser"
	"reflect"
	"strings"
)

type DatabaseAlterations struct {
	Added       []*AddedDatabase
	Modified    []*ModifiedDatabase
	Dropped     []*DroppedDatabase
	Retained    []*RetainedDatabase
	alterations []Alteration
}

func NewDatabaseAlterations(from []*Schema, to []*Schema) *DatabaseAlterations {

	fromMap := map[string]*Schema{}
	fromSet := linkedhashset.New()
	for _, s := range from {
		fromMap[s.Database.DbName] = s
		fromSet.Add(s.Database.DbName)
	}
	toMap := map[string]*Schema{}
	toSet := linkedhashset.New()
	for _, s := range to {
		toMap[s.Database.DbName] = s
		toSet.Add(s.Database.DbName)
	}

	databaseOrder := getDatabaseOrder(from, to)

	var added []*AddedDatabase
	var dropped []*DroppedDatabase
	var modified []*ModifiedDatabase
	var retained []*RetainedDatabase

	for _, v := range difference(fromSet, toSet).Values() {
		s := v.(string)
		t1 := fromMap[s].Tables
		tableAlterations := NewTableAlterations(t1, []*parser.CreateTableStatement{})
		dropped = append(dropped, &DroppedDatabase{
			This:       fromMap[s].Database,
			Tables:     &tableAlterations,
			Sequential: Sequential{databaseOrder[s]},
		})
	}

	for _, v := range difference(toSet, fromSet).Values() {
		s := v.(string)
		t2 := toMap[s].Tables
		tableAlterations := NewTableAlterations([]*parser.CreateTableStatement{}, t2)
		added = append(added, &AddedDatabase{
			This:       toMap[s].Database,
			Tables:     &tableAlterations,
			Sequential: Sequential{databaseOrder[s]},
		})
	}

	for _, v := range intersection(fromSet, toSet).Values() {
		s := v.(string)
		d1 := fromMap[s].Database
		d2 := toMap[s].Database
		t1 := fromMap[s].Tables
		t2 := toMap[s].Tables
		alteredTables := NewTableAlterations(t1, t2)
		if databasesEqual(d1, d2) {
			retained = append(retained, &RetainedDatabase{
				This:       d2,
				Tables:     alteredTables,
				Sequential: Sequential{databaseOrder[s]},
			})
		} else {
			modified = append(modified, &ModifiedDatabase{
				From: d1,
				To:   d2,
				DbOptions: &DatabaseOptionAlterations{
					From: d1.DatabaseOptions,
					To:   d2.DatabaseOptions,
				},
				Tables:     &alteredTables,
				Sequential: Sequential{databaseOrder[s]},
			})
		}
	}

	return &DatabaseAlterations{
		Added:    added,
		Modified: modified,
		Dropped:  dropped,
		Retained: retained,
	}
}

func (r *DatabaseAlterations) Statements() []string {
	ret := []string{}
	for _, a := range r.Alterations() {
		ret = append(ret, a.Statements()...)
	}
	return ret
}

func (r DatabaseAlterations) Diff() []string {
	ret := []string{}
	for _, a := range r.Alterations() {
		ret = append(ret, a.Diff()...)
	}
	return ret
}

func (r DatabaseAlterations) FromString() []string {
	ret := []string{}
	for _, a := range r.Alterations() {
		ret = append(ret, a.FromString()...)
	}
	return ret
}

func (r DatabaseAlterations) ToString() []string {
	ret := []string{}
	for _, a := range r.Alterations() {
		ret = append(ret, a.ToString()...)
	}
	return ret
}

func (r *DatabaseAlterations) Alterations() []Alteration {
	if len(r.alterations) != 0 {
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
	for _, a := range r.Retained {
		alterations = append(alterations, a)
	}
	r.alterations = NewDag(alterations).Sort()
	return r.alterations
}

type AddedDatabase struct {
	This   *parser.CreateDatabaseStatement
	Tables *TableAlterations
	Sequential
	Dependent
	Prefixable
}

func (r AddedDatabase) Statements() []string {
	ret := []string{}
	ret = append(ret, r.This.String())
	ret = append(ret, r.Tables.Statements()...)
	return ret
}

func (r AddedDatabase) Diff() []string {
	ret := []string{}
	// Append "+" to CREATE DATABASE statement
	ret = append(ret, prefix(r.This.String(), "+ "))
	ret = append(ret, r.Tables.Diff()...)
	return ret
}

func (r AddedDatabase) FromString() []string {
	return []string{}
}

func (r AddedDatabase) ToString() []string {
	ret := []string{}
	ret = append(ret, r.This.String())
	ret = append(ret, r.Tables.ToString()...)
	return ret
}

func (r AddedDatabase) Id() string {
	return r.This.DbName
}

type ModifiedDatabase struct {
	From      *parser.CreateDatabaseStatement
	To        *parser.CreateDatabaseStatement
	DbOptions *DatabaseOptionAlterations
	Tables    *TableAlterations
	Sequential
	Dependent
	Prefixable
}

func (r ModifiedDatabase) Statements() []string {
	var ret []string
	for _, s := range r.DbOptions.Statements() {
		ret = append(ret, fmt.Sprintf("ALTER DATABASE `%s` %s;", r.From.DbName, s))
	}
	ret = append(ret, r.Tables.Statements()...)
	return ret
}

func (r ModifiedDatabase) Diff() []string {
	ret := []string{}
	optionDiff := r.DbOptions.Diff()
	databaseOptions := strings.Join(indentDiff(optionDiff, 4), "\n")
	ret = append(ret, fmt.Sprintf("  CREATE DATABASE `%s`%s;",
		r.To.DbName,
		optS(databaseOptions, "\n%s")))
	ret = append(ret, r.Tables.Diff()...)
	return ret
}

func (r ModifiedDatabase) FromString() []string {
	ret := []string{}
	options := r.DbOptions.FromString()
	for i, s := range options {
		options[i] = fmt.Sprintf("%s%s", strings.Repeat(" ", 4), s)
	}
	databaseOptions := strings.Join(options, "\n")
	ret = append(ret, fmt.Sprintf("CREATE DATABASE `%s`%s;",
		r.To.DbName,
		optS(databaseOptions, "\n%s")))
	ret = append(ret, r.Tables.FromString()...)
	return ret
}

func (r ModifiedDatabase) ToString() []string {
	ret := []string{}
	options := r.DbOptions.ToString()
	for i, s := range options {
		options[i] = fmt.Sprintf("%s%s", strings.Repeat(" ", 4), s)
	}
	databaseOptions := strings.Join(options, "\n")
	ret = append(ret, fmt.Sprintf("CREATE DATABASE `%s`%s;",
		r.To.DbName,
		optS(databaseOptions, "\n%s")))
	ret = append(ret, r.Tables.ToString()...)
	return ret
}

func (r ModifiedDatabase) Id() string {
	return r.To.DbName
}

type DroppedDatabase struct {
	This   *parser.CreateDatabaseStatement
	Tables *TableAlterations
	Sequential
	Dependent
	Prefixable
}

func (r DroppedDatabase) Statements() []string {
	return []string{fmt.Sprintf("DROP DATABASE `%s`;", r.This.DbName)}
}

func (r DroppedDatabase) Diff() []string {
	ret := []string{}
	// Prepend "-" for CREATE DATABASE statement
	ret = append(ret, prefix(r.This.String(), "- "))
	ret = append(ret, r.Tables.Diff()...)
	return ret
}

func (r DroppedDatabase) FromString() []string {
	ret := []string{}
	ret = append(ret, r.This.String())
	ret = append(ret, r.Tables.FromString()...)
	return ret
}

func (r DroppedDatabase) ToString() []string {
	return []string{}
}

func (r DroppedDatabase) Id() string {
	return r.This.DbName
}

type RetainedDatabase struct {
	This   *parser.CreateDatabaseStatement
	Tables TableAlterations
	Sequential
	Dependent
	Prefixable
}

func (r RetainedDatabase) Statements() []string {
	return r.Tables.Statements()
}

func (r RetainedDatabase) Diff() []string {
	ret := []string{}
	ret = append(ret, prefix(r.This.String(), "  "))
	ret = append(ret, r.Tables.Diff()...)
	return ret
}

func (r RetainedDatabase) FromString() []string {
	ret := []string{}
	ret = append(ret, r.This.String())
	ret = append(ret, r.Tables.FromString()...)
	return ret
}

func (r RetainedDatabase) ToString() []string {
	ret := []string{}
	ret = append(ret, r.This.String())
	ret = append(ret, r.Tables.ToString()...)
	return ret
}

func (r RetainedDatabase) Id() string {
	return r.This.DbName
}

func getDatabaseOrder(from []*Schema, to []*Schema) map[string]int {
	ret := map[string]int{}
	p1 := 0
	p2 := 0
	seq := 0
	for p1 < len(from) || p2 < len(to) {
		if p1 >= len(from) {
			ret[to[p2].Database.DbName] = seq
			p2 += 1
			seq += 1
			continue
		}
		if p2 >= len(to) {
			if _, ok := ret[from[p1].Database.DbName]; !ok {
				ret[from[p1].Database.DbName] = seq
			}
			p1 += 1
			seq += 1
			continue
		}
		ret[to[p2].Database.DbName] = seq
		if _, ok := ret[from[p1].Database.DbName]; !ok {
			ret[from[p1].Database.DbName] = seq + 1
		}
		p1 += 1
		p2 += 1
		seq += 2
	}
	return ret
}

func databasesEqual(d1 *parser.CreateDatabaseStatement, d2 *parser.CreateDatabaseStatement) bool {
	return d1.DbName == d2.DbName && reflect.DeepEqual(d1.DatabaseOptions, d2.DatabaseOptions)
}

func indentDiff(strs []string, depth int) []string {
	ret := []string{}
	for _, s := range strs {
		ret = append(ret, fmt.Sprintf("%c%s%s", s[0], strings.Repeat(" ", depth+1), s[2:]))
	}
	return ret
}
