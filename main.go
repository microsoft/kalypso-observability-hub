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

package main

import (
	"flag"
	"os"
	"strconv"
	"sync"

	// Import all Kubernetes client auth plugins (e.g. Azure, GCP, OIDC, etc.)
	// to ensure that exec-entrypoint and run can make use of them.
	_ "k8s.io/client-go/plugin/pkg/client/auth"

	kustomizev1 "github.com/fluxcd/kustomize-controller/api/v1beta2"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	httpapi "github.com/microsoft/kalypso-observability-hub/api/http"
	hubv1alpha1 "github.com/microsoft/kalypso-observability-hub/api/v1alpha1"
	"github.com/microsoft/kalypso-observability-hub/controllers"
	db "github.com/microsoft/kalypso-observability-hub/storage/postgres"
	//+kubebuilder:scaffold:imports
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))

	utilruntime.Must(hubv1alpha1.AddToScheme(scheme))
	utilruntime.Must(kustomizev1.AddToScheme(scheme))
	//+kubebuilder:scaffold:scheme
}

func main() {
	var metricsAddr string
	var enableLeaderElection bool
	var probeAddr string
	var httpAPIAddr string
	var enableHTTPAPI bool
	var postgresHost string
	var postgresPort string
	var postgresUser string
	var postgresPassword string
	var postgresDBName string
	var postgresSSLMode string

	flag.StringVar(&metricsAddr, "metrics-bind-address", ":8080", "The address the metric endpoint binds to.")
	flag.StringVar(&probeAddr, "health-probe-bind-address", ":8081", "The address the probe endpoint binds to.")
	flag.StringVar(&httpAPIAddr, "http-api-bind-address", ":8082", "The address the HTTP API endpoint binds to.")
	flag.BoolVar(&enableLeaderElection, "leader-elect", false,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")
	flag.BoolVar(&enableHTTPAPI, "enable-http-api", false, "Enable HTTP API server for external queries.")
	flag.StringVar(&postgresHost, "postgres-host", "localhost", "PostgreSQL host for HTTP API.")
	flag.StringVar(&postgresPort, "postgres-port", "5432", "PostgreSQL port for HTTP API.")
	flag.StringVar(&postgresUser, "postgres-user", "postgres", "PostgreSQL user for HTTP API.")
	flag.StringVar(&postgresPassword, "postgres-password", "", "PostgreSQL password for HTTP API.")
	flag.StringVar(&postgresDBName, "postgres-dbname", "postgres", "PostgreSQL database name for HTTP API.")
	flag.StringVar(&postgresSSLMode, "postgres-sslmode", "disable", "PostgreSQL SSL mode for HTTP API.")

	opts := zap.Options{
		Development: true,
	}
	opts.BindFlags(flag.CommandLine)
	flag.Parse()

	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts)))

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:                 scheme,
		MetricsBindAddress:     metricsAddr,
		Port:                   9443,
		HealthProbeBindAddress: probeAddr,
		LeaderElection:         enableLeaderElection,
		LeaderElectionID:       "a0645f20.kalypso.io",
		// LeaderElectionReleaseOnCancel defines if the leader should step down voluntarily
		// when the Manager ends. This requires the binary to immediately end when the
		// Manager is stopped, otherwise, this setting is unsafe. Setting this significantly
		// speeds up voluntary leader transitions as the new leader don't have to wait
		// LeaseDuration time first.
		//
		// In the default scaffold provided, the program ends immediately after
		// the manager stops, so would be fine to enable this option. However,
		// if you are doing or is intended to do any operation such as perform cleanups
		// after the manager stops then its usage might be unsafe.
		// LeaderElectionReleaseOnCancel: true,
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	if err = (&controllers.DeploymentDescriptorReconciler{
		Client: mgr.GetClient(),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "DeploymentDescriptor")
		os.Exit(1)
	}
	if err = (&controllers.AzureResourceGraphReconciler{
		Client: mgr.GetClient(),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "AzureResourceGraph")
		os.Exit(1)
	}
	if err = (&controllers.ReconcilerReconciler{
		Client: mgr.GetClient(),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "Reconciler")
		os.Exit(1)
	}
	//+kubebuilder:scaffold:builder

	if err := mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up health check")
		os.Exit(1)
	}
	if err := mgr.AddReadyzCheck("readyz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up ready check")
		os.Exit(1)
	}

	// Start HTTP API server if enabled
	if enableHTTPAPI {
		postgresPortInt, err := strconv.Atoi(postgresPort)
		if err != nil {
			setupLog.Error(err, "invalid postgres port")
			os.Exit(1)
		}

		httpAPIPortStr := httpAPIAddr[1:] // Remove the : prefix
		httpAPIPortInt, err := strconv.Atoi(httpAPIPortStr)
		if err != nil {
			setupLog.Error(err, "invalid HTTP API port")
			os.Exit(1)
		}

		dbClient := db.NewPostgresClient(postgresHost, postgresPortInt, postgresUser, postgresPassword, postgresDBName, postgresSSLMode)
		httpAPIServer := httpapi.NewHTTPAPIServer(dbClient, httpAPIPortInt)

		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()
			setupLog.Info("starting HTTP API server", "port", httpAPIPortInt)
			if err := httpAPIServer.Start(); err != nil {
				setupLog.Error(err, "problem running HTTP API server")
			}
		}()
	}

	setupLog.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}
