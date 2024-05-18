package bucket

import (
	"testing"

	"github.com/Evankj/ecgo/ecs/core"
)

func TestCreateEntity(t *testing.T) {

	bucket := NewBucket()

	bucket.CreateEntity()
	bucket.CreateEntity()
	bucket.CreateEntity()

	got := len(bucket.Entities)
	want := 3
	if got != want {
		t.Fatalf("got: %d\nwant: %d", got, want)
	}
}

func TestAddComponentToEntity(t *testing.T) {
	b := NewBucket()

	entity1 := b.CreateEntity()
	entity2 := b.CreateEntity()

	type TestComponent1 struct{}
	type TestComponent2 struct{}

	tc1 := TestComponent1{}
	tc2 := TestComponent2{}

	err := AddComponentToEntityByID(b, entity1, &tc1)
	if err != nil {
		t.Fatal(err)
	}

	err = AddComponentToEntityByID(b, entity1, &tc2)
	if err != nil {
		t.Fatal(err)
	}

	err = AddComponentToEntityByID(b, entity2, &tc1)
	if err != nil {
		t.Fatal(err)
	}

	entries := len(b.ComponentMap[core.TypeId[TestComponent1]()].Entries)
	if entries != 2 {
		t.Fatalf("Expected 2 entry for bucket's component map entries, got %d", entries)
	}

	expectedEntity1Mask := 3
	entity1Mask := b.Entities[entity1].Mask
	if entity1Mask != uint64(expectedEntity1Mask) {
		t.Fatalf("Entity 1 Mask - Expected %d, got %d", expectedEntity1Mask, entity1Mask)
	}

	expectedEntity2Mask := 1
	entity2Mask := b.Entities[entity2].Mask
	if entity2Mask != uint64(expectedEntity2Mask) {
		t.Fatalf("Entity 2 Mask - Expected %d, got %d", expectedEntity2Mask, entity2Mask)
	}

}

func TestDeleteEntityById(t *testing.T) {

	b := NewBucket()

	b.CreateEntity()

	err := b.DeleteEntityById(0)
	if err != nil {
		t.Fatal(err)
	}

	if len(b.UnusedIndexes) != 1 {
		t.Fatalf("Failed to append deleted index to unused indexes")
	}

	if b.Size != 0 {
		t.Fatalf("Failed to update bucket size")
	}

}
