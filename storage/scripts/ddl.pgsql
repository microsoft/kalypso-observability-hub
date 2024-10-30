-- drop table deployment;
-- drop table reconciler;
-- drop table host;
-- drop table deployment_assignment;
-- drop table workload_version;
-- drop table deployment_target;
-- drop table environment;
-- drop table workload;
-- drop table application;
-- drop table workspace;

create table if not exists workspace (
    id serial primary key,
    name varchar(30) not null unique,
    description text,
    created_on timestamp default current_timestamp,
    created_by varchar(30) default current_user,
    updated_on timestamp default current_timestamp,
    updated_by varchar(30) default current_user
);


create table if not exists application (
    id serial primary key,
    name varchar(30) not null,
    description text,
    workspace_id integer not null references workspace(id), 
    created_on timestamp default current_timestamp,
    created_by varchar(30) default current_user,
    updated_on timestamp default current_timestamp,
    updated_by varchar(30) default current_user,
    UNIQUE (workspace_id, name)
);


create table if not exists workload (
    id serial primary key,
    name varchar(30) not null,
    description text,
    source_storage_type varchar(15) default 'git',
    source_endpoint text,
    application_id integer not null references application(id),
    created_on timestamp default current_timestamp,
    created_by varchar(30) default current_user,
    updated_on timestamp default current_timestamp,
    updated_by varchar(30) default current_user,
    unique(application_id, name)
);


create table if not exists environment (
    id serial primary key,
    name varchar(30) not null unique,
    description text,
    created_on timestamp default current_timestamp,
    created_by varchar(30) default current_user,
    updated_on timestamp default current_timestamp,
    updated_by varchar(30) default current_user
);


create table if not exists deployment_target(
    id serial primary key,
    name varchar(100) not null,
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
    unique(workload_id, environment_id, name)
);


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


create table if not exists deployment_assignment(
    id serial primary key,
    deployment_target_id int not null references deployment_target(id),
    workload_version_id int not null references workload_version(id),
    gitops_commit_id varchar(100) not null,
    created_on timestamp default current_timestamp,
    created_by varchar(30) default current_user,
    unique(deployment_target_id, workload_version_id, gitops_commit_id)
);
create table if not exists host(
    id serial primary key,
    name varchar(50) not null unique,
    description text,
    host_type varchar(20) default 'K8s-cluster',
    labels text,
    created_on timestamp default current_timestamp,
    created_by varchar(30) default current_user,
    updated_on timestamp default current_timestamp,
    updated_by varchar(30) default current_user
);


create table if not exists reconciler(
    id serial primary key,
    name varchar(150) not null,
    host_id int not null references host(id),
    description text,
    reconciler_type varchar(20),
    labels text,
    manifests_storage_type varchar(15) default 'git',
    manifests_endpoint text not null,
    created_on timestamp default current_timestamp,
    created_by varchar(30) default current_user,
    updated_on timestamp default current_timestamp,
    updated_by varchar(30) default current_user,
    unique(host_id, name)    
);


create table if not exists deployment(
    id serial primary key,
    gitops_commit_id varchar(100) not null,
    reconciler_id int not null references reconciler(id),
    status varchar(20) not null,
    status_message text,
    created_on timestamp default current_timestamp,
    created_by varchar(30) default current_user,
    updated_on timestamp default current_timestamp,
    updated_by varchar(30) default current_user,
    unique(gitops_commit_id, reconciler_id)    
);

 CREATE USER hub WITH PASSWORD 'c67';
 GRANT pg_read_all_data TO hub;



