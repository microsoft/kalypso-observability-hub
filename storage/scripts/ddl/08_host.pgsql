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


