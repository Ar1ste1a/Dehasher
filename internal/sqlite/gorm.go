package sqlite

import (
	"fmt"
	"go.uber.org/zap"
	"gorm.io/gorm/clause"
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
	if len(results.Results) == 0 {
		return nil
	}

	zap.L().Info("Storing results", zap.Int("count", len(results.Results)))
	db := GetDB()

	// Use batch insert with conflict handling
	const batchSize = 100
	var lastErr error

	// Extract the slice of results
	resultSlice := results.Results

	for i := 0; i < len(resultSlice); i += batchSize {
		end := i + batchSize
		if end > len(resultSlice) {
			end = len(resultSlice)
		}

		batch := resultSlice[i:end]
		// Use Clauses with OnConflict DoNothing to skip conflicts
		err := db.Clauses(clause.OnConflict{DoNothing: true}).CreateInBatches(&batch, batchSize).Error
		if err != nil {
			zap.L().Warn("Error storing some results", zap.Error(err))
			lastErr = err
			// Continue with next batch despite error
		}
	}

	return lastErr
}

func StoreCreds(creds []Creds) error {
	if len(creds) == 0 {
		return nil
	}

	zap.L().Info("Storing credentials", zap.Int("count", len(creds)))
	db := GetDB()

	// Use batch insert with conflict handling
	// This will insert records in batches and continue even if some fail
	const batchSize = 100
	var lastErr error

	for i := 0; i < len(creds); i += batchSize {
		end := i + batchSize
		if end > len(creds) {
			end = len(creds)
		}

		batch := creds[i:end]
		// Use Clauses with OnConflict DoNothing to skip conflicts
		err := db.Clauses(clause.OnConflict{DoNothing: true}).CreateInBatches(&batch, batchSize).Error
		if err != nil {
			zap.L().Warn("Error storing some credentials", zap.Error(err))
			lastErr = err
			// Continue with next batch despite error
		}
	}

	return lastErr
}
