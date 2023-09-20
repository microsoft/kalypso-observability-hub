# Installation guide

## Prerequisites 

- A k8s cluster
- Helm

## Installation

Install Kalypso Observability Hub (KOH) with the following Helm commands:

```sh
helm repo add kalypso-observability-hub https://raw.githubusercontent.com/microsoft/kalypso-observability-hub/gh-pages/ --force-update 
helm upgrade -i kalypso-observability-hub kalypso-observability-hub/kalypso-observability-hub  --create-namespace  -n hub 
```

It installs the following components on your k8s cluster in the `hub` namespace:

- Postgres Database
- KOH Storage API 
- KOH Controller Manager
- Grafana. *Note*, that Grafana component is optional. To omit it add `grafana.enabled` option to the Helm `upgrade` command.

You can check that the installed components are up and running:
```sh
 kubectl get pods -n hub

---
NAME                                                            READY   STATUS    RESTARTS   AGE
kalypso-observability-hub-api-server-56bf784d79-zcrl7           1/1     Running   0          70s
kalypso-observability-hub-controller-manager-6fc48bc875-s9bqt   2/2     Running   0          70s
grafana-6854d5b49c-f86x4                                        2/2     Running   0          70s
postgre-db-0                                                    1/1     Running   0          70s
```  

## Configuring

Kalypso Observability Hub consumes two types of data: Deployment Descriptors and Reconciler statuses. Normally, these entities are coming from different sources, which should be configured differently, 

### Deployment Descriptors

Although, [Deployment Descriptors](../README.md#deployment-descriptor) can be pushed to the KOH cluster simply with a `kubectl` command, normally, they are delivered to the cluster in the pull based GitOps fashion with an operator like Flux. 

Make sure you have Flux installed on the KOH cluster either with any preferable [native Flux methods](https://fluxcd.io/flux/installation/) or with the [Azure Arc GitOps extension](https://learn.microsoft.com/en-us/azure/azure-arc/kubernetes/tutorial-use-gitops-flux2?tabs=azure-cli).    

Create Flux resources to fetch deployment descriptors from a GitOps repository:

```sh
kubectl apply -f - <<EOF
apiVersion: source.toolkit.fluxcd.io/v1beta2
kind: GitRepository
metadata:
  name: deployment-descriptors-sample
  namespace: flux-system
spec:
  interval: 30s
  url: [YOUR REPO] (e.g. https://github.com/kaizentm/kalypso-gitops)
  ref:
    branch: [YOUR BRANCH] (e.g. dev)
---
apiVersion: kustomize.toolkit.fluxcd.io/v1beta2
kind: Kustomization
metadata:
  name: deployment-descriptors-sample
  namespace: flux-system
spec:
  interval: 30s
  targetNamespace: default
  sourceRef:
    kind: GitRepository
    name: deployment-descriptors-sample
  path: [FOLDER WITH DEPLOYMENT DESCRIPTORS IN YOUR REPO]  (e.g. ./deployment-descriptors/samples)
  prune: true
EOF
```

*Note*, if you're using Azure Arc GitOps extension, you can create a configuration above with an [az cli command](https://learn.microsoft.com/en-us/azure/azure-arc/kubernetes/tutorial-use-gitops-flux2?tabs=azure-cli#apply-a-flux-configuration). 

### Reconcilers

[Reconcilers](../README.md#reconciler) can be pushed to the KOH cluster by the workload hosts directly. However, if your workload k8s clusters deploy applications with [Azure Arc GitOps extension](https://learn.microsoft.com/en-us/azure/azure-arc/kubernetes/tutorial-use-gitops-flux2?tabs=azure-cli), they already report the deployment state to Azure Resource Graph (ARG). 

Kalypso Observability Hub has a built in controller to monitor Azure Resource Graph for the deployment state updates. This controller polls ARG and creates [Reconciler](../README.md#reconciler) resources in the KOH cluster automatically. The ARG controller uses a managed identity to authenticate with Azure Resource Graph, so make sure the KOH cluster is configured with a [managed identity](https://learn.microsoft.com/en-us/azure/aks/use-managed-identity).

To configure the ARG controller, create a [AzureResourceGraph](../README.md#arg) resource to the KOH cluster:

```sh
kubectl apply -f - <<EOF
apiVersion: hub.kalypso.io/v1alpha1
kind: AzureResourceGraph
metadata:
  name: azureresourcegraph-sample
spec:
  subscription: [YOUR AZURE SUBSCRIPTION] (e.g 7be1b9e7-57ca-47ff-b5ab-82e7ccb8c611)  
  tenant: [YOUR AZURE TENANT] (e.g. 16b3c013-d300-468d-ac64-7eda0820b6d3)
  managedIdentity: [MANAGED IDENTITY ID] (e.g. 02552706-98f9-4301-a473-017752fc430b)
  interval: 10s
EOF
```

## Monitoring

By default, the [Helm chart](#installation) installs Grafana with the [preconfigured dashboards](../README.md#deployment-reports). To access these dashboards, run the following command:

```sh
kubectl port-forward svc/grafana 3000:3000 -n hub
```
and go to the http://localhost:3000 with your browser.
