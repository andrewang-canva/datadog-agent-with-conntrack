package util

import "sort"

// SortUniqInPlace sorts and remove duplicates from elements in place
// The returned slice is a subslice of elements
func SortUniqInPlace(elements []string) []string {
	if len(elements) < 2 {
		return elements
	}
	size := len(elements)
	if size <= InsertionSortThreshold {
		InsertionSort(elements)
	} else {
		// this will trigger an alloc because sorts uses interface{} internaly
		// which confuses the escape analysis
		sort.Strings(elements)
	}
	return uniqSorted(elements)
}

func DedupInPlace(elements []string) []string {
	if len(elements) < 2 {
		return elements
	}

	m := make(map[string]struct{})
	idx := 0
	for i := range elements {
		if _, exists := m[elements[i]]; !exists {
			m[elements[i]] = struct{}{}
			elements[idx] = elements[i]
			idx++
		}
	}

	return elements[:idx]
}

// uniqSorted remove duplicate elements from the given slice
// the given slice needs to be sorted
func uniqSorted(elements []string) []string {
	j := 0
	for i := 1; i < len(elements); i++ {
		if elements[j] == elements[i] {
			continue
		}
		j++
		elements[j] = elements[i]
	}
	return elements[:j+1]
}
