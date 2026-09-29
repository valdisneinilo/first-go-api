package database

import (
	"database/sql"

	"github.com/valdisneinilo/first-go-api/internal/entity"
)

type UserDB struct {
	DB *sql.DB
}

func NewUserDB(db *sql.DB) *UserDB {
	return &UserDB{
		DB: db,
	}
}

func (u *UserDB) CreateUser(user *entity.User) error {
	stmt, err := u.DB.Prepare("INSERT INTO USUARIOS(ID, NOME, EMAIL, SENHA) VALUES(?, ?, ?, ?)")
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(user.Id, user.Name, user.Email, user.Password)

	if err != nil {
		return err
	}

	return nil
}

func (u *UserDB) FindUserByEmail(email string) (*entity.User, error) {
	var user entity.User
	stmt, err := u.DB.Prepare("SELECT ID, NOME, EMAIL, SENHA FROM USUARIOS WHERE EMAIL=?")
	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	err = stmt.QueryRow(email).Scan(&user.Id, &user.Name, &user.Email, &user.Password)
	if err != nil {
		return nil, err
	}

	return &user, nil

}
