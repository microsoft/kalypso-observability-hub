package postgres

import "database/sql"

type Deployment struct {
	Id             int
	GitopsCommitId string
	ReconcilerId   int
	Status         string
	StatusMessage  string
}

// make sure that the Host implements the Entity interface
var _ Entity = (*Deployment)(nil)

func (d *Deployment) update(conn *sql.DB) (Entity, error) {
	var id int32
	err := conn.QueryRow(`INSERT INTO deployment (gitops_commit_id, reconciler_id, status, status_message) VALUES ($1, $2, $3, $4)
						  ON CONFLICT (gitops_commit_id, reconciler_id) DO
						  UPDATE SET status=$3,
						   	         status_message=$4,
								     updated_on=current_timestamp,
						             updated_by=current_user
						 RETURNING id`, d.GitopsCommitId, d.ReconcilerId, d.Status, d.StatusMessage).Scan(&id)
	if err != nil {
		return nil, err
	}
	d.Id = int(id)
	return d, nil
}

func (d *Deployment) get(conn *sql.DB) (Entity, error) {
	err := conn.QueryRow(`SELECT id, gitops_commit_id, reconciler_id, status, status_message FROM deployment WHERE id=$1`, d.Id).Scan(&d.Id, &d.GitopsCommitId, &d.ReconcilerId, &d.Status, &d.StatusMessage)
	if err != nil {
		return nil, err
	}
	return d, nil
}

func (d *Deployment) getByNaturalKey(conn *sql.DB) (Entity, error) {
	err := conn.QueryRow(`SELECT id, gitops_commit_id, reconciler_id, status, status_message FROM deployment WHERE gitops_commit_id=$1 AND reconciler_id=$2`, d.GitopsCommitId, d.ReconcilerId).Scan(&d.Id, &d.GitopsCommitId, &d.ReconcilerId, &d.Status, &d.StatusMessage)
	if err != nil {
		return nil, err
	}
	return d, nil
}
