// cmd/server/main.go
package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/farmanexo/catalog-service/internal/application/commands"
	"github.com/farmanexo/catalog-service/internal/application/handlers"
	"github.com/farmanexo/catalog-service/internal/application/postprocessors"
	"github.com/farmanexo/catalog-service/internal/application/preprocessors"
	"github.com/farmanexo/catalog-service/internal/application/validators"
	"github.com/farmanexo/catalog-service/internal/presentation/dto/responses"
	"github.com/farmanexo/catalog-service/internal/infrastructure/cache"
	"github.com/farmanexo/catalog-service/internal/infrastructure/messaging"
	"github.com/farmanexo/catalog-service/internal/infrastructure/persistence/postgres"
	"github.com/farmanexo/catalog-service/internal/infrastructure/security"
	"github.com/farmanexo/catalog-service/internal/infrastructure/storage"
	"github.com/farmanexo/catalog-service/internal/presentation/http/controllers"
	"github.com/farmanexo/catalog-service/internal/presentation/http/middlewares"
	"github.com/farmanexo/catalog-service/internal/presentation/http/routes"
	"github.com/farmanexo/catalog-service/pkg/config"
	"github.com/farmanexo/catalog-service/pkg/mediator"

	// Swagger docs
	_ "github.com/farmanexo/catalog-service/docs"

	"go.uber.org/zap"
	pgdriver "gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

// @title           FarmaNexo Catalog Service API
// @version         1.0
// @description     Servicio de catálogo de productos farmacéuticos para FarmaNexo - Microservicio con CQRS y Clean Architecture
// @termsOfService  https://farmanexo.pe/terms

// @contact.name    FarmaNexo API Support
// @contact.url     https://farmanexo.pe/support
// @contact.email   support@farmanexo.pe

// @license.name    Apache 2.0
// @license.url     http://www.apache.org/licenses/LICENSE-2.0.html

// @host            localhost:4003
// @BasePath        /api/v1

// @securityDefinitions.apikey  BearerAuth
// @in                          header
// @name                        Authorization
// @description                 JWT Authorization header using the Bearer scheme. Example: "Bearer {token}"

// @tag.name         Products
// @tag.description  Endpoints de productos farmacéuticos

// @tag.name         Categories
// @tag.description  Endpoints de categorías de productos

// @tag.name         Brands
// @tag.description  Endpoints de marcas farmacéuticas

// @tag.name         Health
// @tag.description  Endpoints de salud del servicio

func main() {
	env := getEnvironment()
	cfg, err := config.LoadConfig(env)
	if err != nil {
		panic(fmt.Sprintf("Error cargando configuración: %v", err))
	}

	logger := initLogger(cfg)
	defer logger.Sync()

	logger.Info("Iniciando Catalog Service",
		zap.String("environment", cfg.Environment),
		zap.Int("port", cfg.Server.Port),
	)

	db := initDatabase(cfg, logger)

	logger.Info("Auto-migration deshabilitado - Usar migraciones manuales")

	// ========================================
	// REPOSITORIOS
	// ========================================
	productRepo := postgres.NewProductRepository(db, logger)
	categoryRepo := postgres.NewCategoryRepository(db, logger)
	brandRepo := postgres.NewBrandRepository(db, logger)
	imageRepo := postgres.NewProductImageRepository(db, logger)

	// ========================================
	// SERVICIOS
	// ========================================
	jwtService := security.NewJWTService(cfg.JWT.Secret, logger)

	// S3 File Storage
	fileStorage, err := storage.NewS3FileStorage(cfg.AWS, logger)
	if err != nil {
		logger.Fatal("Error inicializando S3 FileStorage", zap.Error(err))
	}

	// SQS Event Publisher
	eventPublisher, err := messaging.NewSQSEventPublisher(cfg.AWS, cfg.SQS, logger)
	if err != nil {
		logger.Fatal("Error inicializando SQS EventPublisher", zap.Error(err))
	}

	// Redis Cache
	redisClient, err := cache.NewRedisClient(cfg.Redis, logger)
	if err != nil {
		logger.Fatal("Error inicializando Redis", zap.Error(err))
	}
	defer redisClient.Close()

	cacheService := cache.NewRedisCacheService(redisClient, logger)

	// ========================================
	// MEDIATOR
	// ========================================
	med := mediator.NewMediator()

	// ========================================
	// HANDLERS - Products
	// ========================================
	listProductsHandler := handlers.NewListProductsHandler(productRepo, logger)
	mediator.RegisterHandler(med, listProductsHandler)

	getProductHandler := handlers.NewGetProductHandler(productRepo, cacheService, logger)
	mediator.RegisterHandler(med, getProductHandler)

	searchProductsHandler := handlers.NewSearchProductsHandler(productRepo, logger)
	mediator.RegisterHandler(med, searchProductsHandler)

	createProductHandler := handlers.NewCreateProductHandler(productRepo, categoryRepo, brandRepo, eventPublisher, cacheService, logger)
	mediator.RegisterHandler(med, createProductHandler)

	updateProductHandler := handlers.NewUpdateProductHandler(productRepo, categoryRepo, brandRepo, eventPublisher, cacheService, logger)
	mediator.RegisterHandler(med, updateProductHandler)

	deleteProductHandler := handlers.NewDeleteProductHandler(productRepo, eventPublisher, cacheService, logger)
	mediator.RegisterHandler(med, deleteProductHandler)

	uploadImagesHandler := handlers.NewUploadProductImagesHandler(productRepo, imageRepo, fileStorage, cacheService, cfg.S3.ProductsBucket, logger)
	mediator.RegisterHandler(med, uploadImagesHandler)

	// ========================================
	// HANDLERS - Categories
	// ========================================
	listCategoriesHandler := handlers.NewListCategoriesHandler(categoryRepo, cacheService, logger)
	mediator.RegisterHandler(med, listCategoriesHandler)

	getCategoryHandler := handlers.NewGetCategoryHandler(categoryRepo, logger)
	mediator.RegisterHandler(med, getCategoryHandler)

	listProductsByCategoryHandler := handlers.NewListProductsByCategoryHandler(productRepo, categoryRepo, logger)
	mediator.RegisterHandler(med, listProductsByCategoryHandler)

	createCategoryHandler := handlers.NewCreateCategoryHandler(categoryRepo, cacheService, logger)
	mediator.RegisterHandler(med, createCategoryHandler)

	updateCategoryHandler := handlers.NewUpdateCategoryHandler(categoryRepo, cacheService, logger)
	mediator.RegisterHandler(med, updateCategoryHandler)

	// ========================================
	// HANDLERS - Brands
	// ========================================
	listBrandsHandler := handlers.NewListBrandsHandler(brandRepo, logger)
	mediator.RegisterHandler(med, listBrandsHandler)

	listProductsByBrandHandler := handlers.NewListProductsByBrandHandler(productRepo, brandRepo, logger)
	mediator.RegisterHandler(med, listProductsByBrandHandler)

	createBrandHandler := handlers.NewCreateBrandHandler(brandRepo, logger)
	mediator.RegisterHandler(med, createBrandHandler)

	updateBrandHandler := handlers.NewUpdateBrandHandler(brandRepo, logger)
	mediator.RegisterHandler(med, updateBrandHandler)

	// ========================================
	// VALIDATORS
	// ========================================
	createProductValidator := validators.NewCreateProductValidator()
	mediator.RegisterValidator[commands.CreateProductCommand, responses.ProductResponse](med, createProductValidator)

	createCategoryValidator := validators.NewCreateCategoryValidator()
	mediator.RegisterValidator[commands.CreateCategoryCommand, responses.CategoryResponse](med, createCategoryValidator)

	createBrandValidator := validators.NewCreateBrandValidator()
	mediator.RegisterValidator[commands.CreateBrandCommand, responses.BrandResponse](med, createBrandValidator)

	// ========================================
	// PREPROCESSORS Y POSTPROCESSORS
	// ========================================
	sanitizePreProcessor := preprocessors.NewSanitizeInputPreProcessor(logger)
	med.RegisterPreProcessor(sanitizePreProcessor)

	auditPostProcessor := postprocessors.NewLogAuditPostProcessor(logger)
	med.RegisterPostProcessor(auditPostProcessor)

	logger.Info("Mediator configurado",
		zap.Int("handlers", 16),
		zap.Int("validators", 3),
		zap.Int("preprocessors", 1),
		zap.Int("postprocessors", 1),
	)

	// ========================================
	// MIDDLEWARES
	// ========================================
	authMiddleware := middlewares.NewAuthMiddleware(jwtService, logger)

	// ========================================
	// CONTROLADORES Y RUTAS
	// ========================================
	catalogController := controllers.NewCatalogController(med, logger)
	router := routes.SetupRoutes(catalogController, authMiddleware)

	// ========================================
	// SERVIDOR HTTP
	// ========================================
	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	go func() {
		logger.Info("Servidor HTTP iniciado",
			zap.String("address", server.Addr),
			zap.String("swagger_url", fmt.Sprintf("http://localhost:%d/swagger/index.html", cfg.Server.Port)),
		)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Error iniciando servidor", zap.Error(err))
		}
	}()

	// ========================================
	// GRACEFUL SHUTDOWN
	// ========================================
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Iniciando graceful shutdown...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error("Error en shutdown", zap.Error(err))
	}

	logger.Info("Servidor detenido exitosamente")
}

func getEnvironment() string {
	env := os.Getenv("ENV")
	if env == "" {
		env = "local"
	}
	return env
}

func initLogger(cfg *config.Config) *zap.Logger {
	var logger *zap.Logger
	var err error

	if cfg.IsProduction() || cfg.IsUAT() {
		logger, err = zap.NewProduction()
	} else {
		logger, err = zap.NewDevelopment()
	}

	if err != nil {
		panic(fmt.Sprintf("Error inicializando logger: %v", err))
	}

	return logger
}

func initDatabase(cfg *config.Config, logger *zap.Logger) *gorm.DB {
	gormLogLevel := gormlogger.Silent
	if cfg.IsDevelopment() {
		gormLogLevel = gormlogger.Info
	}

	gormLogger := gormlogger.Default.LogMode(gormLogLevel)

	db, err := gorm.Open(pgdriver.Open(cfg.Database.GetDSN()), &gorm.Config{
		Logger: gormLogger,
	})

	if err != nil {
		logger.Fatal("Error conectando a PostgreSQL",
			zap.Error(err),
			zap.String("host", cfg.Database.Host),
			zap.Int("port", cfg.Database.Port),
		)
	}

	sqlDB, err := db.DB()
	if err != nil {
		logger.Fatal("Error obteniendo SQL DB", zap.Error(err))
	}

	sqlDB.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.Database.ConnMaxLifetime)

	logger.Info("Conexión a PostgreSQL establecida",
		zap.String("host", cfg.Database.Host),
		zap.String("database", cfg.Database.DBName),
	)

	return db
}
