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

	if options.Vin != "" {
		query = applyFilter("vin", options.Vin)
	}

	if options.LicensePlate != "" {
		query = applyFilter("license_plate", options.LicensePlate)
	}

	if options.Address != "" {
		query = applyFilter("address", options.Address)
	}

	if options.Phone != "" {
		query = applyFilter("phone", options.Phone)
	}

	if options.Social != "" {
		query = applyFilter("social", options.Social)
	}

	if options.CryptoCurrencyAddress != "" {
		query = applyFilter("cryptocurrency_address", options.CryptoCurrencyAddress)
	}

	if options.Domain != "" {
		query = applyFilter("url", options.Domain)
	}

	// Apply non-empty field filters
	for _, field := range options.NonEmptyFields {
		switch field {
		case "username":
			query = query.Where("JSON_ARRAY_LENGTH(username) > 0")
		case "email":
			query = query.Where("JSON_ARRAY_LENGTH(email) > 0")
		case "ip_address", "ipaddress", "ip":
			query = query.Where("JSON_ARRAY_LENGTH(ip_address) > 0")
		case "password":
			query = query.Where("JSON_ARRAY_LENGTH(password) > 0")
		case "hashed_password", "hash":
			query = query.Where("JSON_ARRAY_LENGTH(hashed_password) > 0")
		case "name":
			query = query.Where("JSON_ARRAY_LENGTH(name) > 0")
		case "vin":
			query = query.Where("JSON_ARRAY_LENGTH(vin) > 0")
		case "license_plate", "license":
			query = query.Where("JSON_ARRAY_LENGTH(license_plate) > 0")
		case "address":
			query = query.Where("JSON_ARRAY_LENGTH(address) > 0")
		case "phone":
			query = query.Where("JSON_ARRAY_LENGTH(phone) > 0")
		case "social":
			query = query.Where("JSON_ARRAY_LENGTH(social) > 0")
		case "cryptocurrency_address", "crypto":
			query = query.Where("JSON_ARRAY_LENGTH(cryptocurrency_address) > 0")
		case "url", "domain":
			query = query.Where("JSON_ARRAY_LENGTH(url) > 0")
		}
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
