// internal/application/handlers/upload_product_images_handler.go
package handlers

import (
	"bytes"
	"context"
	"fmt"
	"path/filepath"

	"github.com/disintegration/imaging"
	"github.com/farmanexo/catalog-service/internal/application/commands"
	"github.com/farmanexo/catalog-service/internal/domain/entities"
	"github.com/farmanexo/catalog-service/internal/domain/repositories"
	"github.com/farmanexo/catalog-service/internal/domain/services"
	"github.com/farmanexo/catalog-service/internal/presentation/dto/responses"
	"github.com/farmanexo/catalog-service/internal/shared/common"
	"github.com/farmanexo/catalog-service/pkg/mediator"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

const (
	maxImageSize       = 10 * 1024 * 1024 // 10MB
	primaryImageWidth  = 800
	primaryImageHeight = 800
)

var allowedImageTypes = map[string]bool{
	"image/jpeg": true,
	"image/jpg":  true,
	"image/png":  true,
	"image/webp": true,
}

type UploadProductImagesHandler struct {
	productRepo  repositories.ProductRepository
	imageRepo    repositories.ProductImageRepository
	fileStorage  services.FileStorage
	cacheService services.CacheService
	bucket       string
	logger       *zap.Logger
}

func NewUploadProductImagesHandler(
	productRepo repositories.ProductRepository,
	imageRepo repositories.ProductImageRepository,
	fileStorage services.FileStorage,
	cacheService services.CacheService,
	bucket string,
	logger *zap.Logger,
) *UploadProductImagesHandler {
	return &UploadProductImagesHandler{
		productRepo:  productRepo,
		imageRepo:    imageRepo,
		fileStorage:  fileStorage,
		cacheService: cacheService,
		bucket:       bucket,
		logger:       logger,
	}
}

func (h *UploadProductImagesHandler) Handle(ctx context.Context, cmd commands.UploadProductImagesCommand) (*common.ApiResponse[responses.ImageListResponse], error) {
	// Verify product exists
	_, err := h.productRepo.FindByID(ctx, cmd.ProductID)
	if err != nil {
		return common.NotFoundResponse[responses.ImageListResponse]("Producto no encontrado"), nil
	}

	// Validate files
	for _, file := range cmd.Files {
		if !allowedImageTypes[file.ContentType] {
			return common.BadRequestResponse[responses.ImageListResponse]("VAL_004", "Tipo de archivo no permitido. Use: jpg, jpeg, png, webp"), nil
		}
		if file.Size > maxImageSize {
			return common.BadRequestResponse[responses.ImageListResponse]("VAL_005", "El archivo excede el tamaño máximo de 10MB"), nil
		}
	}

	// Get current image count for display_order
	currentCount, _ := h.imageRepo.CountByProductID(ctx, cmd.ProductID)

	var uploadedImages []responses.ImageResponse

	for i, file := range cmd.Files {
		// Resize image
		resizedBuf, contentType, err := h.resizeImage(file.Reader, file.Filename, primaryImageWidth, primaryImageHeight)
		if err != nil {
			h.logger.Error("Error redimensionando imagen", zap.String("filename", file.Filename), zap.Error(err))
			continue
		}

		// Upload to S3
		ext := filepath.Ext(file.Filename)
		key := fmt.Sprintf("products/%s/%s%s", cmd.ProductID, uuid.New().String(), ext)
		url, err := h.fileStorage.Upload(ctx, h.bucket, key, resizedBuf, contentType)
		if err != nil {
			h.logger.Error("Error subiendo imagen a S3", zap.Error(err))
			continue
		}

		// If this is marked as primary, clear existing primaries
		if file.IsPrimary {
			_ = h.imageRepo.ClearPrimaryByProductID(ctx, cmd.ProductID)
		}

		// Save to DB
		image := &entities.ProductImage{
			ID:           uuid.New().String(),
			ProductID:    cmd.ProductID,
			ImageURL:     url,
			IsPrimary:    file.IsPrimary,
			DisplayOrder: int(currentCount) + i,
		}

		if err := h.imageRepo.Create(ctx, image); err != nil {
			h.logger.Error("Error guardando imagen en DB", zap.Error(err))
			continue
		}

		uploadedImages = append(uploadedImages, responses.ImageResponse{
			ID:           image.ID,
			ImageURL:     image.ImageURL,
			IsPrimary:    image.IsPrimary,
			DisplayOrder: image.DisplayOrder,
		})
	}

	// Invalidate product cache
	go func() {
		_ = h.cacheService.Delete(context.Background(), "cache:product:"+cmd.ProductID)
	}()

	response := responses.ImageListResponse{
		Images: uploadedImages,
	}

	h.logger.Info("Imágenes subidas exitosamente",
		zap.String("product_id", cmd.ProductID),
		zap.Int("count", len(uploadedImages)),
	)

	return common.CreatedResponse(response), nil
}

func (h *UploadProductImagesHandler) resizeImage(reader interface{}, filename string, width, height int) (*bytes.Buffer, string, error) {
	var buf bytes.Buffer

	// Read all data from reader
	r, ok := reader.(interface{ Read([]byte) (int, error) })
	if !ok {
		return nil, "", fmt.Errorf("invalid reader type")
	}

	var data bytes.Buffer
	if _, err := data.ReadFrom(r); err != nil {
		return nil, "", fmt.Errorf("error reading image data: %w", err)
	}

	img, err := imaging.Decode(bytes.NewReader(data.Bytes()))
	if err != nil {
		return nil, "", fmt.Errorf("error decoding image: %w", err)
	}

	resized := imaging.Fill(img, width, height, imaging.Center, imaging.Lanczos)

	ext := filepath.Ext(filename)
	var format imaging.Format
	contentType := "image/jpeg"

	switch ext {
	case ".png":
		format = imaging.PNG
		contentType = "image/png"
	default:
		format = imaging.JPEG
		contentType = "image/jpeg"
	}

	if err := imaging.Encode(&buf, resized, format); err != nil {
		return nil, "", fmt.Errorf("error encoding resized image: %w", err)
	}

	return &buf, contentType, nil
}

var _ mediator.RequestHandler[commands.UploadProductImagesCommand, responses.ImageListResponse] = (*UploadProductImagesHandler)(nil)
