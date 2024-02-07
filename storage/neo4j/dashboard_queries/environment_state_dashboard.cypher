MATCH (ws:Workspace)-[:CONTAINS_APP]->(app:Application)-[:CONTAINS_WORKLOAD]->(wl:Workload)-[:HAS_DEPLOYMENT_TARGET]->(dt:DeploymentTarget)<-[:HAS_DEPLOYMENT_TARGET]-(e:Environment)

// Modify filter conditions below
WHERE e.name = "dev" AND dt.name = "functional-test" AND wl.name = "hello"

MATCH (dt)-[da:HAS_DEPLOYMENT_ASSIGNMENT]->(wv:WORKLOAD_VERSION) 
    WITH  ws,app,dt,da,wl,e,wv
        OPTIONAL MATCH (r:Reconciler)<-[:HAS_RECONCILER]-(dt)
        OPTIONAL MATCH (h:Host)-[:HAS_RECONCILER]->(r)-[:HAS_DEPLOYMENT]->(d:Deployment{gitOpsCommitId: da.gitOpsCommitId})
        RETURN 
            ws.name AS `Workspace`,
            app.name AS `Application / use-case`,
            e.name AS `Environment`,
            h.name AS `Cluster`,
            wl.name AS `Workload`,
            dt.name AS `Deployment Target`,
            wv.version AS `Deployed Version`,
            d.gitOpsCommitId AS `GitOps commit`,
            MAX(d.deployedOn) AS `Deployed On`,
            SUM(case when d.status = "success" then 1 else 0 end) AS `Success`,
            SUM(case when d.status = "failure" then 1 else 0 end) AS `Failure`,
            SUM(case when d.status = "in_progress" then 1 else 0 end) AS `In Progress`