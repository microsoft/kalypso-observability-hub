package postgres

import (
	context "context"
	"strconv"
	"testing"
	"time"
)

var (
	host     = "localhost"
	port     = 5432
	user     = "creator"
	password = "c67"
	dbname   = "hub"
	sslmode  = "disable"
)

// Test UpdateWorkspace
func TestUpdateWorkspace(t *testing.T) {
	client := NewPostgresClient(host, port, user, password, dbname, sslmode)
	_, err := client.UpdateWorkspace(context.Background(), &Workspace{
		Name:        strconv.FormatInt(time.Now().UnixNano(), 10),
		Description: "Description",
	})
	if err != nil {
		t.Errorf("UpdateWorkspace() error = %v", err)
		return
	}
}
