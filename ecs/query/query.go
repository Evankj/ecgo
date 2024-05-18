package query

import (
	"fmt"

	"github.com/Evankj/ecgo/ecs/bucket"
	"github.com/Evankj/ecgo/ecs/core"
)

type Query struct {
	Bucket *bucket.Bucket
	Mask   uint64
}

type QueryResult struct {
	Bucket *bucket.Bucket
	Index  core.Size
}

func NewQuery(b *bucket.Bucket) *Query {
	return &Query{
		Bucket: b,
		Mask:   0,
	}
}

func AddComponentToQuery[T interface{}](q *Query) error {
	typeId := core.TypeId[T]()
	component := q.Bucket.ComponentMap[typeId]
	if component == nil {
		return fmt.Errorf("Component with id %s is not registered on this bucket", typeId)
	}

	q.Mask |= component.Mask

	return nil
}

// Runs and returns the list of entity Ids with the components specified in the query
func (q *Query) Execute() []QueryResult {
	entities := []QueryResult{}
	for _, entity := range q.Bucket.Entities {
		if entity.Mask&q.Mask == q.Mask {
			entities = append(entities, QueryResult{
				Bucket: q.Bucket,
				Index:  entity.Index,
			})
		}
	}

	return entities
}

func GetComponentFromQueryResult[T interface{}](q *QueryResult) (*T, error) {
	typeId := core.TypeId[T]()

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
