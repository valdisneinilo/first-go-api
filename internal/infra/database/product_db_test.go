package database

import (
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	internalentity "github.com/valdisneinilo/first-go-api/internal/entity"
	pkgentity "github.com/valdisneinilo/first-go-api/pkg/entity"

	_ "modernc.org/sqlite"
)

func setupProductDB(t *testing.T) *sql.DB {
	t.Helper()

	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)

	_, err = db.Exec(`
		CREATE TABLE PRODUTOS (
			ID TEXT PRIMARY KEY,
			NOME TEXT NOT NULL,
			PRECO REAL NOT NULL,
			CRIADO_EM DATETIME NOT NULL
		)
	`)
	require.NoError(t, err)

	return db
}

func createTestProduct() *internalentity.Product {
	return &internalentity.Product{
		Id:         pkgentity.NewUUID(),
		Name:       "João",
		Price:      2500.00,
		Created_at: time.Now(),
	}
}

func TestCreateProduct(t *testing.T) {
	db := setupProductDB(t)
	defer db.Close()

	productDB := ProductDB{DB: db}

	product := createTestProduct()

	err := productDB.CreateProduct(product)

	require.NoError(t, err)

	var name string
	var price float64

	err = db.QueryRow(
		"SELECT NOME, PRECO FROM PRODUTOS WHERE ID=?",
		product.Id.String(),
	).Scan(&name, &price)

	require.NoError(t, err)
	require.Equal(t, product.Name, name)
	require.Equal(t, product.Price, price)
}

func TestCreateProductDuplicateID(t *testing.T) {
	db := setupProductDB(t)
	defer db.Close()

	productDB := ProductDB{DB: db}

	product := createTestProduct()

	err := productDB.CreateProduct(product)
	require.NoError(t, err)

	err = productDB.CreateProduct(product)

	require.Error(t, err)
}

func TestCreateProductDatabaseError(t *testing.T) {
	db := setupProductDB(t)
	db.Close()

	productDB := ProductDB{DB: db}

	product := createTestProduct()

	err := productDB.CreateProduct(product)

	require.Error(t, err)
}

func TestUpdateProduct(t *testing.T) {
	db := setupProductDB(t)
	defer db.Close()

	productDB := ProductDB{DB: db}

	product := createTestProduct()

	err := productDB.CreateProduct(product)
	require.NoError(t, err)

	product.Name = "Maria"
	product.Price = 3500.00

	err = productDB.UpdateProduct(product)

	require.NoError(t, err)

	result, err := productDB.FindProductById(product.Id.String())

	require.NoError(t, err)
	require.Equal(t, "Maria", result.Name)
	require.Equal(t, 3500.00, result.Price)
}

func TestUpdateProductNotFound(t *testing.T) {
	db := setupProductDB(t)
	defer db.Close()

	productDB := ProductDB{DB: db}

	product := createTestProduct()

	err := productDB.UpdateProduct(product)

	require.NoError(t, err)
}

func TestUpdateProductDatabaseError(t *testing.T) {
	db := setupProductDB(t)
	db.Close()

	productDB := ProductDB{DB: db}

	product := createTestProduct()

	err := productDB.UpdateProduct(product)

	require.Error(t, err)
}

func TestDeleteProduct(t *testing.T) {
	db := setupProductDB(t)
	defer db.Close()

	productDB := ProductDB{DB: db}

	product := createTestProduct()

	err := productDB.CreateProduct(product)
	require.NoError(t, err)

	err = productDB.DeleteProduct(product.Id.String())

	require.NoError(t, err)

	_, err = productDB.FindProductById(product.Id.String())

	require.ErrorIs(t, err, sql.ErrNoRows)
}

func TestDeleteProductNotFound(t *testing.T) {
	db := setupProductDB(t)
	defer db.Close()

	productDB := ProductDB{DB: db}

	err := productDB.DeleteProduct(pkgentity.NewUUID().String())

	require.NoError(t, err)
}

func TestDeleteProductDatabaseError(t *testing.T) {
	db := setupProductDB(t)
	db.Close()

	productDB := ProductDB{DB: db}

	err := productDB.DeleteProduct(pkgentity.NewUUID().String())

	require.Error(t, err)
}

func TestFindProductById(t *testing.T) {
	db := setupProductDB(t)
	defer db.Close()

	productDB := ProductDB{DB: db}

	product := createTestProduct()

	err := productDB.CreateProduct(product)
	require.NoError(t, err)

	result, err := productDB.FindProductById(product.Id.String())

	require.NoError(t, err)
	require.NotNil(t, result)

	require.Equal(t, product.Id.String(), result.Id.String())
	require.Equal(t, product.Name, result.Name)
	require.Equal(t, product.Price, result.Price)
	require.WithinDuration(
		t,
		product.Created_at,
		result.Created_at,
		time.Second,
	)
}

func TestFindProductByIdNotFound(t *testing.T) {
	db := setupProductDB(t)
	defer db.Close()

	productDB := ProductDB{DB: db}

	result, err := productDB.FindProductById(pkgentity.NewUUID().String())

	require.Nil(t, result)
	require.ErrorIs(t, err, sql.ErrNoRows)
}

func TestFindProductByIdDatabaseError(t *testing.T) {
	db := setupProductDB(t)
	db.Close()

	productDB := ProductDB{DB: db}

	result, err := productDB.FindProductById(pkgentity.NewUUID().String())

	require.Nil(t, result)
	require.Error(t, err)
}

func TestFindAllProducts(t *testing.T) {
	db := setupProductDB(t)
	defer db.Close()

	productDB := ProductDB{DB: db}

	products := []*internalentity.Product{
		{
			Id:         pkgentity.NewUUID(),
			Name:       "Produto A",
			Price:      100,
			Created_at: time.Now(),
		},
		{
			Id:         pkgentity.NewUUID(),
			Name:       "Produto B",
			Price:      200,
			Created_at: time.Now().Add(time.Second),
		},
		{
			Id:         pkgentity.NewUUID(),
			Name:       "Produto C",
			Price:      300,
			Created_at: time.Now().Add(2 * time.Second),
		},
	}

	for _, product := range products {
		err := productDB.CreateProduct(product)
		require.NoError(t, err)
	}

	result, err := productDB.FindAllProducts(0, 0, "ASC", "name")

	require.NoError(t, err)
	require.Len(t, result, 3)

	require.Equal(t, "Produto A", result[0].Name)
	require.Equal(t, "Produto B", result[1].Name)
	require.Equal(t, "Produto C", result[2].Name)
}

func TestFindAllProductsOrderByNameDESC(t *testing.T) {
	db := setupProductDB(t)
	defer db.Close()

	productDB := ProductDB{DB: db}

	products := []*internalentity.Product{
		{
			Id:         pkgentity.NewUUID(),
			Name:       "Ana",
			Price:      100,
			Created_at: time.Now(),
		},
		{
			Id:         pkgentity.NewUUID(),
			Name:       "Carlos",
			Price:      200,
			Created_at: time.Now(),
		},
		{
			Id:         pkgentity.NewUUID(),
			Name:       "Bruno",
			Price:      300,
			Created_at: time.Now(),
		},
	}

	for _, product := range products {
		err := productDB.CreateProduct(product)
		require.NoError(t, err)
	}

	result, err := productDB.FindAllProducts(0, 0, "DESC", "name")

	require.NoError(t, err)
	require.Len(t, result, 3)

	require.Equal(t, "Carlos", result[0].Name)
	require.Equal(t, "Bruno", result[1].Name)
	require.Equal(t, "Ana", result[2].Name)
}

func TestFindAllProductsOrderByPrice(t *testing.T) {
	db := setupProductDB(t)
	defer db.Close()

	productDB := ProductDB{DB: db}

	products := []*internalentity.Product{
		{
			Id:         pkgentity.NewUUID(),
			Name:       "Produto A",
			Price:      300,
			Created_at: time.Now(),
		},
		{
			Id:         pkgentity.NewUUID(),
			Name:       "Produto B",
			Price:      100,
			Created_at: time.Now(),
		},
		{
			Id:         pkgentity.NewUUID(),
			Name:       "Produto C",
			Price:      200,
			Created_at: time.Now(),
		},
	}

	for _, product := range products {
		err := productDB.CreateProduct(product)
		require.NoError(t, err)
	}

	result, err := productDB.FindAllProducts(0, 0, "ASC", "price")

	require.NoError(t, err)
	require.Len(t, result, 3)

	require.Equal(t, 100.0, result[0].Price)
	require.Equal(t, 200.0, result[1].Price)
	require.Equal(t, 300.0, result[2].Price)
}

func TestFindAllProductsInvalidOrder(t *testing.T) {
	db := setupProductDB(t)
	defer db.Close()

	productDB := ProductDB{DB: db}

	products := []*internalentity.Product{
		{
			Id:         pkgentity.NewUUID(),
			Name:       "Ana",
			Price:      100,
			Created_at: time.Now(),
		},
		{
			Id:         pkgentity.NewUUID(),
			Name:       "Bruno",
			Price:      200,
			Created_at: time.Now(),
		},
	}

	for _, product := range products {
		err := productDB.CreateProduct(product)
		require.NoError(t, err)
	}

	result, err := productDB.FindAllProducts(0, 0, "INVALID", "name")

	require.NoError(t, err)
	require.Len(t, result, 2)

	require.Equal(t, "Ana", result[0].Name)
	require.Equal(t, "Bruno", result[1].Name)
}

func TestFindAllProductsInvalidField(t *testing.T) {
	db := setupProductDB(t)
	defer db.Close()

	productDB := ProductDB{DB: db}

	products := []*internalentity.Product{
		{
			Id:         pkgentity.NewUUID(),
			Name:       "Produto A",
			Price:      100,
			Created_at: time.Now().Add(time.Second),
		},
		{
			Id:         pkgentity.NewUUID(),
			Name:       "Produto B",
			Price:      200,
			Created_at: time.Now(),
		},
	}

	for _, product := range products {
		err := productDB.CreateProduct(product)
		require.NoError(t, err)
	}

	result, err := productDB.FindAllProducts(0, 0, "ASC", "invalid_field")

	require.NoError(t, err)
	require.Len(t, result, 2)

	// Campo inválido deve cair no padrão CRIADO_EM
	require.Equal(t, "Produto B", result[0].Name)
	require.Equal(t, "Produto A", result[1].Name)
}

func TestFindAllProductsPagination(t *testing.T) {
	db := setupProductDB(t)
	defer db.Close()

	productDB := ProductDB{DB: db}

	for i := 1; i <= 5; i++ {
		product := &internalentity.Product{
			Id:         pkgentity.NewUUID(),
			Name:       "Produto " + string(rune('A'+i-1)),
			Price:      float64(i * 100),
			Created_at: time.Now().Add(time.Duration(i) * time.Second),
		}

		err := productDB.CreateProduct(product)
		require.NoError(t, err)
	}

	result, err := productDB.FindAllProducts(2, 2, "ASC", "price")

	require.NoError(t, err)
	require.Len(t, result, 2)

	require.Equal(t, 300.0, result[0].Price)
	require.Equal(t, 400.0, result[1].Price)
}

func TestFindAllProductsPageWithoutLimit(t *testing.T) {
	db := setupProductDB(t)
	defer db.Close()

	productDB := ProductDB{DB: db}

	for i := 1; i <= 3; i++ {
		product := &internalentity.Product{
			Id:         pkgentity.NewUUID(),
			Name:       "Produto",
			Price:      float64(i * 100),
			Created_at: time.Now(),
		}

		err := productDB.CreateProduct(product)
		require.NoError(t, err)
	}

	result, err := productDB.FindAllProducts(2, 0, "ASC", "price")

	require.NoError(t, err)

	// Como limit <= 0, não aplica paginação
	require.Len(t, result, 3)
}

func TestFindAllProductsEmpty(t *testing.T) {
	db := setupProductDB(t)
	defer db.Close()

	productDB := ProductDB{DB: db}

	result, err := productDB.FindAllProducts(0, 0, "ASC", "name")

	require.NoError(t, err)
	require.Empty(t, result)
}

func TestFindAllProductsDatabaseError(t *testing.T) {
	db := setupProductDB(t)
	db.Close()

	productDB := ProductDB{DB: db}

	result, err := productDB.FindAllProducts(0, 0, "ASC", "name")

	require.Nil(t, result)
	require.Error(t, err)
}
