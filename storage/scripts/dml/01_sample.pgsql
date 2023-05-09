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

insert into environment
(name , description)
values
('qa', 'qa environment');

insert into environment
(name , description)
values
('prod', 'prod environment');

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

insert into workload_version
(version, build_id, build_commit_id, workload_id)
values
('1.0.0', 'build_1', '7e8d9dd3514803b7bee6bd83eed4e2d18858e942', 1);

insert into workload_version
(version, build_id, build_commit_id, workload_id)
values
('1.0.1', 'build_2', '4e8d9dd3514803b7bee6bd83eed4e2d18858e943', 1);

insert into deployment_assignment
(deployment_target_id, workload_version_id, gitops_commit_id)
values
(1, 1, '5e8d9dd3514803b7bee6bd83eed4e2d18858e948');

insert into deployment_assignment
(deployment_target_id, workload_version_id, gitops_commit_id, created_on)
values
(1, 2, '2e8d9dd3514803b7bee6bd83eed4e2d18858e944', (NOW() + interval '1 hour'));

insert into host
(name, description)
values
('Large-cl-us-1', 'Large K8s cluster in the US');

insert into host
(name, description)
values
('Small-cl-eu-1', 'Small K8s cluster in the EU');

insert into host
(name, description, host_type)
values
('Ubuntu-vm-us-1', 'Ubuntu VM in the US', 'Ubuntu-vm');

insert into reconciler
(name, host_id, description, reconciler_type, manifests_storage_type, manifests_endpoint)
values
('hello-world', 1, 'Reconciler for the Hello World application', 'Flux', 'git', 'https://github.com/kaizentm/kalypso-app-gitops/dev/functional-test');

insert into reconciler
(name, host_id, description, reconciler_type, manifests_storage_type, manifests_endpoint)
values
('argo', 2, 'Sample Reconciler', 'Argocd', 'git', 'https://github.com/kaizentm/kalypso-app-gitops/dev/performance-test');

insert into reconciler
(name, host_id, description, reconciler_type, manifests_storage_type, manifests_endpoint)
values
('ansible', 3, 'Sample Reconciler', 'Ansible', 'git', 'https://github.com/kaizentm/kalypso-app-gitops/dev/int-test');

insert into deployment
(gitops_commit_id, reconciler_id, status, status_message)
values
('5e8d9dd3514803b7bee6bd83eed4e2d18858e948',1, 'success', 'Successfully deployed');

insert into deployment
(gitops_commit_id, reconciler_id, status, status_message, created_on)
values
('2e8d9dd3514803b7bee6bd83eed4e2d18858e944',1, 'failure', 'Failed to deploy', (NOW() + interval '1 hour'));

