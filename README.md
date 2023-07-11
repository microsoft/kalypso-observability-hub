# Kalypso Observability Hub

Kalypso Observability Hub is a central storage that contains deployment data with historical information on workload versions and their deployment state across clusters. This data is used by dashboards and alerts for the deployment monitoring purposes, by the CD pipelines, implementing progressive rollout across environments and by various external systems that make decisions basing on the deployment data. 

## Motivation

Platform and Application Dev teams need a deployment observability solution to perform the following activities:

- monitor what application/service versions are deployed to clusters in the environments
- compare environments and see deployment discrepancy (e.g. how my "stage" environment is different from "prod")
- track deployment history per environment, per application/service, per microservice 
- compare desired deployment state to the reality and see deployment drift


## Data flow

![deployment-observability-hub](./docs/images/deployment-observability-hub.png)

## Observability hub abstractions

### Deployment Descriptor
  
#### Example

```yaml
apiVersion: hub.kalypso.io/v1alpha1
kind: DeploymentDescriptor
metadata:
  name: hello-world-functional-test-0-0-1
spec:
  workload:
    name: hello-world
    source:
        repo: https://github.com/kaizentm/hello-world
        branch: main
        path: .
    application:
        name: greeting-service
        workspace:
            name: kaizentm
  deploymentTarget:
    name: functional-test
    environment: dev
    manifests:
        repo: https://github.com/kaizentm/hello-world-gitops
        branch: dev
        path: functional-test
  workloadVersion:
    version: 0.0.1
    build: build-1
    commit: ca9ee9d0ff9ec52b998fdcf64e128c84ddd0e661
    buildTime: "2023-04-27T23:25:05Z"
```

### Cluster

#### Example

### Deployment

#### Example

## Observability hub API

## Deployment reports

See examples of some [deployment reports](./docs/images/DeploymentObservabilityReports.png) that could be built on top of Observability Hub data.

## Storage data model

See [logical data model](./docs/images/DeploymentObservabilityLogicalModel.drawio) of the Observability Hub storage.

## Installation

## Contributing

This project welcomes contributions and suggestions.  Most contributions require you to agree to a
Contributor License Agreement (CLA) declaring that you have the right to, and actually do, grant us
the rights to use your contribution. For details, visit https://cla.opensource.microsoft.com.

When you submit a pull request, a CLA bot will automatically determine whether you need to provide
a CLA and decorate the PR appropriately (e.g., status check, comment). Simply follow the instructions
provided by the bot. You will only need to do this once across all repos using our CLA.

This project has adopted the [Microsoft Open Source Code of Conduct](https://opensource.microsoft.com/codeofconduct/).
For more information see the [Code of Conduct FAQ](https://opensource.microsoft.com/codeofconduct/faq/) or
contact [opencode@microsoft.com](mailto:opencode@microsoft.com) with any additional questions or comments.

## Trademarks

This project may contain trademarks or logos for projects, products, or services. Authorized use of Microsoft 
trademarks or logos is subject to and must follow 
[Microsoft's Trademark & Brand Guidelines](https://www.microsoft.com/en-us/legal/intellectualproperty/trademarks/usage/general).
Use of Microsoft trademarks or logos in modified versions of this project must not cause confusion or imply Microsoft sponsorship.
Any use of third-party trademarks or logos are subject to those third-party's policies.
