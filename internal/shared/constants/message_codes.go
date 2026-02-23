// internal/shared/constants/message_codes.go
package constants

// MessageCode contiene todos los códigos de respuesta del sistema
type MessageCode string

const (
	// Success codes
	CodeSuccess        MessageCode = "SUCCESS_001"
	CodeCreatedSuccess MessageCode = "SUCCESS_002"
	CodeUpdatedSuccess MessageCode = "SUCCESS_003"
	CodeDeletedSuccess MessageCode = "SUCCESS_004"

	// Catalog domain codes
	CodeProductRetrieved   MessageCode = "CAT_001"
	CodeProductCreated     MessageCode = "CAT_002"
	CodeProductUpdated     MessageCode = "CAT_003"
	CodeProductDeleted     MessageCode = "CAT_004"
	CodeProductsListed     MessageCode = "CAT_005"
	CodeCategoryRetrieved  MessageCode = "CAT_006"
	CodeCategoryCreated    MessageCode = "CAT_007"
	CodeCategoryUpdated    MessageCode = "CAT_008"
	CodeCategoriesListed   MessageCode = "CAT_009"
	CodeBrandRetrieved     MessageCode = "CAT_010"
	CodeBrandCreated       MessageCode = "CAT_011"
	CodeBrandUpdated       MessageCode = "CAT_012"
	CodeBrandsListed       MessageCode = "CAT_013"
	CodeImagesUploaded     MessageCode = "CAT_014"
	CodeSearchCompleted    MessageCode = "CAT_015"
	CodeBarcodeRetrieved   MessageCode = "CAT_016"
	CodeInteractionsListed MessageCode = "CAT_017"
	CodeInteractionCreated MessageCode = "CAT_018"
	CodeFBTListed          MessageCode = "CAT_019"
	CodeAvailabilityRetrieved MessageCode = "CAT_020"

	// Validation errors
	CodeValidationError  MessageCode = "VAL_001"
	CodeInvalidFileType  MessageCode = "VAL_004"
	CodeFileTooLarge     MessageCode = "VAL_005"
	CodeRequiredField    MessageCode = "VAL_006"
	CodeInvalidSlug      MessageCode = "VAL_009"
	CodeInvalidSort      MessageCode = "VAL_010"
	CodeInvalidPage      MessageCode = "VAL_011"

	// Authentication errors
	CodeUnauthorized       MessageCode = "AUTH_ERR_001"
	CodeInvalidToken       MessageCode = "AUTH_ERR_002"
	CodeTokenExpired       MessageCode = "AUTH_ERR_003"
	CodeForbidden          MessageCode = "AUTH_ERR_005"

	// Business errors
	CodeProductNotFound    MessageCode = "BUS_001"
	CodeCategoryNotFound   MessageCode = "BUS_002"
	CodeResourceNotFound   MessageCode = "BUS_003"
	CodeBrandNotFound      MessageCode = "BUS_005"
	CodeSlugAlreadyExists  MessageCode = "BUS_006"
	CodeSKUAlreadyExists       MessageCode = "BUS_007"
	CodeInteractionNotFound    MessageCode = "BUS_008"
	CodeBarcodeNotFound        MessageCode = "BUS_009"

	// Rate limiting
	CodeRateLimitExceeded MessageCode = "RATE_001"

	// System errors
	CodeInternalError      MessageCode = "SYS_001"
	CodeDatabaseError      MessageCode = "SYS_002"
	CodeServiceUnavailable MessageCode = "SYS_003"
	CodeStorageError       MessageCode = "SYS_004"
	CodeCacheError         MessageCode = "SYS_005"
)

// MessageDescription contiene las descripciones predefinidas
var MessageDescription = map[MessageCode]string{
	// Success
	CodeSuccess:        "Operación exitosa",
	CodeCreatedSuccess: "Recurso creado exitosamente",
	CodeUpdatedSuccess: "Recurso actualizado exitosamente",
	CodeDeletedSuccess: "Recurso eliminado exitosamente",

	// Catalog
	CodeProductRetrieved:  "Producto obtenido exitosamente",
	CodeProductCreated:    "Producto creado exitosamente",
	CodeProductUpdated:    "Producto actualizado exitosamente",
	CodeProductDeleted:    "Producto eliminado exitosamente",
	CodeProductsListed:    "Productos listados exitosamente",
	CodeCategoryRetrieved: "Categoría obtenida exitosamente",
	CodeCategoryCreated:   "Categoría creada exitosamente",
	CodeCategoryUpdated:   "Categoría actualizada exitosamente",
	CodeCategoriesListed:  "Categorías listadas exitosamente",
	CodeBrandRetrieved:    "Marca obtenida exitosamente",
	CodeBrandCreated:      "Marca creada exitosamente",
	CodeBrandUpdated:      "Marca actualizada exitosamente",
	CodeBrandsListed:      "Marcas listadas exitosamente",
	CodeImagesUploaded:        "Imágenes subidas exitosamente",
	CodeSearchCompleted:       "Búsqueda completada exitosamente",
	CodeBarcodeRetrieved:      "Producto obtenido por código de barras exitosamente",
	CodeInteractionsListed:    "Interacciones listadas exitosamente",
	CodeInteractionCreated:    "Interacción creada exitosamente",
	CodeFBTListed:             "Productos frecuentemente comprados juntos listados exitosamente",
	CodeAvailabilityRetrieved: "Disponibilidad obtenida exitosamente",

	// Validation
	CodeValidationError: "Error de validación",
	CodeInvalidFileType: "Tipo de archivo no permitido. Use: jpg, jpeg, png, webp",
	CodeFileTooLarge:    "El archivo excede el tamaño máximo de 10MB",
	CodeRequiredField:   "Campo requerido",
	CodeInvalidSlug:     "Slug inválido",
	CodeInvalidSort:     "Criterio de ordenamiento inválido",
	CodeInvalidPage:     "Paginación inválida",

	// Auth errors
	CodeUnauthorized: "No autorizado",
	CodeInvalidToken: "Token inválido",
	CodeTokenExpired: "Token expirado",
	CodeForbidden:    "No tiene permisos para esta acción",

	// Business
	CodeProductNotFound:   "Producto no encontrado",
	CodeCategoryNotFound:  "Categoría no encontrada",
	CodeResourceNotFound:  "Recurso no encontrado",
	CodeBrandNotFound:     "Marca no encontrada",
	CodeSlugAlreadyExists:      "El slug ya existe",
	CodeSKUAlreadyExists:       "El SKU ya existe",
	CodeInteractionNotFound:    "Interacción no encontrada",
	CodeBarcodeNotFound:        "Producto no encontrado con el código de barras proporcionado",

	// Rate limiting
	CodeRateLimitExceeded: "Demasiadas solicitudes. Intente nuevamente más tarde",

	// System
	CodeInternalError:      "Error interno del servidor",
	CodeDatabaseError:      "Error de base de datos",
	CodeServiceUnavailable: "Servicio no disponible",
	CodeStorageError:       "Error en servicio de almacenamiento",
	CodeCacheError:         "Error en servicio de caché",
}

// GetDescription retorna la descripción del código
func GetDescription(code MessageCode) string {
	if desc, ok := MessageDescription[code]; ok {
		return desc
	}
	return "Descripción no disponible"
}
