package storage

import (
	"time"

	"github.com/Kpatoc452/container_manager/models"
)

func (p *Postgres) Get(id int) (models.Container, error) {
	var container models.Container

	data := p.conn.QueryRow(ctx, "SELECT * FROM containers WHERE id=$1", id)
	err := data.Scan(&container.Id, &container.Address, &container.LastPing, &container.LastSuccessPing)
	return container, err
}

func (p *Postgres) GetAll() ([]models.Container, error) {
	var containers []models.Container

	rows, err := p.conn.Query(ctx, "SELECT * FROM containers ")
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var container models.Container
		err = rows.Scan(&container.Id, &container.Address, &container.LastPing, &container.LastSuccessPing)
		if err != nil {
			return nil, err
		}

		containers = append(containers, container)
	}

	return containers, err
}

func (p *Postgres) Create(address string) error {
	_, err := p.conn.Exec(ctx, "INSERT INTO containers(address, last_ping, last_success_ping) VALUES($1, $2, $3)", address, time.Date(1970, time.January, 1, 0, 0, 0, 0, time.Now().Location()), time.Date(1970, time.January, 1, 0, 0, 0, 0, time.Now().Location()))
	return err
}

func (p *Postgres) Update(id int, newAddress string) error {
	_, err := p.conn.Exec(ctx, "UPDATE containers SET address=$1 WHERE id=$2", newAddress, id)

	return err
}

func (p *Postgres) Delete(id int) error {
	_, err := p.conn.Exec(ctx, "DELETE FROM containers WHERE id=$1", id)
	return err
}
