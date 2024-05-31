package gorm

import (
	"fmt"

	"gorm.io/gorm"
)

func CreatePaginationQuery(pageNumber, pageSize int, sortBy, sortDirection string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if pageNumber < 1 {
			pageNumber = 1
		}

		switch {
		case pageSize > 1000:
			pageSize = 1000
		case pageSize <= 0:
			pageSize = 10
		}

		if sortBy == "" {
			sortBy = "id"
		}

		if sortDirection == "" {
			sortDirection = "DESC"
		}

		offset := (pageNumber - 1) * pageSize
		return db.Offset(offset).Limit(pageSize).Order(fmt.Sprintf("%s %s", sortBy, sortDirection))
	}
}
