package store

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"
)

type User struct {
	ID        int64
	LoginName string `gorm:"uniqueIndex"`
	Email     string `gorm:"uniqueIndex"`
	Salt      string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

var ErrUnknownUser = errors.New("Unknown user")

func genSalt() string {
	var salt [16]byte
	_, err := io.ReadFull(rand.Reader, salt[:])
	if err != nil {
		panic(fmt.Errorf("random read failed: %v", err))
	}

	return hex.EncodeToString(salt[:])
}

func hashSaltAndPassword(password string, salt string) string {
	hashfunc := sha256.New()
	hashfunc.Write([]byte(salt))
	hashfunc.Write([]byte(password))

	return hex.EncodeToString(hashfunc.Sum(nil))
}

func GetUserByLogin(loginName string) (*User, error) {
	var user User
	res := gDb.Take(&user, "login_name = ?", strings.ToLower(loginName))
	return &user, res.Error
}

func (u *User) CheckPassword(password string) bool {
	hashedPass := hashSaltAndPassword(password, u.Salt)

	if subtle.ConstantTimeCompare([]byte(u.Password), []byte(hashedPass)) == 1 {
		return true
	}

	return false
}

func CreateUser(u *User) error {
	u.LoginName = strings.TrimSpace(u.LoginName)
	u.LoginName = strings.ToLower(u.LoginName)
	u.Email = strings.TrimSpace(u.Email)
	u.Email = strings.ToLower(u.Email)
	u.Salt = genSalt()
	u.Password = hashSaltAndPassword(u.Password, u.Salt)

	res := gDb.Create(u)

	return res.Error
}

func DeleteUserByLogin(loginName string) error {
	loginName = strings.TrimSpace(loginName)
	loginName = strings.ToLower(loginName)

	res := gDb.Where("login_name = ?", loginName).Delete(&User{})
	if res.RowsAffected == 0 {
		return ErrUnknownUser
	}

	return res.Error
}
