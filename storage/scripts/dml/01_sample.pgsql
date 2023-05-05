insert into workspace
(name, description)
values
('kaizentm', 'kaizentm sample workspace');

insert into application
(name, description, workspace_id)
values
('Hello World', 'Hello World sample application', 1);

insert into workload
(name, description, source_storage_type, source_endpoint, application_id)
values
('Hello World', 'Hello World sample workload', 'git', 'https://github.com/kaizentm/kalypso-app-src', 1);

insert into environment
(name , description)
values
('dev', 'dev environment');

insert into deployment_target
(name , description, workload_id, environment_id, manifests_storage_type, manifests_endpoint)
values
('functional test', 'functional test deployment target', 1, 1, 'git', 'https://github.com/kaizentm/kalypso-app-gitops/dev/functional-test');

insert into deployment_target
(name , description, workload_id, environment_id, manifests_storage_type, manifests_endpoint)
values
('performance test', 'performance test deployment target', 1, 1, 'git', 'https://github.com/kaizentm/kalypso-app-gitops/dev/performance-test');

insert into deployment_target
(name , description, workload_id, environment_id, manifests_storage_type, manifests_endpoint)
values
('integration test', 'integration test deployment target', 1, 1, 'git', 'https://github.com/kaizentm/kalypso-app-gitops/dev/int-test');

workload_version(
    id serial primary key,
    version varchar(30) not null,
    build_id varchar(40) not null,
    build_commit_id varchar(40) not null,
    workload_id int not null references workload(id),