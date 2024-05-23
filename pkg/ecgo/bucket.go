package ecgo

import (
	"fmt"
)

type Bucket struct {
	Size          Size
	ComponentMap  map[string]*Component
	UnusedIndexes []Size
	MaxComponents Size
	Entities      []*Entity
}

func NewBucket() *Bucket {
	return &Bucket{
		Size:          0,
		ComponentMap:  make(map[string]*Component),
		UnusedIndexes: []Size{},
		Entities:      []*Entity{},
		MaxComponents: (64 / 8), // 64 bits in bytes TODO: Revisit this to remove this limitation
	}
}

func (b *Bucket) DeleteEntityById(entityId Size) error {
	entity := b.Entities[entityId]
	if entity == nil {
		return fmt.Errorf("No entity with id %d in this bucket", entityId)
	}

	b.Entities[entityId].Mask = 0
	b.UnusedIndexes = append(b.UnusedIndexes, entityId)

	b.Size -= 1

	return nil
}

func (b *Bucket) CreateEntity() Size {
	index, err := Pop(&b.UnusedIndexes)
	if err != nil {
		// unused index slice is empty
		// time to push a new entry onto all of the component arrays
		for _, v := range b.ComponentMap {
			v.Entries = append(v.Entries, nil)
		}

		length := Size(len(b.Entities))

		entity := Entity{
			Index: length,
			Mask:  0,
		}
		b.Entities = append(b.Entities, &entity)

		b.Size += 1

		return entity.Index
	}

	entity := Entity{
		Index: index,
		Mask:  0,
	}

	// reuse a used index
	b.Entities[index] = &entity

	return entity.Index
}

func registerComponent[T any](b *Bucket) error {
	typeId := TypeId[T]()
	if b.ComponentMap[typeId] != nil {
		return fmt.Errorf("Component with type id: %s already registered!", typeId)
	}

	entries := []any{}
	for range len(b.Entities) {
		entries = append(entries, nil)
	}

	mask := uint64(1 << len(b.ComponentMap))

	b.ComponentMap[typeId] = &Component{
		Entries: entries,
		Mask:    mask,
	}

	return nil
}

func AddComponentToEntityByID[T any](b *Bucket, entityId Size, component *T) error {
	typeId := TypeId[T]()
	if b.ComponentMap[typeId] == nil {
		err := registerComponent[T](b)
		if err != nil {
			// This should not happen as we confirmed above
			return err
		}
	}

	c := b.ComponentMap[typeId]

	mask := c.Mask

	entity := b.Entities[entityId]

	if len(c.Entries)-1 < int(entity.Index) {
		return fmt.Errorf("No component entry for this entity's index in this bucket")
	}

	if entity == nil {
		return fmt.Errorf("No entity with provided ID exists in this bucket")
	}

	entity.Mask |= mask

	c.Entries[entity.Index] = component

	return nil
}
