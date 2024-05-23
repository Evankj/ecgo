package ecgo


type ComponentType struct {
	ComponentId string
}

type Component struct {
	Entries []any // slice of "any"s so we can use nil for non-present components
	Mask    uint64
}

func GetComponentType[T any]() ComponentType {
	return ComponentType {
		ComponentId: TypeId[T](),
	}
}
