package client

import (
	context "context"
	"flag"
	"strconv"
	"testing"
	"time"

	pb "github.com/microsoft/kalypso-observability-hub/storage/api"
)

var (
	serverAddr = flag.String("addr", "localhost:50051", "The server address in the format of host:port")
)

// Test UpdateWorkspace
func TestUpdateWorkspace(t *testing.T) {
	client := NewObservabilityStorageClient(*serverAddr)
	_, err := client.UpdateWorkspace(context.Background(), &pb.Workspace{
		Name:        strconv.FormatInt(time.Now().UnixNano(), 10),
		Description: "Description",
	})
	if err != nil {
		t.Errorf("UpdateWorkspace() error = %v", err)
		return
	}
}
