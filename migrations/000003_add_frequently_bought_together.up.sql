-- Productos frecuentemente comprados juntos
CREATE TABLE IF NOT EXISTS catalog.frequently_bought_together (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id UUID NOT NULL REFERENCES catalog.products(id) ON DELETE CASCADE,
    related_product_id UUID NOT NULL REFERENCES catalog.products(id) ON DELETE CASCADE,
    score DECIMAL(5,2) NOT NULL DEFAULT 0.00,
    co_purchase_count INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT uq_fbt_pair UNIQUE (product_id, related_product_id),
    CONSTRAINT chk_fbt_different CHECK (product_id != related_product_id)
);

CREATE INDEX idx_fbt_product_id ON catalog.frequently_bought_together(product_id);
CREATE INDEX idx_fbt_related_product_id ON catalog.frequently_bought_together(related_product_id);
CREATE INDEX idx_fbt_score ON catalog.frequently_bought_together(score DESC);

COMMENT ON TABLE catalog.frequently_bought_together IS 'Relaciones de productos frecuentemente comprados juntos';
COMMENT ON COLUMN catalog.frequently_bought_together.score IS 'Puntuación de relevancia (0-100)';
COMMENT ON COLUMN catalog.frequently_bought_together.co_purchase_count IS 'Cantidad de compras conjuntas registradas';
