// internal/application/handlers/get_product_availability_handler.go
package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/farmanexo/catalog-service/internal/application/queries"
	"github.com/farmanexo/catalog-service/internal/domain/repositories"
	"github.com/farmanexo/catalog-service/internal/domain/services"
	"github.com/farmanexo/catalog-service/internal/presentation/dto/responses"
	"github.com/farmanexo/catalog-service/internal/shared/common"
	"github.com/farmanexo/catalog-service/pkg/mediator"
	"go.uber.org/zap"
)

type GetProductAvailabilityHandler struct {
	productRepo    repositories.ProductRepository
	pharmacyClient services.PharmacyClient
	cacheService   services.CacheService
	logger         *zap.Logger
}

func NewGetProductAvailabilityHandler(
	productRepo repositories.ProductRepository,
	pharmacyClient services.PharmacyClient,
	cacheService services.CacheService,
	logger *zap.Logger,
) *GetProductAvailabilityHandler {
	return &GetProductAvailabilityHandler{
		productRepo:    productRepo,
		pharmacyClient: pharmacyClient,
		cacheService:   cacheService,
		logger:         logger,
	}
}

func (h *GetProductAvailabilityHandler) Handle(ctx context.Context, query queries.GetProductAvailabilityQuery) (*common.ApiResponse[responses.AvailabilityResponse], error) {
	if query.ProductID == "" {
		return common.BadRequestResponse[responses.AvailabilityResponse]("VAL_001", "El ID del producto es requerido"), nil
	}

	// Verify product exists
	product, err := h.productRepo.FindByID(ctx, query.ProductID)
	if err != nil {
		return common.NotFoundResponse[responses.AvailabilityResponse]("Producto no encontrado"), nil
	}

	geo := services.AvailabilityGeo{Lat: query.Latitude, Lng: query.Longitude, RadiusKm: query.RadiusKm}

	// Try cache first (short TTL: 5 minutes). Bucket geo a 3 decimales (~110m)
	// para evitar explosión de claves cuando la geo del browser oscila.
	cacheKey := buildAvailabilityCacheKey(query.ProductID, geo)
	if cached, err := h.cacheService.Get(ctx, cacheKey); err == nil && cached != "" {
		var availResp responses.AvailabilityResponse
		if err := json.Unmarshal([]byte(cached), &availResp); err == nil {
			h.logger.Debug("Disponibilidad obtenida de caché", zap.String("product_id", query.ProductID))
			return common.OkResponse(availResp), nil
		}
	}

	// Call Pharmacy Service
	items, err := h.pharmacyClient.GetProductAvailability(ctx, query.ProductID, geo)
	if err != nil {
		h.logger.Warn("Error consultando disponibilidad en Pharmacy Service",
			zap.String("product_id", query.ProductID),
			zap.Error(err),
		)
		// Return empty list instead of error (graceful degradation)
		availResp := responses.AvailabilityResponse{
			ProductID:   query.ProductID,
			ProductName: product.Name,
			Pharmacies:  []responses.PharmacyAvailability{},
			Total:       0,
		}
		return common.OkResponse(availResp), nil
	}

	pharmacies := make([]responses.PharmacyAvailability, len(items))
	for i, item := range items {
		pharmacies[i] = responses.PharmacyAvailability{
			PharmacyID:       item.PharmacyID,
			PharmacySlug:     item.PharmacySlug,
			PharmacyName:     item.PharmacyName,
			PharmacyDistrict: item.PharmacyDistrict,
			PharmacyAddress:  item.PharmacyAddress,
			Stock:            item.Stock,
			Price:            item.Price,
			IsAvailable:      item.IsAvailable,
			DistanceKm:       item.DistanceKm,
			DistrictAvgPrice: item.DistrictAvgPrice,
			IsOverpriced:     item.IsOverpriced,
			OverpricePct:     item.OverpricePct,
		}
	}

	availResp := responses.AvailabilityResponse{
		ProductID:   query.ProductID,
		ProductName: product.Name,
		Pharmacies:  pharmacies,
		Total:       len(pharmacies),
	}

	// Cache for 5 minutes
	if data, err := json.Marshal(availResp); err == nil {
		go func() {
			_ = h.cacheService.Set(context.Background(), cacheKey, string(data), 5*time.Minute)
		}()
	}

	return common.OkResponse(availResp), nil
}

func buildAvailabilityCacheKey(productID string, geo services.AvailabilityGeo) string {
	if !geo.IsActive() {
		return "cache:product:" + productID + ":availability"
	}
	return fmt.Sprintf("cache:product:%s:availability:geo:%.3f:%.3f:%.1f",
		productID, geo.Lat, geo.Lng, geo.RadiusKm)
}

var _ mediator.RequestHandler[queries.GetProductAvailabilityQuery, responses.AvailabilityResponse] = (*GetProductAvailabilityHandler)(nil)
