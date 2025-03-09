package db

import (
	"database/sql"
	"fmt"
	"log"

	"example.com/rogue/util"
	"example.com/rogue/world"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

type Database struct {
	connection *sql.DB
}

func GetDatabase(host, user, name, password string, port int) Database {
	dbInfo := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable",
		host, port, user, name, password)

	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	return Database{
		connection: db,
	}
}

func (db *Database) InsertChunk(x, y int, chunk world.Chunk) error {
	_, err := db.connection.Exec(`
        INSERT INTO Chunk (x, y, world)
        VALUES ($1, $2, $3)
        ON CONFLICT (x, y) DO NOTHING`,
		x, y, pq.Array(util.Flatten2DArray(chunk.World)),
	)

	if err != nil {
		return fmt.Errorf("failed to insert chunk: %w", err)
	}

	log.Printf("Chunk at (%d, %d) inserted successfully.", x, y)
	return nil
}

func (db *Database) GetChunk(x, y int) (world.Chunk, error) {
	var flatWorld []string
	err := db.connection.QueryRow("SELECT world FROM Chunk WHERE x = $1 AND y = $2", x, y).
		Scan(pq.Array(&flatWorld))

	if err != nil {
		if err == sql.ErrNoRows {
			return world.Chunk{}, nil
		}
		return world.Chunk{}, err
	}

	chunkWorld, err := util.Unflatten2DArray(flatWorld, world.COLS)

	if err != nil {
		return world.Chunk{}, err
	}

	return world.Chunk{
		World: chunkWorld,
	}, nil
}
