package main

import (
	"fmt"
	"log"
	"os"

	"github.com/glebarez/sqlite"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// openDB connects using the DATABASE_PROVIDER env var:
//   - "sqlite" (default): no extra credentials needed, uses DATABASE_PATH (default data/payment.db)
//   - "mysql" | "postgres" | "sqlserver": requires DATABASE_HOST, DATABASE_PORT,
//     DATABASE_USER, DATABASE_PASSWORD, DATABASE_NAME
func openDB() (*gorm.DB, error) {
	provider := envOr("DATABASE_PROVIDER", "sqlite")
	cfg := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	}

	switch provider {
	case "sqlite":
		path := envOr("DATABASE_PATH", "data/payment.db")
		if err := os.MkdirAll(dirOf(path), 0o755); err != nil {
			return nil, fmt.Errorf("mkdir: %w", err)
		}
		return gorm.Open(sqlite.Open(path), cfg)

	case "mysql":
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			requiredEnv("DATABASE_USER"), requiredEnv("DATABASE_PASSWORD"),
			requiredEnv("DATABASE_HOST"), envOr("DATABASE_PORT", "3306"),
			requiredEnv("DATABASE_NAME"))
		return gorm.Open(mysql.Open(dsn), cfg)

	case "postgres":
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
			requiredEnv("DATABASE_HOST"), requiredEnv("DATABASE_USER"),
			requiredEnv("DATABASE_PASSWORD"), requiredEnv("DATABASE_NAME"),
			envOr("DATABASE_PORT", "5432"))
		return gorm.Open(postgres.Open(dsn), cfg)

	case "sqlserver":
		dsn := fmt.Sprintf("sqlserver://%s:%s@%s:%s?database=%s",
			requiredEnv("DATABASE_USER"), requiredEnv("DATABASE_PASSWORD"),
			requiredEnv("DATABASE_HOST"), envOr("DATABASE_PORT", "1433"),
			requiredEnv("DATABASE_NAME"))
		return gorm.Open(sqlserver.Open(dsn), cfg)

	default:
		return nil, fmt.Errorf("unsupported DATABASE_PROVIDER %q (use sqlite, mysql, postgres, or sqlserver)", provider)
	}
}

// migrate runs GORM AutoMigrate ONLY when DATABASE_MIGRATE=true.
// Default is false — without the flag the app will fail loudly if tables are missing,
// so manual schema management (or one-off migrations) is the safe path.
func migrate() error {
	if envOr("DATABASE_MIGRATE", "") != "true" {
		log.Println("migrate: skipped (set DATABASE_MIGRATE=true to run AutoMigrate)")
		return nil
	}
	log.Println("migrate: running AutoMigrate ...")
	return DB.AutoMigrate(&User{}, &Wallet{}, &Provider{})
}

func envOr(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func requiredEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("missing required env var %s for DATABASE_PROVIDER=%s", key, envOr("DATABASE_PROVIDER", "sqlite"))
	}
	return v
}

func dirOf(p string) string {
	for i := len(p) - 1; i >= 0; i-- {
		if p[i] == '/' {
			return p[:i]
		}
	}
	return "."
}
