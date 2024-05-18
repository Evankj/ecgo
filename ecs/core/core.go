package core

import "fmt"

type Size int


func TypeId[T interface{}]() string {
	return fmt.Sprintf("%T", *new(T))
}
