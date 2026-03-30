package database

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	_ "github.com/lib/pq"
	"github.com/thienel/tlog"
	"go.uber.org/zap"

	"github.com/thienel/go-backend-template/internal/ent"
	"github.com/thienel/go-backend-template/pkg/config"
)

var (
	client *ent.Client
	sqlDB  *sql.DB
	once   sync.Once
)

// Init initializes the database connection
func Init(cfg *config.DatabaseConfig) error {
	var initErr error

	once.Do(func() {
		dsn := fmt.Sprintf(
			"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s TimeZone=%s",
			cfg.Host,
			cfg.Port,
			cfg.User,
			cfg.Password,
			cfg.DBName,
			cfg.SSLMode,
			cfg.TimeZone,
		)

		db, err := sql.Open("postgres", dsn)
		sqlDB = db
		if err != nil {
			initErr = fmt.Errorf("failed opening connection to postgres: %w", err)
			return
		}

		db.SetMaxIdleConns(getEnvInt("DB_MAX_IDLE_CONNS", 5))
		db.SetMaxOpenConns(getEnvInt("DB_MAX_OPEN_CONNS", 20))
		db.SetConnMaxLifetime(time.Minute * time.Duration(getEnvInt("DB_CONN_MAX_LIFETIME_MINUTES", 3)))
		db.SetConnMaxIdleTime(time.Minute * 2)

		if err := db.Ping(); err != nil {
			initErr = fmt.Errorf("failed pinging database: %w", err)
			return
		}

		drv := entsql.OpenDB(dialect.Postgres, db)
		client = ent.NewClient(ent.Driver(drv))

		tlog.Info("Database connection established",
			zap.String("host", cfg.Host),
			zap.Int("port", cfg.Port),
			zap.String("database", cfg.DBName),
		)
	})

	return initErr
}

// GetClient returns the ent client instance
func GetClient() *ent.Client {
	return client
}

// GetDB returns the underlying *sql.DB instance
func GetDB() *sql.DB {
	return sqlDB
}

// Close closes the database connection
func Close() error {
	if client != nil {
		return client.Close()
	}
	return nil
}

// WithTx runs callbacks in a transaction
func WithTx(ctx context.Context, fn func(tx *ent.Tx) error) error {
	if client == nil {
		return fmt.Errorf("database not initialized")
	}
	tx, err := client.Tx(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if r := recover(); r != nil {
			_ = tx.Rollback()
			panic(r)
		}
	}()

	if err := fn(tx); err != nil {
		_ = tx.Rollback()
		return err
	}
	return tx.Commit()
}

// getEnvInt reads an integer from an environment variable,
// returning defaultVal if the variable is absent or unparseable.
func getEnvInt(key string, defaultVal int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			return n
		}
	}
	return defaultVal
}
