-- Migration 000004 — campos para poblar el catálogo desde fuentes externas (DIGEMID)
--
-- Contexto: el scraper-service ingesta data del Observatorio DIGEMID
-- (ms-opm.minsa.gob.pe). Esta migración agrega los campos que permiten:
--   1. UPSERT idempotente por (source_product_code, concentration) sin duplicar
--   2. Guardar metadata farmacéutica que hoy no tiene lugar en el schema
--
-- Los campos son aditivos (NOT NULL solo con defaults), zero breaking change.

-- source_product_code: código DIGEMID del producto (ej. 2926 = PARACETAMOL).
-- Es la clave natural que combinada con concentration identifica una variante.
ALTER TABLE catalog.products
    ADD COLUMN IF NOT EXISTS source_product_code INTEGER;

-- form: forma farmacéutica (Tableta, Cápsula, Inyectable, Jarabe, etc.).
-- Distinto de `presentation` que puede incluir tipo de envase.
ALTER TABLE catalog.products
    ADD COLUMN IF NOT EXISTS form VARCHAR(100);

-- registry_number: registro sanitario DIGEMID (ej. EE08928, EN06173).
-- Concepto legal, distinto de sku/barcode comerciales. Único globalmente.
ALTER TABLE catalog.products
    ADD COLUMN IF NOT EXISTS registry_number VARCHAR(50);

-- manufacturer: laboratorio fabricante (ej. "TITAN LABORATORIES PVT. LTD.").
-- No-FK porque una marca puede tener múltiples fabricantes.
ALTER TABLE catalog.products
    ADD COLUMN IF NOT EXISTS manufacturer VARCHAR(500);

-- Índice único parcial para UPSERT por fuente DIGEMID.
-- Permite que otras fuentes (con source_product_code NULL) convivan sin conflicto.
CREATE UNIQUE INDEX IF NOT EXISTS uq_products_source_digemid
    ON catalog.products (source_product_code, concentration)
    WHERE source_product_code IS NOT NULL;

-- Índice no-único en source_product_code para lookups rápidos.
CREATE INDEX IF NOT EXISTS idx_products_source_product_code
    ON catalog.products (source_product_code)
    WHERE source_product_code IS NOT NULL;

-- Índice en registry_number (único pero parcial — otras fuentes pueden tenerlo NULL).
CREATE UNIQUE INDEX IF NOT EXISTS uq_products_registry_number
    ON catalog.products (registry_number)
    WHERE registry_number IS NOT NULL;

-- Índice en form para filtros frecuentes (ej. "solo Tabletas").
CREATE INDEX IF NOT EXISTS idx_products_form
    ON catalog.products (form)
    WHERE form IS NOT NULL;
