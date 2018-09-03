package store

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"strings"
)

type User struct {
	ID        int64
	LoginName string
	Email     string
	Salt      string
	Password  string
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

	stmt, err := gDb.Prepare("SELECT ID, LOGINNAME, EMAIL, SALT, PASSWORD FROM USER WHERE LOGINNAME=?")
	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	rs, err := stmt.Query(loginName)
	if err != nil {
		return nil, err
	}

	var user User
	if rs.Next() {
		err = rs.Scan(&user.ID, &user.LoginName, &user.Email, user.Salt, &user.Password)
		rs.Close()
		if err != nil {
			return nil, err
		}
	} else {
		return nil, ErrUnknownUser
	}

	return &user, nil
}

func (u *User) CheckPassword(password string) bool {
	hashedPass := hashSaltAndPassword(password, u.Salt)
	if hashedPass == u.Password {
		return true
	}

	return false
}

func CreateUser(u *User) error {
	u.LoginName = strings.TrimSpace(u.LoginName)
	u.Email = strings.TrimSpace(u.Email)
	u.Salt = genSalt()
	u.Password = hashSaltAndPassword(u.Password, u.Salt)

	stmt, err := gDb.Prepare("INSERT INTO USER(LOGINNAME, EMAIL, SALT, PASSWORD) VALUES(?, ?, ?, ?)")
	if err != nil {
		return err
	}

	defer stmt.Close()

	res, err := stmt.Exec(u.LoginName, u.Email, u.Salt, u.Password)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	u.ID = id

	return nil
}
