package store

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

type MccMnc struct {
	Mcc string `json:"mcc" gorm:"index:idx_mccmnc,unique"`
	Mnc string `json:"mnc" gorm:"index:idx_mccmnc,unique"`
}

func (v *MccMnc) String() string {
	return "MCC:" + v.Mcc + " MNC:" + v.Mcc
}

type E212Country struct {
	ID        int
	Name      string    `gorm:"uniqueIndex"`
	CreatedAt time.Time `gorm:"<-:create"`
	UpdatedAt time.Time
}

type E212Entry struct {
	ID   int
	Code MccMnc `json:"code" gorm:"embedded"`
	//Country   string `json:"country"`
	E212CountryID int
	E212Country   E212Country
	Operator      string    `json:"operator"`
	CreatedAt     time.Time `gorm:"<-:create"`
	UpdatedAt     time.Time
}

func isAsciiDigits(s string) bool {
	const digits = "0123456789"
	for _, r := range s {
		if !strings.ContainsRune(digits, r) {
			return false
		}
	}
	return len(s) > 0
}
func (e *E212Entry) Validate() error {
	if len(e.E212Country.Name) == 0 {
		return errors.New("Invalid Country")
	}
	if len(e.Operator) == 0 {
		return errors.New("Invalid Operator")
	}
	if !isAsciiDigits(e.Code.Mcc) {
		return errors.New("Invalid MCC")
	}
	if !isAsciiDigits(e.Code.Mnc) {
		return errors.New("Invalid Mnc")
	}

	return nil
}

func NewE212Entry(mcc, mnc, country, operator string) *E212Entry {

	return &E212Entry{
		Code: MccMnc{
			Mcc: strings.TrimSpace(mcc),
			Mnc: strings.TrimSpace(mnc)},
		E212Country: E212Country{Name: strings.TrimSpace(country)},
		Operator:    strings.TrimSpace(operator)}

}

func E212GetByMccMnc(mccMnc *MccMnc) (*E212Entry, error) {
	var entry E212Entry
	res := gDb.Joins("E212Country").Take(&entry, "mcc = ? AND mnc = ?", mccMnc.Mcc, mccMnc.Mnc)

	return &entry, res.Error

}

func E212GetAll() ([]E212Entry, error) {
	var entries []E212Entry
	res := gDb.Joins("E212Country").Find(&entries)
	fmt.Printf("%s\n", entries[0].E212Country.Name)

	return entries, res.Error
}

func E212GetByMcc(mcc string) ([]E212Entry, error) {
	var entries []E212Entry

	res := gDb.Joins("E212Country").Find(&entries, "MCC=?", mcc)
	if len(entries) == 0 {
		return nil, ErrEntryMissing
	}

	return entries, res.Error

}

func E212Add(e *E212Entry) error {
	country, err := E212GetCountryByName(e.E212Country.Name)
	if err == nil {
		e.E212Country = *country
		e.E212CountryID = country.ID
	}

	res := gDb.Create(e)

	return res.Error
}

func E212Update(e *E212Entry) error {
	if e.ID == 0 {
		return ErrEntryMissing
	}
	country, err := E212GetCountryByName(e.E212Country.Name)
	if err == nil {
		e.E212Country = *country
		e.E212CountryID = country.ID
	}

	res := gDb.Save(e)

	if res.RowsAffected == 0 {
		return ErrEntryMissing
	}

	return res.Error
}

func E212DeleteById(id int) error {
	if id == 0 {
		return ErrEntryMissing
	}

	res := gDb.Delete(&E212Entry{}, id)

	if res.RowsAffected == 0 {
		return ErrEntryMissing
	}

	return res.Error
}

func E212Remove(e *MccMnc) error {
	res := gDb.Where("mcc = ? and mnc = ?", e.Mcc, e.Mnc).Delete(&E212Entry{})
	if res.RowsAffected == 0 {
		return ErrEntryMissing
	}
	return res.Error
}

func E212GetCountryByName(name string) (*E212Country, error) {
	var entry E212Country
	res := gDb.Take(&entry, "name = ?", name)

	return &entry, res.Error

}
