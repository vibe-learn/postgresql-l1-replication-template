// Package main is the postgresql lesson `l1_replication_basics` homework scaffold for Vibe Learn.
//
// Задача: primary+2 standby: health-check, выбор кандидата по replay_lsn, программный failover.
// Реализуй функции ниже — сигнатуры и тестовая поверхность фиксированы;
// CI (.github/workflows/ci.yml) гоняет `go vet` и `go test ./...`.
// Подробности и критерии приёмки — в README.md.
//
// Драйвер: github.com/jackc/pgx/v5 (+ pgxpool). DATABASE_URL — DSN из env.
package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Latencies — собранные перцентили для бенчмарка запроса.
type Latencies struct{ P50, P95, P99 time.Duration }

// StandbyInfo — строка из pg_stat_replication для выбора кандидата на promote.
type StandbyInfo struct {
	ClientAddr string
	ReplayLSN  string
	State      string
}

// ----- config -----

// envOr returns the env var for `key` if set, else `fallback`.
func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// DatabaseURL — DSN PostgreSQL. Дефолт совпадает с docker-compose.yml.
func DatabaseURL() string {
	return envOr("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable")
}

// Connect открывает пул pgx из DATABASE_URL.
func Connect(ctx context.Context) (*pgxpool.Pool, error) {
	return pgxpool.New(ctx, DatabaseURL())
}

// ----- TODO #1: ReplicationStatus -----
//
// SELECT client_addr, replay_lsn, state FROM pg_stat_replication
func ReplicationStatus(ctx context.Context, primary *pgxpool.Pool) ([]StandbyInfo, error) {
	// TODO: implement
	panic("ReplicationStatus: not implemented")
}

// ----- TODO #2: PickCandidate -----
//
// чистая функция: выбрать standby с максимальным replay_lsn (наименьший lag)
func PickCandidate(standbys []StandbyInfo) (StandbyInfo, error) {
	// TODO: implement
	panic("PickCandidate: not implemented")
}

// ----- TODO #3: Promote -----
//
// pg_promote() на выбранном standby; дождаться pg_is_in_recovery()=false
func Promote(ctx context.Context, standby *pgxpool.Pool) error {
	// TODO: implement
	panic("Promote: not implemented")
}

// _refs keeps imports live while the TODO bodies are unimplemented stubs.
// Удали эту функцию, когда реализуешь TODO выше.
var _refs = []any{
	Latencies{},
	StandbyInfo{},
	time.Second,
}

// ----- main entry -----

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	log.Printf("Vibe Learn — postgresql lesson %s scaffold up", "l1_replication_basics")
	log.Printf("DATABASE_URL: %s", DatabaseURL())
	log.Printf("Реализуй TODO-функции, затем `go test ./...`. README.md содержит задачу.")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Graceful shutdown so `go run .` is interactive — Ctrl-C exits cleanly.
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigCh
		log.Printf("shutdown signal received")
		cancel()
	}()
	<-ctx.Done()
}
