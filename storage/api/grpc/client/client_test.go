package grpcClient

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
