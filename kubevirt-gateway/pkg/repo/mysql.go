package repo

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// InitDBConfig Init DB info, dataSource: test.db
func NewMysqlClient(dataSource string) (*gorm.DB, error) {
	// db, err := gorm.Open("mysql", dataSource)
	db, err := gorm.Open(mysql.Open(dataSource), &gorm.Config{DisableForeignKeyConstraintWhenMigrating: true})
	if err != nil {
		return nil, fmt.Errorf("open Mysql datasource err, %v", err)
	}
	return db, nil
}
