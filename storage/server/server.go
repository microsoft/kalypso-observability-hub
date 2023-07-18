package main

import (
	context "context"
	"flag"
	"fmt"
	"log"
	"net"

	pb "github.com/microsoft/kalypso-observability-hub/storage/api"
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
	log.Printf("Received: %v", workspace)

	workspace, err := s.dbClient.UpdateWorkspace(ctx, &db.Workspace{
		Name:        workspace.Name,
		Description: workspace.Description,
	})
	if err != nil {
		return nil, err
	}

	log.Printf("Updated: %v", workspace)
	//return the workspace
	return workspace, nil
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
