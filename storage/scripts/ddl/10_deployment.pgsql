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


