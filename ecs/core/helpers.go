package core


import (
	"fmt"
)

func Pop[T interface{}](s *[]T) (T, error) {

	sVal := *s

	length := len(sVal)

	if length <= 0 {
		var zero T
		return zero, fmt.Errorf("cannot pop from an empty slice")
	}

	lastIndex := length - 1
	elem := sVal[lastIndex]
	*s = append(sVal[:lastIndex], sVal[lastIndex+1:]...)
	return elem, nil
}
