package db

import (
	"flag"

	"github.com/Opsi/sparschwein/util"
)

// Config holds the database configuration values
type Config struct {
	// Host of the database
	Host string
	// Port of the database
	Port int
	// User to use for connecting to the database
	User string
	// Password of the user
	Password string
	// Database name
	Database string
}

// AddFlags adds the database flags and returns the configuration struct to
// be filled with the values from the flags.
func AddFlags() *Config {
	dbConfig := &Config{}

	flag.StringVar(
		&dbConfig.Host,
		"db-host",
		util.LookupStringEnv("DB_HOST", "localhost"),
		"Database host (default: localhost)")
	flag.IntVar(
		&dbConfig.Port,
		"db-port",
		util.LookupIntEnv("DB_PORT", 5432),
		"Database port (default: 5432)")
	flag.StringVar(
		&dbConfig.User,
		"db-user",
		util.LookupStringEnv("DB_USER", "postgres"),
		"Database user (default: postgres)")
	flag.StringVar(
		&dbConfig.Password,
		"db-password",
		util.LookupStringEnv("DB_PASSWORD", ""),
		"Database password")
	flag.StringVar(
		&dbConfig.Database,
		"db-name",
		util.LookupStringEnv("DB_NAME", "sparschwein"),
		"Database name (default: sparschwein)")

	return dbConfig
}
