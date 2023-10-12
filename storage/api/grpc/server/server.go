package main

import (
	context "context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"

	pb "github.com/microsoft/kalypso-observability-hub/storage/api/grpc/proto"
	db "github.com/microsoft/kalypso-observability-hub/storage/postgres"
	"google.golang.org/grpc"

	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

var (
	port             *string
	postgresHost     *string
	postgresPort     *string
	postgresUser     *string
	postgresPassword *string
	postgresDbName   *string
	postgresSslmode  *string

	portInt         int
	postgresPortInt int
)

type storageApiServer struct {
	pb.UnimplementedStorageApiServer
	dbClient db.DBClient
}

// make sure that the server implements the interface
var _ pb.StorageApiServer = (*storageApiServer)(nil)

func (s *storageApiServer) UpdateWorkspace(ctx context.Context, workspace *pb.Workspace) (*pb.Workspace, error) {
	log.Printf("Received Workspace: %v", workspace)

	ws, err := s.dbClient.Update(ctx, &db.Workspace{
		Name:        workspace.Name,
		Description: workspace.Description,
	})
	if err != nil {
		return nil, err
	}
	workspace_entity := ws.(*db.Workspace)

	log.Printf("Updated Workspace: %v", workspace_entity)
	//return the workspace
	return &pb.Workspace{
		Id:          int32(workspace_entity.Id),
		Name:        workspace_entity.Name,
		Description: workspace_entity.Description,
	}, nil
}

func (s *storageApiServer) UpdateApplication(ctx context.Context, application *pb.Application) (*pb.Application, error) {
	log.Printf("Received Application: %v", application)

	app, err := s.dbClient.Update(ctx, &db.Application{
		Name:        application.Name,
		Description: application.Description,
		WorkspaceId: int(application.WorkspaceId),
	})
	if err != nil {
		return nil, err
	}
	application_entity := app.(*db.Application)

	log.Printf("Updated Application: %v", application_entity)
	//return the application
	return &pb.Application{
		Id:          int32(application_entity.Id),
		Name:        application_entity.Name,
		Description: application_entity.Description,
		WorkspaceId: int32(application_entity.WorkspaceId),
	}, nil
}

func (s *storageApiServer) UpdateWorkload(ctx context.Context, workload *pb.Workload) (*pb.Workload, error) {
	log.Printf("Received Workload: %v", workload)

	wl, err := s.dbClient.Update(ctx, &db.Workload{
		Name:              workload.Name,
		Description:       workload.Description,
		SourceStorageType: workload.SourceStorageType,
		SourceEndpoint:    workload.SourceEndpoint,
		ApplicationId:     int(workload.ApplicationId),
	})
	if err != nil {
		return nil, err
	}
	workload_entity := wl.(*db.Workload)

	log.Printf("Updated Workload: %v", workload_entity)
	//return the workload
	return &pb.Workload{
		Id:                int32(workload_entity.Id),
		Name:              workload_entity.Name,
		Description:       workload_entity.Description,
		SourceStorageType: workload_entity.SourceStorageType,
		SourceEndpoint:    workload_entity.SourceEndpoint,
		ApplicationId:     int32(workload_entity.ApplicationId),
	}, nil
}

func (s *storageApiServer) UpdateEnvironment(ctx context.Context, environment *pb.Environment) (*pb.Environment, error) {
	log.Printf("Received Environment: %v", environment)

	env, err := s.dbClient.Update(ctx, &db.Environment{
		Name:        environment.Name,
		Description: environment.Description,
	})
	if err != nil {
		return nil, err
	}
	environment_entity := env.(*db.Environment)

	log.Printf("Updated Environment: %v", environment_entity)
	//return the environment
	return &pb.Environment{
		Id:          int32(environment_entity.Id),
		Name:        environment_entity.Name,
		Description: environment_entity.Description,
	}, nil
}

// Update DeploymentTarget
func (s *storageApiServer) UpdateDeploymentTarget(ctx context.Context, deploymentTarget *pb.DeploymentTarget) (*pb.DeploymentTarget, error) {
	log.Printf("Received DeploymentTarget: %v", deploymentTarget)

	dt, err := s.dbClient.Update(ctx, &db.DeploymentTarget{
		Name:                 deploymentTarget.Name,
		Description:          deploymentTarget.Description,
		WorkloadId:           int(deploymentTarget.WorkloadId),
		EnvironmentId:        int(deploymentTarget.EnvironmentId),
		Labels:               deploymentTarget.Labels,
		ManifestsStorageType: deploymentTarget.ManifestsStorageType,
		ManifestsEndpoint:    deploymentTarget.ManifestsEndpoint,
	})
	if err != nil {
		return nil, err
	}
	deploymentTarget_entity := dt.(*db.DeploymentTarget)

	log.Printf("Updated DeploymentTarget: %v", deploymentTarget_entity)
	//return the deploymentTarget
	return &pb.DeploymentTarget{
		Id:                   int32(deploymentTarget_entity.Id),
		Name:                 deploymentTarget_entity.Name,
		Description:          deploymentTarget_entity.Description,
		WorkloadId:           int32(deploymentTarget_entity.WorkloadId),
		EnvironmentId:        int32(deploymentTarget_entity.EnvironmentId),
		Labels:               deploymentTarget_entity.Labels,
		ManifestsStorageType: deploymentTarget_entity.ManifestsStorageType,
		ManifestsEndpoint:    deploymentTarget_entity.ManifestsEndpoint,
	}, nil
}

// Update WorkloadVersion
func (s *storageApiServer) UpdateWorkloadVersion(ctx context.Context, workloadVersion *pb.WorkloadVersion) (*pb.WorkloadVersion, error) {
	log.Printf("Received WorkloadVersion: %v", workloadVersion)

	wv, err := s.dbClient.Update(ctx, &db.WorkloadVersion{
		WorkloadId:    int(workloadVersion.WorkloadId),
		Version:       workloadVersion.Version,
		BuildId:       workloadVersion.BuildId,
		BuildCommitId: workloadVersion.BuildCommitId,
	})
	if err != nil {
		return nil, err
	}
	workloadVersion_entity := wv.(*db.WorkloadVersion)

	log.Printf("Updated WorkloadVersion: %v", workloadVersion_entity)
	//return the workloadVersion
	return &pb.WorkloadVersion{
		Id:            int32(workloadVersion_entity.Id),
		WorkloadId:    int32(workloadVersion_entity.WorkloadId),
		Version:       workloadVersion_entity.Version,
		BuildId:       workloadVersion_entity.BuildId,
		BuildCommitId: workloadVersion_entity.BuildCommitId,
	}, nil
}

// Update DeploymentAssignment
func (s *storageApiServer) UpdateDeploymentAssignment(ctx context.Context, deploymentAssignment *pb.DeploymentAssignment) (*pb.DeploymentAssignment, error) {
	log.Printf("Received DeploymentAssignment: %v", deploymentAssignment)

	da, err := s.dbClient.Update(ctx, &db.DeploymentAssignment{
		DeploymentTargetId: int(deploymentAssignment.DeploymentTargetId),
		WorkloadVersionId:  int(deploymentAssignment.WorkloadVersionId),
		GitopsCommitId:     deploymentAssignment.GitopsCommitId,
	})
	if err != nil {
		return nil, err
	}
	deploymentAssignment_entity := da.(*db.DeploymentAssignment)

	log.Printf("Updated DeploymentAssignment: %v", deploymentAssignment_entity)
	//return the deploymentAssignment
	return &pb.DeploymentAssignment{
		Id:                 int32(deploymentAssignment_entity.Id),
		DeploymentTargetId: int32(deploymentAssignment_entity.DeploymentTargetId),
		WorkloadVersionId:  int32(deploymentAssignment_entity.WorkloadVersionId),
		GitopsCommitId:     deploymentAssignment_entity.GitopsCommitId,
	}, nil
}

// Update Host
func (s *storageApiServer) UpdateHost(ctx context.Context, host *pb.Host) (*pb.Host, error) {
	log.Printf("Received Host: %v", host)

	h, err := s.dbClient.Update(ctx, &db.Host{
		Name:        host.Name,
		Description: host.Description,
		HostType:    host.HostType,
		Labels:      host.Labels,
	})
	if err != nil {
		return nil, err
	}
	host_entity := h.(*db.Host)

	log.Printf("Updated Host: %v", host_entity)
	//return the host
	return &pb.Host{
		Id:          int32(host_entity.Id),
		Name:        host_entity.Name,
		Description: host_entity.Description,
		HostType:    host_entity.HostType,
		Labels:      host_entity.Labels,
	}, nil
}

// Update Reconciler
func (s *storageApiServer) UpdateReconciler(ctx context.Context, reconciler *pb.Reconciler) (*pb.Reconciler, error) {
	log.Printf("Received Reconciler: %v", reconciler)

	r, err := s.dbClient.Update(ctx, &db.Reconciler{
		Name:                 reconciler.Name,
		Description:          reconciler.Description,
		HostId:               int(reconciler.HostId),
		ReconcilerType:       reconciler.ReconcilerType,
		Labels:               reconciler.Labels,
		ManifestsStorageType: reconciler.ManifestsStorageType,
		ManifestsEndpoint:    reconciler.ManifestsEndpoint,
	})
	if err != nil {
		return nil, err
	}
	reconciler_entity := r.(*db.Reconciler)

	log.Printf("Updated Reconciler: %v", reconciler_entity)
	//return the reconciler
	return &pb.Reconciler{
		Id:                   int32(reconciler_entity.Id),
		Name:                 reconciler_entity.Name,
		Description:          reconciler_entity.Description,
		HostId:               int32(reconciler_entity.HostId),
		ReconcilerType:       reconciler_entity.ReconcilerType,
		Labels:               reconciler_entity.Labels,
		ManifestsStorageType: reconciler_entity.ManifestsStorageType,
		ManifestsEndpoint:    reconciler_entity.ManifestsEndpoint,
	}, nil
}

// Update Deployment
func (s *storageApiServer) UpdateDeployment(ctx context.Context, deployment *pb.Deployment) (*pb.Deployment, error) {
	log.Printf("Received Deployment: %v", deployment)

	d, err := s.dbClient.Update(ctx, &db.Deployment{
		GitopsCommitId: deployment.GitopsCommitId,
		ReconcilerId:   int(deployment.ReconcilerId),
		Status:         deployment.Status,
		StatusMessage:  deployment.StatusMessage,
	})
	if err != nil {
		return nil, err
	}
	deployment_entity := d.(*db.Deployment)

	log.Printf("Updated Deployment: %v", deployment_entity)
	//return the deployment
	return &pb.Deployment{
		Id:             int32(deployment_entity.Id),
		GitopsCommitId: deployment_entity.GitopsCommitId,
		ReconcilerId:   int32(deployment_entity.ReconcilerId),
		Status:         deployment_entity.Status,
		StatusMessage:  deployment_entity.StatusMessage,
	}, nil
}

// Get Deployment Target
func (s *storageApiServer) GetDeploymentTarget(ctx context.Context, deploymentTargetSearch *pb.DeploymentTargetSearch) (*pb.DeploymentTarget, error) {
	log.Printf("Received DeploymentTargetSearch: %v", deploymentTargetSearch)

	//Get Environment by natural key
	env, err := s.dbClient.GetByNaturalKey(ctx, &db.Environment{
		Name: deploymentTargetSearch.EnvironmentName,
	})
	if err != nil {
		return nil, err
	}

	//Get Workspace by natural key
	ws, err := s.dbClient.GetByNaturalKey(ctx, &db.Workspace{
		Name: deploymentTargetSearch.WorkspaceName,
	})
	if err != nil {
		return nil, err
	}

	//Get Application by natural key
	app, err := s.dbClient.GetByNaturalKey(ctx, &db.Application{
		Name:        deploymentTargetSearch.ApplicationName,
		WorkspaceId: ws.(*db.Workspace).Id,
	})
	if err != nil {
		return nil, err
	}

	//Get Workload by natural key
	wl, err := s.dbClient.GetByNaturalKey(ctx, &db.Workload{
		Name:          deploymentTargetSearch.WorkloadName,
		ApplicationId: app.(*db.Application).Id,
	})
	if err != nil {
		return nil, err
	}

	//Get DeploymentTarget by natural key
	dt, err := s.dbClient.GetByNaturalKey(ctx, &db.DeploymentTarget{
		Name:          deploymentTargetSearch.DeploymentTargetName,
		EnvironmentId: env.(*db.Environment).Id,
		WorkloadId:    wl.(*db.Workload).Id,
	})
	if err != nil {
		return nil, err
	}

	deploymentTarget_entity := dt.(*db.DeploymentTarget)

	log.Printf("Got DeploymentTarget: %v", deploymentTarget_entity)
	//return the deploymentTarget
	return &pb.DeploymentTarget{
		Id:                   int32(deploymentTarget_entity.Id),
		Name:                 deploymentTarget_entity.Name,
		Description:          deploymentTarget_entity.Description,
		WorkloadId:           int32(deploymentTarget_entity.WorkloadId),
		EnvironmentId:        int32(deploymentTarget_entity.EnvironmentId),
		Labels:               deploymentTarget_entity.Labels,
		ManifestsStorageType: deploymentTarget_entity.ManifestsStorageType,
		ManifestsEndpoint:    deploymentTarget_entity.ManifestsEndpoint,
	}, nil
}

func newStorageApiServer(dbClient db.DBClient) *storageApiServer {
	return &storageApiServer{dbClient: dbClient}
}

func getEnv(key string) *string {
	// if ther is an env variable with the key, return it
	if value, ok := os.LookupEnv(key); ok {
		return &value
	} else {
		// throw an error if there is no env variable with the key
		log.Fatalf("environment variable %s not set", key)
		return nil
	}
}

func readConfigValuesFromEnv() {
	port = getEnv("PORT")
	postgresHost = getEnv("POSTGRES_HOST")
	postgresPort = getEnv("POSTGRES_PORT")
	postgresUser = getEnv("POSTGRES_USER")
	postgresPassword = getEnv("POSTGRES_PASSWORD")
	postgresDbName = getEnv("POSTGRES_DBNAME")
	postgresSslmode = getEnv("POSTGRES_SSL_MODE")

	var err error
	portInt, err = strconv.Atoi(*port)
	if err != nil {
		log.Fatalf("failed to convert port to int: %v", err)
	}
	postgresPortInt, err = strconv.Atoi(*postgresPort)
	if err != nil {
		log.Fatalf("failed to convert postgres port to int: %v", err)
	}

}

func main() {
	flag.Parse()
	readConfigValuesFromEnv()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", portInt))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()

	grpc_health_v1.RegisterHealthServer(grpcServer, health.NewServer())

	dbClient := db.NewPostgresClient(*postgresHost, postgresPortInt, *postgresUser, *postgresPassword, *postgresDbName, *postgresSslmode)
	pb.RegisterStorageApiServer(grpcServer, newStorageApiServer(dbClient))
	//log starting the server
	log.Printf("Starting server on port %d", portInt)

	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
