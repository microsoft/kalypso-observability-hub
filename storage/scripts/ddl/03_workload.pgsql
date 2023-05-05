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


