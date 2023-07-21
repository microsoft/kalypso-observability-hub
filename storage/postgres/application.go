package postgres

import "database/sql"

type Application struct {
	Id          int
	Name        string
	Description string
	WorkspaceId int
}

// make sure that the Application implements the Entity interface
var _ Entity = (*Application)(nil)

func (app *Application) update(conn *sql.DB) (Entity, error) {
	var id int32
	err := conn.QueryRow(`INSERT INTO application (name, description, workspace_id) VALUES ($1, $2, $3)
						 ON CONFLICT (workspace_id, name) DO
						 UPDATE SET description=$2,
								   updated_on=current_timestamp,
						           updated_by=current_user
						 RETURNING id`, app.Name, app.Description, app.WorkspaceId).Scan(&id)
	if err != nil {
		return nil, err
	}
	app.Id = int(id)
	return app, nil
}
