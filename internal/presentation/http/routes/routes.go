// internal/presentation/http/routes/routes.go
package routes

import (
	"net/http"

	"github.com/farmanexo/catalog-service/internal/presentation/http/controllers"
	"github.com/farmanexo/catalog-service/internal/presentation/http/middlewares"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	httpSwagger "github.com/swaggo/http-swagger"
)

// SetupRoutes configura todas las rutas del servicio
func SetupRoutes(
	catalogController *controllers.CatalogController,
	authMiddleware *middlewares.AuthMiddleware,
) *chi.Mux {
	r := chi.NewRouter()

	// ========================================
	// MIDDLEWARES GLOBALES
	// ========================================

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// CORS
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "https://farmanexo.pe"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Use(middlewares.CorrelationID)

	// ========================================
	// SWAGGER DOCUMENTATION
	// ========================================

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:4003/swagger/doc.json"),
	))

	// ========================================
	// HEALTH CHECK
	// ========================================

	r.Get("/health", catalogController.HealthCheck)
	r.Get("/", catalogController.HealthCheck)

	// ========================================
	// API ROUTES - VERSION 1
	// ========================================

	r.Route("/api/v1", func(r chi.Router) {
		// ========================================
		// PRODUCTOS - ENDPOINTS PÚBLICOS
		// ========================================
		r.Route("/products", func(r chi.Router) {
			r.Get("/", catalogController.ListProducts)
			r.Post("/search", catalogController.SearchProducts)
			r.Get("/barcode/{barcode}", catalogController.GetProductByBarcode)
			r.Get("/{id}", catalogController.GetProduct)
			r.Get("/{id}/interactions", catalogController.ListDrugInteractions)
			r.Get("/{id}/frequently-bought-together", catalogController.ListFrequentlyBoughtTogether)
			r.Get("/{id}/availability", catalogController.GetProductAvailability)

			// Endpoints protegidos (Admin)
			r.Group(func(r chi.Router) {
				r.Use(authMiddleware.RequireAuth)
				r.Use(authMiddleware.RequireAdmin)

				r.Post("/", catalogController.CreateProduct)
				r.Put("/{id}", catalogController.UpdateProduct)
				r.Delete("/{id}", catalogController.DeleteProduct)
				r.Put("/{id}/images", catalogController.UploadProductImages)
				r.Post("/interactions", catalogController.CreateDrugInteraction)
			})
		})

		// ========================================
		// CATEGORÍAS - ENDPOINTS PÚBLICOS
		// ========================================
		r.Route("/categories", func(r chi.Router) {
			r.Get("/", catalogController.ListCategories)
			r.Get("/{id}", catalogController.GetCategory)
			r.Get("/{id}/products", catalogController.ListProductsByCategory)

			// Endpoints protegidos (Admin)
			r.Group(func(r chi.Router) {
				r.Use(authMiddleware.RequireAuth)
				r.Use(authMiddleware.RequireAdmin)

				r.Post("/", catalogController.CreateCategory)
				r.Put("/{id}", catalogController.UpdateCategory)
			})
		})

		// ========================================
		// MARCAS - ENDPOINTS PÚBLICOS
		// ========================================
		r.Route("/brands", func(r chi.Router) {
			r.Get("/", catalogController.ListBrands)
			r.Get("/{id}/products", catalogController.ListProductsByBrand)

			// Endpoints protegidos (Admin)
			r.Group(func(r chi.Router) {
				r.Use(authMiddleware.RequireAuth)
				r.Use(authMiddleware.RequireAdmin)

				r.Post("/", catalogController.CreateBrand)
				r.Put("/{id}", catalogController.UpdateBrand)
			})
		})
	})

	// ========================================
	// API ROUTES - VERSION 2 (Futuro)
	// ========================================

	r.Route("/api/v2", func(r chi.Router) {
		r.Get("/status", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("API v2 - Próximamente"))
		})
	})

	return r
}
