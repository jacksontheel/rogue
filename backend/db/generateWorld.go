package db

import (
	"fmt"

	"example.com/rogue/util"
	"example.com/rogue/world"
)

type chunkExits struct {
	north *[]int
	east  *[]int
	south *[]int
	west  *[]int
}

func pointerSliceToRandom(s *[]int, max int) *[]int {
	if s == nil {
		ret := util.RandIntSlice(0, 5, 1, max)
		return &ret
	}

	return s
}

func pointerSliceToSlice(s *[]int) []int {
	if s == nil {
		return make([]int, 0)
	}

	return *s
}

func GenerateWorld(db Database, radius int) {
	exitsMap := make(map[string]chunkExits)

	fmt.Printf("%v\n", exitsMap["hello"])

	for x := -radius; x <= radius; x++ {
		for y := -radius; y <= radius; y++ {
			exits := chunkExits{
				north: pointerSliceToRandom(exitsMap[fmt.Sprintf("%d,%d", x, y-1)].south, world.COLS-1),
				east:  pointerSliceToRandom(exitsMap[fmt.Sprintf("%d,%d", x+1, y)].west, world.ROWS-1),
				south: pointerSliceToRandom(exitsMap[fmt.Sprintf("%d,%d", x, y+1)].north, world.COLS-1),
				west:  pointerSliceToRandom(exitsMap[fmt.Sprintf("%d,%d", x-1, y)].east, world.ROWS-1),
			}

			exitsMap[fmt.Sprintf("%d,%d", x, y)] = exits

			err := db.InsertChunk(
				x,
				y,
				world.GenerateChunk(
					pointerSliceToSlice(exits.north),
					pointerSliceToSlice(exits.east),
					pointerSliceToSlice(exits.south),
					pointerSliceToSlice(exits.west),
				),
			)

			if err != nil {
				fmt.Println(err.Error())
			}
		}
	}
	fmt.Println("done")
}
