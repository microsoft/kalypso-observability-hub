package postgres

import "database/sql"

type Environment struct {
	Id          int
	Name        string
	Description string
}

// make sure that the Environment implements the Entity interface
var _ Entity = (*Environment)(nil)

func (en *Environment) update(conn *sql.DB) (Entity, error) {
	var id int32
	err := conn.QueryRow(`INSERT INTO environment (name, description) VALUES ($1, $2)
						 ON CONFLICT (name) DO
						 UPDATE SET description=$2,
								   updated_on=current_timestamp,
						           updated_by=current_user
						 RETURNING id`, en.Name, en.Description).Scan(&id)
	if err != nil {
		return nil, err
	}
	en.Id = int(id)
	return en, nil
}

func (en *Environment) get(conn *sql.DB) (Entity, error) {
	err := conn.QueryRow(`SELECT id, name, description FROM environment WHERE id=$1`, en.Id).Scan(&en.Id, &en.Name, &en.Description)
	if err != nil {
		return nil, err
	}
	return en, nil
}

func (en *Environment) getByNaturalKey(conn *sql.DB) (Entity, error) {
	err := conn.QueryRow(`SELECT id, name, description FROM environment WHERE name=$1`, en.Name).Scan(&en.Id, &en.Name, &en.Description)
	if err != nil {
		return nil, err
	}
	return en, nil
}
