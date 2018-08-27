package store

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/mattn/go-sqlite3"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB
var ErrEntryExists = errors.New("Entry already exists")
var ErrEntryMissing = errors.New("Entry missing")

type MccMnc struct {
	Mcc string `json:"mcc"`
	Mnc string `json:"nnc"`
}

func (v *MccMnc) String() string {
	return "MCC:" + v.Mcc + " MNC:" + v.Mcc
}

type MccMncEntry struct {
	MccMNC   MccMnc `json:"mccmnc"`
	Country  string `json:"country"`
	Operator string `json:"operator"`
}

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
	if db != nil {
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

	db = tmpdb
	return nil
}

func NewMccMncEntry(mcc, mnc, country, operator string) *MccMncEntry {

	return &MccMncEntry{
		MccMNC: MccMnc{
			Mcc: strings.TrimSpace(mcc),
			Mnc: strings.TrimSpace(mnc)},
		Country:  strings.TrimSpace(country),
		Operator: strings.TrimSpace(operator)}

}

func GetByMccMnc(mccMnc *MccMnc) (*MccMncEntry, error) {
	stmt, err := db.Prepare("SELECT MCC, MNC, COUNTRY, OPERATOR FROM E212 WHERE MCC=? AND MNC=?")
	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	rs, err := stmt.Query(mccMnc.Mcc, mccMnc.Mnc)
	if err != nil {
		return nil, err
	}

	var entry MccMncEntry
	if rs.Next() {
		err = rs.Scan(&entry.MccMNC.Mcc, &entry.MccMNC.Mnc, &entry.Country, &entry.Operator)
		rs.Close()
		if err != nil {
			return nil, err
		}
	} else {
		return nil, ErrEntryMissing
	}

	return &entry, nil

}

func GetAll() ([]MccMncEntry, error) {
	var entries []MccMncEntry

	stmt, err := db.Prepare("SELECT MCC, MNC, COUNTRY, OPERATOR FROM E212")
	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	rs, err := stmt.Query()
	if err != nil {
		return nil, err
	}

	for rs.Next() {
		var entry MccMncEntry

		err = rs.Scan(&entry.MccMNC.Mcc, &entry.MccMNC.Mnc, &entry.Country, &entry.Operator)
		entries = append(entries, entry)
		if err != nil {
			rs.Close()
			return entries, err
		}
	}

	return entries, err
}

func GetByMcc(mcc string) ([]MccMncEntry, error) {
	var entries []MccMncEntry

	stmt, err := db.Prepare("SELECT MCC, MNC, COUNTRY, OPERATOR FROM E212 WHERE MCC=?")
	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	rs, err := stmt.Query(mcc)
	if err != nil {
		return nil, err
	}

	for rs.Next() {
		var entry MccMncEntry

		err = rs.Scan(&entry.MccMNC.Mcc, &entry.MccMNC.Mnc, &entry.Country, &entry.Operator)
		entries = append(entries, entry)
		if err != nil {
			rs.Close()
			return entries, err
		}
	}

	return entries, err
}

func Add(e *MccMncEntry) error {

	stmt, err := db.Prepare("INSERT INTO E212(MCC, MNC, COUNTRY, OPERATOR) VALUES(?, ?, ?, ?)")
	if err != nil {
		return err
	}

	defer stmt.Close()
	_, err = stmt.Exec(e.MccMNC.Mcc, e.MccMNC.Mnc, e.Country, e.Operator)
	if s3err, ok := err.(sqlite3.Error); ok {
		if s3err.Code == sqlite3.ErrConstraint {
			return ErrEntryExists
		}
	}
	return err
}

func Remove(e *MccMnc) error {

	stmt, err := db.Prepare("DELETE FROM  E212 WHERE MCC=? and MNC=?")
	if err != nil {
		return err
	}

	defer stmt.Close()
	r, err := stmt.Exec(e.Mcc, e.Mnc)
	if err != nil {
		return err
	}
	if cnt, _ := r.RowsAffected(); cnt == 0 {
		return ErrEntryMissing
	}

	return err
}
