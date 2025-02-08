package storage

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5"
)

var ctx = context.Background()

type Postgres struct {
	conn *pgx.Conn
}

func MustNew() *Postgres {
	connStr := "postgres://postgres:manager@localhost:8081/postgres"
	conn, err := pgx.Connect(ctx, connStr)
	if err != nil {
		log.Println("Connection to postgress Error!")

		panic(err)
	}

	p := Postgres{conn}
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
	return &p
}
