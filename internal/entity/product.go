package entity

import (
	"errors"
	"time"

	"github.com/valdisneinilo/first-go-api/pkg/entity"
)

type Product struct {
	Id         entity.ID `json:"id"`
	Name       string    `json:"name"`
	Price      float64   `json:"price"`
	Created_at time.Time `json:"created_at"`
}

var (
	ErrorProductInvalidId       = errors.New("Id is not null")
	ErrorProductNameIsRequired  = errors.New("Name is required")
	ErrorProductPriceIsRequired = errors.New("Price is required")
	ErrorProductInvalidPrice    = errors.New("Invalid price")
)

func NewProduct(name string, price float64) (*Product, error) {
	product := &Product{
		Id:         entity.NewUUID(),
		Name:       name,
		Price:      price,
		Created_at: time.Now(),
	}

	err := product.ProductValidation()
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (p *Product) ProductValidation() error {
	if p.Name == "" {
		return ErrorProductNameIsRequired
	}

	if p.Price == 0 {
		return ErrorProductPriceIsRequired
	}

	if p.Price < 0 {
		return ErrorProductInvalidPrice
	}

	if _, err := entity.ParseID(p.Id.String()); err != nil {
		return ErrorProductInvalidId
	}

	return nil
}
