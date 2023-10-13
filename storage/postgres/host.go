package postgres

import "database/sql"

type Host struct {
	Id          int
	Name        string
	Description string
	HostType    string
	Labels      string
}

// make sure that the Host implements the Entity interface
var _ Entity = (*Host)(nil)

func (host *Host) update(conn *sql.DB) (Entity, error) {
	var id int32
	err := conn.QueryRow(`INSERT INTO host (name, description, host_type, labels) VALUES ($1, $2, $3, $4)
						 ON CONFLICT (name) DO
						 UPDATE SET description=$2,
						 		    host_type=$3,
									labels=$4,
								    updated_on=current_timestamp,
						            updated_by=current_user
						 RETURNING id`, host.Name, host.Description, host.HostType, host.Labels).Scan(&id)
	if err != nil {
		return nil, err
	}
	host.Id = int(id)
	return host, nil
}

func (host *Host) get(conn *sql.DB) (Entity, error) {
	err := conn.QueryRow(`SELECT id, name, description, host_type, labels FROM host WHERE id=$1`, host.Id).Scan(&host.Id, &host.Name, &host.Description, &host.HostType, &host.Labels)
	if err != nil {
		return nil, err
	}
	return host, nil
}

func (host *Host) getByNaturalKey(conn *sql.DB) (Entity, error) {
	err := conn.QueryRow(`SELECT id, name, description, host_type, labels FROM host WHERE name=$1`, host.Name).Scan(&host.Id, &host.Name, &host.Description, &host.HostType, &host.Labels)
	if err != nil {
		return nil, err
	}
	return host, nil
}
