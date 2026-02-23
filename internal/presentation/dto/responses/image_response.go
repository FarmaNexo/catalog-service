// internal/presentation/dto/responses/image_response.go
package responses

// ImageResponse DTO de respuesta de imagen
type ImageResponse struct {
	ID           string `json:"id"`
	ImageURL     string `json:"image_url"`
	IsPrimary    bool   `json:"is_primary"`
	DisplayOrder int    `json:"display_order"`
}

// ImageListResponse DTO de respuesta de lista de imágenes
type ImageListResponse struct {
	Images []ImageResponse `json:"images"`
}

// EmptyResponse DTO vacío para respuestas sin datos
type EmptyResponse struct{}
