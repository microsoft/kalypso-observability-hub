package grpcClient

import (
	context "context"
	"log"

	pb "github.com/microsoft/kalypso-observability-hub/storage/api/grpc/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ObservabilityStorageGrpcClient interface {
	UpdateWorkspace(ctx context.Context, in *pb.Workspace) (*pb.Workspace, error)
	UpdateApplication(ctx context.Context, in *pb.Application) (*pb.Application, error)
	UpdateWorkload(ctx context.Context, in *pb.Workload) (*pb.Workload, error)
}

type observabilityStorageGrpcClient struct {
	serverAddr string
}

// make sure that the client implements the interface
var _ ObservabilityStorageGrpcClient = (*observabilityStorageGrpcClient)(nil)

func NewObservabilityStorageGrpcClient(serverAddr string) ObservabilityStorageGrpcClient {
	return &observabilityStorageGrpcClient{serverAddr: serverAddr}
}

func (c *observabilityStorageGrpcClient) getConnection() (*grpc.ClientConn, error) {
	conn, err := grpc.Dial(c.serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
		return nil, err
	}
	return conn, nil
}

func (c *observabilityStorageGrpcClient) UpdateWorkspace(ctx context.Context, in *pb.Workspace) (*pb.Workspace, error) {
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

func (c *observabilityStorageGrpcClient) UpdateApplication(ctx context.Context, in *pb.Application) (*pb.Application, error) {
	conn, err := c.getConnection()
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	client := pb.NewStorageApiClient(conn)
	app, err := client.UpdateApplication(ctx, in)
	if err != nil {
		return nil, err
	}
	return app, nil
}

func (c *observabilityStorageGrpcClient) UpdateWorkload(ctx context.Context, in *pb.Workload) (*pb.Workload, error) {
	conn, err := c.getConnection()
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	client := pb.NewStorageApiClient(conn)
	wl, err := client.UpdateWorkload(ctx, in)
	if err != nil {
		return nil, err
	}
	return wl, nil
}
