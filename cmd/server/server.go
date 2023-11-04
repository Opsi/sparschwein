package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/Opsi/sparschwein/db"
	"github.com/Opsi/sparschwein/server"
	"github.com/Opsi/sparschwein/util"
	"github.com/joho/godotenv"
)

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
	_ = db.AddFlags()
	flag.Parse()

	if err := logConfig.InitSlogDefault(); err != nil {
		return fmt.Errorf("init slog: %w", err)
	}

	return server.ListenAndServe()
}
