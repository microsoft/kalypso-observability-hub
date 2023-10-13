package postgres

import "database/sql"

type WorkloadVersion struct {
	Id            int
	Version       string
	BuildId       string
	BuildCommitId string
	WorkloadId    int
}

// make sure that the WorkloadVersion implements the Entity interface
var _ Entity = (*WorkloadVersion)(nil)

func (wv *WorkloadVersion) update(conn *sql.DB) (Entity, error) {
	var id int32
	err := conn.QueryRow(`INSERT INTO workload_version (version, build_id, build_commit_id, workload_id) VALUES ($1, $2, $3, $4)
						 ON CONFLICT (workload_id, version) DO
						 UPDATE SET build_id=$2,
						 			build_commit_id=$3
						 RETURNING id`, wv.Version, wv.BuildId, wv.BuildCommitId, wv.WorkloadId).Scan(&id)
	if err != nil {
		return nil, err
	}
	wv.Id = int(id)
	return wv, nil
}

func (wv *WorkloadVersion) get(conn *sql.DB) (Entity, error) {
	err := conn.QueryRow(`SELECT id, version, build_id, build_commit_id, workload_id FROM workload_version WHERE id=$1`, wv.Id).Scan(&wv.Id, &wv.Version, &wv.BuildId, &wv.BuildCommitId, &wv.WorkloadId)
	if err != nil {
		return nil, err
	}
	return wv, nil
}

func (wv *WorkloadVersion) getByNaturalKey(conn *sql.DB) (Entity, error) {
	err := conn.QueryRow(`SELECT id, version, build_id, build_commit_id, workload_id FROM workload_version WHERE workload_id=$1 AND version=$2`, wv.WorkloadId, wv.Version).Scan(&wv.Id, &wv.Version, &wv.BuildId, &wv.BuildCommitId, &wv.WorkloadId)
	if err != nil {
		return nil, err
	}
	return wv, nil
}
