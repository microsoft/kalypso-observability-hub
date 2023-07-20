package postgres

import "database/sql"

type Workspace struct {
	Id          int
	Name        string
	Description string
}

// make sure that the Workspace implements the Entity interface
var _ Entity = (*Workspace)(nil)

func (ws *Workspace) update(conn *sql.DB) (Entity, error) {
	var id int32
	err := conn.QueryRow(`INSERT INTO workspace (name, description) VALUES ($1, $2)
						 ON CONFLICT (name) DO
						 UPDATE SET description=$2,
								   updated_on=current_timestamp,
						           updated_by=current_user
						 RETURNING id`, ws.Name, ws.Description).Scan(&id)
	if err != nil {
		return nil, err
	}
	ws.Id = int(id)
	return ws, nil
}
