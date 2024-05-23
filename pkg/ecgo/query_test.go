package ecgo

import (
	"fmt"
	"testing"

)

func TestQuery(t *testing.T) {
	b := NewBucket()

	type TestPosComponent struct {
		x int
		y int
	}
	type TestRotComponent struct {
		angle float32
	}

	entityId := b.CreateEntity()

	err := AddComponentToEntityByID(b, entityId, &TestPosComponent{
		x: 10,
		y: 20,
	})
	if err != nil {
		t.Fatal(err)
	}
	err = AddComponentToEntityByID(b, entityId, &TestRotComponent{
		angle: 45.0,
	})
	if err != nil {
		t.Fatal(err)
	}

	entityId = b.CreateEntity()
	err = AddComponentToEntityByID(b, entityId, &TestPosComponent{
		x: 0,
		y: 0,
	})
	if err != nil {
		t.Fatal(err)
	}

	q := NewQuery(b)
	q.WithComponents(
		GetComponentType[TestPosComponent](),
		GetComponentType[TestRotComponent](),
	)

	res := q.Execute()

  if len(res) != 1 {
    t.Fatalf("Expected one result for this query")
  }

	pos, err := GetComponentFromQueryResult[TestPosComponent](&res[0])
	if err != nil || pos.x != 10 {
		t.Fatalf("Failed to get pos component from entity")
	}

	rot, err := GetComponentFromQueryResult[TestRotComponent](&res[0])
	if err != nil || rot.angle != 45.0 {
		t.Fatalf("Failed to get rot component from entity")
	}

	q = NewQuery(b)
	err = q.WithComponents(GetComponentType[TestPosComponent]())
	if err != nil {
		t.Fatal(err)
	}
	err = q.WithoutComponents(GetComponentType[TestRotComponent]())
	if err != nil {
		t.Fatal(err)
	}

	res = q.Execute()

	fmt.Println(len(res))

	if len(res) != 1 {
		t.Fatalf("Expected one result for this query!")
	}

}

func TestBulkQuery(t *testing.T) {
	b := NewBucket()

	type TestPosComponent struct {
		x int
		y int
	}

	type TestRotComponent struct {
		angle float32
	}

	for index := range 10 {

		entityId := b.CreateEntity()

		err := AddComponentToEntityByID(b, entityId, &TestPosComponent{
			x: index,
			y: index,
		})
		if err != nil {
			t.Fatal(err)
		}

		err = AddComponentToEntityByID(b, entityId, &TestRotComponent{
			angle: float32(index),
		})
		if err != nil {
			t.Fatal(err)
		}
	}

	q := NewQuery(b)
	err := q.WithComponents(GetComponentType[TestPosComponent]())
	if err != nil {
		t.Fatal(err)
	}

	res := q.Execute()
	fmt.Println(len(res))

	for index, queryResult := range res {
		pos, err := GetComponentFromQueryResult[TestPosComponent](&queryResult)
		if err != nil {
			t.Fatalf("Failed to get pos component from entity - index: %d", index)
		}
		if pos.x != index || pos.y != index {
			t.Fatalf("Failed to get pos component from entity - index: %d", index)
		}

		rot, err := GetComponentFromQueryResult[TestRotComponent](&queryResult)

		if err != nil {
			t.Fatalf("Failed to get rot component from entity - index: %d", index)
		}
		if rot.angle != float32(index) {
			t.Fatalf("Failed to get rot component from entity - index: %d", index)
		}
	}

}
