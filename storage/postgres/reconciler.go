package postgres

import (
	"database/sql"
)

type Reconciler struct {
	Id                   int
	Name                 string
	HostId               int
	Description          string
	ReconcilerType       string
	Labels               string
	ManifestsStorageType string
	ManifestsEndpoint    string
}

type StatusStats struct {
	Success    int32
	Failed     int32
	InProgress int32
}

// make sure that the Host implements the Entity interface
var _ Entity = (*Reconciler)(nil)

func (r *Reconciler) update(conn *sql.DB) (Entity, error) {
	var id int32
	err := conn.QueryRow(`INSERT INTO reconciler (name, host_id, description, reconciler_type, labels, manifests_storage_type, manifests_endpoint) 
						  VALUES ($1, $2, $3, $4, $5, $6, $7)
						  ON CONFLICT (host_id, name) DO
						  UPDATE SET description=$3,
						             reconciler_type=$4,
									 labels=$5,
									 manifests_storage_type=$6,
									 manifests_endpoint=$7,
								     updated_on=current_timestamp,
						             updated_by=current_user
						 RETURNING id`, r.Name, r.HostId, r.Description, r.ReconcilerType, r.Labels, r.ManifestsStorageType, r.ManifestsEndpoint).Scan(&id)
	if err != nil {
		return nil, err
	}
	r.Id = int(id)
	return r, nil
}

func (r *Reconciler) get(conn *sql.DB) (Entity, error) {
	err := conn.QueryRow(`SELECT id, name, host_id, description, reconciler_type, labels, manifests_storage_type, manifests_endpoint FROM reconciler WHERE id=$1`, r.Id).Scan(&r.Id, &r.Name, &r.HostId, &r.Description, &r.ReconcilerType, &r.Labels, &r.ManifestsStorageType, &r.ManifestsEndpoint)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (r *Reconciler) getByNaturalKey(conn *sql.DB) (Entity, error) {
	err := conn.QueryRow(`SELECT id, name, host_id, description, reconciler_type, labels, manifests_storage_type, manifests_endpoint FROM reconciler WHERE host_id=$1 AND name=$2`, r.HostId, r.Name).Scan(&r.Id, &r.Name, &r.HostId, &r.Description, &r.ReconcilerType, &r.Labels, &r.ManifestsStorageType, &r.ManifestsEndpoint)
	if err != nil {
		return nil, err
	}
	return r, nil
}

var _ QueryFunc = CountByManifestsEndpoint

func CountByManifestsEndpoint(conn *sql.DB, args ...interface{}) (interface{}, error) {
	var count int32

	manifest_endpoint := args[0].(string)
	err := conn.QueryRow(`SELECT count(1) FROM reconciler WHERE manifests_endpoint like $1`, manifest_endpoint).Scan(&count)
	if err != nil {
		return nil, err
	}
	return count, nil
}

var _ QueryFunc = CountByStatuses

func CountByStatuses(conn *sql.DB, args ...interface{}) (interface{}, error) {

	status := StatusStats{}

	manifest_endpoint := args[0].(string)
	gitops_commit_id := args[1].(string)
	err := conn.QueryRow(`SELECT count(case when status='success' then 1 end) as success,
								count(case when status='failure' then 1 end) as failed,
								count(case when status='in_progress' then 1 end) as in_progress
							FROM reconciler r, deployment d
							WHERE r.id = d.reconciler_id and
								r.manifests_endpoint like $1 and
								d.gitops_commit_id = $2`,
		manifest_endpoint, gitops_commit_id).Scan(&status.Success, &status.Failed, &status.InProgress)
	if err != nil {
		return nil, err
	}
	return status, nil
}

var _ QueryFunc = GetByStatus

func GetByStatus(conn *sql.DB, args ...interface{}) (interface{}, error) {
	manifest_endpoint := args[0].(string)
	gitops_commit_id := args[1].(string)
	status := args[2].(string)
	reconcilers := make([]map[string]string, 0)

	rows, err := conn.Query(`SELECT r.name, d.status_message
							FROM reconciler r, deployment d
							WHERE r.id = d.reconciler_id and
								r.manifests_endpoint like $1 and
								d.gitops_commit_id = $2 and
								d.status = $3`,
		manifest_endpoint, gitops_commit_id, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var reconciler_name string
		var reconciler_status_message string

		err := rows.Scan(&reconciler_name, &reconciler_status_message)
		if err != nil {
			return nil, err
		}
		reconcilers = append(reconcilers, map[string]string{"name": reconciler_name, "status_message": reconciler_status_message})
	}
	return reconcilers, nil

}
