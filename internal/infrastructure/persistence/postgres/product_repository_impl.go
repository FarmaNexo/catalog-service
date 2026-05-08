// internal/infrastructure/persistence/postgres/product_repository_impl.go
package postgres

import (
	"context"
	"fmt"
	"math"
	"regexp"
	"strings"

	"github.com/farmanexo/catalog-service/internal/domain/entities"
	"github.com/farmanexo/catalog-service/internal/domain/repositories"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// ProductRepositoryImpl implementación PostgreSQL del repositorio de productos
type ProductRepositoryImpl struct {
	db     *gorm.DB
	logger *zap.Logger
}

func NewProductRepository(db *gorm.DB, logger *zap.Logger) *ProductRepositoryImpl {
	return &ProductRepositoryImpl{db: db, logger: logger}
}

func (r *ProductRepositoryImpl) Create(ctx context.Context, product *entities.Product) error {
	return r.db.WithContext(ctx).Create(product).Error
}

func (r *ProductRepositoryImpl) FindByID(ctx context.Context, id string) (*entities.Product, error) {
	var product entities.Product
	err := r.db.WithContext(ctx).
		Where("id = ? AND deleted_at IS NULL", id).
		Preload("Category").
		Preload("Brand").
		Preload("Images", func(db *gorm.DB) *gorm.DB {
			return db.Order("display_order ASC")
		}).
		First(&product).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *ProductRepositoryImpl) FindByIDWithDeleted(ctx context.Context, id string) (*entities.Product, error) {
	var product entities.Product
	err := r.db.WithContext(ctx).
		Where("id = ?", id).
		Preload("Category").
		Preload("Brand").
		Preload("Images", func(db *gorm.DB) *gorm.DB {
			return db.Order("display_order ASC")
		}).
		First(&product).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *ProductRepositoryImpl) FindBySlug(ctx context.Context, slug string) (*entities.Product, error) {
	var product entities.Product
	err := r.db.WithContext(ctx).Where("slug = ? AND deleted_at IS NULL", slug).First(&product).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *ProductRepositoryImpl) FindBySKU(ctx context.Context, sku string) (*entities.Product, error) {
	var product entities.Product
	err := r.db.WithContext(ctx).Where("sku = ? AND deleted_at IS NULL", sku).First(&product).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *ProductRepositoryImpl) FindByBarcode(ctx context.Context, barcode string) (*entities.Product, error) {
	var product entities.Product
	err := r.db.WithContext(ctx).
		Where("barcode = ? AND deleted_at IS NULL", barcode).
		Preload("Category").
		Preload("Brand").
		Preload("Images", func(db *gorm.DB) *gorm.DB {
			return db.Order("display_order ASC")
		}).
		First(&product).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

// FindBySourceCode busca un producto por (source_product_code, concentration).
// Se usa desde pharmacy-service (HTTP) para resolver inventory events.
func (r *ProductRepositoryImpl) FindBySourceCode(ctx context.Context, sourceProductCode int, concentration string) (*entities.Product, error) {
	var product entities.Product
	err := r.db.WithContext(ctx).
		Where("source_product_code = ? AND concentration = ? AND deleted_at IS NULL", sourceProductCode, concentration).
		First(&product).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *ProductRepositoryImpl) Update(ctx context.Context, product *entities.Product) error {
	return r.db.WithContext(ctx).Save(product).Error
}

// UpsertBySource hace INSERT ... ON CONFLICT (source_product_code, concentration) DO UPDATE.
// Idempotente: el mismo evento procesado dos veces deja la fila igual (UPDATE de
// los mismos valores). Retorna el id (UUID) del producto resultante.
//
// Reglas de merge:
//   - name → siempre se sobrescribe con el valor del evento (canonical_name DIGEMID).
//   - registry_number / manufacturer / form / etc → COALESCE(EXCLUDED, existing) =
//     no se borra un valor previo si el evento trae vacío.
//   - requires_prescription → OR lógico (si alguna fuente lo marcó true, queda true).
//   - sku → registry_number si existe; sino se fabrica DIGEMID-{code}-{form}-{slug-conc}.
//   - slug → incluye source_product_code para evitar colisiones entre marcas con
//     mismo (name, concentration, form).
func (r *ProductRepositoryImpl) UpsertBySource(ctx context.Context, p repositories.ProductUpsertParams) (string, error) {
	slug := buildProductSlugWithSource(p.CanonicalName, p.Concentration, p.SourceFormCode, p.SourceProductCode)
	manufacturer := firstNonEmptyStr(p.Manufacturer, p.Holder)

	var sku string
	if p.RegistryNumber != "" {
		sku = p.RegistryNumber
	} else {
		sku = fmt.Sprintf("DIGEMID-%d-%s-%s", p.SourceProductCode, p.SourceFormCode, slugify(p.Concentration))
	}

	var resultID string
	err := r.db.WithContext(ctx).Raw(`
		INSERT INTO catalog.products (
			name, slug, description, active_ingredient,
			presentation, concentration, form, registry_number,
			manufacturer, source_product_code,
			requires_prescription, sku, is_active, created_at, updated_at
		) VALUES (
			?, ?, '', NULLIF(?, ''),
			NULLIF(?, ''), ?, NULLIF(?, ''), NULLIF(?, ''),
			NULLIF(?, ''), ?,
			?, NULLIF(?, ''), true, NOW(), NOW()
		)
		ON CONFLICT (source_product_code, concentration)
		WHERE source_product_code IS NOT NULL
		DO UPDATE SET
			name                  = EXCLUDED.name,
			active_ingredient     = COALESCE(EXCLUDED.active_ingredient, catalog.products.active_ingredient),
			presentation          = COALESCE(EXCLUDED.presentation, catalog.products.presentation),
			form                  = COALESCE(EXCLUDED.form, catalog.products.form),
			registry_number       = COALESCE(EXCLUDED.registry_number, catalog.products.registry_number),
			manufacturer          = COALESCE(EXCLUDED.manufacturer, catalog.products.manufacturer),
			requires_prescription = EXCLUDED.requires_prescription OR catalog.products.requires_prescription,
			sku                   = COALESCE(EXCLUDED.sku, catalog.products.sku),
			updated_at            = NOW()
		RETURNING id::text
	`,
		p.CanonicalName, slug, p.ActiveIngredient,
		p.Presentation, p.Concentration, p.Form, p.RegistryNumber,
		manufacturer, p.SourceProductCode,
		p.RequiresPrescription, sku,
	).Scan(&resultID).Error
	if err != nil {
		return "", err
	}
	if resultID == "" {
		return "", fmt.Errorf("UPSERT no retornó id (source_product_code=%d, concentration=%s)", p.SourceProductCode, p.Concentration)
	}
	return resultID, nil
}

// ----- Helpers internos para UpsertBySource -----

var slugProductReplacer = regexp.MustCompile(`[^a-z0-9]+`)

func slugify(s string) string {
	s = strings.ToLower(s)
	s = slugProductReplacer.ReplaceAllString(s, "-")
	return strings.Trim(s, "-")
}

// buildProductSlugWithSource agrega source_product_code al slug para asegurar
// unicidad entre marcas con mismo (name, concentration, form).
func buildProductSlugWithSource(name, concent, formCode string, sourceCode int) string {
	base := slugify(strings.Join([]string{name, concent, formCode}, "-"))
	return fmt.Sprintf("%s-%d", base, sourceCode)
}

func firstNonEmptyStr(vals ...string) string {
	for _, v := range vals {
		if v != "" {
			return v
		}
	}
	return ""
}

func (r *ProductRepositoryImpl) SoftDelete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).
		Model(&entities.Product{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"deleted_at": gorm.Expr("CURRENT_TIMESTAMP"),
			"is_active":  false,
		}).Error
}

func (r *ProductRepositoryImpl) List(ctx context.Context, page, limit int, sort string) (*repositories.PaginatedResult, error) {
	var products []entities.Product
	var total int64

	query := r.db.WithContext(ctx).Model(&entities.Product{}).Where("deleted_at IS NULL")

	query.Count(&total)

	query = r.applySort(query, sort)
	offset := (page - 1) * limit

	err := query.
		Offset(offset).
		Limit(limit).
		Preload("Category").
		Preload("Brand").
		Preload("Images", func(db *gorm.DB) *gorm.DB {
			return db.Order("display_order ASC")
		}).
		Find(&products).Error
	if err != nil {
		return nil, err
	}

	return &repositories.PaginatedResult{
		Products:   products,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: int(math.Ceil(float64(total) / float64(limit))),
	}, nil
}

func (r *ProductRepositoryImpl) Search(ctx context.Context, params repositories.ProductSearchParams) (*repositories.PaginatedResult, error) {
	var products []entities.Product
	var total int64

	query := r.db.WithContext(ctx).Model(&entities.Product{}).Where("deleted_at IS NULL")

	if params.Query != "" {
		searchPattern := "%" + params.Query + "%"
		query = query.Where("name ILIKE ? OR active_ingredient ILIKE ?", searchPattern, searchPattern)
	}

	if params.CategoryID != "" {
		query = query.Where("category_id = ?", params.CategoryID)
	}

	if params.BrandID != "" {
		query = query.Where("brand_id = ?", params.BrandID)
	}

	// Búsqueda por DCI: case-insensitive exacta. Permite encontrar alternativas terapéuticas
	// (productos con el mismo principio activo, sean genéricos o de marca).
	if params.ActiveIngredient != "" {
		query = query.Where("UPPER(active_ingredient) = UPPER(?)", params.ActiveIngredient)
	}

	// Excluir un producto específico de los resultados (HU-015: no mostrar el mismo producto
	// como su propia alternativa).
	if params.ExcludeID != "" {
		query = query.Where("id <> ?", params.ExcludeID)
	}

	if params.RequiresPrescription != nil {
		query = query.Where("requires_prescription = ?", *params.RequiresPrescription)
	}

	if params.IsActive != nil {
		query = query.Where("is_active = ?", *params.IsActive)
	}

	query.Count(&total)

	query = r.applySort(query, params.Sort)
	offset := (params.Page - 1) * params.Limit

	err := query.
		Offset(offset).
		Limit(params.Limit).
		Preload("Category").
		Preload("Brand").
		Preload("Images", func(db *gorm.DB) *gorm.DB {
			return db.Order("display_order ASC")
		}).
		Find(&products).Error
	if err != nil {
		return nil, err
	}

	return &repositories.PaginatedResult{
		Products:   products,
		Total:      total,
		Page:       params.Page,
		Limit:      params.Limit,
		TotalPages: int(math.Ceil(float64(total) / float64(params.Limit))),
	}, nil
}

func (r *ProductRepositoryImpl) FindByCategoryID(ctx context.Context, categoryID string, page, limit int) (*repositories.PaginatedResult, error) {
	return r.Search(ctx, repositories.ProductSearchParams{
		CategoryID: categoryID,
		Page:       page,
		Limit:      limit,
		Sort:       "name_asc",
	})
}

func (r *ProductRepositoryImpl) FindByBrandID(ctx context.Context, brandID string, page, limit int) (*repositories.PaginatedResult, error) {
	return r.Search(ctx, repositories.ProductSearchParams{
		BrandID: brandID,
		Page:    page,
		Limit:   limit,
		Sort:    "name_asc",
	})
}

func (r *ProductRepositoryImpl) applySort(query *gorm.DB, sort string) *gorm.DB {
	switch sort {
	case "name_asc":
		return query.Order("name ASC")
	case "name_desc":
		return query.Order("name DESC")
	case "created_at_asc":
		return query.Order("created_at ASC")
	case "created_at_desc":
		return query.Order("created_at DESC")
	default:
		return query.Order("name ASC")
	}
}

// Compile-time interface check
var _ repositories.ProductRepository = (*ProductRepositoryImpl)(nil)
