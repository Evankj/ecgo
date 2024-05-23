package ecgo

import "fmt"

type Size int

// TODO: Use this in place of other mask variables so it can be expanded without refactoring other code
type Mask uint64

func TypeId[T any]() string {
	return fmt.Sprintf("%T", *new(T))
}
