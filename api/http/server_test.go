package http

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func TestHealthHandler(t *testing.T) {
	server := &HTTPAPIServer{}

	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.healthHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var response map[string]string
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Errorf("could not parse response: %v", err)
	}

	if response["status"] != "healthy" {
		t.Errorf("handler returned unexpected body: got %v want %v", response["status"], "healthy")
	}
}

func TestGetDeploymentStateHandlerMissingParams(t *testing.T) {
	server := &HTTPAPIServer{}

	req, err := http.NewRequest("GET", "/api/v1/deployment-state", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.getDeploymentStateHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}

	expectedBody := "manifests_endpoint and commit_id query parameters are required\n"
	if rr.Body.String() != expectedBody {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expectedBody)
	}
}

func TestAPIRoutes(t *testing.T) {
	server := &HTTPAPIServer{}
	router := mux.NewRouter()

	// API routes
	api := router.PathPrefix("/api/v1").Subrouter()

	// Health check
	router.HandleFunc("/health", server.healthHandler).Methods("GET")

	// Deployment state endpoint
	api.HandleFunc("/deployment-state", server.getDeploymentStateHandler).Methods("GET")

	// Environment endpoints
	api.HandleFunc("/environments", server.listEnvironmentsHandler).Methods("GET")
	api.HandleFunc("/environments/{name}", server.getEnvironmentHandler).Methods("GET")

	// Deployment endpoints
	api.HandleFunc("/deployments", server.listDeploymentsHandler).Methods("GET")
	api.HandleFunc("/deployments/{id}", server.getDeploymentHandler).Methods("GET")

	testCases := []struct {
		method         string
		path           string
		expectedStatus int
	}{
		{"GET", "/health", http.StatusOK},
		{"GET", "/api/v1/deployment-state", http.StatusBadRequest}, // Missing params
		{"GET", "/api/v1/environments", http.StatusOK},
		{"GET", "/api/v1/deployments", http.StatusOK},
	}

	for _, tc := range testCases {
		req, err := http.NewRequest(tc.method, tc.path, nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		if status := rr.Code; status != tc.expectedStatus {
			t.Errorf("route %s %s returned wrong status code: got %v want %v", tc.method, tc.path, status, tc.expectedStatus)
		}
	}
}
