package sqlite

import (
	"fmt"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// QueryResults queries the database for results based on the provided options
func QueryResults(options *DBOptions) ([]Result, error) {
	db := GetDB()
	var results []Result
	query := db.Model(&Result{})

	// Apply filters based on the provided options
	query = applyFilters(query, options)

	// Apply limit
	if options.Limit > 0 {
		query = query.Limit(options.Limit)
	}

	// Execute the query
	if err := query.Find(&results).Error; err != nil {
		zap.L().Error("query_results",
			zap.String("message", "failed to query results"),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to query results: %w", err)
	}

	return results, nil
}

// applyFilters applies filters to the query based on the provided options
func applyFilters(query *gorm.DB, options *DBOptions) *gorm.DB {
	// Helper function to apply filter based on exact match setting
	applyFilter := func(field, value string) *gorm.DB {
		if value == "" {
			return query
		}

		if options.ExactMatch {
			return query.Where(field+" = ?", value)
		} else {
			return query.Where(field+" LIKE ?", "%"+value+"%")
		}
	}

	// Apply filters for each field if provided
	if options.Email != "" {
		query = applyFilter("email", options.Email)
	}

	if options.Username != "" {
		query = applyFilter("username", options.Username)
	}

	if options.IPAddress != "" {
		query = applyFilter("ip_address", options.IPAddress)
	}

	if options.Password != "" {
		query = applyFilter("password", options.Password)
	}

	if options.HashedPassword != "" {
		query = applyFilter("hashed_password", options.HashedPassword)
	}

	if options.Name != "" {
		query = applyFilter("name", options.Name)
	}

	return query
}

// GetResultsCount returns the count of results matching the provided options
func GetResultsCount(options *DBOptions) (int64, error) {
	db := GetDB()
	var count int64
	query := db.Model(&Result{})

	// Apply filters based on the provided options
	query = applyFilters(query, options)

	// Count the results
	if err := query.Count(&count).Error; err != nil {
		zap.L().Error("get_results_count",
			zap.String("message", "failed to count results"),
			zap.Error(err),
		)
		return 0, fmt.Errorf("failed to count results: %w", err)
	}

	return count, nil
}
