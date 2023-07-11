package data

import (
	"brillinkmicros/pkg"
	"gorm.io/gorm"
)

func ApplyBlDataScope(dsi *pkg.DataScopeInfo) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if dsi.AccessType == pkg.DataScopeFullAccess {
			return db.Where("create_by IN (?)", dsi.AccessibleIds)
		}
		return db
	}
}
