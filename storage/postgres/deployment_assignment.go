package postgres

import "database/sql"

type DeploymentAssignment struct {
	Id                 int
	DeploymentTargetId int
	WorkloadVersionId  int
	GitopsCommitId     string
}

// make sure that the DeploymentAssignment implements the Entity interface
var _ Entity = (*DeploymentAssignment)(nil)

func (da *DeploymentAssignment) update(conn *sql.DB) (Entity, error) {
	var id int32
	err := conn.QueryRow(`INSERT INTO deployment_assignment (deployment_target_id, workload_version_id, gitops_commit_id) VALUES ($1, $2, $3)	                     
						 ON CONFLICT (deployment_target_id, workload_version_id) DO
						 UPDATE SET gitops_commit_id=$3
						 RETURNING id`, da.DeploymentTargetId, da.WorkloadVersionId, da.GitopsCommitId).Scan(&id)
	if err != nil {
		return nil, err
	}
	da.Id = int(id)
	return da, nil
}
