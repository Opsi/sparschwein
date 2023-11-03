package db

import (
	"flag"
	"fmt"

	"github.com/Opsi/sparschwein/util"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
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

func (c Config) Open() (*sqlx.DB, error) {
	// Construct the connection string.
	// SSL mode 'disable' is not recommended for production use.
	dataSourceName := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		c.Host, c.Port, c.User, c.Password, c.Database)

	return sqlx.Open("postgres", dataSourceName)
}
