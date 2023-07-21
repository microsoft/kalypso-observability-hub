package main

import (
	context "context"
	"flag"
	"fmt"
	"log"
	"net"

	pb "github.com/microsoft/kalypso-observability-hub/storage/api/grpc/proto"
	db "github.com/microsoft/kalypso-observability-hub/storage/postgres"
	"google.golang.org/grpc"
)

var (
	port             = flag.Int("port", 50051, "The server port")
	postgresHost     = flag.String("postgresHost", "localhost", "Postgres Host")
	postgresPort     = flag.Int("postgresPort", 5432, "Postgres Port")
	postgresUser     = flag.String("postgresUser", "creator", "Postgres User")
	postgresPassword = flag.String("postgresPassword", "c67", "Postgres Password")
	postgresDbName   = flag.String("postgresDbName", "hub", "Postgres Db Name")
	postgresSslmode  = flag.String("postgresSslmode", "disable", "Postgres Db Name")
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

func newStorageApiServer(dbClient db.DBClient) *storageApiServer {
	return &storageApiServer{dbClient: dbClient}
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	dbClient := db.NewPostgresClient(*postgresHost, *postgresPort, *postgresUser, *postgresPassword, *postgresDbName, *postgresSslmode)
	pb.RegisterStorageApiServer(grpcServer, newStorageApiServer(dbClient))
	//log starting the server
	log.Printf("Starting server on port %d", *port)

	grpcServer.Serve(lis)

}
