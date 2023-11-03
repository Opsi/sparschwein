package main

import (
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

	dbConn, err := dbConfig.Open()
	if err != nil {
		return fmt.Errorf("open connection: %w", err)
	}
	defer dbConn.Close()

	// Check the connection
	err = dbConn.Ping()
	if err != nil {
		return fmt.Errorf("ping db: %w", err)
	}

	slog.Info("successfully connected")

	_, err = dbConn.Exec(string(seedData))
	if err != nil {
		return fmt.Errorf("exec seed script: %w", err)
	}

	slog.Info("successfully seeded")
	return nil
}
