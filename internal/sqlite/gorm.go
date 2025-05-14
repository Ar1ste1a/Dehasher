package sqlite

import (
	"fmt"
	"go.uber.org/zap"
	"os"
	"path/filepath"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// InitDB initializes the database connection
func InitDB(dbDir string) (*gorm.DB, error) {
	zap.L().Info("Initializing database")

	// Create directory if it doesn't exist
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		zap.L().Error("Failed to create database directory", zap.Error(err))
		return nil, fmt.Errorf("failed to create database directory: %w", err)
	}

	dbPath := filepath.Join(dbDir, "dehashed.sqlite")
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		zap.L().Error("Failed to connect to database", zap.Error(err))
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Auto migrate your models
	err = db.AutoMigrate(&Result{}, &Creds{}, QueryOptions{}, Creds{})
	if err != nil {
		zap.L().Error("Failed to migrate database", zap.Error(err))
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	DB = db
	return db, nil
}

// GetDB returns the database connection
func GetDB() *gorm.DB {
	if DB == nil {
		zap.L().Error("database not initialized")
		fmt.Println("sqlite database not initialized")
		os.Exit(1)
	}
	return DB
}

func StoreResults(results DehashedResults) error {
	db := GetDB()
	return db.Create(&results).Error
}
