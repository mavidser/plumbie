package models

import (
	"fmt"
	"regexp"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

var (
	usernamePattern = "^[a-z][a-z0-9_]*$"
)

type User struct {
	ID       int64  `xorm:"id pk autoincr"`
	Username string `xorm:"UNIQUE NOT NULL"`
	Password string `xorm:"NOT NULL"`
}

func (u *User) hashPassword() error {
	password, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	u.Password = string(password)
	return err
}

func (u *User) validateUsername() error {
	re := regexp.MustCompile("^[a-zA-Z0-9_]*$")
	if len(u.Username) > 32 {
		return fmt.Errorf("Username cannot be longer than 32 characters")
	}
	if re.MatchString(u.Username) {
		return nil
	}
	u.Username = strings.ToLower(u.Username)
	return fmt.Errorf("Username must start with a letter and can contain only letters, digits and underscores")
}

func (u *User) Create() error {
	sess := x.NewSession()
	defer sess.Close()
	if err := sess.Begin(); err != nil {
		return err
	}

	u.validateUsername()
	u.hashPassword()

	if _, err := sess.Insert(u); err != nil {
		return err
	}
	return sess.Commit()
}

func (u *User) verifyPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}

func UserLogin(username string, password string) (*User, error) {
	user := User{
		Username: strings.ToLower(username),
	}

	hasUser, err := x.Get(&user)
	if err != nil {
		return nil, err
	}
	if !hasUser {
		return nil, fmt.Errorf("User doesn't exist")
	}

	if err := user.verifyPassword(password); err != nil {
		return nil, err
	}

	return &user, nil
}
