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


