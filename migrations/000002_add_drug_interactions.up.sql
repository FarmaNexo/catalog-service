-- Interacciones medicamentosas
CREATE TABLE IF NOT EXISTS catalog.drug_interactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id UUID NOT NULL REFERENCES catalog.products(id) ON DELETE CASCADE,
    interacts_with_product_id UUID NOT NULL REFERENCES catalog.products(id) ON DELETE CASCADE,
    severity VARCHAR(20) NOT NULL CHECK (severity IN ('leve', 'moderada', 'grave')),
    description TEXT NOT NULL,
    recommendation TEXT,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT uq_drug_interaction UNIQUE (product_id, interacts_with_product_id)
);

CREATE INDEX idx_drug_interactions_product_id ON catalog.drug_interactions(product_id);
CREATE INDEX idx_drug_interactions_interacts_with ON catalog.drug_interactions(interacts_with_product_id);
CREATE INDEX idx_drug_interactions_severity ON catalog.drug_interactions(severity);

COMMENT ON TABLE catalog.drug_interactions IS 'Interacciones medicamentosas entre productos';
COMMENT ON COLUMN catalog.drug_interactions.severity IS 'Severidad: leve, moderada, grave';
COMMENT ON COLUMN catalog.drug_interactions.recommendation IS 'Recomendación para el profesional de salud';
