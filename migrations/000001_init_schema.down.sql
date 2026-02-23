-- migrations/000001_init_schema.down.sql
-- Revertir migración inicial del Catalog Service

DROP TABLE IF EXISTS catalog.product_images;
DROP TABLE IF EXISTS catalog.products;
DROP TABLE IF EXISTS catalog.brands;
DROP TABLE IF EXISTS catalog.categories;
DROP SCHEMA IF EXISTS catalog;
