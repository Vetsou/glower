package database

import (
	"glower/database/model"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Handle *gorm.DB

func setupDatabase() {
	dsn := os.Getenv("DB_DSN")

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true,
	}), &gorm.Config{})

	if err != nil {
		log.Fatal("Error setting up database: " + err.Error())
	}

	conn, err := db.DB()
	if err != nil {
		log.Fatal("Error setting up database connection pool: " + err.Error())
	}

	conn.SetConnMaxIdleTime(5)
	conn.SetMaxIdleConns(10)
	conn.SetConnMaxLifetime(5 * time.Minute)

	Handle = db
}

func Init() {
	setupDatabase()

	err := Handle.AutoMigrate(&model.Flower{}, &model.Inventory{}, &model.User{}, &model.Cart{}, &model.CartItem{})
	if err != nil {
		log.Fatal("Error during DB auto migrate: " + err.Error())
	}
}
