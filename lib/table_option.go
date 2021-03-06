package lib

import (
	"fmt"
	"github.com/kota65535/alternator/parser"
	"strconv"
)

type TableOptionAlterations struct {
	From *parser.TableOptions
	To   *parser.TableOptions
	Sequential
	Dependent
	Prefixable
}

func NewTableOptionAlterations(from *parser.TableOptions, to *parser.TableOptions) TableOptionAlterations {
	// Do not show AUTO_INCREMENT of the remote schema if local schema does not have it
	if to.AutoIncrement == "" {
		from.AutoIncrement = ""
	}
	// Unset AUTO_INCREMENT if the local one is smaller than the remote one
	toAi, err1 := strconv.Atoi(to.AutoIncrement)
	fromAi, err2 := strconv.Atoi(from.AutoIncrement)
	if err1 == nil && err2 == nil && toAi < fromAi {
		to.AutoIncrement = ""
		from.AutoIncrement = ""
	}

	return TableOptionAlterations{
		From: from,
		To:   to,
	}
}

func (r *TableOptionAlterations) Alterations() []Alteration {
	return []Alteration{r}
}

func (r TableOptionAlterations) Statements() []string {
	from := r.From.MapWithDefault()
	to := r.To.MapWithDefault()

	keys := to.Keys()
	ret := []string{}
	for _, k := range keys {
		cur, curOk := to.Get(k)
		old, _ := from.Get(k)
		if curOk && cur != old {
			ret = append(ret, fmt.Sprintf("%s = %s", k, cur))
		}
	}
	return ret
}

func (r TableOptionAlterations) Diff() []string {
	from := r.From.MapWithDefault()
	to := r.To.MapWithDefault()

	keys := to.Keys()
	ret := []string{}
	for _, k := range keys {
		cur, curOk := to.Get(k)
		old, oldOk := from.Get(k)
		if curOk {
			if !oldOk {
				ret = append(ret, fmt.Sprintf("+ %s = %s", k, cur))
			} else if cur != old {
				ret = append(ret, fmt.Sprintf("~ %s = %s\t-> %s = %s", k, old, k, cur))
			} else if _, ok := r.To.Map().Get(k); ok {
				ret = append(ret, fmt.Sprintf("  %s = %s", k, cur))
			}
		}
	}
	ret = parser.Align(ret)
	return ret
}

func (r TableOptionAlterations) FromString() []string {
	return r.From.Strings()
}

func (r TableOptionAlterations) ToString() []string {
	return r.To.Strings()
}

func (r TableOptionAlterations) Equivalent() bool {
	return len(r.Diff()) == 0
}

func (r TableOptionAlterations) Id() string {
	return "TableOptionAlterations"
}
