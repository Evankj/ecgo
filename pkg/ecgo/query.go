package ecgo

import (
	"fmt"

)

type Query struct {
	Bucket      *Bucket
	WithMask    uint64
	WithoutMask uint64
}

type QueryResult struct {
	Bucket *Bucket
	Index  Size
}

func NewQuery(b *Bucket) *Query {
	return &Query{
		Bucket:      b,
		WithMask:    0,
		WithoutMask: 0,
	}
}

func (q *Query) WithComponents(t ...ComponentType) error {
	for _, c := range t {
		component := q.Bucket.ComponentMap[c.ComponentId]
		if component == nil {
			return fmt.Errorf("Component type %s not regisered in this query's bucket.", c.ComponentId)
		}
		q.WithMask |= component.Mask
	}

	return nil
}

func (q *Query) WithoutComponents(t ...ComponentType) error {
	for _, c := range t {
		component := q.Bucket.ComponentMap[c.ComponentId]
		if component == nil {
			return fmt.Errorf("Component type %s not regisered in this query's bucket.", c.ComponentId)
		}
		q.WithoutMask |= component.Mask
	}

	return nil
}

func AddComponentToQuery[T any](q *Query) error {
	typeId := TypeId[T]()
	component := q.Bucket.ComponentMap[typeId]
	if component == nil {
		return fmt.Errorf("Component with id %s is not registered on this bucket", typeId)
	}

	q.WithMask |= component.Mask

	return nil
}

// Runs and returns the list of entity Ids with the components specified in the query
func (q *Query) Execute() []QueryResult {
	entities := []QueryResult{}
	for _, entity := range q.Bucket.Entities {
		if (entity.Mask&q.WithMask == q.WithMask) && (entity.Mask&q.WithoutMask == q.WithoutMask) {
			entities = append(entities, QueryResult{
				Bucket: q.Bucket,
				Index:  entity.Index,
			})
		}
	}

	return entities
}

func GetComponentFromQueryResult[T any](q *QueryResult) (*T, error) {
	typeId := TypeId[T]()

	componentMap := q.Bucket.ComponentMap[typeId]
	if componentMap == nil {
		return nil, fmt.Errorf("No component map registered for this type")
	}

	entry := componentMap.Entries[q.Index]
	if entry == nil {
		return nil, fmt.Errorf("No component for this index")
	}

	component, ok := (entry).(*T)
	if !ok {
		return nil, fmt.Errorf("Component does not match T - failed assertion")
	}

	return component, nil
}
