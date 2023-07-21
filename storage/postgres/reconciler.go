package postgres

import "database/sql"

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
