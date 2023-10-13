create table if not exists reconciler(
    id serial primary key,
    name varchar(150) not null unique,
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


