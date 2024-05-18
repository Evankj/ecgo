package component

type Component struct {
	Entries []interface{} // slice of pointers to "any"s which means we can use nil for non-present components
	Mask    uint64
}
