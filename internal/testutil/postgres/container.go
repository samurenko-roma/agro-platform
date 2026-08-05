//go:build integration

package postgres

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type TestDB struct {
	Pool *pgxpool.Pool
	DSN  string
}

// migrationsRoot — относительный путь от этого файла до папки migrations.
const migrationsRoot = "../../../migrations"

func NewTestDB(t *testing.T, modules ...string) *TestDB {
	t.Helper()
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "postgis/postgis:16-3.5",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_DB":       "agrodb_test",
			"POSTGRES_USER":     "agro",
			"POSTGRES_PASSWORD": "agro123",
		},
		WaitingFor: wait.ForLog("database system is ready to accept connections").
			WithOccurrence(2).
			WithStartupTimeout(60 * time.Second),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Fatalf("start postgres container: %v", err)
	}
	t.Cleanup(func() { _ = container.Terminate(context.Background()) })

	host, err := container.Host(ctx)
	if err != nil {
		t.Fatalf("container host: %v", err)
	}
	port, err := container.MappedPort(ctx, "5432")
	if err != nil {
		t.Fatalf("container port: %v", err)
	}

	dsn := fmt.Sprintf("postgres://agro:agro123@%s:%s/agrodb_test?sslmode=disable", host, port.Port())

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		t.Fatalf("create pool: %v", err)
	}
	t.Cleanup(pool.Close)

	if err := waitForPing(ctx, pool); err != nil {
		t.Fatalf("ping: %v", err)
	}

	for _, module := range modules {
		if err := applyMigrations(ctx, pool, module); err != nil {
			t.Fatalf("migrations %q: %v", module, err)
		}
	}

	return &TestDB{Pool: pool, DSN: dsn}
}

func waitForPing(ctx context.Context, pool *pgxpool.Pool) error {
	deadline := time.Now().Add(15 * time.Second)
	var lastErr error
	for time.Now().Before(deadline) {
		if lastErr = pool.Ping(ctx); lastErr == nil {
			return nil
		}
		time.Sleep(300 * time.Millisecond)
	}
	return lastErr
}

func applyMigrations(ctx context.Context, pool *pgxpool.Pool, module string) error {
	dir := filepath.Join(migrationsRoot, module)
	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("read dir %s: %w", dir, err)
	}
	var files []string
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".up.sql") {
			files = append(files, e.Name())
		}
	}
	sort.Strings(files)
	for _, f := range files {
		content, err := os.ReadFile(filepath.Join(dir, f))
		if err != nil {
			return fmt.Errorf("read %s: %w", f, err)
		}
		if _, err := pool.Exec(ctx, string(content)); err != nil {
			return fmt.Errorf("exec %s: %w", f, err)
		}
	}
	return nil
}
