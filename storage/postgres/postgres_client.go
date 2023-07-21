package postgres

import (
	context "context"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type Entity interface {
	update(conn *sql.DB) (Entity, error)
}

type DBClient interface {
	Update(ctx context.Context, enity Entity) (Entity, error)
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

func (c *postgresClient) Update(ctx context.Context, entity Entity) (Entity, error) {
	conn, err := c.getConnection()
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	entity, err = entity.update(conn)
	if err != nil {
		log.Fatalf("fail to update entity: %v", err)
		return nil, err
	}
	return entity, nil
}
