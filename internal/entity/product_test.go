package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var name = "Notebook"
var price = 2500.00

func TestNewProduct(t *testing.T) {

	p, err := NewProduct(name, price)

	assert.Nil(t, err)
	assert.NotNil(t, p)

	assert.NotEmpty(t, p.Id)
	assert.NotEmpty(t, p.Name)
	assert.NotEmpty(t, p.Price)
	assert.NotEmpty(t, p.Created_at)

	assert.Equal(t, p.Name, name)
	assert.Equal(t, p.Price, price)
}

func TestProductValidation(t *testing.T) {
	p, err := NewProduct(name, price)

	assert.Nil(t, err)
	assert.NotNil(t, p)
	assert.Nil(t, p.ProductValidation())
}

func TestProductNameRequired(t *testing.T) {
	p, err := NewProduct("", price)

	assert.Nil(t, p)
	assert.Equal(t, ErrorProductNameIsRequired, err)
}

func TestProductPriceRequired(t *testing.T) {
	p, err := NewProduct(name, 0)

	assert.Nil(t, p)
	assert.Equal(t, ErrorProductPriceIsRequired, err)
}

func TestProductInvalidPrice(t *testing.T) {
	p, err := NewProduct(name, -10)

	assert.Nil(t, p)
	assert.Equal(t, ErrorProductInvalidPrice, err)
}
