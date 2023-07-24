package client

import (
	context "context"
	"errors"
	"log"
	"os"

	pb "github.com/microsoft/kalypso-observability-hub/storage/api/grpc/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ObservabilityStorageGrpcClient interface {
	UpdateWorkspace(ctx context.Context, in *pb.Workspace) (*pb.Workspace, error)
	UpdateApplication(ctx context.Context, in *pb.Application) (*pb.Application, error)
	UpdateWorkload(ctx context.Context, in *pb.Workload) (*pb.Workload, error)
	UpdateEnvironment(ctx context.Context, in *pb.Environment) (*pb.Environment, error)
	UpdateDeploymentTarget(ctx context.Context, in *pb.DeploymentTarget) (*pb.DeploymentTarget, error)
	UpdateWorkloadVersion(ctx context.Context, in *pb.WorkloadVersion) (*pb.WorkloadVersion, error)
	UpdateDeploymentAssignment(ctx context.Context, in *pb.DeploymentAssignment) (*pb.DeploymentAssignment, error)
	UpdateHost(ctx context.Context, in *pb.Host) (*pb.Host, error)
	UpdateReconciler(ctx context.Context, in *pb.Reconciler) (*pb.Reconciler, error)
	UpdateDeployment(ctx context.Context, in *pb.Deployment) (*pb.Deployment, error)
}

type observabilityStorageGrpcClient struct {
	serverAddr string
}

// make sure that the client implements the interface
var _ ObservabilityStorageGrpcClient = (*observabilityStorageGrpcClient)(nil)

func NewObservabilityStorageGrpcClient(serverAddr string) ObservabilityStorageGrpcClient {

	return &observabilityStorageGrpcClient{serverAddr: serverAddr}
}

func GetObservabilityStorageGrpcClient() (ObservabilityStorageGrpcClient, error) {
	serverAddr := os.Getenv("STORAGE_SERVICE_ADDRESS")
	if serverAddr == "" {
		return nil, errors.New("STORAGE_SERVICE_ADDRESS environment variable not set")
	}

	return NewObservabilityStorageGrpcClient(serverAddr), nil

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

func (c *observabilityStorageGrpcClient) UpdateEnvironment(ctx context.Context, in *pb.Environment) (*pb.Environment, error) {
	conn, err := c.getConnection()
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	client := pb.NewStorageApiClient(conn)
	env, err := client.UpdateEnvironment(ctx, in)
	if err != nil {
		return nil, err
	}
	return env, nil
}

func (c *observabilityStorageGrpcClient) UpdateDeploymentTarget(ctx context.Context, in *pb.DeploymentTarget) (*pb.DeploymentTarget, error) {
	conn, err := c.getConnection()
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	client := pb.NewStorageApiClient(conn)
	dt, err := client.UpdateDeploymentTarget(ctx, in)
	if err != nil {
		return nil, err
	}
	return dt, nil
}

func (c *observabilityStorageGrpcClient) UpdateWorkloadVersion(ctx context.Context, in *pb.WorkloadVersion) (*pb.WorkloadVersion, error) {
	conn, err := c.getConnection()
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	client := pb.NewStorageApiClient(conn)
	wlv, err := client.UpdateWorkloadVersion(ctx, in)
	if err != nil {
		return nil, err
	}
	return wlv, nil
}

func (c *observabilityStorageGrpcClient) UpdateDeploymentAssignment(ctx context.Context, in *pb.DeploymentAssignment) (*pb.DeploymentAssignment, error) {
	conn, err := c.getConnection()
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	client := pb.NewStorageApiClient(conn)
	da, err := client.UpdateDeploymentAssignment(ctx, in)
	if err != nil {
		return nil, err
	}
	return da, nil
}

func (c *observabilityStorageGrpcClient) UpdateHost(ctx context.Context, in *pb.Host) (*pb.Host, error) {
	conn, err := c.getConnection()
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	client := pb.NewStorageApiClient(conn)
	host, err := client.UpdateHost(ctx, in)
	if err != nil {
		return nil, err
	}
	return host, nil
}

func (c *observabilityStorageGrpcClient) UpdateReconciler(ctx context.Context, in *pb.Reconciler) (*pb.Reconciler, error) {
	conn, err := c.getConnection()
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	client := pb.NewStorageApiClient(conn)
	rec, err := client.UpdateReconciler(ctx, in)
	if err != nil {
		return nil, err
	}
	return rec, nil
}

func (c *observabilityStorageGrpcClient) UpdateDeployment(ctx context.Context, in *pb.Deployment) (*pb.Deployment, error) {
	conn, err := c.getConnection()
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	client := pb.NewStorageApiClient(conn)
	dep, err := client.UpdateDeployment(ctx, in)
	if err != nil {
		return nil, err
	}
	return dep, nil
}
