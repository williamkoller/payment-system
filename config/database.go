package config

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"

	migratePG "github.com/golang-migrate/migrate/v4/database/postgres"
)

func NewDatabaseConnection() *gorm.DB {
	config, err := LoadDatabaseConfiguration()
	if err != nil {
		panic(err)
	}

	sqlInfo := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.Username, config.Password, config.Database,
	)

	dbLogger := gormLogger.New(
		log.New(os.Stdout, "", log.LstdFlags),
		gormLogger.Config{IgnoreRecordNotFoundError: true},
	)

	db, err := gorm.Open(postgres.Open(sqlInfo), &gorm.Config{
		TranslateError: true,
		Logger:         dbLogger,
	})

	log.Println("✅ Connected to DB")
	return db
}

func RunMigrations(gormDB *gorm.DB, migrationsPath string) {
	sqlDB, err := gormDB.DB()
	if err != nil {
		log.Fatal("failed to get raw database from GORM: ", err)
	}

	driver, err := migratePG.WithInstance(sqlDB, &migratePG.Config{})
	if err != nil {
		log.Fatal("migration WithInstance failed: ", err)
	}

	if migrationsPath == "" {
		wd, err := os.Getwd()
		if err != nil {
			log.Fatal("could not get working directory: ", err)
		}
		migrationsPath = "file://" + filepath.Join(wd, "db/migrations")
	}

	m, err := migrate.NewWithDatabaseInstance(migrationsPath, "postgres", driver)
	if err != nil {
		log.Fatal("migration NewWithDatabaseInstance failed: ", err)
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Println("✅ No new migrations to apply.")
			return
		}
		log.Fatal("migration Up failed: ", err)
	}

	log.Println("✅ Migrations completed successfully.")
}
