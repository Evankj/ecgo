package entity

import "github.com/Evankj/ecgo/ecs/core"

type Entity struct {
	Index core.Size
	Mask  uint64
}
