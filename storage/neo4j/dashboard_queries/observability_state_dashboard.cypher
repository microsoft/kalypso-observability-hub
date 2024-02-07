MATCH (ws:Workspace)-[:CONTAINS_APP]->(app:Application)-[:CONTAINS_WORKLOAD]->(wl:Workload)-[:HAS_DEPLOYMENT_TARGET]->(dt:DeploymentTarget)<-[:HAS_DEPLOYMENT_TARGET]-(e:Environment) 

MATCH (dt)-[da:HAS_DEPLOYMENT_ASSIGNMENT]->(wv:WORKLOAD_VERSION)
// RETURN ws,app,wl,dt,wv,e

MATCH (r:Reconciler)<-[:HAS_RECONCILER]-(dt)
MATCH (h:Host)-[:HAS_RECONCILER]->(r)-[:HAS_DEPLOYMENT]->(d:Deployment{gitOpsCommitId: da.gitOpsCommitId})
// RETURN ws,app,wl,dt,wv,e,r,d,h
WITH ws,app,wl,dt,wv,e,r,d,h
ORDER BY wv.createdOn desc, d.createdOn desc
// RETURN ws,app,wl,dt,wv,e,r,d,h
WITH ws,app,wl,dt,e,r,h, COLLECT(wv)[0] AS latestVersion, COLLECT(d)[0] AS latestDeployment
// RETURN ws,app,wl,dt,e,r,h, latestVersion, latestDeployment


// tabular
MATCH (ws:Workspace)-[:CONTAINS_APP]->(app:Application)-[:CONTAINS_WORKLOAD]->(wl:Workload)-[:HAS_DEPLOYMENT_TARGET]->(dt:DeploymentTarget)<-[:HAS_DEPLOYMENT_TARGET]-(e:Environment) 

MATCH (dt)-[da:HAS_DEPLOYMENT_ASSIGNMENT]->(wv:WORKLOAD_VERSION)
// RETURN ws,app,wl,dt,wv,e

MATCH (r:Reconciler)<-[:HAS_RECONCILER]-(dt)
MATCH (h:Host)-[:HAS_RECONCILER]->(r)-[:HAS_DEPLOYMENT]->(d:Deployment{gitOpsCommitId: da.gitOpsCommitId})
// RETURN ws,app,wl,dt,wv,e,r,d,h
WITH ws,app,wl,dt,wv,e,r,d,h
ORDER BY wv.createdOn desc, d.createdOn desc
// RETURN ws,app,wl,dt,wv,e,r,d,h
WITH ws,app,wl,dt,e,r,h, COLLECT(wv)[0] AS latestVersion, COLLECT(d)[0] AS latestDeployment
// RETURN ws,app,wl,dt,e,r,h, latestVersion, latestDeployment
RETURN
    ws.name AS `Workspace`,
    app.name AS `Application / use-case`,
    e.name AS `Environment`,
    wl.name AS `Workload`,                
    dt.name AS `Deployment Target`,
    latestVersion.version AS `Scheduled Version`,
    latestDeployment.status AS `Status`,
    r.name AS `Host`