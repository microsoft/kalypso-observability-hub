MATCH (ws:Workspace)-[:CONTAINS_APP]->(app:Application)-[:CONTAINS_WORKLOAD]->(wl:Workload)-[:HAS_DEPLOYMENT_TARGET]->(dt:DeploymentTarget)<-[:HAS_DEPLOYMENT_TARGET]-(e:Environment) 

//Modify filter conditions in below line
WHERE e.name = "dev" AND dt.name = "functional-test" AND wl.name = "hello"

MATCH (dt)-[da:HAS_DEPLOYMENT_ASSIGNMENT]->(wv:WORKLOAD_VERSION) 
    WITH  ws,app,dt,da,wl,e,wv ORDER BY wv.createdOn desc limit 1
        OPTIONAL MATCH (r:Reconciler)<-[:HAS_RECONCILER]-(dt)
        OPTIONAL MATCH (h:Host)-[:HAS_RECONCILER]->(r)-[:HAS_DEPLOYMENT]->(d:Deployment{gitOpsCommitId: da.gitOpsCommitId})
        RETURN 
            ws.name AS `Workspace`,
            app.name AS `Application / use-case`,
            e.name AS `Environment`,
            h.name AS `Cluster`,
            r.name AS `Host`,
            wl.name AS `Workload`,
            dt.name AS `Deployment Target`,
            wv.version AS `Deployed Version`,
            d.gitOpsCommitId AS `GitOps commit`,
            d.deployedOn AS `Deployed On`,
            d.status AS `Status`