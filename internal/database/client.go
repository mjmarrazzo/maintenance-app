package database

import (
	"context"
	_ "embed"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Client struct {
	pool *pgxpool.Pool
}

var once sync.Once
var db *Client

func InitPool() (*Client, error) {
	var err error
	var pool *pgxpool.Pool

	once.Do(func() {
		connString := fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			os.Getenv("DB_HOST"),
			os.Getenv("DB_PORT"),
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_NAME"),
		)

		pool, err = pgxpool.New(context.Background(), connString)
		if err != nil {
			log.Printf("Error creating connection pool: %v", err)
			return
		}

		err = pool.Ping(context.Background())
		if err != nil {
			log.Printf("Error pinging database: %v", err)
			return
		}

		db = &Client{pool: pool}
	})

	return db, err
}

func (db *Client) Pool() *pgxpool.Pool {
	return db.pool
}

func (db *Client) Close() {
	if db.pool != nil {
		db.pool.Close()
	}
}
