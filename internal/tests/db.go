package tests

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"testing"

	"github.com/docker/go-connections/nat"
	"github.com/pressly/goose/v3"
	"github.com/testcontainers/testcontainers-go"
	testPostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"github.com/yeralin-munar/tt-go-json-fernet/config"
	"github.com/yeralin-munar/tt-go-json-fernet/internal/data/postgres"

	_ "github.com/jackc/pgx/v5/stdlib"
)

var (
	user          = "test"
	database      = "test"
	port          = nat.Port("5432")
	WaitCondition = wait.ForAll(
		wait.ForLog("PostgreSQL init process complete; ready for start up."),
		wait.ForExec([]string{"psql", "-U", user, "-d", database, "-c", "SELECT pg_sleep(.5);"}),
		wait.ForListeningPort(port),
		wait.ForLog("database system is ready to accept connections"),
	)
)

func WithExportPorts(exportPorts []string) testcontainers.CustomizeRequestOption {
	return func(req *testcontainers.GenericContainerRequest) {
		req.ContainerRequest.ExposedPorts = exportPorts
		req.ExposedPorts = exportPorts
	}
}

func WithWaitStrategy(waitStrategy wait.Strategy) testcontainers.CustomizeRequestOption {
	return func(req *testcontainers.GenericContainerRequest) {
		req.WaitingFor = waitStrategy
	}
}

func GetTestPostgresConnection(
	t *testing.T,
	migrationPath string,
) (
	*config.Data,
	*testPostgres.PostgresContainer,
) {
	ctx := context.Background()

	postgresContainer, err := testPostgres.RunContainer(
		ctx,
		testcontainers.WithImage("docker.io/postgres:15.2-alpine"),
		testPostgres.WithDatabase("test"),
		testPostgres.WithUsername("test"),
		testPostgres.WithPassword("test"),
		WithExportPorts([]string{port.Port() + "/tcp"}),
		// WithWaitStrategy(WaitCondition),
		// testcontainers.WithWaitStrategy(
		// 	wait.ForLog("database system is ready to accept connections").
		// 		WithOccurrence(2).
		// 		WithStartupTimeout(5*time.Second)),
		testcontainers.WithWaitStrategy(WaitCondition),
	)
	if err != nil {
		t.Fatalf("failed to start container: %s", err)
	}

	// Clean up the container
	t.Cleanup(func() {
		log.Println("terminating PostgreSQL container")
		if err := postgresContainer.Terminate(ctx); err != nil {
			log.Fatalf("failed to terminate container: %s", err)
		}
	})

	ip, err := postgresContainer.Host(ctx)
	if err != nil {
		t.Fatal(err)
	}

	mappedPort, err := postgresContainer.MappedPort(ctx, port)
	if err != nil {
		t.Fatal(err)
	}

	// Connect to the database
	cfg :=
		config.Data{
			Database: config.Database{
				Host:     ip,
				Port:     mappedPort.Int(),
				User:     "test",
				Password: "test",
				Name:     "test",
			},
		}

	// Run goose migrations
	if err := runMigrations(&cfg, migrationPath); err != nil {
		t.Fatalf("failed to run migrations: %s", err)
	}

	// Create a snapshot of the database to restore later
	log.Println("creating snapshot of the database")
	err = postgresContainer.Snapshot(ctx, testPostgres.WithSnapshotName("test-snapshot"))
	if err != nil {
		t.Fatal(err)
	}

	return &cfg, postgresContainer

}

func runMigrations(
	cfg *config.Data,
	migrationPath string,
) error {
	migrationConn, err := sql.Open("pgx", postgres.GenerateDBURL(cfg))
	if err != nil {
		return fmt.Errorf("failed to open migration connection: %s", err)
	}

	defer migrationConn.Close()

	err = goose.SetDialect("postgres")
	if err != nil {
		return fmt.Errorf("failed to set goose dialect: %s", err)
	}

	err = goose.Up(migrationConn, migrationPath)
	if err != nil {
		return fmt.Errorf("failed to run migrations: %s", err)
	}

	return nil
}
