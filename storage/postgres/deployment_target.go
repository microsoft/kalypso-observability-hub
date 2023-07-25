package postgres

import "database/sql"

type DeploymentTarget struct {
	Id                   int
	Name                 string
	Description          string
	WorkloadId           int
	EnvironmentId        int
	Labels               string
	ManifestsStorageType string
	ManifestsEndpoint    string
}

// make sure that the DeploymentTarget implements the Entity interface
var _ Entity = (*DeploymentTarget)(nil)

func (dt *DeploymentTarget) update(conn *sql.DB) (Entity, error) {
	var id int32
	err := conn.QueryRow(`INSERT INTO deployment_target 
	                     (name, description, workload_id, environment_id, labels, manifests_storage_type, manifests_endpoint) 
						 VALUES ($1, $2, $3, $4, $5, $6, $7)
						 ON CONFLICT (workload_id, environment_id, name) DO
						 UPDATE SET description=$2,
									labels=$5, 
									manifests_storage_type=$6, 
									manifests_endpoint=$7,
								    updated_on=current_timestamp,
						            updated_by=current_user
						 RETURNING id`, dt.Name, dt.Description, dt.WorkloadId, dt.EnvironmentId, dt.Labels, dt.ManifestsStorageType, dt.ManifestsEndpoint).Scan(&id)
	if err != nil {
		return nil, err
	}
	dt.Id = int(id)
	return dt, nil
}
