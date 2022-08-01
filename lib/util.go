package lib

import (
	"fmt"
	"github.com/emirpasic/gods/sets/hashset"
	"github.com/emirpasic/gods/sets/linkedhashset"
	"strings"
)

func optS(s1 string, s2 string) string {
	if s1 == "" {
		return ""
	}
	return fmt.Sprintf(s2, s1)
}

func Find[E any](arr []E, f func(E) bool) int {
	for i, a := range arr {
		if f(a) {
			return i
		}
	}
	return -1
}

func RemoveIf[E any](arr []E, f func(E) bool) []E {
	ret := []E{}
	for _, a := range arr {
		if !f(a) {
			ret = append(ret, a)
		}
	}
	return ret
}

func Index[E comparable](arr []E, v E) int {
	for i, a := range arr {
		if a == v {
			return i
		}
	}
	return -1
}

func IndexIf[E comparable](arr []E, f func(E) bool) int {
	for i, a := range arr {
		if f(a) {
			return i
		}
	}
	return -1
}

func Contains[E comparable](arr []E, v E) bool {
	return Index(arr, v) >= 0
}

func ContainsIf[E comparable](arr []E, f func(E) bool) bool {
	return IndexIf(arr, f) >= 0
}

func Replace[E comparable](arr []E, from E, to E) []E {
	ret := []E{}
	for _, v := range arr {
		if v == from {
			ret = append(ret, to)
		} else {
			ret = append(ret, v)
		}
	}
	return ret
}

func diffArray[E any](a1 []E, a2 []E) []E {
	ret := []E{}
	if a1 == nil {
		return nil
	}
	if a2 == nil {
		return a1
	}
	s2Set := hashset.New()
	for _, v := range a2 {
		s2Set.Add(v)
	}
	for _, v := range a1 {
		if !s2Set.Contains(v) {
			ret = append(ret, v)
		}
	}
	return ret
}

func intersection(s1 *linkedhashset.Set, s2 *linkedhashset.Set) *linkedhashset.Set {
	result := linkedhashset.New()

	if s1.Size() <= s2.Size() {
		for _, item := range s1.Values() {
			if s2.Contains(item) {
				result.Add(item)
			}
		}
	} else {
		for _, item := range s2.Values() {
			if s1.Contains(item) {
				result.Add(item)
			}
		}
	}
	return result
}

func difference(s1 *linkedhashset.Set, s2 *linkedhashset.Set) *linkedhashset.Set {
	result := linkedhashset.New()

	for _, item := range s1.Values() {
		if !s2.Contains(item) {
			result.Add(item)
		}
	}
	return result
}

func arraysEqual[E comparable](a []E, b []E) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

func prefix(s string, p string) string {
	ret := []string{}
	for _, a := range strings.Split(s, "\n") {
		ret = append(ret, p+a)
	}
	return strings.Join(ret, "\n")
}
