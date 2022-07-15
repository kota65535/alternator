package parser

import (
	"fmt"
	"reflect"
	"strings"
)

func findFirstIdentifier(s string) string {
	isIn := false
	name := ""
	for _, c := range s {
		if c == '`' {
			if isIn {
				return name
			}
			isIn = !isIn
		} else if isIn {
			name = name + string(c)
		}
	}
	return ""
}

func compactJoin(array []string, delim string) string {
	var ret []string
	for _, e := range array {
		if e != "" {
			ret = append(ret, e)
		}
	}
	return strings.Join(ret, delim)
}

func optB(b bool, s string) string {
	if !b {
		return ""
	}
	return s
}

func optS(s1 string, s2 string) string {
	if s1 == "" {
		return ""
	}
	return fmt.Sprintf(s2, s1)
}

func JoinS(elems []string, sep string, enclose string) string {
	var enclosed []string
	for _, e := range elems {
		enclosed = append(enclosed, enclose+e+enclose)
	}
	return strings.Join(enclosed, sep)
}

func JoinT[T fmt.Stringer](elems []T, sep string, enclose string) string {
	var enclosed []string
	for _, e := range elems {
		enclosed = append(enclosed, enclose+e.String()+enclose)
	}
	return strings.Join(enclosed, sep)
}

func Align(lines []string) []string {
	var matrix [][]string
	for _, l := range lines {
		matrix = append(matrix, strings.Split(strings.TrimRight(l, "\t\n "), "\t"))
	}
	nRows := len(lines)
	if nRows == 0 {
		return nil
	}
	nMaxCols := 0
	for i := 0; i < nRows; i++ {
		if nMaxCols < len(matrix[i]) {
			nMaxCols = len(matrix[i])
		}
	}

	for i := 0; i < nMaxCols; i++ {
		maxLen := 0
		for j := 0; j < nRows; j++ {
			if i < len(matrix[j]) {
				str := matrix[j][i]
				if maxLen < len(str) {
					maxLen = len(str)
				}
			}
		}
		for j := 0; j < nRows; j++ {
			if i >= len(matrix[j])-1 {
				continue
			}
			str := matrix[j][i]
			matrix[j][i] = str + strings.Repeat(" ", maxLen-len(str))
		}
	}

	var ret []string

	for i := 0; i < nRows; i++ {
		ret = append(ret, strings.Join(matrix[i], " "))

	}

	return ret
}

func structDifference[T any](i1 T, i2 T) T {
	t1 := reflect.TypeOf(i1)
	ret := reflect.New(t1).Elem()
	v1 := reflect.ValueOf(i1)
	v2 := reflect.ValueOf(i2)
	for i := 0; i < t1.NumField(); i++ {
		f1 := v1.Field(i)
		f2 := v2.Field(i)
		if f1.Kind() == reflect.String {
			f1str := f1.Interface().(string)
			f2str := f2.Interface().(string)
			if f1str != "" && f1str != f2str {
				ret.Field(i).Set(f1)
			}
		}
	}
	return ret.Interface().(T)
}

func stripSequentialSpaces(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

func addIndent(str []string, indent int) []string {
	ret := []string{}
	for _, s := range str {
		ret = append(ret, strings.Repeat(" ", indent)+s)
	}
	return ret
}
