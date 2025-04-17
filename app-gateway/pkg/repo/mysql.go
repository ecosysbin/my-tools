package repo

import (
	"github.com/go-errors/errors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// InitDBConfig Init DB info, dataSource: test.db
func NewMysqlClient(dataSource string) (*gorm.DB, error) {
	// db, err := gorm.Open("mysql", dataSource)
	db, err := gorm.Open(mysql.Open(dataSource), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		Logger:                                   logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, errors.Errorf("open Mysql datasource err, %v", err)
	}
	return db, nil
}
