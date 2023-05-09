-- Environment State
select dt.id deployment_target_id, dt.name deployment_target_name,
       e.name environemnt_name,
       wk.name workspace_name,
       a.name application_name, 
       w.name workload_name,
       wv.version workload_version,
       deployments.deployed_workload_version,
       deployments.status,
       deployments.reconcilers
from deployment_target dt
     left outer join (select dt.id deployment_target_id, depl.status, depl.deployed_workload_version, count(1) reconcilers
     from reconciler r
          left outer join (select d.reconciler_id, d.status, wv.version deployed_workload_version
           from deployment d,
                deployment_assignment da,
                workload_version wv
            where d.gitops_commit_id = da.gitops_commit_id 
              and da.workload_version_id = wv.id
              and d.created_on = (select max(created_on) from deployment where reconciler_id = d.reconciler_id)
          ) depl on r.id = depl.reconciler_id,
          deployment_target dt
     where r.manifests_storage_type = dt.manifests_storage_type
          and r.manifests_endpoint = dt.manifests_endpoint
    group by deployment_target_id, status, deployed_workload_version) deployments on dt.id = deployments.deployment_target_id,    
     deployment_assignment da,
     workload_version wv,
     workload w,
     application a,
     workspace wk,
     environment e
where dt.id = da.deployment_target_id
    and da.workload_version_id = wv.id
    and wv.workload_id = w.id
    and w.application_id = a.id
    and a.workspace_id = wk.id
    and dt.environment_id = e.id    
    and da.created_on = (select max(created_on) from deployment_assignment where deployment_target_id = dt.id);


-- Host State
select h.Name host_name, 
       wk.name workspace_name,
       a.name application_name, 
       w.name workload_name,
       dt.name deployment_target_name,
       dt.id deployment_target_id, depl.status, depl.deployed_workload_version, depl.gitops_commit_id, depl.deployed_on
from reconciler r
    left outer join (select d.reconciler_id, d.status, wv.version deployed_workload_version, d.gitops_commit_id, d.created_on deployed_on
    from deployment d,
        deployment_assignment da,
        workload_version wv
    where d.gitops_commit_id = da.gitops_commit_id 
        and da.workload_version_id = wv.id
        and d.created_on = (select max(created_on) from deployment where reconciler_id = d.reconciler_id)
    ) depl on r.id = depl.reconciler_id,
    deployment_target dt,
    host h,
    workspace wk,
    workload w,
    application a,
    environment e
where r.manifests_storage_type = dt.manifests_storage_type
    and r.manifests_endpoint = dt.manifests_endpoint
    and r.host_id = h.id
    and dt.workload_id = w.id
    and w.application_id = a.id
    and a.workspace_id = wk.id 
    and dt.environment_id = e.id;


-- Enviroment History
select h.Name host_name, 
       wk.name workspace_name,
       a.name application_name, 
       w.name workload_name,
       dt.name deployment_target_name,
       dt.id deployment_target_id,        
       depl.deployed_workload_version, 
       depl.gitops_commit_id,
       sum (case depl.status
        when 'success' THEN 1
        else 0
       end) success_count,
       sum (case depl.status
        when 'failure' THEN 1
        else 0
       end) failure_count,
       sum (case depl.status
        when 'in_progress' THEN 1
        else 0
       end) in_progress_count
from reconciler r
    left outer join (select d.reconciler_id, d.status, wv.version deployed_workload_version, d.gitops_commit_id, d.created_on deployed_on
    from deployment d,
        deployment_assignment da,
        workload_version wv
    where d.gitops_commit_id = da.gitops_commit_id 
        and da.workload_version_id = wv.id        
    ) depl on r.id = depl.reconciler_id,
    deployment_target dt,
    host h,
    workspace wk,
    workload w,
    application a,
    environment e
where r.manifests_storage_type = dt.manifests_storage_type
    and r.manifests_endpoint = dt.manifests_endpoint
    and r.host_id = h.id
    and dt.workload_id = w.id
    and w.application_id = a.id
    and a.workspace_id = wk.id 
    and dt.environment_id = e.id
group by h.Name, wk.name, a.name, w.name, dt.name, dt.id, depl.deployed_workload_version, depl.gitops_commit_id;    

-- Environment Discrepancies
select dt.id deployment_target_id, dt.name deployment_target_name,       
       wk.name workspace_name,
       a.name application_name, 
       w.name workload_name,
       max(case e.name
        when 'dev' THEN wv.version
        else null
       end) dev_workload_version,
       max(case e.name
        when 'dev' THEN deployments.deployed_workload_version
        else null
       end) dev_deployed_workload_version,
       max(case e.name
        when 'dev' THEN successfull_reconcilers
        else null
       end) dev_successfull_reconcilers,
       max(case e.name
        when 'dev' THEN reconcilers
        else null
       end) dev_reconcilers,

       max(case e.name
        when 'qa' THEN wv.version
        else null
       end) qa_workload_version,
       max(case e.name
        when 'qa' THEN deployments.deployed_workload_version
        else null
       end) qa_deployed_workload_version,
       max(case e.name
        when 'qa' THEN successfull_reconcilers
        else null
       end) qa_successfull_reconcilers,
       max(case e.name
        when 'qa' THEN reconcilers
        else null
       end) qa_reconcilers,

       max(case e.name
        when 'prod' THEN wv.version
        else null
       end) prod_workload_version,
       max(case e.name
        when 'prod' THEN deployments.deployed_workload_version
        else null
       end) prod_deployed_workload_version,
       max(case e.name
        when 'prod' THEN successfull_reconcilers
        else null
       end) prod_successfull_reconcilers,
       max(case e.name
        when 'prod' THEN reconcilers
        else null
       end) prod_reconcilers       

from deployment_target dt
     left outer join (select dt.id deployment_target_id, depl.deployed_workload_version, 
                        sum(
                        case depl.status
                            when 'success' THEN 1
                            else 0
                        end
                     ) successfull_reconcilers,
                     count(1) reconcilers
     from reconciler r
          left outer join (select d.reconciler_id, d.status, wv.version deployed_workload_version
           from deployment d,
                deployment_assignment da,
                workload_version wv
            where d.gitops_commit_id = da.gitops_commit_id 
              and da.workload_version_id = wv.id
              and d.created_on = (select max(created_on) from deployment where reconciler_id = d.reconciler_id )
          ) depl on r.id = depl.reconciler_id,
          deployment_target dt
     where r.manifests_storage_type = dt.manifests_storage_type
          and r.manifests_endpoint = dt.manifests_endpoint
    group by deployment_target_id, deployed_workload_version) deployments on dt.id = deployments.deployment_target_id,    
     deployment_assignment da,
     workload_version wv,
     workload w,
     application a,
     workspace wk,
     environment e
where dt.id = da.deployment_target_id
    and da.workload_version_id = wv.id
    and wv.workload_id = w.id
    and w.application_id = a.id
    and a.workspace_id = wk.id
    and dt.environment_id = e.id    
    and da.created_on = (select max(created_on) from deployment_assignment where deployment_target_id = dt.id)
group by dt.id,  deployment_target_name,
       workspace_name,
       application_name, 
       workload_name    

