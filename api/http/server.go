package http

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	db "github.com/microsoft/kalypso-observability-hub/storage/postgres"
)

type HTTPAPIServer struct {
	dbClient db.DBClient
	server   *http.Server
	port     int
}

type DeploymentStateResponse struct {
	TotalSubscribers           int32                `json:"total_subscribers"`
	TotalSucceededSubscribers  int32                `json:"total_succeeded_subscribers"`
	TotalFailedSubscribers     int32                `json:"total_failed_subscribers"`
	TotalInProgressSubscribers int32                `json:"total_in_progress_subscribers"`
	SucceededSubscribers       []SubscriberResponse `json:"succeeded_subscribers"`
	FailedSubscribers          []SubscriberResponse `json:"failed_subscribers"`
	InProgressSubscribers      []SubscriberResponse `json:"in_progress_subscribers"`
}

type SubscriberResponse struct {
	Name          string `json:"name"`
	StatusMessage string `json:"status_message"`
}

type DeploymentResponse struct {
	ID             int    `json:"id"`
	GitopsCommitID string `json:"gitops_commit_id"`
	ReconcilerID   int    `json:"reconciler_id"`
	Status         string `json:"status"`
	StatusMessage  string `json:"status_message"`
}

type EnvironmentResponse struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type DeploymentTargetResponse struct {
	ID                   int    `json:"id"`
	Name                 string `json:"name"`
	Description          string `json:"description"`
	WorkloadID           int    `json:"workload_id"`
	EnvironmentID        int    `json:"environment_id"`
	Labels               string `json:"labels"`
	ManifestsStorageType string `json:"manifests_storage_type"`
	ManifestsEndpoint    string `json:"manifests_endpoint"`
}

func NewHTTPAPIServer(dbClient db.DBClient, port int) *HTTPAPIServer {
	return &HTTPAPIServer{
		dbClient: dbClient,
		port:     port,
	}
}

func (s *HTTPAPIServer) Start() error {
	router := mux.NewRouter()

	// API routes
	api := router.PathPrefix("/api/v1").Subrouter()

	// Health check
	router.HandleFunc("/health", s.healthHandler).Methods("GET")

	// Deployment state endpoint
	api.HandleFunc("/deployment-state", s.getDeploymentStateHandler).Methods("GET")

	// Environment endpoints
	api.HandleFunc("/environments", s.listEnvironmentsHandler).Methods("GET")
	api.HandleFunc("/environments/{name}", s.getEnvironmentHandler).Methods("GET")

	// Deployment endpoints
	api.HandleFunc("/deployments", s.listDeploymentsHandler).Methods("GET")
	api.HandleFunc("/deployments/{id}", s.getDeploymentHandler).Methods("GET")

	s.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", s.port),
		Handler: router,
	}

	log.Printf("Starting HTTP API server on port %d", s.port)
	return s.server.ListenAndServe()
}

func (s *HTTPAPIServer) Stop() error {
	if s.server != nil {
		return s.server.Shutdown(context.Background())
	}
	return nil
}

func (s *HTTPAPIServer) healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
}

func (s *HTTPAPIServer) getDeploymentStateHandler(w http.ResponseWriter, r *http.Request) {
	// Get query parameters
	manifestsEndpoint := r.URL.Query().Get("manifests_endpoint")
	commitID := r.URL.Query().Get("commit_id")

	if manifestsEndpoint == "" || commitID == "" {
		http.Error(w, "manifests_endpoint and commit_id query parameters are required", http.StatusBadRequest)
		return
	}

	// Query deployment state using existing logic from storage server
	deploymentState := &DeploymentStateResponse{
		TotalSubscribers:           0,
		TotalSucceededSubscribers:  0,
		TotalFailedSubscribers:     0,
		TotalInProgressSubscribers: 0,
		SucceededSubscribers:       []SubscriberResponse{},
		FailedSubscribers:          []SubscriberResponse{},
		InProgressSubscribers:      []SubscriberResponse{},
	}

	// Get total subscribers by manifests endpoint
	totalSubscribers, err := s.dbClient.Query(r.Context(), db.CountByManifestsEndpoint, manifestsEndpoint)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error querying total subscribers: %v", err), http.StatusInternalServerError)
		return
	}

	deploymentState.TotalSubscribers = totalSubscribers.(int32)

	if deploymentState.TotalSubscribers > 0 {
		// Get deployment statuses by manifests endpoint and gitops commit ID
		deploymentStatuses, err := s.dbClient.Query(r.Context(), db.CountByStatuses, manifestsEndpoint, commitID)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error querying deployment statuses: %v", err), http.StatusInternalServerError)
			return
		}

		status := deploymentStatuses.(db.StatusStats)
		deploymentState.TotalSucceededSubscribers = status.Success
		deploymentState.TotalFailedSubscribers = status.Failed
		deploymentState.TotalInProgressSubscribers = status.InProgress

		if deploymentState.TotalFailedSubscribers > 0 {
			failedDeployments, err := s.dbClient.Query(r.Context(), db.GetByStatus, manifestsEndpoint, commitID, "failure")
			if err != nil {
				http.Error(w, fmt.Sprintf("Error querying failed deployments: %v", err), http.StatusInternalServerError)
				return
			}
			for _, deployment := range failedDeployments.([]map[string]string) {
				subscriber := SubscriberResponse{
					Name:          deployment["name"],
					StatusMessage: deployment["status_message"],
				}
				deploymentState.FailedSubscribers = append(deploymentState.FailedSubscribers, subscriber)
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(deploymentState)
}

func (s *HTTPAPIServer) listEnvironmentsHandler(w http.ResponseWriter, r *http.Request) {
	// This is a simplified implementation - in a real scenario, we'd need a query to list all environments
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "List environments endpoint - implementation needed"})
}

func (s *HTTPAPIServer) getEnvironmentHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	if name == "" {
		http.Error(w, "Environment name is required", http.StatusBadRequest)
		return
	}

	// Get environment by natural key (name)
	env, err := s.dbClient.GetByNaturalKey(r.Context(), &db.Environment{Name: name})
	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
			http.Error(w, "Environment not found", http.StatusNotFound)
			return
		}
		http.Error(w, fmt.Sprintf("Error querying environment: %v", err), http.StatusInternalServerError)
		return
	}

	environment := env.(*db.Environment)
	response := EnvironmentResponse{
		ID:          environment.Id,
		Name:        environment.Name,
		Description: environment.Description,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (s *HTTPAPIServer) listDeploymentsHandler(w http.ResponseWriter, r *http.Request) {
	// This is a simplified implementation - in a real scenario, we'd need a query to list deployments
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "List deployments endpoint - implementation needed"})
}

func (s *HTTPAPIServer) getDeploymentHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	if idStr == "" {
		http.Error(w, "Deployment ID is required", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid deployment ID", http.StatusBadRequest)
		return
	}

	// Get deployment by ID
	dep, err := s.dbClient.Get(r.Context(), &db.Deployment{Id: id})
	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
			http.Error(w, "Deployment not found", http.StatusNotFound)
			return
		}
		http.Error(w, fmt.Sprintf("Error querying deployment: %v", err), http.StatusInternalServerError)
		return
	}

	deployment := dep.(*db.Deployment)
	response := DeploymentResponse{
		ID:             deployment.Id,
		GitopsCommitID: deployment.GitopsCommitId,
		ReconcilerID:   deployment.ReconcilerId,
		Status:         deployment.Status,
		StatusMessage:  deployment.StatusMessage,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
