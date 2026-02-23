// internal/presentation/http/controllers/catalog_controller.go
package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/farmanexo/catalog-service/internal/application/commands"
	"github.com/farmanexo/catalog-service/internal/application/queries"
	"github.com/farmanexo/catalog-service/internal/presentation/dto/requests"
	"github.com/farmanexo/catalog-service/internal/presentation/dto/responses"
	"github.com/farmanexo/catalog-service/internal/presentation/http/middlewares"
	"github.com/farmanexo/catalog-service/internal/shared/common"
	"github.com/farmanexo/catalog-service/pkg/mediator"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

// CatalogController controlador HTTP del catálogo
type CatalogController struct {
	mediator *mediator.Mediator
	logger   *zap.Logger
}

func NewCatalogController(med *mediator.Mediator, logger *zap.Logger) *CatalogController {
	return &CatalogController{mediator: med, logger: logger}
}

// respondJSON helper para escribir respuesta JSON
func (c *CatalogController) respondJSON(w http.ResponseWriter, response interface{}) {
	statusCode := http.StatusOK

	if resp, ok := response.(interface{ GetHttpStatus() *int }); ok {
		if httpStatus := resp.GetHttpStatus(); httpStatus != nil {
			statusCode = *httpStatus
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		c.logger.Error("Error codificando respuesta JSON", zap.Error(err))
	}
}

// HealthCheck godoc
// @Summary      Health check del servicio
// @Description  Retorna el estado del servicio
// @Tags         Health
// @Produce      json
// @Success      200  {object}  map[string]interface{}  "Servicio saludable"
// @Router       /health [get]
func (c *CatalogController) HealthCheck(w http.ResponseWriter, r *http.Request) {
	type HealthResponse struct {
		Status  string `json:"status" example:"healthy"`
		Service string `json:"service" example:"catalog-service"`
		Version string `json:"version" example:"1.0.0"`
	}

	health := HealthResponse{
		Status:  "healthy",
		Service: "catalog-service",
		Version: "1.0.0",
	}

	c.respondJSON(w, common.OkResponse(health))
}

// ========================================
// PRODUCTOS - ENDPOINTS PÚBLICOS
// ========================================

// ListProducts godoc
// @Summary      Listar productos
// @Description  Retorna lista paginada de productos activos
// @Tags         Products
// @Accept       json
// @Produce      json
// @Param        page   query    int     false  "Página"          default(1)
// @Param        limit  query    int     false  "Límite por página" default(20)
// @Param        sort   query    string  false  "Ordenamiento"    default(name_asc)
// @Success      200  {object}  common.ApiResponse[responses.ProductListResponse]
// @Failure      500  {object}  common.ApiResponse[responses.ProductListResponse]
// @Router       /api/v1/products [get]
func (c *CatalogController) ListProducts(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	sort := r.URL.Query().Get("sort")

	query := queries.ListProductsQuery{
		Page:  page,
		Limit: limit,
		Sort:  sort,
	}

	response, _ := mediator.Send[queries.ListProductsQuery, responses.ProductListResponse](r.Context(), c.mediator, query)
	c.respondJSON(w, response)
}

// GetProduct godoc
// @Summary      Obtener producto por ID
// @Description  Retorna el detalle de un producto
// @Tags         Products
// @Accept       json
// @Produce      json
// @Param        id   path     string  true  "Product ID"
// @Success      200  {object}  common.ApiResponse[responses.ProductResponse]
// @Failure      404  {object}  common.ApiResponse[responses.ProductResponse]
// @Router       /api/v1/products/{id} [get]
func (c *CatalogController) GetProduct(w http.ResponseWriter, r *http.Request) {
	productID := chi.URLParam(r, "id")

	_, isAdmin := middlewares.GetUserRoleFromContext(r.Context())

	query := queries.GetProductQuery{
		ID:      productID,
		IsAdmin: isAdmin,
	}

	response, _ := mediator.Send[queries.GetProductQuery, responses.ProductResponse](r.Context(), c.mediator, query)
	c.respondJSON(w, response)
}

// SearchProducts godoc
// @Summary      Búsqueda avanzada de productos
// @Description  Busca productos por nombre, ingrediente activo, categoría, marca
// @Tags         Products
// @Accept       json
// @Produce      json
// @Param        body  body     requests.SearchProductsRequest  true  "Parámetros de búsqueda"
// @Success      200  {object}  common.ApiResponse[responses.ProductListResponse]
// @Failure      400  {object}  common.ApiResponse[responses.ProductListResponse]
// @Router       /api/v1/products/search [post]
func (c *CatalogController) SearchProducts(w http.ResponseWriter, r *http.Request) {
	var req requests.SearchProductsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		c.respondJSON(w, common.BadRequestResponse[responses.ProductListResponse]("VAL_001", "Body inválido"))
		return
	}

	query := queries.SearchProductsQuery{
		Query:                req.Query,
		CategoryID:           req.CategoryID,
		BrandID:              req.BrandID,
		RequiresPrescription: req.RequiresPrescription,
		Page:                 req.Page,
		Limit:                req.Limit,
		Sort:                 req.Sort,
	}

	response, _ := mediator.Send[queries.SearchProductsQuery, responses.ProductListResponse](r.Context(), c.mediator, query)
	c.respondJSON(w, response)
}

// GetProductByBarcode godoc
// @Summary      Obtener producto por código de barras
// @Description  Busca un producto por su código de barras (cacheado, 1hr TTL)
// @Tags         Products
// @Produce      json
// @Param        barcode  path     string  true  "Código de barras"
// @Success      200  {object}  common.ApiResponse[responses.ProductResponse]
// @Failure      404  {object}  common.ApiResponse[responses.ProductResponse]
// @Router       /api/v1/products/barcode/{barcode} [get]
func (c *CatalogController) GetProductByBarcode(w http.ResponseWriter, r *http.Request) {
	barcode := chi.URLParam(r, "barcode")

	query := queries.GetProductByBarcodeQuery{
		Barcode: barcode,
	}

	response, _ := mediator.Send[queries.GetProductByBarcodeQuery, responses.ProductResponse](r.Context(), c.mediator, query)
	c.respondJSON(w, response)
}

// ========================================
// PRODUCTOS - ENDPOINTS ADMIN
// ========================================

// CreateProduct godoc
// @Summary      Crear producto
// @Description  Crea un nuevo producto (requiere rol admin)
// @Tags         Products
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body     requests.CreateProductRequest  true  "Datos del producto"
// @Success      201  {object}  common.ApiResponse[responses.ProductResponse]
// @Failure      400  {object}  common.ApiResponse[responses.ProductResponse]
// @Failure      401  {object}  common.ApiResponse[responses.ProductResponse]
// @Failure      403  {object}  common.ApiResponse[responses.ProductResponse]
// @Failure      409  {object}  common.ApiResponse[responses.ProductResponse]
// @Router       /api/v1/products [post]
func (c *CatalogController) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var req requests.CreateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		c.respondJSON(w, common.BadRequestResponse[responses.ProductResponse]("VAL_001", "Body inválido"))
		return
	}

	cmd := commands.CreateProductCommand{
		Name:                 req.Name,
		Slug:                 req.Slug,
		Description:          req.Description,
		ActiveIngredient:     req.ActiveIngredient,
		Presentation:         req.Presentation,
		Concentration:        req.Concentration,
		RequiresPrescription: req.RequiresPrescription,
		CategoryID:           req.CategoryID,
		BrandID:              req.BrandID,
		SKU:                  req.SKU,
		Barcode:              req.Barcode,
	}

	response, _ := mediator.Send[commands.CreateProductCommand, responses.ProductResponse](r.Context(), c.mediator, cmd)
	c.respondJSON(w, response)
}

// UpdateProduct godoc
// @Summary      Actualizar producto
// @Description  Actualiza un producto existente (requiere rol admin)
// @Tags         Products
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path     string                          true  "Product ID"
// @Param        body  body     requests.UpdateProductRequest   true  "Datos del producto"
// @Success      200  {object}  common.ApiResponse[responses.ProductResponse]
// @Failure      400  {object}  common.ApiResponse[responses.ProductResponse]
// @Failure      404  {object}  common.ApiResponse[responses.ProductResponse]
// @Router       /api/v1/products/{id} [put]
func (c *CatalogController) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	productID := chi.URLParam(r, "id")

	var req requests.UpdateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		c.respondJSON(w, common.BadRequestResponse[responses.ProductResponse]("VAL_001", "Body inválido"))
		return
	}

	cmd := commands.UpdateProductCommand{
		ID:                   productID,
		Name:                 req.Name,
		Slug:                 req.Slug,
		Description:          req.Description,
		ActiveIngredient:     req.ActiveIngredient,
		Presentation:         req.Presentation,
		Concentration:        req.Concentration,
		RequiresPrescription: req.RequiresPrescription,
		CategoryID:           req.CategoryID,
		BrandID:              req.BrandID,
		SKU:                  req.SKU,
		Barcode:              req.Barcode,
		IsActive:             req.IsActive,
	}

	response, _ := mediator.Send[commands.UpdateProductCommand, responses.ProductResponse](r.Context(), c.mediator, cmd)
	c.respondJSON(w, response)
}

// DeleteProduct godoc
// @Summary      Eliminar producto (soft delete)
// @Description  Marca un producto como eliminado (requiere rol admin)
// @Tags         Products
// @Produce      json
// @Security     BearerAuth
// @Param        id   path     string  true  "Product ID"
// @Success      200  {object}  common.ApiResponse[responses.EmptyResponse]
// @Failure      404  {object}  common.ApiResponse[responses.EmptyResponse]
// @Router       /api/v1/products/{id} [delete]
func (c *CatalogController) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	productID := chi.URLParam(r, "id")

	cmd := commands.DeleteProductCommand{ID: productID}

	response, _ := mediator.Send[commands.DeleteProductCommand, responses.EmptyResponse](r.Context(), c.mediator, cmd)
	c.respondJSON(w, response)
}

// UploadProductImages godoc
// @Summary      Subir imágenes de producto
// @Description  Sube una o más imágenes para un producto (requiere rol admin, max 10MB por imagen)
// @Tags         Products
// @Accept       multipart/form-data
// @Produce      json
// @Security     BearerAuth
// @Param        id      path     string  true   "Product ID"
// @Param        images  formData file    true   "Archivos de imagen"
// @Param        primary formData string  false  "Índice de imagen principal (0-based)"
// @Success      201  {object}  common.ApiResponse[responses.ImageListResponse]
// @Failure      400  {object}  common.ApiResponse[responses.ImageListResponse]
// @Failure      404  {object}  common.ApiResponse[responses.ImageListResponse]
// @Router       /api/v1/products/{id}/images [put]
func (c *CatalogController) UploadProductImages(w http.ResponseWriter, r *http.Request) {
	productID := chi.URLParam(r, "id")

	if err := r.ParseMultipartForm(50 << 20); err != nil { // 50MB max total
		c.respondJSON(w, common.BadRequestResponse[responses.ImageListResponse]("VAL_001", "Error procesando formulario multipart"))
		return
	}

	files := r.MultipartForm.File["images"]
	if len(files) == 0 {
		c.respondJSON(w, common.BadRequestResponse[responses.ImageListResponse]("VAL_006", "Al menos una imagen es requerida"))
		return
	}

	primaryIndex, _ := strconv.Atoi(r.FormValue("primary"))

	var imageFiles []commands.ImageFile
	for i, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			continue
		}

		imageFiles = append(imageFiles, commands.ImageFile{
			Reader:      file,
			Filename:    fileHeader.Filename,
			ContentType: fileHeader.Header.Get("Content-Type"),
			Size:        fileHeader.Size,
			IsPrimary:   i == primaryIndex,
		})
	}

	cmd := commands.UploadProductImagesCommand{
		ProductID: productID,
		Files:     imageFiles,
	}

	response, _ := mediator.Send[commands.UploadProductImagesCommand, responses.ImageListResponse](r.Context(), c.mediator, cmd)
	c.respondJSON(w, response)
}

// ========================================
// DISPONIBILIDAD EN FARMACIAS
// ========================================

// GetProductAvailability godoc
// @Summary      Disponibilidad en farmacias
// @Description  Consulta la disponibilidad y precio de un producto en las farmacias registradas (cacheado, 5min TTL)
// @Tags         Products
// @Produce      json
// @Param        id   path     string  true  "Product ID"
// @Success      200  {object}  common.ApiResponse[responses.AvailabilityResponse]
// @Failure      404  {object}  common.ApiResponse[responses.AvailabilityResponse]
// @Router       /api/v1/products/{id}/availability [get]
func (c *CatalogController) GetProductAvailability(w http.ResponseWriter, r *http.Request) {
	productID := chi.URLParam(r, "id")

	query := queries.GetProductAvailabilityQuery{
		ProductID: productID,
	}

	response, _ := mediator.Send[queries.GetProductAvailabilityQuery, responses.AvailabilityResponse](r.Context(), c.mediator, query)
	c.respondJSON(w, response)
}

// ========================================
// PRODUCTOS FRECUENTEMENTE COMPRADOS JUNTOS
// ========================================

// ListFrequentlyBoughtTogether godoc
// @Summary      Productos frecuentemente comprados juntos
// @Description  Retorna productos que se compran frecuentemente junto con el producto dado (cacheado, 6hr TTL)
// @Tags         Products
// @Produce      json
// @Param        id     path     string  true   "Product ID"
// @Param        limit  query    int     false  "Límite de resultados" default(10)
// @Success      200  {object}  common.ApiResponse[responses.FBTListResponse]
// @Failure      404  {object}  common.ApiResponse[responses.FBTListResponse]
// @Router       /api/v1/products/{id}/frequently-bought-together [get]
func (c *CatalogController) ListFrequentlyBoughtTogether(w http.ResponseWriter, r *http.Request) {
	productID := chi.URLParam(r, "id")
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	query := queries.ListFBTQuery{
		ProductID: productID,
		Limit:     limit,
	}

	response, _ := mediator.Send[queries.ListFBTQuery, responses.FBTListResponse](r.Context(), c.mediator, query)
	c.respondJSON(w, response)
}

// ========================================
// INTERACCIONES MEDICAMENTOSAS
// ========================================

// ListDrugInteractions godoc
// @Summary      Listar interacciones medicamentosas
// @Description  Retorna las interacciones medicamentosas de un producto (cacheado, 24hr TTL)
// @Tags         Products
// @Produce      json
// @Param        id   path     string  true  "Product ID"
// @Success      200  {object}  common.ApiResponse[responses.InteractionListResponse]
// @Failure      404  {object}  common.ApiResponse[responses.InteractionListResponse]
// @Router       /api/v1/products/{id}/interactions [get]
func (c *CatalogController) ListDrugInteractions(w http.ResponseWriter, r *http.Request) {
	productID := chi.URLParam(r, "id")

	query := queries.ListDrugInteractionsQuery{
		ProductID: productID,
	}

	response, _ := mediator.Send[queries.ListDrugInteractionsQuery, responses.InteractionListResponse](r.Context(), c.mediator, query)
	c.respondJSON(w, response)
}

// CreateDrugInteraction godoc
// @Summary      Crear interacción medicamentosa
// @Description  Registra una nueva interacción entre dos productos (requiere rol admin)
// @Tags         Products
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body     requests.CreateInteractionRequest  true  "Datos de la interacción"
// @Success      201  {object}  common.ApiResponse[responses.InteractionResponse]
// @Failure      400  {object}  common.ApiResponse[responses.InteractionResponse]
// @Failure      404  {object}  common.ApiResponse[responses.InteractionResponse]
// @Failure      409  {object}  common.ApiResponse[responses.InteractionResponse]
// @Router       /api/v1/products/interactions [post]
func (c *CatalogController) CreateDrugInteraction(w http.ResponseWriter, r *http.Request) {
	var req requests.CreateInteractionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		c.respondJSON(w, common.BadRequestResponse[responses.InteractionResponse]("VAL_001", "Body inválido"))
		return
	}

	cmd := commands.CreateDrugInteractionCommand{
		ProductID:              req.ProductID,
		InteractsWithProductID: req.InteractsWithProductID,
		Severity:               req.Severity,
		Description:            req.Description,
		Recommendation:         req.Recommendation,
	}

	response, _ := mediator.Send[commands.CreateDrugInteractionCommand, responses.InteractionResponse](r.Context(), c.mediator, cmd)
	c.respondJSON(w, response)
}

// ========================================
// CATEGORÍAS - ENDPOINTS PÚBLICOS
// ========================================

// ListCategories godoc
// @Summary      Listar categorías
// @Description  Retorna todas las categorías activas con sus hijos
// @Tags         Categories
// @Produce      json
// @Success      200  {object}  common.ApiResponse[responses.CategoryListResponse]
// @Router       /api/v1/categories [get]
func (c *CatalogController) ListCategories(w http.ResponseWriter, r *http.Request) {
	query := queries.ListCategoriesQuery{}
	response, _ := mediator.Send[queries.ListCategoriesQuery, responses.CategoryListResponse](r.Context(), c.mediator, query)
	c.respondJSON(w, response)
}

// GetCategory godoc
// @Summary      Obtener categoría por ID
// @Description  Retorna una categoría con sus hijos
// @Tags         Categories
// @Produce      json
// @Param        id   path     string  true  "Category ID"
// @Success      200  {object}  common.ApiResponse[responses.CategoryResponse]
// @Failure      404  {object}  common.ApiResponse[responses.CategoryResponse]
// @Router       /api/v1/categories/{id} [get]
func (c *CatalogController) GetCategory(w http.ResponseWriter, r *http.Request) {
	categoryID := chi.URLParam(r, "id")
	query := queries.GetCategoryQuery{ID: categoryID}
	response, _ := mediator.Send[queries.GetCategoryQuery, responses.CategoryResponse](r.Context(), c.mediator, query)
	c.respondJSON(w, response)
}

// ListProductsByCategory godoc
// @Summary      Productos por categoría
// @Description  Retorna lista paginada de productos de una categoría
// @Tags         Categories
// @Produce      json
// @Param        id     path     string  true   "Category ID"
// @Param        page   query    int     false  "Página"          default(1)
// @Param        limit  query    int     false  "Límite por página" default(20)
// @Success      200  {object}  common.ApiResponse[responses.ProductListResponse]
// @Failure      404  {object}  common.ApiResponse[responses.ProductListResponse]
// @Router       /api/v1/categories/{id}/products [get]
func (c *CatalogController) ListProductsByCategory(w http.ResponseWriter, r *http.Request) {
	categoryID := chi.URLParam(r, "id")
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	query := queries.ListProductsByCategoryQuery{
		CategoryID: categoryID,
		Page:       page,
		Limit:      limit,
	}

	response, _ := mediator.Send[queries.ListProductsByCategoryQuery, responses.ProductListResponse](r.Context(), c.mediator, query)
	c.respondJSON(w, response)
}

// ========================================
// CATEGORÍAS - ENDPOINTS ADMIN
// ========================================

// CreateCategory godoc
// @Summary      Crear categoría
// @Description  Crea una nueva categoría (requiere rol admin)
// @Tags         Categories
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body     requests.CreateCategoryRequest  true  "Datos de la categoría"
// @Success      201  {object}  common.ApiResponse[responses.CategoryResponse]
// @Failure      400  {object}  common.ApiResponse[responses.CategoryResponse]
// @Router       /api/v1/categories [post]
func (c *CatalogController) CreateCategory(w http.ResponseWriter, r *http.Request) {
	var req requests.CreateCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		c.respondJSON(w, common.BadRequestResponse[responses.CategoryResponse]("VAL_001", "Body inválido"))
		return
	}

	cmd := commands.CreateCategoryCommand{
		Name:         req.Name,
		Slug:         req.Slug,
		Description:  req.Description,
		ParentID:     req.ParentID,
		ImageURL:     req.ImageURL,
		DisplayOrder: req.DisplayOrder,
	}

	response, _ := mediator.Send[commands.CreateCategoryCommand, responses.CategoryResponse](r.Context(), c.mediator, cmd)
	c.respondJSON(w, response)
}

// UpdateCategory godoc
// @Summary      Actualizar categoría
// @Description  Actualiza una categoría existente (requiere rol admin)
// @Tags         Categories
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path     string                          true  "Category ID"
// @Param        body  body     requests.UpdateCategoryRequest  true  "Datos de la categoría"
// @Success      200  {object}  common.ApiResponse[responses.CategoryResponse]
// @Failure      404  {object}  common.ApiResponse[responses.CategoryResponse]
// @Router       /api/v1/categories/{id} [put]
func (c *CatalogController) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	categoryID := chi.URLParam(r, "id")

	var req requests.UpdateCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		c.respondJSON(w, common.BadRequestResponse[responses.CategoryResponse]("VAL_001", "Body inválido"))
		return
	}

	cmd := commands.UpdateCategoryCommand{
		ID:           categoryID,
		Name:         req.Name,
		Slug:         req.Slug,
		Description:  req.Description,
		ParentID:     req.ParentID,
		ImageURL:     req.ImageURL,
		IsActive:     req.IsActive,
		DisplayOrder: req.DisplayOrder,
	}

	response, _ := mediator.Send[commands.UpdateCategoryCommand, responses.CategoryResponse](r.Context(), c.mediator, cmd)
	c.respondJSON(w, response)
}

// ========================================
// MARCAS - ENDPOINTS PÚBLICOS
// ========================================

// ListBrands godoc
// @Summary      Listar marcas
// @Description  Retorna todas las marcas activas
// @Tags         Brands
// @Produce      json
// @Success      200  {object}  common.ApiResponse[responses.BrandListResponse]
// @Router       /api/v1/brands [get]
func (c *CatalogController) ListBrands(w http.ResponseWriter, r *http.Request) {
	query := queries.ListBrandsQuery{}
	response, _ := mediator.Send[queries.ListBrandsQuery, responses.BrandListResponse](r.Context(), c.mediator, query)
	c.respondJSON(w, response)
}

// ListProductsByBrand godoc
// @Summary      Productos por marca
// @Description  Retorna lista paginada de productos de una marca
// @Tags         Brands
// @Produce      json
// @Param        id     path     string  true   "Brand ID"
// @Param        page   query    int     false  "Página"          default(1)
// @Param        limit  query    int     false  "Límite por página" default(20)
// @Success      200  {object}  common.ApiResponse[responses.ProductListResponse]
// @Failure      404  {object}  common.ApiResponse[responses.ProductListResponse]
// @Router       /api/v1/brands/{id}/products [get]
func (c *CatalogController) ListProductsByBrand(w http.ResponseWriter, r *http.Request) {
	brandID := chi.URLParam(r, "id")
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	query := queries.ListProductsByBrandQuery{
		BrandID: brandID,
		Page:    page,
		Limit:   limit,
	}

	response, _ := mediator.Send[queries.ListProductsByBrandQuery, responses.ProductListResponse](r.Context(), c.mediator, query)
	c.respondJSON(w, response)
}

// ========================================
// MARCAS - ENDPOINTS ADMIN
// ========================================

// CreateBrand godoc
// @Summary      Crear marca
// @Description  Crea una nueva marca (requiere rol admin)
// @Tags         Brands
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body     requests.CreateBrandRequest  true  "Datos de la marca"
// @Success      201  {object}  common.ApiResponse[responses.BrandResponse]
// @Failure      400  {object}  common.ApiResponse[responses.BrandResponse]
// @Router       /api/v1/brands [post]
func (c *CatalogController) CreateBrand(w http.ResponseWriter, r *http.Request) {
	var req requests.CreateBrandRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		c.respondJSON(w, common.BadRequestResponse[responses.BrandResponse]("VAL_001", "Body inválido"))
		return
	}

	cmd := commands.CreateBrandCommand{
		Name:        req.Name,
		Slug:        req.Slug,
		Description: req.Description,
		LogoURL:     req.LogoURL,
		Website:     req.Website,
		Country:     req.Country,
	}

	response, _ := mediator.Send[commands.CreateBrandCommand, responses.BrandResponse](r.Context(), c.mediator, cmd)
	c.respondJSON(w, response)
}

// UpdateBrand godoc
// @Summary      Actualizar marca
// @Description  Actualiza una marca existente (requiere rol admin)
// @Tags         Brands
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path     string                       true  "Brand ID"
// @Param        body  body     requests.UpdateBrandRequest  true  "Datos de la marca"
// @Success      200  {object}  common.ApiResponse[responses.BrandResponse]
// @Failure      404  {object}  common.ApiResponse[responses.BrandResponse]
// @Router       /api/v1/brands/{id} [put]
func (c *CatalogController) UpdateBrand(w http.ResponseWriter, r *http.Request) {
	brandID := chi.URLParam(r, "id")

	var req requests.UpdateBrandRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		c.respondJSON(w, common.BadRequestResponse[responses.BrandResponse]("VAL_001", "Body inválido"))
		return
	}

	cmd := commands.UpdateBrandCommand{
		ID:          brandID,
		Name:        req.Name,
		Slug:        req.Slug,
		Description: req.Description,
		LogoURL:     req.LogoURL,
		Website:     req.Website,
		Country:     req.Country,
		IsActive:    req.IsActive,
	}

	response, _ := mediator.Send[commands.UpdateBrandCommand, responses.BrandResponse](r.Context(), c.mediator, cmd)
	c.respondJSON(w, response)
}
