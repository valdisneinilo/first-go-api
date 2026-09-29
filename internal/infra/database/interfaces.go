package database

import (
	internalEntity "github.com/valdisneinilo/first-go-api/internal/entity"
)

type DBUserInterface interface {
	CreateUser(user *internalEntity.User) error
	FindUserByEmail(email string) (*internalEntity.User, error)
}

type DBProductInterface interface {
	CreateProduct(product *internalEntity.Product) error
	UpdateProduct(product *internalEntity.Product) error
	DeleteProduct(id string) error
	FindProductById(id string) (*internalEntity.Product, error)
	FindAllProducts(page, limit int, order, field string) ([]internalEntity.Product, error)
}
