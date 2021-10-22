package store

import (
	"errors"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var gDb *gorm.DB
var ErrEntryExists = errors.New("Entry already exists")
var ErrEntryMissing = errors.New("Entry missing")

func Init(file string) error {
	if gDb != nil {
		panic("Database is already opened")
	}
	tmpdb, err := gorm.Open(sqlite.Open(file), &gorm.Config{})
	if err != nil {
		return err
	}
	tmpdb.AutoMigrate(&E212Entry{}, &User{})
	gDb = tmpdb

	return nil
}
