package database

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql" // for the mysql connector
	"github.com/jmoiron/sqlx"
)

// MySQLConnection func for connection to MySQL database.
func MySQLConnection() (*sqlx.DB, error) {
	maxConn, _ := strconv.Atoi(os.Getenv("FGONBOARD_DB_MAX_CONNECTIONS"))
	maxIdleConn, _ := strconv.Atoi(os.Getenv("FGONBOARD_DB_MAX_IDLE_CONNECTIONS"))
	maxLifetimeConn, _ := strconv.Atoi(os.Getenv("FGONBOARD_DB_MAX_LIFETIME_CONNECTIONS"))

	databaseURL := os.Getenv("FGONBOARD_DB_SERVER_URL")
	databaseURL = strings.ReplaceAll(databaseURL, "\"", "")

	db, err := sqlx.Connect("mysql", databaseURL+"?charset=utf8mb4&collation=utf8mb4_unicode_ci")
	if err != nil {
		return nil, fmt.Errorf("error, not connected to database, %w", err)
	}

	// MaxOpenConnection: the default is 0 (unlimited)
	// MaxIdleConns: defaultMaxIdleConns = 2
	// ConnMaxLifetime: 0, connections are reused forever

	db.SetMaxOpenConns(maxConn)
	db.SetMaxIdleConns(maxIdleConn)
	db.SetConnMaxLifetime(time.Duration(maxLifetimeConn))

	// Try to ping database.
	if errPing := db.Ping(); errPing != nil {
		defer db.Close() // close database connection
		return nil, fmt.Errorf("error, not sent ping to database, %w", errPing)
	}

	return db, nil
}
