package main

import (
	"log"

	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

var DB *gorm.DB

func main() {
	_ = godotenv.Load()
	log.Println("payment-gateway starting")

	// --- database (DATABASE_PROVIDER driven) ---
	db, err := openDB()
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
