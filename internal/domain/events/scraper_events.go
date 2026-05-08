// Package events — eventos del bus interno y eventos consumidos del scraper.
//
// scraper_events.go define el schema "consumer-side" de los eventos
// publicados por scraper-service en la cola compartida farmanexo-{env}-scraper-events.
// catalog-service solo procesa PRODUCT_DISCOVERED; ignora los otros 2 tipos.
package events

import (
	"encoding/json"
	"time"
)

// Tipos de eventos del scraper. Las 3 constantes están aquí para que el
// switch del consumer las pueda nombrar; solo PRODUCT_DISCOVERED tiene
// handler en catalog-service.
const (
	ScraperEventProductDiscovered   = "PRODUCT_DISCOVERED"
	ScraperEventPharmacyDiscovered  = "PHARMACY_DISCOVERED"
	ScraperEventInventoryDiscovered = "INVENTORY_DISCOVERED"
)

// ScraperEvent envelope común. Schema idéntico al publisher en scraper-service.
type ScraperEvent struct {
	EventType string          `json:"event_type"`
	SourceID  string          `json:"source_id"`
	Timestamp time.Time       `json:"timestamp"`
	Data      json.RawMessage `json:"data"`
}

// ProductDiscoveredData — payload de PRODUCT_DISCOVERED que catalog-service
// proyecta a un row de catalog.products via UPSERT por
// (source_product_code, concentration).
type ProductDiscoveredData struct {
	SourceProductCode    int    `json:"source_product_code"`
	CanonicalName        string `json:"canonical_name"`
	ActiveIngredient     string `json:"active_ingredient,omitempty"`
	Concentration        string `json:"concentration"`
	Form                 string `json:"form,omitempty"`
	SourceFormCode       string `json:"source_form_code,omitempty"`
	Presentation         string `json:"presentation,omitempty"`
	RegistryNumber       string `json:"registry_number,omitempty"`
	Manufacturer         string `json:"manufacturer,omitempty"`
	Holder               string `json:"holder,omitempty"`
	RequiresPrescription bool   `json:"requires_prescription"`
}
