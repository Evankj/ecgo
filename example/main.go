package main

import (
	"fmt"
	"math/rand/v2"

	"github.com/Evankj/ecgo/ecs/bucket"
	"github.com/Evankj/ecgo/ecs/core"
	"github.com/Evankj/ecgo/ecs/query"
	rl "github.com/gen2brain/raylib-go/raylib"
)

const MAX_MOVE_SPEED = 0.5
const WINDOW_WIDTH = 800
const WINDOW_HEIGHT = 450

const MOVER_COUNT = 1

type Position rl.Vector2

type Velocity rl.Vector2

type Speed float32

func generateWorld() *bucket.Bucket {
	return bucket.NewBucket()
}

func createPlayer(world *bucket.Bucket) *core.Size {
	player := world.CreateEntity()

	bucket.AddComponentToEntityByID(world, player, &Position{
		X: 0,
		Y: 0,
	})

	bucket.AddComponentToEntityByID(world, player, &Velocity{
		X: rand.Float32(),
		Y: rand.Float32(),
	})

	return &player
}

func createRandomMover(world *bucket.Bucket) *core.Size {
	entity := world.CreateEntity()

	err := bucket.AddComponentToEntityByID(world, entity, &Position{
		X: 0,
		Y: 0,
	})

	if err != nil {
		fmt.Println("Error adding position component")
		panic(err)
	}

	err = bucket.AddComponentToEntityByID(world, entity, &Velocity{
		X: rand.Float32(),
		Y: rand.Float32(),
	})
	if err != nil {
		fmt.Println("Error adding velocity component")
		panic(err)
	}

	speed := Speed(rand.Float32() * MAX_MOVE_SPEED)
	err = bucket.AddComponentToEntityByID(world, entity, &speed)
	if err != nil {
		fmt.Println("Error adding speed component")
		panic(err)
	}

	return &entity
}

func main() {
	rl.InitWindow(WINDOW_WIDTH, WINDOW_HEIGHT, "raylib [core] example - basic window")
	defer rl.CloseWindow()

	rl.SetTargetFPS(60)

	world := generateWorld()

	moverQuery := query.NewQuery(world)
	err := query.AddComponentToQuery[Position](moverQuery)
	if err != nil {
		fmt.Println("Error adding position component to query")
	}
	query.AddComponentToQuery[Speed](moverQuery)
	if err != nil {
		fmt.Println("Error adding speed component to query")
	}
	query.AddComponentToQuery[Velocity](moverQuery)
	if err != nil {
		fmt.Println("Error adding velocity component to query")
	}

	for range MOVER_COUNT {
		createRandomMover(world)
	}

	for !rl.WindowShouldClose() {

		targetPos := rl.Vector2{
			X: float32(rl.GetMouseX()),
			Y: float32(rl.GetMouseY()),
		}

		fps := rl.GetFPS()
		fmt.Println(fps)

		// Execute query each frame in case we add/remove movers etc.
		movers := moverQuery.Execute()

		rl.BeginDrawing()

		for _, mover := range movers {
			pos, _ := query.GetComponentFromQueryResult[Position](&mover)
			speed, _ := query.GetComponentFromQueryResult[Speed](&mover)

			dir := rl.Vector2Subtract(targetPos, rl.Vector2(*pos))

      move := rl.Vector2Scale(dir, float32(*speed))
			pos.X += move.X
			pos.Y += move.Y

			rl.DrawRectangle(int32(pos.X), int32(pos.Y), 10, 10, rl.Red)
		}

		rl.ClearBackground(rl.RayWhite)
		rl.EndDrawing()
	}
}
