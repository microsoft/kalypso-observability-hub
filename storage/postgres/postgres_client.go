package postgres

import (
	context "context"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	pb "github.com/microsoft/kalypso-observability-hub/storage/api"
)

type Workspace struct {
	Id          int
	Name        string
	Description string
}

type DBClient interface {
	UpdateWorkspace(ctx context.Context, ws *Workspace) (*pb.Workspace, error)
}

type postgresClient struct {
	connectionString string
}

// make sure that the client implements the interface
var _ DBClient = (*postgresClient)(nil)

func NewPostgresClient(host string, port int, user, password, dbname, sslmode string) DBClient {
	return &postgresClient{connectionString: fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode)}
}

func (c *postgresClient) getConnection() (*sql.DB, error) {
	conn, err := sql.Open("postgres", c.connectionString)
	if err != nil {
		log.Fatalf("fail to connect to database: %v", err)
		return nil, err
	}
	return conn, nil
}

func (c *postgresClient) UpdateWorkspace(ctx context.Context, ws *Workspace) (*pb.Workspace, error) {
	conn, err := c.getConnection()
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	var id int32
	err = conn.QueryRow("INSERT INTO workspace (name, description) VALUES ($1, $2) RETURNING id", ws.Name, ws.Description).Scan(&id)
	if err != nil {
		return nil, err
	}

	return &pb.Workspace{Id: id, Name: ws.Name, Description: ws.Description}, nil
}
