package database

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	internalentity "github.com/valdisneinilo/first-go-api/internal/entity"
	pkgentity "github.com/valdisneinilo/first-go-api/pkg/entity"
	_ "modernc.org/sqlite"
)

func setupUserDB(t *testing.T) *sql.DB {
	t.Helper()

	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)

	_, err = db.Exec(`
		CREATE TABLE USUARIOS (
			ID TEXT PRIMARY KEY,
			NOME TEXT NOT NULL,
			EMAIL TEXT NOT NULL,
			SENHA TEXT NOT NULL
		)
	`)
	require.NoError(t, err)

	return db
}

func TestCreateUser(t *testing.T) {
	db := setupUserDB(t)
	defer db.Close()

	userDB := UserDB{
		DB: db,
	}

	user := &internalentity.User{
		Id:       pkgentity.NewUUID(),
		Name:     "João",
		Email:    "joao@email.com",
		Password: "123456",
	}

	err := userDB.CreateUser(user)

	require.NoError(t, err)
}

func TestFindUserByEmail(t *testing.T) {
	db := setupUserDB(t)
	defer db.Close()

	userDB := UserDB{
		DB: db,
	}

	user := &internalentity.User{
		Id:       pkgentity.NewUUID(),
		Name:     "João",
		Email:    "joao@email.com",
		Password: "123456",
	}

	err := userDB.CreateUser(user)
	require.NoError(t, err)

	foundUser, err := userDB.FindUserByEmail("joao@email.com")
	require.NoError(t, err)

	assert.Equal(t, user.Id, foundUser.Id)
	assert.Equal(t, user.Name, foundUser.Name)
	assert.Equal(t, user.Email, foundUser.Email)
}

func TestFindUserByEmailNotFound(t *testing.T) {
	db := setupUserDB(t)
	defer db.Close()

	userDB := UserDB{
		DB: db,
	}

	_, err := userDB.FindUserByEmail("naoexiste@email.com")

	assert.Error(t, err)
}
