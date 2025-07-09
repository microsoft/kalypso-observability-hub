/*
Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	azidentity "github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	armkubernetesconfiguration "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/kubernetesconfiguration/armkubernetesconfiguration"
	armresourcegraph "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resourcegraph/armresourcegraph"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	hubv1alpha1 "github.com/microsoft/kalypso-observability-hub/api/v1alpha1"
	grpcClient "github.com/microsoft/kalypso-observability-hub/storage/api/grpc/client"
	pb "github.com/microsoft/kalypso-observability-hub/storage/api/grpc/proto"
)

// AzureResourceGraphReconciler reconciles a AzureResourceGraph object
type AzureResourceGraphReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=hub.kalypso.io,resources=azureresourcegraphs,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=hub.kalypso.io,resources=azureresourcegraphs/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=hub.kalypso.io,resources=azureresourcegraphs/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the AzureResourceGraph object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.13.1/pkg/reconcile
func (r *AzureResourceGraphReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	reqLogger := log.FromContext(ctx)
	reqLogger.Info("=== Reconciling Azure Resource Graph  ===")

	// Fetch the AzureResourceGraph instance
	arg := &hubv1alpha1.AzureResourceGraph{}
	err := r.Get(ctx, req.NamespacedName, arg)
	if err != nil {
		ignroredNotFound := client.IgnoreNotFound(err)
		if ignroredNotFound != nil {
			reqLogger.Error(err, "Failed to get Azure Resource Graph")
		}
		return ctrl.Result{}, ignroredNotFound
	}

	// Check if the resource is being deleted
	if !arg.DeletionTimestamp.IsZero() {
		return ctrl.Result{}, nil
	}

	//get Azure Credentials
	cred, err := r.getAzureCredentials(arg)
	if err != nil {
		return r.manageFailure(ctx, reqLogger, arg, err, "Failed to get Azure Credentials")
	}

	argClient, err := r.getARGClient(cred)
	if err != nil {
		return r.manageFailure(ctx, reqLogger, arg, err, "Failed to get Azure Resource Graph client")
	}

	//create arm client for flux configurations
	armClient, err := r.getARMClient(cred, arg.Spec.Subscription)
	if err != nil {
		return r.manageFailure(ctx, reqLogger, arg, err, "Failed to get Azure Resource Manager client")
	}

	// Get GitOps Reconcilers Data
	reconcilersData, err := r.getGitOpsReconcilersData(ctx, cred, argClient, arg.Spec.Subscription, reqLogger)
	if err != nil {
		return r.manageFailure(ctx, reqLogger, arg, err, "Failed to get GitOps Reconcilers Data")
	}

	// Add WO Reconcilers Data
	woReconcilersData, err := r.getWoReconcilersData(ctx, argClient, armClient, arg.Spec.Subscription, reqLogger)
	if err != nil {
		return r.manageFailure(ctx, reqLogger, arg, err, "Failed to get WO Reconcilers Data")
	}
	// Append WO Reconcilers Data to GitOps Reconcilers Data
	reconcilersData = append(reconcilersData, woReconcilersData...)

	// Get list of all all reconcilers with the label set to the name of the AzureResourceGraph
	reconcilerList := &hubv1alpha1.ReconcilerList{}
	err = r.List(ctx, reconcilerList, client.MatchingLabels(map[string]string{"azure-resource-graph": arg.Name}))
	if err != nil {
		return r.manageFailure(ctx, reqLogger, arg, err, "Failed to get Reconcilers")
	}

	// Garbage collect reconcilers
	err = r.garbageCollectReconcilers(ctx, arg, reconcilerList, reconcilersData, reqLogger)
	if err != nil {
		return r.manageFailure(ctx, reqLogger, arg, err, "Failed to garbage collect Reconcilers")
	}

	// iterate over the list of reconcilers from the Azure Resource Graph and create/update them
	for _, argReconcilerData := range reconcilersData {

		// Create or update the reconciler
		reconciler, err := r.createOrUpdateReconciler(ctx, arg, reconcilerList, argReconcilerData)
		if err != nil {
			return r.manageFailure(ctx, reqLogger, arg, err, "Failed to create or update Reconciler")
		}
		// log the created or updated reconciler
		reqLogger.Info(fmt.Sprintf("Created or Updated Reconciler: " + fmt.Sprint(reconciler) + "\n"))
	}

	condition := metav1.Condition{
		Type:   "Ready",
		Status: metav1.ConditionTrue,
		Reason: "ArgReconciled",
	}
	meta.SetStatusCondition(&arg.Status.Conditions, condition)

	updateErr := r.Status().Update(ctx, arg)
	if updateErr != nil {
		reqLogger.Info("Error when updating status.")
		return ctrl.Result{RequeueAfter: time.Second * 3}, updateErr
	}

	return ctrl.Result{RequeueAfter: arg.Spec.Interval.Duration}, nil
}

// Garbage collect Arc reconcilers
func (r *AzureResourceGraphReconciler) garbageCollectReconcilers(ctx context.Context, arg *hubv1alpha1.AzureResourceGraph, reconcilerList *hubv1alpha1.ReconcilerList, reconcilersData []hubv1alpha1.ReconcilerSpec, logger logr.Logger) error {
	// iterate over the list of reconcilers and delete the ones that are not in the list of reconcilers from the Azure Resource Graph
	for _, reconciler := range reconcilerList.Items {
		if reconciler.Spec.Type == hubv1alpha1.ReconcilerTypeArc {
			found := false
			for _, argReconciler := range reconcilersData {
				if argReconciler.HostName == reconciler.Spec.HostName && argReconciler.ReconcilerName == reconciler.Spec.ReconcilerName {
					found = true
					break
				}
			}
			if !found {
				err := r.Delete(ctx, &reconciler)
				if err != nil {
					return err
				}
				// log the deleted reconciler
				logger.Info(fmt.Sprintf("Deleted Reconciler: " + fmt.Sprint(reconciler) + "\n"))
			}

		}

	}

	return nil
}

// Create or Update Reconciler
func (r *AzureResourceGraphReconciler) createOrUpdateReconciler(ctx context.Context, arg *hubv1alpha1.AzureResourceGraph, reconcilerList *hubv1alpha1.ReconcilerList, reconcilerData hubv1alpha1.ReconcilerSpec) (*hubv1alpha1.Reconciler, error) {
	reconcilerFullName := fmt.Sprint(reconcilerData.HostName + "-" + reconcilerData.ReconcilerName)
	found := false
	var reconciler *hubv1alpha1.Reconciler = nil

	for _, existingReconciler := range reconcilerList.Items {
		if existingReconciler.Name == reconcilerFullName {
			found = true
			reconciler = &existingReconciler
			break
		}
	}

	if !found {
		reconciler := &hubv1alpha1.Reconciler{
			ObjectMeta: metav1.ObjectMeta{
				Name:      reconcilerFullName,
				Namespace: arg.Namespace,
				Labels: map[string]string{
					"azure-resource-graph": arg.Name,
				},
			},
			Spec:   reconcilerData,
			Status: hubv1alpha1.ReconcilerStatus{},
		}

		// set the owner reference to the AzureResourceGraph
		err := ctrl.SetControllerReference(arg, reconciler, r.Scheme)
		if err != nil {
			return nil, err
		}

		err = r.Create(ctx, reconciler)
		if err != nil {
			return nil, err
		}
	} else {
		reconciler.Spec = reconcilerData
		err := r.Update(ctx, reconciler)
		if err != nil {
			return nil, err
		}
	}

	return reconciler, nil
}

// get a list of ReconcilerSpec from the Azure Resource Graph
// at this point only git/kustomize is supported. Helm is not supported.
func (r *AzureResourceGraphReconciler) getFluxReconcilersData(ctx context.Context, fluxConfigClient *armkubernetesconfiguration.FluxConfigurationsClient, fluxConfigs []interface{}, logger logr.Logger) ([]hubv1alpha1.ReconcilerSpec, error) {
	// Iterate over the results and create a slice of ReconcilerSpec
	var reconcilerData []hubv1alpha1.ReconcilerSpec

	storageClient, err := grpcClient.GetObservabilityStorageGrpcClient()
	if err != nil {
		return nil, err
	}

	for _, fluxConfig := range fluxConfigs {
		if fluxConfig != nil {

			fluxConfigMap := fluxConfig.(map[string]interface{})

			// Get Reconciler Name
			fluxConfigName := fluxConfigMap["name"].(string)

			// Get the resource group name
			resourceGroup := fluxConfigMap["resourceGroup"].(string)
			// Get the id
			id := fluxConfigMap["id"].(string)
			clusterName := strings.Split(id, "/")[8]

			propeties := fluxConfigMap["properties"].(map[string]interface{})

			gitRepository := propeties["gitRepository"].(map[string]interface{})
			// Get the repo
			repo := gitRepository["url"].(string)

			repositoryRef := gitRepository["repositoryRef"].(map[string]interface{})
			// Get the branch
			branch := repositoryRef["branch"].(string)

			kustomizations := propeties["kustomizations"].(map[string]interface{})
			path := ""
			// take first entry in kustomizations map
			// itereate over kustomizations map
			for _, kustomization := range kustomizations {
				// Get the path
				path = kustomization.(map[string]interface{})["path"].(string)
				break
			}
			// Get the endpoint as repo/branch/path
			endpoint := repo + "/" + branch + "/" + path

			sourceSyncedCommitId := propeties["sourceSyncedCommitId"]
			gitOpsCommitId := sourceSyncedCommitId.(string)

			statusMessage := ""
			sourceComplianceState := propeties["complianceState"]
			if sourceComplianceState != nil {
				statusMessage = r.getStatusMessage(sourceComplianceState.(string), propeties["statuses"].([]interface{}))
			}

			reconciler := r.createReconciler(sourceComplianceState.(string), statusMessage, gitOpsCommitId,
				resourceGroup, clusterName, fluxConfigName, endpoint, hubv1alpha1.ReconcilerTypeArc)

			reconcilerData = append(reconcilerData, reconciler)

			cheildReconcilerData, err := r.getReconcilersDataFromChildKalypsoObjects(ctx, storageClient,
				resourceGroup, clusterName, fluxConfigName, fluxConfigClient, fluxConfigs, logger)
			if err != nil {
				return nil, err
			}
			reconcilerData = append(reconcilerData, cheildReconcilerData...)

		}
	}
	return reconcilerData, nil

}
func (r *AzureResourceGraphReconciler) createReconciler(status string, statusMessage string, gitOpsCommitId string,
	resourceGroup string, clusterName string, reconcilerName string, manifestsEndpoint string, reconcilerType hubv1alpha1.ReconcilerType) hubv1alpha1.ReconcilerSpec {
	reconcilerStatus := r.translateComplianceState(status)
	clusterName = strings.ToLower(clusterName)

	deployment := hubv1alpha1.Deployment{
		GitOpsCommitId: gitOpsCommitId,
		Status:         reconcilerStatus,
		StatusMessage:  statusMessage,
	}

	// Create the reconciler spec
	reconciler := hubv1alpha1.ReconcilerSpec{
		HostName:             clusterName,
		ReconcilerName:       reconcilerName,
		Type:                 reconcilerType,
		ManifestsStorageType: hubv1alpha1.Git,
		ManifestsEndpoint:    manifestsEndpoint,
		Deployment:           deployment,
	}

	return reconciler

}

func (r *AzureResourceGraphReconciler) getReconcilersDataFromChildKalypsoObjects(ctx context.Context, storageClient grpcClient.ObservabilityStorageGrpcClient, resourceGroup string, clusterName string, fluxConfigName string, fluxConfigClient *armkubernetesconfiguration.FluxConfigurationsClient, fluxConfigs []interface{}, logger logr.Logger) ([]hubv1alpha1.ReconcilerSpec, error) {
	var reconcilerData []hubv1alpha1.ReconcilerSpec

	res, err := fluxConfigClient.Get(ctx, resourceGroup, "Microsoft.Kubernetes", "connectedClusters", clusterName, fluxConfigName, nil)
	if err != nil {
		res, err = fluxConfigClient.Get(ctx, resourceGroup, "Microsoft.ContainerService", "managedClusters", clusterName, fluxConfigName, nil)
		if err != nil {
			return nil, err
		}
	}

	fluxConfigurationDetal := res.FluxConfiguration
	// iteretae over the statuses and log them
	for _, status := range fluxConfigurationDetal.Properties.Statuses {
		if *status.Kind != "Kustomization" {
			continue
		}

		//TODO Update Kalypso: name deployment target as workload.deploymentTarget or without workload at all
		// expected flux resource name format: env.workspace.application.workload.deploymentTarget[.clusterType]
		nameParts := strings.Split(*status.Name, ".")
		if len(nameParts) < 5 {
			continue
		}
		environmentName := nameParts[0]
		workspace := nameParts[1]
		application := nameParts[2]
		workloadName := nameParts[3]
		deploymentTargetName := nameParts[4]

		dt, err := storageClient.GetDeploymentTarget(ctx, &pb.DeploymentTargetSearch{
			WorkloadName:         workloadName,
			DeploymentTargetName: deploymentTargetName,
			EnvironmentName:      environmentName,
			WorkspaceName:        workspace,
			ApplicationName:      application,
		})
		if err != nil {
			//log  workspace, application, workloadName and deploymentTargetName
			logger.Info("Could not find deployment target", "workspace", workspace, "application", application, "workloadName", workloadName, "deploymentTargetName", deploymentTargetName)
			continue
		}

		manifestsEndpoint := dt.ManifestsEndpoint
		statusMessage := ""
		gitOpsCommitId := ""
		for _, statusCondition := range status.StatusConditions {
			if statusCondition.Message != nil {
				statusConditionMessage := *statusCondition.Message
				shaIndex := strings.Index(statusConditionMessage, "sha1:")
				if shaIndex > 0 {
					gitOpsCommitId = statusConditionMessage[shaIndex : shaIndex+45]
				}
				statusMessage += statusConditionMessage
			}

		}

		reconciler := r.createReconciler(string(*status.ComplianceState), statusMessage, gitOpsCommitId,
			resourceGroup, clusterName, *status.Name, manifestsEndpoint, hubv1alpha1.ReconcilerTypeFlux)
		reconcilerData = append(reconcilerData, reconciler)

	}
	return reconcilerData, nil
}

// Translate Compliance State
func (r *AzureResourceGraphReconciler) translateComplianceState(complianceState string) hubv1alpha1.DeploymentStatusType {

	switch complianceState {
	case "Compliant":
		return hubv1alpha1.DeploymentStatusSuccess
	case "Succeeded":
		return hubv1alpha1.DeploymentStatusSuccess
	case "Non-Compliant":
		return hubv1alpha1.DeploymentStatusFailed
	case "Noncompliant":
		return hubv1alpha1.DeploymentStatusFailed
	case "Failed":
		return hubv1alpha1.DeploymentStatusFailed
	default:
		return hubv1alpha1.DeploymentStatusPending
	}

}

// Get Status Message
func (r *AzureResourceGraphReconciler) getStatusMessage(complianceState string, statuses []interface{}) string {
	// iterate over statuses and concatenate message for the records with complianceState
	var statusMessage string
	for _, status := range statuses {
		statusMap := status.(map[string]interface{})
		if statusMap["complianceState"] == complianceState {
			statusMessage += statusMap["message"].(string) + "\n"
		}
	}
	return statusMessage
}

// Get Acxure Credentials
func (r *AzureResourceGraphReconciler) getAzureCredentials(arg *hubv1alpha1.AzureResourceGraph) (*azidentity.DefaultAzureCredential, error) {
	// find secret by name
	secret := &corev1.Secret{}
	err := r.Get(context.TODO(), types.NamespacedName{Name: arg.Spec.SecretRef, Namespace: arg.Namespace}, secret)
	if err != nil {
		return nil, err
	}

	//set the environment variables from the secret
	if err := os.Setenv("AZURE_TENANT_ID", arg.Spec.Tenant); err != nil {
		return nil, err
	}
	if err := os.Setenv("AZURE_SUBSCRIPTION_ID", arg.Spec.Subscription); err != nil {
		return nil, err
	}
	if err := os.Setenv("AZURE_CLIENT_SECRET", string(secret.Data["AZURE_CLIENT_SECRET"])); err != nil {
		return nil, err
	}
	if err := os.Setenv("AZURE_CLIENT_ID", string(secret.Data["AZURE_CLIENT_ID"])); err != nil {
		return nil, err
	}

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return nil, err
	}
	return cred, nil
}

// Get ARG client
func (r *AzureResourceGraphReconciler) getARGClient(cred *azidentity.DefaultAzureCredential) (*armresourcegraph.Client, error) {

	clientFactory, err := armresourcegraph.NewClientFactory(cred, nil)
	if err != nil {
		return nil, err
	}
	return clientFactory.NewClient(), nil
}

// Get ARM client
func (r *AzureResourceGraphReconciler) getARMClient(cred *azidentity.DefaultAzureCredential, subscriptionId string) (*armresources.Client, error) {
	clientFactory, err := armresources.NewClientFactory(subscriptionId, cred, nil)
	if err != nil {
		return nil, err
	}
	return clientFactory.NewClient(), nil
}

// Get Flux Config Client
func (r *AzureResourceGraphReconciler) getFluxConfigClient(cred *azidentity.DefaultAzureCredential, subscriptionId string) (*armkubernetesconfiguration.FluxConfigurationsClient, error) {
	clientFactory, err := armkubernetesconfiguration.NewClientFactory(subscriptionId, cred, nil)
	if err != nil {
		return nil, err
	}
	return clientFactory.NewFluxConfigurationsClient(), nil

}

func (r *AzureResourceGraphReconciler) getWoTargets(ctx context.Context, argClient *armresourcegraph.Client, subscription string) (armresourcegraph.ClientResourcesResponse, error) {
	query := "where type == 'microsoft.edge/targets'"
	resultFormatObjectArray := armresourcegraph.ResultFormatObjectArray
	return argClient.Resources(ctx, armresourcegraph.QueryRequest{
		Subscriptions: []*string{&subscription},
		Query:         &query,
		Options:       &armresourcegraph.QueryRequestOptions{ResultFormat: &resultFormatObjectArray},
	}, nil)
}

// Get Flux Configurations
func (r *AzureResourceGraphReconciler) getFluxConfigurations(ctx context.Context, argClient *armresourcegraph.Client, subscription string) (armresourcegraph.ClientResourcesResponse, error) {
	query := "kubernetesconfigurationresources | where type == 'microsoft.kubernetesconfiguration/fluxconfigurations'"
	resultFormatObjectArray := armresourcegraph.ResultFormatObjectArray

	// Set options
	RequestOptions := armresourcegraph.QueryRequestOptions{
		ResultFormat: &resultFormatObjectArray,
	}

	// Create the query request
	Request := armresourcegraph.QueryRequest{
		Subscriptions: []*string{&subscription},
		Query:         &query,
		Options:       &RequestOptions,
	}

	return argClient.Resources(ctx, Request, nil)
}

// Get Reconcilers Data from Workload Orchetration service
// Since there is not API to get all deployed instances we have to iterate over all deployment descriptors and targets and check if the instance exist one-by-one.
// This is not the most efficient way, but it is the only way until the Workload Orchestration service provides an API to get all instances.
func (r *AzureResourceGraphReconciler) getWoReconcilersData(ctx context.Context, argClient *armresourcegraph.Client, armClient *armresources.Client, subscription string, logger logr.Logger) ([]hubv1alpha1.ReconcilerSpec, error) {
	reconcilerData := []hubv1alpha1.ReconcilerSpec{}
	// get a list of deployment descriptors
	deploymentDescriptors := &hubv1alpha1.DeploymentDescriptorList{}
	err := r.List(ctx, deploymentDescriptors)
	if err != nil {
		return nil, fmt.Errorf("failed to list deployment descriptors: %w", err)
	}

	//iterate over a list of WO targets
	woTargets, err := r.getWoTargets(ctx, argClient, subscription)
	if err != nil {
		return nil, fmt.Errorf("failed to get WO targets: %w", err)
	}

	// Iterate over a list of deployment descriptors (hubv1alpha1.DeploymentDescriptor)
	for _, deploymentDescriptor := range deploymentDescriptors.Items {
		deploymentTargetName := deploymentDescriptor.Spec.DeploymentTarget.Name

		gitOpsCommitId := deploymentDescriptor.Status.GitOpsCommitId
		descriptorDeploymentTarget := deploymentDescriptor.Spec.DeploymentTarget
		endpoint := fmt.Sprintf("%s/%s/%s", descriptorDeploymentTarget.Manifests.Repo, descriptorDeploymentTarget.Manifests.Branch, descriptorDeploymentTarget.Manifests.Path)
		workloadVsersion := deploymentDescriptor.Spec.WorkloadVersion.Version

		for _, woTarget := range woTargets.Data.([]interface{}) {
			woTargetMap := woTarget.(map[string]interface{})
			//get targetName
			targetName := woTargetMap["name"].(string)
			// Get custom location as extendedLocation.name
			// extendedLocation is a map with a key "name"
			extendedLocation := woTargetMap["extendedLocation"].(map[string]interface{})
			extendedLocationName := extendedLocation["name"].(string)
			// extract the last section starting with "/" as a locationName
			extendedLocationNameSections := strings.Split(extendedLocationName, "/")
			locationName := ""
			if len(extendedLocationNameSections) > 0 {
				locationName = extendedLocationNameSections[len(extendedLocationNameSections)-1]
			}

			instanceId := fmt.Sprintf("%v/solutions/%v/instances/%v", woTargetMap["id"], deploymentTargetName, deploymentTargetName)
			//get the instance with getAzureResourceById
			instance, err := r.getAzureResourceById(ctx, armClient, instanceId)
			if err != nil {
				// check if the error is not found
				if !strings.Contains(err.Error(), "404") {
					logger.Info("Instance not found for deployment target", "deploymentTargetName", deploymentTargetName, "instanceId", instanceId)
				}
				continue
			}

			logger.Info("Found an instance for deployment target", "deploymentTargetName", deploymentTargetName, "instanceId", instanceId)
			//Check the status as properties.status.status
			properties := instance.Properties.(map[string]interface{})
			solutionVersionId := properties["solutionVersionId"].(string)
			//extract vsersion from solutionVersionId string like "something-0.0.126.2 -> 0.0.126.2
			solutionVersionParts := strings.Split(solutionVersionId, "-")
			solutionVersion := ""
			if len(solutionVersionParts) > 1 {
				solutionVersion = solutionVersionParts[len(solutionVersionParts)-1]
			}
			// remove last section of solutionVersion if it contains a dot 0.0.126.2 -> 0.0.126
			if strings.Contains(solutionVersion, ".") {
				solutionVersionParts = strings.Split(solutionVersion, ".")
				if len(solutionVersionParts) > 1 {
					solutionVersion = strings.Join(solutionVersionParts[:len(solutionVersionParts)-1], ".")
				}
			}

			if workloadVsersion != solutionVersion {
				logger.Info("Skipping reconciler creation for deployment descriptor", "deploymentDescriptorName", deploymentDescriptor.Name,
					"deploymentTargetName", deploymentTargetName, "solutionVersion", solutionVersion,
					"workloadVsersion", workloadVsersion, "targetName", targetName)
				continue
			}

			statusProperties := properties["status"].(map[string]interface{})
			status := statusProperties["status"].(string)

			// Handle targetStatuses which might be a slice or string
			var statusMessage string
			if targetStatuses, exists := statusProperties["targetStatuses"]; exists {
				statusMessage = fmt.Sprintf("%v", targetStatuses)
			}

			resourceGroup := woTargetMap["resourceGroup"].(string)

			// create or update reconciler
			reconciler := r.createReconciler(status, statusMessage, gitOpsCommitId,
				resourceGroup, locationName, targetName, endpoint, hubv1alpha1.ReconcilerTypeWo)

			reconcilerData = append(reconcilerData, reconciler)
		}
	}
	return reconcilerData, nil
}

// get azure resource by id using armresources.Client
// This function is not used in the current implementation, but it can be used to get Azure
// resources by their ID if needed in the future.
func (r *AzureResourceGraphReconciler) getAzureResourceById(ctx context.Context, armClient *armresources.Client, resourceId string) (*armresources.GenericResource, error) {
	resource, err := armClient.GetByID(ctx, resourceId, "2025-01-01-preview", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get Azure resource by ID: %w", err)
	}
	return &resource.GenericResource, nil
}

func (r *AzureResourceGraphReconciler) getGitOpsReconcilersData(ctx context.Context, cred *azidentity.DefaultAzureCredential, argClient *armresourcegraph.Client, subscription string, logger logr.Logger) ([]hubv1alpha1.ReconcilerSpec, error) {
	// Get Flux Configurations
	fluxConfigs, err := r.getFluxConfigurations(ctx, argClient, subscription)
	if err != nil {
		return nil, err
	}

	// Get Flux Config Client
	fluxConfigClient, err := r.getFluxConfigClient(cred, subscription)
	if err != nil {
		return nil, err
	}

	// Get Reconcilers Data
	reconcilersData, err := r.getFluxReconcilersData(ctx, fluxConfigClient, fluxConfigs.Data.([]interface{}), logger)
	if err != nil {
		return nil, err
	}

	return reconcilersData, nil
}

// Gracefully handle errors
func (r *AzureResourceGraphReconciler) manageFailure(ctx context.Context, logger logr.Logger, arg *hubv1alpha1.AzureResourceGraph, err error, message string) (ctrl.Result, error) {
	logger.Error(err, message)

	//crerate a condition
	condition := metav1.Condition{
		Type:    "Ready",
		Status:  metav1.ConditionFalse,
		Reason:  "UpdateFailed",
		Message: err.Error(),
	}

	meta.SetStatusCondition(&arg.Status.Conditions, condition)

	updateErr := r.Status().Update(ctx, arg)
	if updateErr != nil {
		logger.Info("Error when updating status. Requeued")
		return ctrl.Result{RequeueAfter: time.Second * 3}, updateErr
	}
	return ctrl.Result{}, err
}

// SetupWithManager sets up the controller with the Manager.
func (r *AzureResourceGraphReconciler) SetupWithManager(mgr ctrl.Manager) error {

	return ctrl.NewControllerManagedBy(mgr).
		For(&hubv1alpha1.AzureResourceGraph{}, builder.WithPredicates(predicate.GenerationChangedPredicate{})).
		Owns(&hubv1alpha1.Reconciler{}, builder.WithPredicates(predicate.GenerationChangedPredicate{})).
		Complete(r)
}

//TODO:
// perhaps handle only GitOPs extensions in ARG that correlate with deployment descriptors and ignore the rest
