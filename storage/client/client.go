package client

import (
	context "context"
	"log"

	pb "github.com/microsoft/kalypso-observability-hub/storage/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ObservabilityStorageClient interface {
	UpdateWorkspace(ctx context.Context, in *pb.Workspace) (*pb.Workspace, error)
}

type observabilityStorageClient struct {
	serverAddr string
}

// make sure that the client implements the interface
var _ ObservabilityStorageClient = (*observabilityStorageClient)(nil)

func NewObservabilityStorageClient(serverAddr string) ObservabilityStorageClient {
	return &observabilityStorageClient{serverAddr: serverAddr}
}

func (c *observabilityStorageClient) getConnection() (*grpc.ClientConn, error) {
	conn, err := grpc.Dial(c.serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
		return nil, err
	}
	return conn, nil
}

func (c *observabilityStorageClient) UpdateWorkspace(ctx context.Context, in *pb.Workspace) (*pb.Workspace, error) {
	conn, err := c.getConnection()
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	client := pb.NewStorageApiClient(conn)
	ws, err := client.UpdateWorkspace(ctx, in)
	if err != nil {
		return nil, err
	}
	return ws, nil
}
