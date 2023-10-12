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

func (app *Application) get(conn *sql.DB) (Entity, error) {
	err := conn.QueryRow(`SELECT id, name, description, workspace_id FROM application WHERE id=$1`, app.Id).Scan(&app.Id, &app.Name, &app.Description, &app.WorkspaceId)
	if err != nil {
		return nil, err
	}
	return app, nil
}

func (app *Application) getByNaturalKey(conn *sql.DB) (Entity, error) {
	err := conn.QueryRow(`SELECT id, name, description, workspace_id FROM application WHERE workspace_id=$1 AND name=$2`, app.WorkspaceId, app.Name).Scan(&app.Id, &app.Name, &app.Description, &app.WorkspaceId)
	if err != nil {
		return nil, err
	}
	return app, nil
}
