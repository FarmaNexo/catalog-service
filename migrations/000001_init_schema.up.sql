-- migrations/000001_init_schema.up.sql
-- Migración inicial del Catalog Service

-- Crear esquema dedicado para catalog
CREATE SCHEMA IF NOT EXISTS catalog;

-- Habilitar extensión UUID
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Categorías (árbol jerárquico)
CREATE TABLE IF NOT EXISTS catalog.categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(255) NOT NULL UNIQUE,
    description TEXT,
    parent_id UUID REFERENCES catalog.categories(id),
    image_url VARCHAR(500),
    is_active BOOLEAN NOT NULL DEFAULT true,
    display_order INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE INDEX idx_categories_parent_id ON catalog.categories(parent_id);
CREATE INDEX idx_categories_slug ON catalog.categories(slug);

-- Marcas farmacéuticas
CREATE TABLE IF NOT EXISTS catalog.brands (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL UNIQUE,
    slug VARCHAR(255) NOT NULL UNIQUE,
    description TEXT,
    logo_url VARCHAR(500),
    website VARCHAR(500),
    country VARCHAR(100),
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE INDEX idx_brands_slug ON catalog.brands(slug);

-- Productos
CREATE TABLE IF NOT EXISTS catalog.products (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(500) NOT NULL,
    slug VARCHAR(500) NOT NULL UNIQUE,
    description TEXT,
    active_ingredient VARCHAR(500),
    presentation VARCHAR(255),
    concentration VARCHAR(100),
    requires_prescription BOOLEAN NOT NULL DEFAULT false,
    category_id UUID REFERENCES catalog.categories(id),
    brand_id UUID REFERENCES catalog.brands(id),
    sku VARCHAR(100) UNIQUE,
    barcode VARCHAR(100),
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_products_category_id ON catalog.products(category_id);
CREATE INDEX idx_products_brand_id ON catalog.products(brand_id);
CREATE INDEX idx_products_slug ON catalog.products(slug);
CREATE INDEX idx_products_active_ingredient ON catalog.products(active_ingredient);
CREATE INDEX idx_products_is_active ON catalog.products(is_active);

-- Imágenes de productos
CREATE TABLE IF NOT EXISTS catalog.product_images (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id UUID NOT NULL REFERENCES catalog.products(id) ON DELETE CASCADE,
    image_url VARCHAR(500) NOT NULL,
    is_primary BOOLEAN NOT NULL DEFAULT false,
    display_order INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE INDEX idx_product_images_product_id ON catalog.product_images(product_id);

-- Comentarios en tablas
COMMENT ON TABLE catalog.categories IS 'Categorías de productos farmacéuticos (jerárquicas)';
COMMENT ON TABLE catalog.brands IS 'Marcas/laboratorios farmacéuticos';
COMMENT ON TABLE catalog.products IS 'Catálogo de productos farmacéuticos';
COMMENT ON TABLE catalog.product_images IS 'Imágenes de productos';

-- Comentarios en columnas importantes
COMMENT ON COLUMN catalog.categories.parent_id IS 'FK a categoría padre, NULL si es raíz';
COMMENT ON COLUMN catalog.categories.slug IS 'Slug URL-friendly único';
COMMENT ON COLUMN catalog.products.active_ingredient IS 'Principio activo del medicamento';
COMMENT ON COLUMN catalog.products.presentation IS 'Forma: Tabletas, Jarabe, Cápsulas, etc.';
COMMENT ON COLUMN catalog.products.concentration IS 'Concentración: 500mg, 10mg/5ml, etc.';
COMMENT ON COLUMN catalog.products.requires_prescription IS 'Si requiere receta médica';
COMMENT ON COLUMN catalog.products.deleted_at IS 'Soft delete timestamp';
COMMENT ON COLUMN catalog.product_images.is_primary IS 'Solo una imagen puede ser principal por producto';
