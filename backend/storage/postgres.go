package storage

import (
	"context"
	"log"
	"time"

	// "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ctx = context.Background()

type Postgres struct {
	conn *pgxpool.Pool
}

func MustNew() *Postgres {
	time.Sleep(20 * time.Second)
	connStr := "postgres://postgres:manager@localhost:8081/postgres"

	conn, err := pgxpool.New(ctx, connStr)

	if err != nil || conn == nil {
		log.Println("Connection to postgress Error!")

		panic(err)
	}

	log.Println("Connected to db")

	err = conn.Ping(ctx)
	if err != nil {
		log.Println("error ping db")
		panic(err)
	}

	p := &Postgres{conn}
	querry := `CREATE TABLE IF NOT EXISTS containers (
			id SERIAL PRIMARY KEY,
			address VARCHAR(25) NOT NULL,
			last_ping TIMESTAMP WITH TIME ZONE,
			last_success_ping TIMESTAMP WITH TIME ZONE
			)`
	_, err = p.conn.Exec(ctx, querry)
	if err != nil {
		log.Println("Error create containers table")

		panic(err)
	}
	return p
}
