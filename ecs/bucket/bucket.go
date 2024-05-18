package bucket

import (
	"fmt"

	"github.com/Evankj/ecgo/ecs/component"
	"github.com/Evankj/ecgo/ecs/core"
	"github.com/Evankj/ecgo/ecs/entity"
)

type Bucket struct {
	Size          core.Size
	ComponentMap  map[string]*component.Component
	UnusedIndexes []core.Size
	MaxComponents core.Size
	Entities      []*entity.Entity
}

func NewBucket() *Bucket {
	return &Bucket{
		Size:          0,
		ComponentMap:  make(map[string]*component.Component),
		UnusedIndexes: []core.Size{},
		Entities:      []*entity.Entity{},
		MaxComponents: (64 / 8), // 64 bits in bytes TODO: Revisit this to remove this limitation
	}
}

func (b *Bucket) DeleteEntityById(entityId core.Size) error {
	entity := b.Entities[entityId]
	if entity == nil {
		return fmt.Errorf("No entity with id %d in this bucket", entityId)
	}

	b.Entities[entityId].Mask = 0
	b.UnusedIndexes = append(b.UnusedIndexes, entityId)

	b.Size -= 1

	return nil
}

func (b *Bucket) CreateEntity() core.Size {
	index, err := core.Pop(&b.UnusedIndexes)
	if err != nil {
		// unused index slice is empty
		// time to push a new entry onto all of the component arrays
		for _, v := range b.ComponentMap {
			v.Entries = append(v.Entries, nil)
		}

		length := core.Size(len(b.Entities))

		entity := entity.Entity{
			Index: length,
			Mask:  0,
		}
		b.Entities = append(b.Entities, &entity)

		b.Size += 1

		return entity.Index
	}

	entity := entity.Entity{
		Index: index,
		Mask:  0,
	}

	// reuse a used index
	b.Entities[index] = &entity

	return entity.Index
}


func registerComponent[T interface{}](b *Bucket) error {
	typeId := core.TypeId[T]()
	if b.ComponentMap[typeId] != nil {
		return fmt.Errorf("Component with type id: %s already registered!", typeId)
	}

	entries := []interface{}{}
  for range len(b.Entities) {
    entries = append(entries, nil)
  }

	mask := uint64(1 << len(b.ComponentMap))

	b.ComponentMap[typeId] = &component.Component{
		Entries: entries,
		Mask:    mask,
	}

	return nil
}

func AddComponentToEntityByID[T interface{}](b *Bucket, entityId core.Size, component *T) error {
	typeId := core.TypeId[T]()
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

  if len(c.Entries) - 1 < int(entity.Index)  {
    return fmt.Errorf("No component entry for this entity's index in this bucket")
  }

  if entity == nil {
    return fmt.Errorf("No entity with provided ID exists in this bucket")
  }

	entity.Mask |= mask

  c.Entries[entity.Index] = component

	return nil
}
