package data

import (
	"brillinkmicros/common"
	"gorm.io/gorm"
)

func ApplyBlDataScope(dsi *common.DataScopeInfo) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if dsi.AccessType == common.DataScopeFullAccess {
			return db.Where("create_by IN (?)", dsi.AccessibleIds)
		}
		return db
	}
}
