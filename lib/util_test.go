package lib

import (
	"fmt"
	"sort"
	"testing"
)

func TestSort(t *testing.T) {

	elems := map[string]int{
		"b": 1,
		"a": 2,
		"d": 3,
	}

	arr := [][]string{{"a"}, {"b"}, {"d", "b"}, {"a", "b", "d"}, {"b", "d"}, {"b", "d", "a"}}

	sort.SliceStable(arr, func(i, j int) bool {
		if len(arr[i]) != len(arr[j]) {
			return len(arr[i]) < len(arr[j])
		}
		ln := len(arr[i])
		for a := 0; a < ln; a++ {
			if elems[arr[i][a]] != elems[arr[j][a]] {
				return elems[arr[i][a]] < elems[arr[j][a]]
			}
		}
		return true
	})

	fmt.Println(arr)
}

type I interface {
	Get() int
}

type Impl1 struct {
}

func (r Impl1) Get() int {
	return 1
}

type Impl2 struct {
}

func (r Impl2) Get() int {
	return 2
}

func TestAlterationaa(t *testing.T) {

	arr := []I{}

	arr = append(arr, &Impl1{})
	arr = append(arr, &Impl2{})

}
