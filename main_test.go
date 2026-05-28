package main

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"
)

// TestDatabaseURLDefault — pure unit test. No PostgreSQL required.
func TestDatabaseURLDefault(t *testing.T) {
	if os.Getenv("DATABASE_URL") == "" && !strings.HasPrefix(DatabaseURL(), "postgres://") {
		t.Errorf("DatabaseURL() default = %q, want postgres:// DSN", DatabaseURL())
	}
}

// TestEnvOverride — env должен перекрывать дефолт.
func TestEnvOverride(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://u:p@h:5432/db")
	if DatabaseURL() != "postgres://u:p@h:5432/db" {
		t.Errorf("DatabaseURL() = %q, env override ignored", DatabaseURL())
	}
}

// TestIntegration — требует запущенный PostgreSQL (docker compose up -d).
// SKIPPED по умолчанию; убери skip через PG_INTEGRATION=1.
func TestIntegration(t *testing.T) {
	if os.Getenv("PG_INTEGRATION") == "" {
		t.Skip("set PG_INTEGRATION=1 and run `docker compose up -d` to enable")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	pool, err := Connect(ctx)
	if err != nil {
		t.Fatalf("connect failed: %v", err)
	}
	defer pool.Close()
	if err := pool.Ping(ctx); err != nil {
		t.Fatalf("ping failed: %v", err)
	}
	// TODO: вызови свои реализованные функции и проверь поведение урока «Streaming replication, WAL и failover».
}
