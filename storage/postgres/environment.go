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
