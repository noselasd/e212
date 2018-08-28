package store

import (
	"strings"

	"github.com/mattn/go-sqlite3"
)

type MccMnc struct {
	Mcc string `json:"mcc"`
	Mnc string `json:"nnc"`
}

func (v *MccMnc) String() string {
	return "MCC:" + v.Mcc + " MNC:" + v.Mcc
}

type E212Entry struct {
	Code     MccMnc `json:"code"`
	Country  string `json:"country"`
	Operator string `json:"operator"`
}

func NewE212Entry(mcc, mnc, country, operator string) *E212Entry {

	return &E212Entry{
		Code: MccMnc{
			Mcc: strings.TrimSpace(mcc),
			Mnc: strings.TrimSpace(mnc)},
		Country:  strings.TrimSpace(country),
		Operator: strings.TrimSpace(operator)}

}

func E212GetByMccMnc(mccMnc *MccMnc) (*E212Entry, error) {
	stmt, err := gDb.Prepare("SELECT MCC, MNC, COUNTRY, OPERATOR FROM E212 WHERE MCC=? AND MNC=?")
	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	rs, err := stmt.Query(mccMnc.Mcc, mccMnc.Mnc)
	if err != nil {
		return nil, err
	}

	var entry E212Entry
	if rs.Next() {
		err = rs.Scan(&entry.Code.Mcc, &entry.Code.Mnc, &entry.Country, &entry.Operator)
		rs.Close()
		if err != nil {
			return nil, err
		}
	} else {
		return nil, ErrEntryMissing
	}

	return &entry, nil

}

func E212GetAll() ([]E212Entry, error) {
	var entries []E212Entry

	stmt, err := gDb.Prepare("SELECT MCC, MNC, COUNTRY, OPERATOR FROM E212")
	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	rs, err := stmt.Query()
	if err != nil {
		return nil, err
	}

	for rs.Next() {
		var entry E212Entry

		err = rs.Scan(&entry.Code.Mcc, &entry.Code.Mnc, &entry.Country, &entry.Operator)
		entries = append(entries, entry)
		if err != nil {
			rs.Close()
			return entries, err
		}
	}

	return entries, err
}

func E212GetByMcc(mcc string) ([]E212Entry, error) {
	var entries []E212Entry

	stmt, err := gDb.Prepare("SELECT MCC, MNC, COUNTRY, OPERATOR FROM E212 WHERE MCC=?")
	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	rs, err := stmt.Query(mcc)
	if err != nil {
		return nil, err
	}

	for rs.Next() {
		var entry E212Entry

		err = rs.Scan(&entry.Code.Mcc, &entry.Code.Mnc, &entry.Country, &entry.Operator)
		entries = append(entries, entry)
		if err != nil {
			rs.Close()
			return entries, err
		}
	}

	return entries, err
}

func E212Add(e *E212Entry) error {

	stmt, err := gDb.Prepare("INSERT INTO E212(MCC, MNC, COUNTRY, OPERATOR) VALUES(?, ?, ?, ?)")
	if err != nil {
		return err
	}

	defer stmt.Close()
	_, err = stmt.Exec(e.Code.Mcc, e.Code.Mnc, e.Country, e.Operator)
	if s3err, ok := err.(sqlite3.Error); ok {
		if s3err.Code == sqlite3.ErrConstraint {
			return ErrEntryExists
		}
	}
	return err
}

func E212Remove(e *MccMnc) error {

	stmt, err := gDb.Prepare("DELETE FROM  E212 WHERE MCC=? and MNC=?")
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
