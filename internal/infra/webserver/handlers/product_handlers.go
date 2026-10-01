package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/valdisneinilo/first-go-api/internal/dto"
	"github.com/valdisneinilo/first-go-api/internal/entity"
	entityPkg "github.com/valdisneinilo/first-go-api/pkg/entity"

	"github.com/valdisneinilo/first-go-api/internal/infra/database"
)

type ProductHandler struct {
	ProductDB database.DBProductInterface
}

func NewProductHandler(db database.DBProductInterface) *ProductHandler {
	return &ProductHandler{
		ProductDB: db,
	}
}

// Create Product godoc
// @Summary      Create Product
// @Description  Creates a new product
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        request body dto.CreateProductInput true "Product data"
// @Success      201
// @Failure      400 {object} dto.Error
// @Failure      500 {object} dto.Error
// @Router       /products [post]
// @Security     ApiKeyAuth
func (h ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {

	var product dto.CreateProductInput
	err := json.NewDecoder(r.Body).Decode(&product)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad Request: \n" + err.Error()))
		return
	}

	p, err := entity.NewProduct(product.Name, product.Price)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad Request: \n" + err.Error()))
		return
	}

	err = h.ProductDB.CreateProduct(p)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Erro interno no servidor: \n" + err.Error()))
		return
	}

	w.WriteHeader(http.StatusCreated)

}

// Get Product godoc
// @Summary      Get Product
// @Description  Retrieves a product by ID
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        id path string true "Product ID"
// @Success      200 {object} entity.Product
// @Failure      400 {object} dto.Error
// @Failure      404 {object} dto.Error
// @Failure      500 {object} dto.Error
// @Router       /products/{id} [get]
// @Security     ApiKeyAuth
func (h ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad Request: \n" + "Parâmetro 'id' não foi passado"))
		return
	}

	product, err := h.ProductDB.FindProductById(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Not Found: \n" + err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(product)
}

// Update Product godoc
// @Summary      Update Product
// @Description  Updates an existing product
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        id path string true "Product ID"
// @Param        request body entity.Product true "Product data"
// @Success      200
// @Failure      400 {object} dto.Error
// @Failure      404 {object} dto.Error
// @Failure      500 {object} dto.Error
// @Router       /products/{id} [patch]
// @Security     ApiKeyAuth
func (h ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad Request: \n" + "Parâmetro 'id' não foi passado"))
		return
	}

	var product entity.Product

	err := json.NewDecoder(r.Body).Decode(&product)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad Request: \n" + err.Error()))
		return
	}

	product.Id, err = entityPkg.ParseID(id)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad Request: \n" + "Parâmetro 'id' é inválido"))
		return
	}

	currentProduct, err := h.ProductDB.FindProductById(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Produto não encontrado: \n" + err.Error()))
		return
	}

	if product.Name == "" {
		product.Name = currentProduct.Name
	}

	if product.Price == 0 {
		product.Price = currentProduct.Price
	}

	err = h.ProductDB.UpdateProduct(&product)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Erro interno no servidor: \n" + err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Delete Product godoc
// @Summary      Delete Product
// @Description  Deletes a product by ID
// @Tags         products
// @Produce      json
// @Param        id path string true "Product ID"
// @Success      200
// @Failure      400 {object} dto.Error
// @Failure      404 {object} dto.Error
// @Failure      500 {object} dto.Error
// @Router       /products/{id} [delete]
// @Security     ApiKeyAuth
func (h ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(
			w,
			"Parâmetro 'id' não foi passado na requisição",
			http.StatusBadRequest,
		)
		return
	}

	_, err := h.ProductDB.FindProductById(id)
	if err != nil {
		http.Error(
			w,
			"Produto não encntrado",
			http.StatusNotFound,
		)
		return
	}

	err = h.ProductDB.DeleteProduct(id)
	if err != nil {
		http.Error(
			w,
			"Internal Server Error: \n"+err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Get All Products godoc
// @Summary      Get All Products
// @Description  Retrieves a paginated list of products
// @Tags         products
// @Produce      json
// @Param        page query int false "Page number"
// @Param        limit query int false "Number of products per page"
// @Param        order query string false "Sort direction"
// @Param        field query string false "Field used for sorting"
// @Success      200 {array} entity.Product
// @Failure      500 {object} dto.Error
// @Router       /products [get]
// @Security     ApiKeyAuth
func (h ProductHandler) GetAllProducts(w http.ResponseWriter, r *http.Request) {
	page := r.URL.Query().Get("page")
	limit := r.URL.Query().Get("limit")
	order := r.URL.Query().Get("order")
	field := r.URL.Query().Get("field")

	pageInt, err := strconv.Atoi(page)
	if err != nil {
		pageInt = 0
	}

	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		limitInt = 0
	}

	products, err := h.ProductDB.FindAllProducts(pageInt, limitInt, order, field)
	if err != nil {
		http.Error(
			w,
			"Internal Server Error\n "+err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(products)
}
