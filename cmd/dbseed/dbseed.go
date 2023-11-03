package main

import (
	"database/sql"
	"embed"
	"fmt"
	"log/slog"
	"os"

	"github.com/Opsi/sparschwein/util"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq" // Import the PostgreSQL driver
)

//go:embed seed.sql
var seedFile embed.FS

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run() error {
	if err := godotenv.Load(); err != nil {
		return fmt.Errorf("load .env: %w", err)
	}
	if err := util.InitSlogDefault(); err != nil {
		return fmt.Errorf("init slog default: %w", err)
	}
	var (
		host     = util.LookupStringEnv("DB_HOST", "localhost")
		port     = util.LookupIntEnv("DB_PORT", 5432)
		user     = util.LookupStringEnv("DB_USER", "postgres")
		password = util.LookupStringEnv("DB_PASSWORD", "postgres")
		dbname   = util.LookupStringEnv("DB_NAME", "sparschwein")
	)

	seedData, err := seedFile.ReadFile("seed.sql")
	if err != nil {
		return fmt.Errorf("read seed file: %w", err)
	}

	// Construct the connection string.
	// SSL mode 'disable' is not recommended for production use.
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	// Open the connection
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return fmt.Errorf("open connection: %w", err)
	}
	defer db.Close()

	// Check the connection
	err = db.Ping()
	if err != nil {
		return fmt.Errorf("ping db: %w", err)
	}

	slog.Info("successfully connected")

	_, err = db.Exec(string(seedData))
	if err != nil {
		return fmt.Errorf("exec seed script: %w", err)
	}

	slog.Info("successfully seeded")
	return nil
}
