create table if not exists deployment_target(
    id serial primary key,
    name varchar(30) not null,
    description text,
    workload_id int not null references workload(id),
    environment_id int not null references environment(id),
    labels text,
    manifests_storage_type varchar(15) default 'git',
    manifests_endpoint text not null unique,
    created_on timestamp default current_timestamp,
    created_by varchar(30) default current_user,
    updated_on timestamp default current_timestamp,
    updated_by varchar(30) default current_user,
    unique(workload_id, name),
);


