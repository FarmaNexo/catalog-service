// internal/application/commands/upload_product_images_command.go
package commands

import "io"

// UploadProductImagesCommand comando para subir imágenes de producto
type UploadProductImagesCommand struct {
	ProductID string
	Files     []ImageFile
}

// ImageFile representa un archivo de imagen a subir
type ImageFile struct {
	Reader      io.Reader
	Filename    string
	ContentType string
	Size        int64
	IsPrimary   bool
}

func (c UploadProductImagesCommand) GetName() string {
	return "UploadProductImagesCommand"
}
