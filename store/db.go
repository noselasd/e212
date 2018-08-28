package store

import (
	"database/sql"
	"errors"

	//init sqlite2 driver
	_ "github.com/mattn/go-sqlite3"
)

var gDb *sql.DB
var ErrEntryExists = errors.New("Entry already exists")
var ErrEntryMissing = errors.New("Entry missing")

var initSQL = [...]string{
	`CREATE TABLE IF NOT EXISTS E212 (
		ID INTEGER PRIMARY KEY AUTOINCREMENT,
		MCC VARCHAR(3) NOT NULL,
		MNC VARCHAR(3) NOT NULL,
		COUNTRY VARCHAR(255) NOT NULL,
		OPERATOR VARCHAR(255) NOT NULL,
		CONSTRAINT MCCMNC_UNIQUE UNIQUE(MCC, MNC)
	) `,
}

func Init(file string) error {
	if gDb != nil {
		panic("Database is already opened")
	}

	tmpdb, err := sql.Open("sqlite3", file)
	if err != nil {
		return err
	}

	for i := range initSQL {
		_, err = tmpdb.Exec(initSQL[i])
		if err != nil {
			tmpdb.Close()
			return err
		}
	}

	gDb = tmpdb
	return nil
}
