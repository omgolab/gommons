package gccollections

import (
	"sort"

	"golang.org/x/exp/constraints"
)

func SortMapKeysByValue[K comparable, V constraints.Ordered](m map[K]V, isAscending bool) []K {
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}

	if isAscending {
		sort.SliceStable(keys, func(i, j int) bool {
			return m[keys[i]] < m[keys[j]]
		})
	} else {
		sort.SliceStable(keys, func(i, j int) bool {
			return m[keys[i]] > m[keys[j]]
		})
	}

	return keys
}

func MergeMaps[M ~map[K]V, K comparable, V any](src ...M) M {
	merged := make(M)
	for _, m := range src {
		for k, v := range m {
			merged[k] = v
		}
	}
	return merged
}
