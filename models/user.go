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

const UserSchema = `
CREATE TABLE user_acc (
  id serial,
  username varchar(255) UNIQUE NOT NULL,
  password varchar(255) NOT NULL,
  PRIMARY KEY (id)
);
`

type User struct {
	ID       int
	Username string
	Password string
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
	u.validateUsername()
	u.hashPassword()

	stmt := "INSERT INTO user_acc (username, password) VALUES (:username, :password)"
	_, err := db.NamedExec(stmt, u)
	return err
}

func (u *User) verifyPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}

func UserLogin(username string, password string) (*User, error) {
	user := User{}

	err := db.Get(&user, "SELECT * from user_acc WHERE username = $1", username)
	if err != nil {
		return nil, err
	}

	if err := user.verifyPassword(password); err != nil {
		return nil, err
	}

	return &user, nil
}
