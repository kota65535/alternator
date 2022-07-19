package lib

import (
	"fmt"
	"github.com/kota65535/alternator/parser"
)

type DatabaseOptionAlterations struct {
	From *parser.DatabaseOptions
	To   *parser.DatabaseOptions
	Sequential
	Dependent
	Prefixable
}

func (r *DatabaseOptionAlterations) Alterations() []Alteration {
	return []Alteration{r}
}

func (r DatabaseOptionAlterations) Statements() []string {
	from := r.From.Map()
	to := r.To.Map()
	from.Put("DEFAULT CHARACTER SET", r.From.ActualDefaultCharset())
	from.Put("DEFAULT COLLATE", r.From.ActualDefaultCollate())
	to.Put("DEFAULT CHARACTER SET", r.To.ActualDefaultCharset())
	to.Put("DEFAULT COLLATE", r.To.ActualDefaultCollate())

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

func (r DatabaseOptionAlterations) Diff() []string {
	from := r.From.Map()
	to := r.To.Map()

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
			} else {
				ret = append(ret, fmt.Sprintf("  %s = %s", k, cur))
			}
		}
	}
	ret = parser.Align(ret)
	return ret
}

func (r DatabaseOptionAlterations) Id() string {
	return "DatabaseOptionAlterations"
}
