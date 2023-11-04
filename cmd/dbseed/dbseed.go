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

	dbConn, err := dbConfig.OpenPingedConnection()
	if err != nil {
		return fmt.Errorf("open database: %w", err)
	}
	defer dbConn.Close()

	slog.Info("successfully connected to database")

	_, err = dbConn.Exec(string(seedData))
	if err != nil {
		return fmt.Errorf("exec seed script: %w", err)
	}

	slog.Info("successfully seeded database")
	return nil
}
