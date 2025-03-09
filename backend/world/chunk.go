package world

import (
	"math/rand/v2"

	"example.com/rogue/util"
)

type Chunk struct {
	World [][]string `json:"world"`
}

type room struct {
	height int
	width  int

	rowOffset int
	colOffset int
}

const ROOM_SIZE = 12
const ROWS = 36
const COLS = 48

func GenerateChunk(northHallways, eastHallways, southHallways, westHallways []int) Chunk {
	chunk := Chunk{
		World: make([][]string, ROWS),
	}
	for i := range chunk.World {
		chunk.World[i] = make([]string, COLS)
	}

	rooms := make([][]room, ROWS/ROOM_SIZE)
	for i := range rooms {
		rooms[i] = make([]room, COLS/ROOM_SIZE)
	}

	for rowN := 0; rowN < ROWS; rowN += ROOM_SIZE {
		for colN := 0; colN < COLS; colN += ROOM_SIZE {
			roomHeight := util.RandRange(4, ROOM_SIZE-2)
			roomWidth := util.RandRange(4, ROOM_SIZE-2)

			rowOffset := util.RandRange(1, ROOM_SIZE-roomHeight-1)

			colOffset := util.RandRange(1, ROOM_SIZE-roomWidth-1)

			for y := rowN + rowOffset; y <= rowN+rowOffset+roomHeight; y++ {
				for x := colN + colOffset; x <= colN+colOffset+roomWidth; x++ {
					chunk.World[y][x] = MaterialStoneFloor
				}
			}

			rooms[rowN/ROOM_SIZE][colN/ROOM_SIZE] = room{
				height:    roomHeight,
				width:     roomWidth,
				rowOffset: rowOffset,
				colOffset: colOffset,
			}
		}

	}

	chunk.addHallways(rooms)
	chunk.addNorthCorridorToChunk(rooms, northHallways)
	chunk.addEastCorridorToChunk(rooms, eastHallways)
	chunk.addSouthCorridorToChunk(rooms, southHallways)
	chunk.addWestCorridorToChunk(rooms, westHallways)
	return chunk
}

func (c *Chunk) addHallways(rooms [][]room) {
	connectPercentage := float32(util.RandRange(5, 11)) / 10

	for rowN := 0; rowN < ROWS/ROOM_SIZE; rowN += 1 {
		for colN := 0; colN < COLS/ROOM_SIZE; colN += 1 {
			room := rooms[rowN][colN]

			if colN+1 < COLS/ROOM_SIZE {
				rightNeighbor := rooms[rowN][colN+1]

				hallwayStart := max(room.rowOffset, rightNeighbor.rowOffset) + 1
				hallwayEnd := min(room.rowOffset+room.height, rightNeighbor.rowOffset+rightNeighbor.height) - 1
				hallwayHeight := hallwayEnd - hallwayStart

				connectRoll := rand.Float32()

				if connectRoll < connectPercentage && hallwayHeight > 0 {
					hallwayHeight -= util.RandRange(0, hallwayHeight)
					hallwayWidth := ROOM_SIZE - room.colOffset - room.width + rightNeighbor.colOffset
					rowOffset := hallwayStart + rowN*ROOM_SIZE
					colOffset := room.colOffset + room.width + colN*ROOM_SIZE

					for y := rowOffset; y < rowOffset+hallwayHeight; y++ {
						for x := colOffset; x < colOffset+hallwayWidth; x++ {
							c.World[y][x] = MaterialStoneFloor
						}
					}
				}
			}

			if rowN+1 < ROWS/ROOM_SIZE {
				downNeighbor := rooms[rowN+1][colN]

				hallwayStart := max(room.colOffset, downNeighbor.colOffset) + 1
				hallwayEnd := min(room.colOffset+room.width, downNeighbor.colOffset+downNeighbor.width) - 1
				hallwayWidth := hallwayEnd - hallwayStart

				connectRoll := rand.Float32()

				if connectRoll < connectPercentage && hallwayWidth > 0 {
					hallwayWidth -= util.RandRange(0, hallwayWidth)
					hallwayHeight := ROOM_SIZE - room.rowOffset - room.height + downNeighbor.rowOffset
					rowOffset := room.rowOffset + room.height + rowN*ROOM_SIZE
					colOffset := hallwayStart + colN*ROOM_SIZE

					for y := rowOffset; y < rowOffset+hallwayHeight; y++ {
						for x := colOffset; x < colOffset+hallwayWidth; x++ {
							c.World[y][x] = MaterialStoneFloor
						}
					}
				}
			}
		}
	}
}

func (c *Chunk) addWestCorridorToChunk(rooms [][]room, requiredYs []int) {
	for _, y := range requiredYs {
		roomRow := y / ROOM_SIZE
		room := rooms[roomRow][0]

		for x := 0; x <= room.colOffset; x++ {
			c.World[y][x] = MaterialStoneFloor
		}

		for j := y - 1; j >= roomRow*ROOM_SIZE+room.rowOffset+room.height; j-- {
			c.World[j][room.colOffset] = MaterialStoneFloor
		}

		for j := y + 1; j <= roomRow*ROOM_SIZE+room.rowOffset; j++ {
			c.World[j][room.colOffset] = MaterialStoneFloor
		}
	}
}

func (c *Chunk) addEastCorridorToChunk(rooms [][]room, requiredYs []int) {
	for _, y := range requiredYs {
		roomRow := y / ROOM_SIZE
		room := rooms[roomRow][3]

		for x := COLS - 1; x >= COLS-(ROOM_SIZE-room.colOffset)+room.width; x-- {
			c.World[y][x] = MaterialStoneFloor
		}

		for j := y - 1; j >= roomRow*ROOM_SIZE+room.rowOffset+room.height; j-- {
			c.World[j][ROOM_SIZE*3+room.colOffset+room.width] = MaterialStoneFloor
		}

		for j := y + 1; j <= roomRow*ROOM_SIZE+room.rowOffset; j++ {
			c.World[j][ROOM_SIZE*3+room.colOffset+room.width] = MaterialStoneFloor
		}
	}
}

func (c *Chunk) addNorthCorridorToChunk(rooms [][]room, requiredXs []int) {
	for _, x := range requiredXs {
		roomCol := x / ROOM_SIZE
		room := rooms[0][roomCol]

		for y := 0; y <= room.rowOffset; y++ {
			c.World[y][x] = MaterialStoneFloor
		}

		for i := x - 1; i >= roomCol*ROOM_SIZE+room.colOffset+room.width; i-- {
			c.World[room.rowOffset][i] = MaterialStoneFloor
		}

		for i := x + 1; i <= roomCol*ROOM_SIZE+room.colOffset; i++ {
			c.World[room.rowOffset][i] = MaterialStoneFloor
		}
	}
}

func (c *Chunk) addSouthCorridorToChunk(rooms [][]room, requiredXs []int) {
	for _, x := range requiredXs {
		roomCol := x / ROOM_SIZE
		room := rooms[2][roomCol]

		for y := ROWS - 1; y >= ROWS-(ROOM_SIZE-room.rowOffset)+room.height; y-- {
			c.World[y][x] = MaterialStoneFloor
		}

		for i := x - 1; i >= roomCol*ROOM_SIZE+room.colOffset+room.width; i-- {
			c.World[ROOM_SIZE*2+room.rowOffset+room.height][i] = MaterialStoneFloor
		}

		for i := x + 1; i <= roomCol*ROOM_SIZE+room.colOffset; i++ {
			c.World[ROOM_SIZE*2+room.rowOffset+room.height][i] = MaterialStoneFloor
		}
	}
}
