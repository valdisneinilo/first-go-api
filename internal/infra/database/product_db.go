package database

import (
	"database/sql"
	"fmt"
	"strings"

	internalEntity "github.com/valdisneinilo/first-go-api/internal/entity"
)

type ProductDB struct {
	DB *sql.DB
}

func NewProductDB(db *sql.DB) *ProductDB {
	return &ProductDB{DB: db}
}

func (p *ProductDB) CreateProduct(product *internalEntity.Product) error {
	stmt, err := p.DB.Prepare("INSERT INTO PRODUTOS(ID, NOME, PRECO, CRIADO_EM) VALUES(?,?,?,?)")
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(product.Id, product.Name, product.Price, product.Created_at)

	if err != nil {
		return err
	}
	return nil
}

func (p *ProductDB) UpdateProduct(product *internalEntity.Product) error {
	stmt, err := p.DB.Prepare("UPDATE PRODUTOS SET NOME=?, PRECO=? WHERE ID=?")
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(product.Name, product.Price, product.Id)

	if err != nil {
		return err
	}
	return nil
}

func (p *ProductDB) DeleteProduct(id string) error {
	stmt, err := p.DB.Prepare("DELETE FROM PRODUTOS WHERE ID=?")
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(id)

	if err != nil {
		return err
	}
	return nil
}

func (p *ProductDB) FindProductById(id string) (*internalEntity.Product, error) {
	var product internalEntity.Product

	stmt, err := p.DB.Prepare("SELECT ID, NOME, PRECO, CRIADO_EM FROM PRODUTOS WHERE ID=?")
	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	err = stmt.QueryRow(id).Scan(&product.Id, &product.Name, &product.Price, &product.Created_at)
	if err != nil {
		return nil, err
	}

	return &product, nil
}

func (p *ProductDB) FindAllProducts(page, limit int, order, field string) ([]internalEntity.Product, error) {
	var rows *sql.Rows
	var err error
	var products []internalEntity.Product
	var query string
	order = strings.ToUpper(order)

	sortFieldMap := map[string]string{
		"name":       "NOME",
		"price":      "PRECO",
		"created_at": "CRIADO_EM",
	}

	column, ok := sortFieldMap[field]

	if order != "ASC" && order != "DESC" {
		order = "ASC"
	}

	if !ok {
		column = "CRIADO_EM"
	}

	if page > 0 && limit > 0 {
		offset := (page - 1) * limit
		query = fmt.Sprintf("SELECT ID, NOME, PRECO, CRIADO_EM FROM PRODUTOS ORDER BY %s %s LIMIT ? OFFSET ?", column, order)
		rows, err = p.DB.Query(query, limit, offset)

	} else {
		query = fmt.Sprintf("SELECT ID, NOME, PRECO, CRIADO_EM FROM PRODUTOS ORDER BY %s %s", column, order)
		rows, err = p.DB.Query(query)
	}

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var product internalEntity.Product

		err = rows.Scan(&product.Id, &product.Name, &product.Price, &product.Created_at)
		if err != nil {
			return nil, err
		}

		products = append(products, product)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	defer rows.Close()

	return products, nil
}
