// internal/application/queries/get_product_by_barcode_query.go
package queries

// GetProductByBarcodeQuery consulta para obtener un producto por código de barras
type GetProductByBarcodeQuery struct {
	Barcode string `json:"barcode"`
}

func (q GetProductByBarcodeQuery) GetName() string {
	return "GetProductByBarcodeQuery"
}
