-- Revierte 000004.

DROP INDEX IF EXISTS catalog.idx_products_form;
DROP INDEX IF EXISTS catalog.uq_products_registry_number;
DROP INDEX IF EXISTS catalog.idx_products_source_product_code;
DROP INDEX IF EXISTS catalog.uq_products_source_digemid;

ALTER TABLE catalog.products
    DROP COLUMN IF EXISTS manufacturer,
    DROP COLUMN IF EXISTS registry_number,
    DROP COLUMN IF EXISTS form,
    DROP COLUMN IF EXISTS source_product_code;
