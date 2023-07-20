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
