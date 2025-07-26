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
						  ON CONFLICT (deployment_target_id, workload_version_id, gitops_commit_id) DO
						  UPDATE SET deployment_target_id=$1,
									 workload_version_id=$2, 
									 gitops_commit_id=$3 
						  RETURNING id`, da.DeploymentTargetId, da.WorkloadVersionId, da.GitopsCommitId).Scan(&id)
	if err != nil {
		return nil, err
	}
	da.Id = int(id)
	return da, nil
}

func (da *DeploymentAssignment) get(conn *sql.DB) (Entity, error) {
	err := conn.QueryRow(`SELECT id, deployment_target_id, workload_version_id, gitops_commit_id FROM deployment_assignment WHERE id=$1`, da.Id).Scan(&da.Id, &da.DeploymentTargetId, &da.WorkloadVersionId, &da.GitopsCommitId)
	if err != nil {
		return nil, err
	}
	return da, nil
}

func (da *DeploymentAssignment) getByNaturalKey(conn *sql.DB) (Entity, error) {
	err := conn.QueryRow(`SELECT id, deployment_target_id, workload_version_id, gitops_commit_id FROM deployment_assignment WHERE deployment_target_id=$1 AND workload_version_id=$2 AND gitops_commit_id=$3`, da.DeploymentTargetId, da.WorkloadVersionId, da.GitopsCommitId).Scan(&da.Id, &da.DeploymentTargetId, &da.WorkloadVersionId, &da.GitopsCommitId)
	if err != nil {
		return nil, err
	}
	return da, nil
}

// GetByDeploymentTargetAndCommit is a custom query function to find deployment assignment
// by deployment target ID and git commit ID (without requiring workload version ID)
var _ QueryFunc = GetByDeploymentTargetAndCommit

func GetByDeploymentTargetAndCommit(conn *sql.DB, args ...interface{}) (interface{}, error) {
	deploymentTargetId := args[0].(int)
	gitopsCommitId := args[1].(string)

	var da DeploymentAssignment
	err := conn.QueryRow(
		`SELECT id, deployment_target_id, workload_version_id, gitops_commit_id 
		 FROM deployment_assignment 
		 WHERE deployment_target_id=$1 AND gitops_commit_id=$2`,
		deploymentTargetId, gitopsCommitId).Scan(&da.Id, &da.DeploymentTargetId, &da.WorkloadVersionId, &da.GitopsCommitId)
	if err != nil {
		return nil, err
	}
	return &da, nil
}
