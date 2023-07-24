package client

import (
	context "context"
	"flag"
	"testing"

	pb "github.com/microsoft/kalypso-observability-hub/storage/api/grpc/proto"
)

var (
	serverAddr = flag.String("addr", "localhost:50051", "The server address in the format of host:port")
)

// Test UpdateWorkspace
func TestUpdateWorkspace(t *testing.T) {
	client := NewObservabilityStorageGrpcClient(*serverAddr)
	_, err := client.UpdateWorkspace(context.Background(), &pb.Workspace{
		Name:        "Name",
		Description: "Description",
	})
	if err != nil {
		t.Errorf("UpdateWorkspace() error = %v", err)
		return
	}
}

// Test UpdateApplication
func TestUpdateApplication(t *testing.T) {
	client := NewObservabilityStorageGrpcClient(*serverAddr)
	_, err := client.UpdateApplication(context.Background(), &pb.Application{
		Name:        "Name",
		Description: "Description",
		WorkspaceId: 1,
	})
	if err != nil {
		t.Errorf("UpdateApplication() error = %v", err)
		return
	}
}

// Test UpdateWorkload
func TestUpdateWorkload(t *testing.T) {
	client := NewObservabilityStorageGrpcClient(*serverAddr)
	_, err := client.UpdateWorkload(context.Background(), &pb.Workload{
		Name:              "Name",
		Description:       "Description",
		SourceStorageType: "git",
		SourceEndpoint:    "http://github.com",
		ApplicationId:     1,
	})
	if err != nil {
		t.Errorf("UpdateWorkload() error = %v", err)
		return
	}
}

// Test UpdateEnvironment
func TestUpdateEnvironment(t *testing.T) {
	client := NewObservabilityStorageGrpcClient(*serverAddr)
	_, err := client.UpdateEnvironment(context.Background(), &pb.Environment{
		Name:        "Name",
		Description: "Description",
	})
	if err != nil {
		t.Errorf("UpdateEnvironment() error = %v", err)
		return
	}
}

// Test UpdateDeploymentTarget
func TestUpdateDeploymentTarget(t *testing.T) {
	client := NewObservabilityStorageGrpcClient(*serverAddr)
	_, err := client.UpdateDeploymentTarget(context.Background(), &pb.DeploymentTarget{
		Name:                 "Name",
		Description:          "Description",
		WorkloadId:           1,
		EnvironmentId:        1,
		Labels:               "good",
		ManifestsStorageType: "git",
		ManifestsEndpoint:    "http://github.com",
	})
	if err != nil {
		t.Errorf("UpdateDeploymentTarget() error = %v", err)
		return
	}
}

// Test UpdateWorkloadVersion
func TestUpdateWorkloadVersion(t *testing.T) {
	client := NewObservabilityStorageGrpcClient(*serverAddr)
	_, err := client.UpdateWorkloadVersion(context.Background(), &pb.WorkloadVersion{
		Version:       "1.0.0",
		WorkloadId:    1,
		BuildCommitId: "1234567890",
		BuildId:       "123",
	})
	if err != nil {
		t.Errorf("UpdateWorkloadVersion() error = %v", err)
		return
	}
}

// Test UpdateDeploymentAssignment
func TestUpdateDeploymentAssignment(t *testing.T) {
	client := NewObservabilityStorageGrpcClient(*serverAddr)
	_, err := client.UpdateDeploymentAssignment(context.Background(), &pb.DeploymentAssignment{
		DeploymentTargetId: 1,
		WorkloadVersionId:  1,
		GitopsCommitId:     "1234567890",
	})
	if err != nil {
		t.Errorf("UpdateDeploymentAssignment() error = %v", err)
		return
	}
}

// Test UpdateHost
func TestUpdateHost(t *testing.T) {
	client := NewObservabilityStorageGrpcClient(*serverAddr)
	_, err := client.UpdateHost(context.Background(), &pb.Host{
		Name:        "Name",
		Description: "Description",
		HostType:    "HostType",
		Labels:      "Labels",
	})
	if err != nil {
		t.Errorf("UpdateHost() error = %v", err)
		return
	}
}

// Test UpdateReconciler
func TestUpdateReconciler(t *testing.T) {
	client := NewObservabilityStorageGrpcClient(*serverAddr)
	_, err := client.UpdateReconciler(context.Background(), &pb.Reconciler{
		Name:                 "Name",
		Description:          "Description",
		HostId:               1,
		ReconcilerType:       "Flux",
		Labels:               "Labels",
		ManifestsStorageType: "git",
		ManifestsEndpoint:    "http://github.com",
	})
	if err != nil {
		t.Errorf("UpdateReconciler() error = %v", err)
		return
	}
}

// Test UpdateDeployment
func TestUpdateDeployment(t *testing.T) {
	client := NewObservabilityStorageGrpcClient(*serverAddr)
	_, err := client.UpdateDeployment(context.Background(), &pb.Deployment{
		GitopsCommitId: "1234567890",
		ReconcilerId:   1,
		Status:         "success",
		StatusMessage:  "Successfully deployed",
	})
	if err != nil {
		t.Errorf("UpdateDeployment() error = %v", err)
		return
	}
}
