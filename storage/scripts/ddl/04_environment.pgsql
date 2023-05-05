create table if not exists environment (
    id serial primary key,
    name varchar(30) not null unique,
    description text,
    created_on timestamp default current_timestamp,
    created_by varchar(30) default current_user,
    updated_on timestamp default current_timestamp,
    updated_by varchar(30) default current_user
);


