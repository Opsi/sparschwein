package main

import (
	"database/sql"
	"embed"
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/Opsi/sparschwein/db"
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

	// init and parse flags
	logConfig := util.AddLogFlags()
	dbConfig := db.AddFlags()
	flag.Parse()

	if err := logConfig.InitSlogDefault(); err != nil {
		return fmt.Errorf("init slog: %w", err)
	}

	seedData, err := seedFile.ReadFile("seed.sql")
	if err != nil {
		return fmt.Errorf("read seed file: %w", err)
	}

	// Construct the connection string.
	// SSL mode 'disable' is not recommended for production use.
	psqlInfo := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		dbConfig.Host, dbConfig.Port, dbConfig.User, dbConfig.Password, dbConfig.Database)

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
