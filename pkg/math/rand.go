package gmath

import (
	"math/rand"
)

type Int3264 interface {
	~int32 | ~int64 | ~int
}

func GetARandNumber[T Int3264](min T, max T, inclusive bool) T {
	if inclusive {
		max++
	}

	// detect if the type is int32 or int64
	if _, ok := any(min).(int32); ok {
		return T(rand.Int31n(int32(max-min)) + int32(min))
	}

	if _, ok := any(min).(int64); ok {
		return T(rand.Int63n(int64(max-min)) + int64(min))
	}

	return T(rand.Intn(int(max-min)) + int(min))
}
