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
	"time"

	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	"github.com/go-logr/logr"
	hubv1alpha1 "github.com/microsoft/kalypso-observability-hub/api/v1alpha1"
	grpcClient "github.com/microsoft/kalypso-observability-hub/storage/api/grpc/client"
	pb "github.com/microsoft/kalypso-observability-hub/storage/api/grpc/proto"
)

// ReconcilerReconciler reconciles a Reconciler object
type ReconcilerReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=hub.kalypso.io,resources=reconcilers,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=hub.kalypso.io,resources=reconcilers/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=hub.kalypso.io,resources=reconcilers/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Reconciler object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.13.1/pkg/reconcile
func (r *ReconcilerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	reqLogger := log.FromContext(ctx)
	reqLogger.Info("=== Reconciling Reconciler  ===")

	// Fetch the Reconciler instance
	reconciler := &hubv1alpha1.Reconciler{}
	err := r.Get(ctx, req.NamespacedName, reconciler)
	if err != nil {
		ignroredNotFound := client.IgnoreNotFound(err)
		if ignroredNotFound != nil {
			reqLogger.Error(err, "Failed to get Azure Resource Graph")
		}
		return ctrl.Result{}, ignroredNotFound
	}

	// Check if the resource is being deleted
	if !reconciler.ObjectMeta.DeletionTimestamp.IsZero() {
		return ctrl.Result{}, nil
	}

	storageClient, err := grpcClient.GetObservabilityStorageGrpcClient()
	if err != nil {
		return r.manageFailure(ctx, reqLogger, reconciler, err, "Failed to get storage client")
	}

	// Update Host
	host, err := storageClient.UpdateHost(ctx, &pb.Host{
		Name:        reconciler.Spec.HostName,
		Description: reconciler.Spec.HostName,
	})
	if err != nil {
		return r.manageFailure(ctx, reqLogger, reconciler, err, "Failed to update host")
	}

	// Update Reconciler
	rc, err := storageClient.UpdateReconciler(ctx, &pb.Reconciler{
		HostId:               host.Id,
		Name:                 reconciler.Spec.ReconcilerName,
		Description:          reconciler.Spec.ReconcilerName,
		ReconcilerType:       reconciler.Spec.Type,
		ManifestsStorageType: (string)(reconciler.Spec.ManifestsStorageType),
		ManifestsEndpoint:    reconciler.Spec.ManifestsEndpoint,
	})
	if err != nil {
		return r.manageFailure(ctx, reqLogger, reconciler, err, "Failed to update reconciler")
	}

	// Update Deployment
	reconcilerDeployment := reconciler.Spec.Deployment
	_, err = storageClient.UpdateDeployment(ctx, &pb.Deployment{
		ReconcilerId:   rc.Id,
		GitopsCommitId: reconcilerDeployment.GitOpsCommitId,
		Status:         (string)(reconcilerDeployment.Status),
		StatusMessage:  reconcilerDeployment.StatusMessage,
	})
	if err != nil {
		return r.manageFailure(ctx, reqLogger, reconciler, err, "Failed to update deployment")
	}

	condition := metav1.Condition{
		Type:   "Ready",
		Status: metav1.ConditionTrue,
		Reason: "Succeeded",
	}
	meta.SetStatusCondition(&reconciler.Status.Conditions, condition)

	updateErr := r.Status().Update(ctx, reconciler)
	if updateErr != nil {
		reqLogger.Info("Error when updating status.")
		return ctrl.Result{RequeueAfter: time.Second * 3}, updateErr
	}

	return ctrl.Result{}, nil
}

// Gracefully handle errors
func (r *ReconcilerReconciler) manageFailure(ctx context.Context, logger logr.Logger, reconciler *hubv1alpha1.Reconciler, err error, message string) (ctrl.Result, error) {
	logger.Error(err, message)

	//crerate a condition
	condition := metav1.Condition{
		Type:    "Ready",
		Status:  metav1.ConditionFalse,
		Reason:  "UpdateFailed",
		Message: err.Error(),
	}

	meta.SetStatusCondition(&reconciler.Status.Conditions, condition)

	updateErr := r.Status().Update(ctx, reconciler)
	if updateErr != nil {
		logger.Info("Error when updating status. Requeued")
		return ctrl.Result{RequeueAfter: time.Second * 3}, updateErr
	}
	return ctrl.Result{}, err
}

// SetupWithManager sets up the controller with the Manager.
func (r *ReconcilerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&hubv1alpha1.Reconciler{}, builder.WithPredicates(predicate.GenerationChangedPredicate{})).
		Complete(r)
}
