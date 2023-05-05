create table if not exists deployment_assignment(
    id serial primary key,
    deployment_target_id int not null references deployment_target(id),
    workload_version_id int not null references workload_version(id),
    gitops_commit_id varchar(40) not null,
    created_on timestamp default current_timestamp,
    created_by varchar(30) default current_user,
    unique(deployment_target_id, workload_version_id)
);


