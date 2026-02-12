package pqsqldb

import (
	"context"
	"fmt"
	"log"
	"order-bot-mgmt-svc/internal/config"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

var testConfig config.Db

func mustStartPostgresContainer() (config.Db, func(context.Context, ...testcontainers.TerminateOption) error, error) {
	var (
		dbName = "database"
		dbPwd  = "password"
		dbUser = "user"
	)

	dbContainer, err := postgres.Run(
		context.Background(),
		"postgres:latest",
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPwd),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		return config.Db{}, nil, fmt.Errorf("mustStartPostgresContainer: %w", err)
	}

	dbHost, err := dbContainer.Host(context.Background())
	if err != nil {
		return config.Db{}, dbContainer.Terminate, fmt.Errorf("mustStartPostgresContainer: %w", err)
	}

	dbPort, err := dbContainer.MappedPort(context.Background(), "5432/tcp")
	if err != nil {
		return config.Db{}, dbContainer.Terminate, fmt.Errorf("mustStartPostgresContainer: %w", err)
	}

	return config.Db{
		Database: dbName,
		Password: dbPwd,
		Username: dbUser,
		Host:     dbHost,
		Port:     dbPort.Port(),
		Schema:   "public",
	}, dbContainer.Terminate, err
}

func TestMain(m *testing.M) {
	dbConfig, teardown, err := mustStartPostgresContainer()
	if err != nil {
		log.Fatalf("could not start postgres container: %v", err)
	}
	testConfig = dbConfig

	m.Run()

	if teardown != nil && teardown(context.Background()) != nil {
		log.Fatalf("could not teardown postgres container: %v", err)
	}
}

func TestNew(t *testing.T) {
	srv, err := New(testConfig)
	if err != nil {
		t.Fatalf("New() returned error: %v", err)
	}
	if srv == nil {
		t.Fatal("New() returned nil")
	}
}

func TestHealth(t *testing.T) {
	srv, err := New(testConfig)
	if err != nil {
		t.Fatalf("New() returned error: %v", err)
	}

	stats, err := srv.Health()
	if err != nil {
		t.Fatalf("Health() returned error: %v", err)
	}

	if stats["status"] != "up" {
		t.Fatalf("expected status to be up, got %s", stats["status"])
	}

	if _, ok := stats["error"]; ok {
		t.Fatalf("expected error not to be present")
	}

	if stats["message"] != "It's healthy" {
		t.Fatalf("expected message to be 'It's healthy', got %s", stats["message"])
	}
}

func TestClose(t *testing.T) {
	srv, err := New(testConfig)
	if err != nil {
		t.Fatalf("New() returned error: %v", err)
	}

	if srv.Close() != nil {
		t.Fatalf("expected Close() to return nil")
	}
}
