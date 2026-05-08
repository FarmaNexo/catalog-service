// internal/infrastructure/clients/pharmacy_client_impl.go
package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	neturl "net/url"
	"strconv"
	"time"

	"github.com/farmanexo/catalog-service/internal/domain/services"
	"go.uber.org/zap"
)

// PharmacyClientImpl implementación HTTP del cliente de Pharmacy Service
type PharmacyClientImpl struct {
	baseURL    string
	httpClient *http.Client
	logger     *zap.Logger
}

func NewPharmacyClient(baseURL string, logger *zap.Logger) *PharmacyClientImpl {
	return &PharmacyClientImpl{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 3 * time.Second,
		},
		logger: logger,
	}
}

// pharmacyInventoryResponse estructura de respuesta del Pharmacy Service
type pharmacyInventoryResponse struct {
	Meta struct {
		Result bool `json:"resultado"`
	} `json:"meta"`
	Data struct {
		Items []struct {
			PharmacyID       string   `json:"pharmacy_id"`
			PharmacySlug     string   `json:"pharmacy_slug"`
			PharmacyName     string   `json:"pharmacy_name"`
			PharmacyDistrict string   `json:"pharmacy_district"`
			PharmacyAddress  string   `json:"pharmacy_address"`
			ProductID        string   `json:"product_id"`
			Stock            int      `json:"stock"`
			Price            float64  `json:"price"`
			IsAvailable      bool     `json:"is_available"`
			DistanceKm       *float64 `json:"distance_km,omitempty"`
			DistrictAvgPrice *float64 `json:"district_avg_price,omitempty"`
			IsOverpriced     bool     `json:"is_overpriced"`
			OverpricePct     *float64 `json:"overprice_pct,omitempty"`
		} `json:"items"`
	} `json:"datos"`
}

func (c *PharmacyClientImpl) GetProductAvailability(ctx context.Context, productID string, geo services.AvailabilityGeo) ([]services.PharmacyInventoryItem, error) {
	url := fmt.Sprintf("%s/api/v1/pharmacies/inventory/product/%s", c.baseURL, productID)
	if geo.IsActive() {
		// pharmacy-service acepta lat/lng/radius_km opcionales (HU-014).
		q := neturl.Values{}
		q.Set("lat", strconv.FormatFloat(geo.Lat, 'f', -1, 64))
		q.Set("lng", strconv.FormatFloat(geo.Lng, 'f', -1, 64))
		if geo.RadiusKm > 0 {
			q.Set("radius_km", strconv.FormatFloat(geo.RadiusKm, 'f', -1, 64))
		}
		url = url + "?" + q.Encode()
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creando request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		c.logger.Warn("Error consultando Pharmacy Service",
			zap.String("product_id", productID),
			zap.Error(err),
		)
		return nil, fmt.Errorf("error consultando pharmacy service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.logger.Warn("Pharmacy Service respondió con error",
			zap.String("product_id", productID),
			zap.Int("status_code", resp.StatusCode),
		)
		return []services.PharmacyInventoryItem{}, nil
	}

	var pharmacyResp pharmacyInventoryResponse
	if err := json.NewDecoder(resp.Body).Decode(&pharmacyResp); err != nil {
		return nil, fmt.Errorf("error decodificando respuesta: %w", err)
	}

	items := make([]services.PharmacyInventoryItem, len(pharmacyResp.Data.Items))
	for i, item := range pharmacyResp.Data.Items {
		items[i] = services.PharmacyInventoryItem{
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

	return items, nil
}

// Compile-time interface check
var _ services.PharmacyClient = (*PharmacyClientImpl)(nil)
