package store

import (
	"errors"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var gDb *gorm.DB
var ErrEntryExists = errors.New("Entry already exists")
var ErrEntryMissing = errors.New("Entry missing")

func Init(file string, traceSql bool) error {
	if gDb != nil {
		panic("Database is already opened")
	}
	var config *gorm.Config
	if traceSql {
		config = &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		}
	} else {
		config = &gorm.Config{}
	}
	tmpdb, err := gorm.Open(sqlite.Open(file), config)

	if err != nil {
		return err
	}

	tmpdb.AutoMigrate(&E212Entry{}, &User{})
	gDb = tmpdb

	return nil
}
