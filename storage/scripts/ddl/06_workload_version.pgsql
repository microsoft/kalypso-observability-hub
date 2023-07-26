create table if not exists workload_version(
    id serial primary key,
    version varchar(30) not null,
    build_id varchar(40),
    build_commit_id varchar(40) not null,
    workload_id int not null references workload(id),
    created_on timestamp default current_timestamp,
    created_by varchar(30) default current_user,
    unique(workload_id, version)
);


