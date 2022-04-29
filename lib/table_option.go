package lib

import (
	"fmt"
	"github.com/kota65535/alternator/parser"
)

type TableOptionAlterations struct {
	From *parser.TableOptions
	To   *parser.TableOptions
	Sequential
	Dependent
	Prefixable
}

func NewTableOptionAlterations(from *parser.TableOptions, to *parser.TableOptions) TableOptionAlterations {
	// Unset AUTO_INCREMENT if it is smaller than current value
	if to.AutoIncrement < from.AutoIncrement {
		to.AutoIncrement = ""
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
	from := r.From.Map()
	to := r.To.Map()
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
	from := r.From.Map()
	to := r.To.Map()
	keys := to.Keys()
	ret := []string{}
	for _, k := range keys {
		cur, curOk := to.Get(k)
		old, oldOk := from.Get(k)
		if curOk {
			if !oldOk {
				ret = append(ret, fmt.Sprintf("+ %s = %s\t", k, cur))
			} else if cur != old {
				ret = append(ret, fmt.Sprintf("~ %s = %s\t-> %s = %s", k, old, k, cur))
			}
		}
	}
	return ret
}

func (r TableOptionAlterations) Equivalent() bool {
	return len(r.Diff()) == 0
}

func (r TableOptionAlterations) Id() string {
	return "TableOptionAlterations"
}
