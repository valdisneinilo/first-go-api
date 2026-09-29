package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewUser(t *testing.T) {
	user, err := NewUser("Valdisnei", "valdisneinilo@gmail.com", "abc123")

	assert.Nil(t, err)
	assert.NotNil(t, user)
	assert.NotEmpty(t, user.Id)
	assert.NotEmpty(t, user.Password)
	assert.Equal(t, "Valdisnei", user.Name)
	assert.Equal(t, "valdisneinilo@gmail.com", user.Email)
}

func TestValidatePassword(t *testing.T) {
	user, err := NewUser("Valdisnei", "valdisneinilo@gmail.com", "abc123")
	assert.Nil(t, err)

	err = user.ValidatePassword("abc123")
	assert.Nil(t, err)

	err = user.ValidatePassword("abc")
	assert.Error(t, err)

	assert.NotEqual(t, user.Password, "abc123")
}

func TestValidate(t *testing.T) {
	user, err := NewUser("Valdisnei", "valdisneinilo@gmail.com", "abc123")

	assert.Nil(t, err)
	assert.NotNil(t, user)
	assert.Nil(t, user.UserValidate())
}

func TestValidateUserNameRequired(t *testing.T) {
	user, err := NewUser("", "valdisneinilo@gmail.com", "abc123")

	assert.Nil(t, user)
	assert.Equal(t, ErrorUserNameIsRequired, err)
}

func TestValidateUserEmailRequired(t *testing.T) {
	user, err := NewUser("Valdisnei", "", "abc123")

	assert.Nil(t, user)
	assert.Equal(t, ErrorUserEmailIsRequired, err)
}

func TestValidateUserPasswordRequired(t *testing.T) {
	user, err := NewUser("Valdisnei", "valdisneinilo@gmail.com", "")

	assert.Nil(t, user)
	assert.Equal(t, ErrorUserPasswordIsRequired, err)
}
