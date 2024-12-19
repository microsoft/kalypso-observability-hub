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
	get(conn *sql.DB) (Entity, error)
	getByNaturalKey(conn *sql.DB) (Entity, error)
}

type QueryFunc func(conn *sql.DB, args ...interface{}) (interface{}, error)

type DBClient interface {
	Update(ctx context.Context, enity Entity) (Entity, error)
	Get(ctx context.Context, enity Entity) (Entity, error)
	GetByNaturalKey(ctx context.Context, enity Entity) (Entity, error)
	Query(ctx context.Context, query QueryFunc, args ...interface{}) (interface{}, error)
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
		log.Printf("fail to connect to database: %v", err)
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
		log.Printf("fail to update entity: %v", err)
		return nil, err
	}
	return entity, nil
}

func (c *postgresClient) Get(ctx context.Context, entity Entity) (Entity, error) {
	conn, err := c.getConnection()
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	entity, err = entity.get(conn)
	if err != nil {
		log.Printf("fail to get entity: %v", err)
		return nil, err
	}
	return entity, nil
}

func (c *postgresClient) GetByNaturalKey(ctx context.Context, entity Entity) (Entity, error) {
	conn, err := c.getConnection()
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	entity, err = entity.getByNaturalKey(conn)
	if err != nil {
		log.Printf("fail to get entity by natural key: %v", err)
		return nil, err
	}
	return entity, nil
}

func (c *postgresClient) Query(ctx context.Context, query QueryFunc, args ...interface{}) (interface{}, error) {
	conn, err := c.getConnection()
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	entities, err := query(conn, args...)
	if err != nil {
		log.Printf("fail to query entities: %v", err)
		return nil, err
	}
	return entities, nil
}
