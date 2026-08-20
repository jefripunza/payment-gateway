package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func main() {
	// .env loaded at the very start; real env vars take precedence
	_ = godotenv.Load()
	log.Println("payment-gateway starting")

	// --- database ---
	dbPath := envOr("PAYMENT_DB_PATH", "data/payment.db")
	if err := os.MkdirAll(dirOf(dbPath), 0o755); err != nil {
		log.Fatalf("mkdir: %v", err)
	}
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		log.Fatalf("open db: %v", err)
	}
	DB = db
	if err := migrate(); err != nil {
		log.Fatalf("migrate: %v", err)
	}
	seedAdmin()

	// --- http server ---
	addr := envOr("PAYMENT_BE", ":3005")
	app := newApp()
	log.Printf("listening on %s", addr)
	if err := app.Listen(addr); err != nil {
		log.Fatalf("listen: %v", err)
	}
}

func envOr(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func dirOf(p string) string {
	for i := len(p) - 1; i >= 0; i-- {
		if p[i] == '/' {
			return p[:i]
		}
	}
	return "."
}
