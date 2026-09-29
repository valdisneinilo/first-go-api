package entity

import (
	"errors"

	"github.com/valdisneinilo/first-go-api/pkg/entity"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id       entity.ID `json:"id"`
	Name     string    `json:"name"`
	Email    string    `json:"email"`
	Password string    `json:"-"`
}

var (
	ErrorUserInvalidId          = errors.New("Invalid ID")
	ErrorUserNameIsRequired     = errors.New("Name is required")
	ErrorUserEmailIsRequired    = errors.New("Email is required")
	ErrorUserPasswordIsRequired = errors.New("Password is required")
)

func NewUser(name, email, password string) (*User, error) {

	user := &User{
		Id:       entity.NewUUID(),
		Name:     name,
		Email:    email,
		Password: password,
	}

	err := user.UserValidate()
	if err != nil {
		return nil, err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user.Password = string(hash)
	return user, nil
}

func (user *User) ValidatePassword(password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	return err
}

func (u *User) UserValidate() error {

	if u.Name == "" {
		return ErrorUserNameIsRequired
	}

	if u.Email == "" {
		return ErrorUserEmailIsRequired
	}

	if u.Password == "" {
		return ErrorUserPasswordIsRequired
	}

	if _, err := entity.ParseID(u.Id.String()); err != nil {
		return ErrorUserInvalidId
	}

	return nil
}
