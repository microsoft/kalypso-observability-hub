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
	"time"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	kustomizev1 "github.com/fluxcd/kustomize-controller/api/v1beta2"
	"github.com/microsoft/kalypso-observability-hub/api/v1alpha1"
	hubv1alpha1 "github.com/microsoft/kalypso-observability-hub/api/v1alpha1"
	grpcClient "github.com/microsoft/kalypso-observability-hub/storage/api/grpc/client"
	pb "github.com/microsoft/kalypso-observability-hub/storage/api/grpc/proto"
)

const (
	FluxKustomizationNameLabel      = "kustomize.toolkit.fluxcd.io/name"
	FluxKustomizationNamespaceLabel = "kustomize.toolkit.fluxcd.io/namespace"
)

// DeploymentDescriptorReconciler reconciles a DeploymentDescriptor object
type DeploymentDescriptorReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=hub.kalypso.io,resources=deploymentdescriptors,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=hub.kalypso.io,resources=deploymentdescriptors/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=hub.kalypso.io,resources=deploymentdescriptors/finalizers,verbs=update
//+kubebuilder:rbac:groups=kustomize.toolkit.fluxcd.io,resources=kustomizations,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the DeploymentDescriptor object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.13.1/pkg/reconcile
func (r *DeploymentDescriptorReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	reqLogger := log.FromContext(ctx)
	reqLogger.Info("=== Reconciling Deployment Descriptor  ===")

	// Fetch the DeploymentDescriptor instance
	deploymentDescriptor := &hubv1alpha1.DeploymentDescriptor{}
	err := r.Get(ctx, req.NamespacedName, deploymentDescriptor)
	if err != nil {
		ignroredNotFound := client.IgnoreNotFound(err)
		if ignroredNotFound != nil {
			reqLogger.Error(err, "Failed to get Deployment Descriptor")
		}
		return ctrl.Result{}, ignroredNotFound
	}

	// Check if the resource is being deleted
	if !deploymentDescriptor.ObjectMeta.DeletionTimestamp.IsZero() {
		return ctrl.Result{}, nil
	}

	storageClient, err := grpcClient.GetObservabilityStorageGrpcClient()
	if err != nil {
		return r.manageFailure(ctx, reqLogger, deploymentDescriptor, err, "Failed to get storage client")
	}

	descriptorWorkload := deploymentDescriptor.Spec.Workload
	descriptorApplication := descriptorWorkload.Application
	descriptorWorkspace := descriptorApplication.Workspace

	// Update Workspace
	ws, err := storageClient.UpdateWorkspace(ctx, &pb.Workspace{
		Name:        descriptorWorkspace.Name,
		Description: descriptorWorkspace.Name,
	})
	if err != nil {
		return r.manageFailure(ctx, reqLogger, deploymentDescriptor, err, "Failed to update workspace")
	}

	// Update Application
	app, err := storageClient.UpdateApplication(ctx, &pb.Application{
		Name:        descriptorApplication.Name,
		Description: descriptorApplication.Name,
		WorkspaceId: ws.Id,
	})
	if err != nil {
		return r.manageFailure(ctx, reqLogger, deploymentDescriptor, err, "Failed to update application")
	}

	// Update Workload
	wkl, err := storageClient.UpdateWorkload(ctx, &pb.Workload{
		Name:              descriptorWorkload.Name,
		Description:       descriptorWorkload.Name,
		SourceStorageType: v1alpha1.GitStorageType,
		SourceEndpoint:    fmt.Sprintf("%s/%s/%s", descriptorWorkload.Source.Repo, descriptorWorkload.Source.Branch, descriptorWorkload.Source.Path),
		ApplicationId:     app.Id,
	})
	if err != nil {
		return r.manageFailure(ctx, reqLogger, deploymentDescriptor, err, "Failed to update workload")
	}

	descriptorDeploymentTarget := deploymentDescriptor.Spec.DeploymentTarget
	descriptorEnvironment := descriptorDeploymentTarget.Environment

	// Update Environment
	env, err := storageClient.UpdateEnvironment(ctx, &pb.Environment{
		Name:        descriptorEnvironment,
		Description: descriptorEnvironment,
	})
	if err != nil {
		return r.manageFailure(ctx, reqLogger, deploymentDescriptor, err, "Failed to update environment")
	}

	// Update Deployment Target
	dt, err := storageClient.UpdateDeploymentTarget(ctx, &pb.DeploymentTarget{
		Name:                 descriptorDeploymentTarget.Name,
		Description:          descriptorDeploymentTarget.Name,
		EnvironmentId:        env.Id,
		WorkloadId:           wkl.Id,
		ManifestsStorageType: v1alpha1.GitStorageType,
		ManifestsEndpoint:    fmt.Sprintf("%s/%s/%s", descriptorDeploymentTarget.Manifests.Repo, descriptorDeploymentTarget.Manifests.Branch, descriptorDeploymentTarget.Manifests.Path),
	})
	if err != nil {
		return r.manageFailure(ctx, reqLogger, deploymentDescriptor, err, "Failed to update deployment target")
	}

	// Update Workload Version
	descriptorWorkloadVersion := deploymentDescriptor.Spec.WorkloadVersion
	wv, err := storageClient.UpdateWorkloadVersion(ctx, &pb.WorkloadVersion{
		Version:       descriptorWorkloadVersion.Version,
		WorkloadId:    wkl.Id,
		BuildId:       descriptorWorkloadVersion.Build,
		BuildCommitId: descriptorWorkloadVersion.Commit,
	})
	if err != nil {
		return r.manageFailure(ctx, reqLogger, deploymentDescriptor, err, "Failed to update workload version")
	}

	commit_id, err := r.getCommitFromFluxKustomization(deploymentDescriptor)
	if err != nil {
		return r.manageFailure(ctx, reqLogger, deploymentDescriptor, err, "Failed to get commit from flux kustomization")
	}

	// Update Deployment Assignment
	_, err = storageClient.UpdateDeploymentAssignment(ctx, &pb.DeploymentAssignment{
		DeploymentTargetId: dt.Id,
		WorkloadVersionId:  wv.Id,
		GitopsCommitId:     *commit_id,
	})
	if err != nil {
		return r.manageFailure(ctx, reqLogger, deploymentDescriptor, err, "Failed to update deployment assignment")
	}

	condition := metav1.Condition{
		Type:   "Ready",
		Status: metav1.ConditionTrue,
		Reason: "DeploymentDescriptorReconciled",
	}
	meta.SetStatusCondition(&deploymentDescriptor.Status.Conditions, condition)

	updateErr := r.Status().Update(ctx, deploymentDescriptor)
	if updateErr != nil {
		reqLogger.Info("Error when updating status.")
		return ctrl.Result{RequeueAfter: time.Second * 3}, updateErr
	}

	return ctrl.Result{}, nil
}

// Get Commit from Flux Kustomization
func (r *DeploymentDescriptorReconciler) getCommitFromFluxKustomization(deploymentDescriptor *hubv1alpha1.DeploymentDescriptor) (*string, error) {
	//get the flux kustomization name from the deployment descriptor
	fluxKustomizationName := deploymentDescriptor.Labels[FluxKustomizationNameLabel]
	if fluxKustomizationName == "" {
		return nil, fmt.Errorf("Flux Kustomization name not found in the deployment descriptor")
	}

	fluxKustomizationNamespace := deploymentDescriptor.Labels[FluxKustomizationNamespaceLabel]
	if fluxKustomizationNamespace == "" {
		return nil, fmt.Errorf("Flux Kustomization namespace not found in the deployment descriptor")
	}

	//get the flux kustomization
	fluxKustomization := &kustomizev1.Kustomization{}
	err := r.Get(context.Background(), client.ObjectKey{
		Name:      fluxKustomizationName,
		Namespace: fluxKustomizationNamespace,
	}, fluxKustomization)
	if err != nil {
		return nil, err
	}

	//get the commit from the flux kustomization
	commit := fluxKustomization.Status.LastAppliedRevision

	return &commit, nil

}

// Gracefully handle errors
func (r *DeploymentDescriptorReconciler) manageFailure(ctx context.Context, logger logr.Logger, deploymentDescriptor *hubv1alpha1.DeploymentDescriptor, err error, message string) (ctrl.Result, error) {
	logger.Error(err, message)

	//crerate a condition
	condition := metav1.Condition{
		Type:    "Ready",
		Status:  metav1.ConditionFalse,
		Reason:  "UpdateFailed",
		Message: err.Error(),
	}

	meta.SetStatusCondition(&deploymentDescriptor.Status.Conditions, condition)

	updateErr := r.Status().Update(ctx, deploymentDescriptor)
	if updateErr != nil {
		logger.Info("Error when updating status. Requeued")
		return ctrl.Result{RequeueAfter: time.Second * 3}, updateErr
	}
	return ctrl.Result{}, err
}

// SetupWithManager sets up the controller with the Manager.
func (r *DeploymentDescriptorReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&hubv1alpha1.DeploymentDescriptor{}).
		WithEventFilter(predicate.Or(predicate.GenerationChangedPredicate{}, predicate.LabelChangedPredicate{})).
		Complete(r)
}

// TODO: condider other storages rather than git
// Think on : Commi-Id-> hash
// DeaploymentDescriptor has an optional hash in the manifeats section (if not specified take it from Flux like it is now)
// Simplify repo/branch/path to just endpoint
