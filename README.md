# Catalog Service

Microservicio de catalogo de productos farmaceuticos para FarmaNexo. Gestiona productos, categorias, marcas e imagenes de productos. Incluye busqueda avanzada, caching con Redis y publicacion de eventos via SQS.

## Inicio Rapido

### Prerequisitos
- Go 1.25+
- PostgreSQL 16
- Redis 7
- LocalStack (desarrollo local)
- Docker & Docker Compose

### Instalacion
```bash
# Clonar repositorio
git clone <url>
cd services/catalog-service

# Instalar dependencias
go mod download

# Configurar ambiente local
cp configs/config.development.yaml configs/config.local.yaml
# Editar configs/config.local.yaml con tus credenciales

# Crear base de datos
docker exec -it farmanexo-postgres psql -U admin -c "CREATE DATABASE catalog_db;"

# Ejecutar migraciones
make migrate-up

# Ejecutar servicio
make dev
```

Swagger UI disponible en: http://localhost:4003/swagger/index.html

## Endpoints

### Publicos

**GET /api/v1/products** - Listar productos (paginado)
```bash
curl "http://localhost:4003/api/v1/products?page=1&limit=20"
```

**GET /api/v1/products/{id}** - Detalle de producto (cacheado, 1hr TTL)
```bash
curl http://localhost:4003/api/v1/products/{id}
```

Respuesta (200):
```json
{
  "meta": {
    "mensajes": [{ "codigo": "CAT_001", "mensaje": "Producto obtenido exitosamente", "tipo": "exito" }],
    "idTransaccion": "...",
    "resultado": true,
    "timestamp": "20260222 103000"
  },
  "datos": {
    "id": "uuid",
    "name": "Aspirina 500mg",
    "slug": "aspirina-500mg",
    "description": "Analgesico y antipiretico",
    "active_ingredient": "Acido acetilsalicilico",
    "presentation": "Tabletas",
    "concentration": "500mg",
    "requires_prescription": false,
    "category": { "id": "uuid", "name": "Analgesicos" },
    "brand": { "id": "uuid", "name": "Bayer" },
    "sku": "ASP-500-TAB",
    "barcode": "7750000000123",
    "images": [
      { "id": "uuid", "image_url": "http://...", "is_primary": true, "display_order": 0 }
    ]
  }
}
```

**POST /api/v1/products/search** - Busqueda avanzada
```bash
curl -X POST http://localhost:4003/api/v1/products/search \
  -H "Content-Type: application/json" \
  -d '{
    "query": "aspirina",
    "category_id": "uuid",
    "brand_id": "uuid",
    "requires_prescription": false,
    "page": 1,
    "limit": 20
  }'
```

**GET /api/v1/categories** - Listar categorias (cacheado, 6hr TTL)
```bash
curl http://localhost:4003/api/v1/categories
```

**GET /api/v1/categories/{id}** - Detalle de categoria
```bash
curl http://localhost:4003/api/v1/categories/{id}
```

**GET /api/v1/categories/{id}/products** - Productos por categoria
```bash
curl "http://localhost:4003/api/v1/categories/{id}/products?page=1&limit=20"
```

**GET /api/v1/brands** - Listar marcas
```bash
curl http://localhost:4003/api/v1/brands
```

**GET /api/v1/brands/{id}/products** - Productos por marca
```bash
curl "http://localhost:4003/api/v1/brands/{id}/products?page=1&limit=20"
```

**GET /api/v1/products/barcode/{barcode}** - Buscar producto por codigo de barras (cacheado, 1hr TTL)
```bash
curl http://localhost:4003/api/v1/products/barcode/7750000000123
```

**GET /api/v1/products/{id}/interactions** - Interacciones medicamentosas (cacheado, 24hr TTL)
```bash
curl http://localhost:4003/api/v1/products/{id}/interactions
```

Respuesta (200):
```json
{
  "meta": { "resultado": true },
  "datos": {
    "interactions": [
      {
        "id": "uuid",
        "product_id": "uuid",
        "product_name": "Aspirina 500mg",
        "interacts_with_product_id": "uuid",
        "interacts_with_product_name": "Warfarina 5mg",
        "severity": "grave",
        "description": "Aumenta el riesgo de sangrado",
        "recommendation": "Evitar uso concomitante"
      }
    ],
    "total": 1
  }
}
```

**GET /api/v1/products/{id}/frequently-bought-together** - Productos frecuentemente comprados juntos (cacheado, 6hr TTL)
```bash
curl "http://localhost:4003/api/v1/products/{id}/frequently-bought-together?limit=10"
```

**GET /api/v1/products/{id}/availability** - Disponibilidad en farmacias (cacheado, 5min TTL)
```bash
curl http://localhost:4003/api/v1/products/{id}/availability
```

Respuesta (200):
```json
{
  "meta": { "resultado": true },
  "datos": {
    "product_id": "uuid",
    "product_name": "Aspirina 500mg",
    "pharmacies": [
      {
        "pharmacy_id": "uuid",
        "pharmacy_name": "Farmacia San Pablo",
        "stock": 50,
        "price": 25.50,
        "is_available": true
      }
    ],
    "total_pharmacies": 1
  }
}
```

### Admin (requieren JWT + rol admin)

**POST /api/v1/products** - Crear producto
```bash
curl -X POST http://localhost:4003/api/v1/products \
  -H "Authorization: Bearer {token}" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Aspirina 500mg",
    "description": "Analgesico y antipiretico",
    "active_ingredient": "Acido acetilsalicilico",
    "presentation": "Tabletas",
    "concentration": "500mg",
    "requires_prescription": false,
    "category_id": "uuid",
    "brand_id": "uuid",
    "sku": "ASP-500-TAB",
    "barcode": "7750000000123"
  }'
```

**PUT /api/v1/products/{id}** - Actualizar producto
```bash
curl -X PUT http://localhost:4003/api/v1/products/{id} \
  -H "Authorization: Bearer {token}" \
  -H "Content-Type: application/json" \
  -d '{"name": "Aspirina 500mg Actualizada"}'
```

**DELETE /api/v1/products/{id}** - Eliminar producto (soft delete)
```bash
curl -X DELETE http://localhost:4003/api/v1/products/{id} \
  -H "Authorization: Bearer {token}"
```

**PUT /api/v1/products/{id}/images** - Subir imagenes
```bash
curl -X PUT http://localhost:4003/api/v1/products/{id}/images \
  -H "Authorization: Bearer {token}" \
  -F "images=@image1.jpg" \
  -F "images=@image2.jpg" \
  -F "primary=0"
```

- Max total: 50 MB
- Max por imagen: 10 MB
- Formatos: JPEG, PNG, GIF, WebP
- Redimensionado: max 800x800

**POST /api/v1/categories** - Crear categoria
```bash
curl -X POST http://localhost:4003/api/v1/categories \
  -H "Authorization: Bearer {token}" \
  -H "Content-Type: application/json" \
  -d '{"name": "Analgesicos", "description": "Medicamentos para el dolor"}'
```

**PUT /api/v1/categories/{id}** - Actualizar categoria
```bash
curl -X PUT http://localhost:4003/api/v1/categories/{id} \
  -H "Authorization: Bearer {token}" \
  -H "Content-Type: application/json" \
  -d '{"name": "Analgesicos y Antipiréticos"}'
```

**POST /api/v1/brands** - Crear marca
```bash
curl -X POST http://localhost:4003/api/v1/brands \
  -H "Authorization: Bearer {token}" \
  -H "Content-Type: application/json" \
  -d '{"name": "Bayer", "country": "Alemania", "website": "https://bayer.com"}'
```

**PUT /api/v1/brands/{id}** - Actualizar marca
```bash
curl -X PUT http://localhost:4003/api/v1/brands/{id} \
  -H "Authorization: Bearer {token}" \
  -H "Content-Type: application/json" \
  -d '{"name": "Bayer Healthcare"}'
```

**POST /api/v1/products/interactions** - Crear interaccion medicamentosa
```bash
curl -X POST http://localhost:4003/api/v1/products/interactions \
  -H "Authorization: Bearer {token}" \
  -H "Content-Type: application/json" \
  -d '{
    "product_id": "uuid-producto-1",
    "interacts_with_product_id": "uuid-producto-2",
    "severity": "grave",
    "description": "Aumenta el riesgo de sangrado",
    "recommendation": "Evitar uso concomitante"
  }'
```

Severidades: `leve`, `moderada`, `grave`

### Health Check

**GET /health** - Estado del servicio
```bash
curl http://localhost:4003/health
```

## Arquitectura

- **Puerto:** 4003
- **Base de datos:** `catalog_db`
- **Schema:** `catalog`
- **Patron:** Clean Architecture + CQRS + MediatR

### Capas

```
internal/
  domain/           Entidades (Product, Category, Brand, ProductImage), interfaces
  application/      Commands, Queries, Handlers, Validators, Pre/Post processors
  infrastructure/   PostgreSQL (GORM), Redis, S3, SQS, JWT
  presentation/     Controllers, Middlewares, Routes, DTOs
  shared/           ApiResponse[T], Constants, Errors
pkg/
  mediator/         Mediator CQRS generico con pipeline
  config/           Carga de configuracion por ambiente (Viper)
```

### Flujo de un request

```
HTTP Request
  -> Chi Router
    -> [Middlewares: RequestID, RealIP, Logger, Recoverer, CORS, CorrelationID]
    -> [AuthMiddleware + RequireAdmin (si es admin endpoint)]
    -> Controller
      -> Mediator.Send(Command/Query)
        -> Validator
        -> PreProcessor (SanitizeInput)
        -> Handler (+ Redis cache para queries)
        -> PostProcessor (LogAudit)
      <- ApiResponse[T]
    <- JSON Response
```

## Configuracion

### Variables de Entorno

```yaml
# configs/config.local.yaml
environment: local

server:
  host: 0.0.0.0
  port: 4003
  read_timeout: 15s
  write_timeout: 15s
  idle_timeout: 60s

database:
  host: localhost
  port: 5432
  user: admin
  password: admin
  db_name: catalog_db
  schema: catalog
  sslmode: disable
  max_open_conns: 25
  max_idle_conns: 5
  conn_max_lifetime: 5m

jwt:
  secret: "dev-super-secret-key-change-in-production-min-32-chars"
  access_token_duration: 15m
  issuer: "farmanexo-catalog-service"

redis:
  host: localhost
  port: 6379
  password: farmanexo2026
  db: 0
  max_retries: 3
  pool_size: 10

aws:
  region: us-east-1
  endpoint: "http://localhost:4566"

sqs:
  catalog_events_queue_url: "http://sqs.us-east-1.localhost.localstack.cloud:4566/000000000000/farmanexo-catalog-events"

s3:
  products_bucket: "farmanexo-products"

services:
  pharmacy_service:
    base_url: "http://localhost:4004"

log:
  level: debug
  encoding: console
```

### Variables de entorno requeridas (ambientes desplegados)

| Variable | Descripcion |
|---|---|
| `DB_HOST` | Host de PostgreSQL |
| `DB_USER` | Usuario de PostgreSQL |
| `DB_PASSWORD` | Password de PostgreSQL |
| `JWT_SECRET` | Secret para validar JWT (mismo que Auth Service) |
| `REDIS_HOST` | Host de Redis |
| `REDIS_PASSWORD` | Password de Redis |
| `AWS_REGION` | Region de AWS |
| `SQS_CATALOG_EVENTS_QUEUE_URL` | URL cola SQS |
| `PHARMACY_SERVICE_URL` | URL base del Pharmacy Service |

## Infraestructura

### PostgreSQL
- **Database:** `catalog_db`
- **Schema:** `catalog`
- **Tablas:** `products`, `categories`, `brands`, `product_images`, `drug_interactions`, `frequently_bought_together`
- **Soft delete:** Productos usan campo `deleted_at`

### Redis
- **Uso:** Cache de productos, categorias, interacciones, FBT y disponibilidad
- **Keys:**
  - `cache:product:{id}` - TTL: 1 hora
  - `cache:product:barcode:{barcode}` - TTL: 1 hora
  - `cache:catalog:categories:all` - TTL: 6 horas
  - `cache:product:{id}:interactions` - TTL: 24 horas
  - `cache:product:{id}:fbt` - TTL: 6 horas
  - `cache:product:{id}:availability` - TTL: 5 minutos
  - `cache:search:*` - Invalidado en create/update/delete
- **Invalidacion:** Automatica al crear/actualizar/eliminar productos y categorias

### S3
- **Bucket:** `farmanexo-products`
- **Path:** `products/{product_id}/{uuid}.{ext}`
- **Procesamiento:** Resize a max 800x800

### SQS
- **Publica en:** `farmanexo-catalog-events`
- **Eventos:** `PRODUCT_CREATED`, `PRODUCT_UPDATED`, `PRODUCT_DELETED`
- **No consume** eventos de otros servicios

## Eventos

| Evento | Cola | Trigger |
|---|---|---|
| `PRODUCT_CREATED` | `farmanexo-catalog-events` | Creacion de producto |
| `PRODUCT_UPDATED` | `farmanexo-catalog-events` | Actualizacion de producto |
| `PRODUCT_DELETED` | `farmanexo-catalog-events` | Eliminacion (soft delete) de producto |

Formato:
```json
{
  "event_type": "PRODUCT_CREATED",
  "product_id": "uuid",
  "timestamp": "2026-02-22T12:00:00Z",
  "metadata": {
    "source": "catalog-service",
    "version": "1.0"
  }
}
```

## Testing
```bash
# Unit tests
make test

# Tests con coverage
make test-coverage

# Generar mocks
make gen-mocks
```

## Comandos Utiles
```bash
# Desarrollo
make dev              # Ejecutar en modo desarrollo
make build            # Compilar binario a bin/catalog-service
make swagger          # Generar documentacion Swagger

# Base de datos
make migrate-up       # Aplicar migraciones pendientes
make migrate-down     # Revertir ultima migracion
make migrate-create NAME=nombre  # Crear nueva migracion

# Calidad
make lint             # Ejecutar golangci-lint
make format           # Formatear codigo con goimports

# Docker
make docker-build     # Construir imagen Docker
make docker-run       # Ejecutar container
```

## Dependencias

### Principales
- `github.com/go-chi/chi/v5` - HTTP router
- `gorm.io/gorm` - ORM
- `github.com/redis/go-redis/v9` - Redis client
- `github.com/aws/aws-sdk-go-v2` - AWS SDK (S3, SQS)
- `github.com/golang-jwt/jwt/v5` - JWT validation
- `github.com/disintegration/imaging` - Procesamiento de imagenes
- `go.uber.org/zap` - Structured logging
- `github.com/spf13/viper` - Configuracion
- `github.com/swaggo/swag` - Swagger
- `golang.org/x/text` - Generacion de slugs (remocion de acentos)

### Completas
Ver `go.mod`

## Documentacion Adicional

- [CLAUDE.md](./CLAUDE.md) - Contexto para Claude AI
- [INFRASTRUCTURE.md](./INFRASTRUCTURE.md) - Detalle de infraestructura
- [Swagger UI](http://localhost:4003/swagger/index.html) - API docs interactiva

## Estructura de Directorios

```
catalog-service/
  cmd/server/main.go                    Punto de entrada con DI
  configs/                              YAML por ambiente (5 archivos)
  migrations/                           SQL (golang-migrate, schema: catalog)
  internal/
    application/
      commands/                         CreateProduct, UpdateProduct, DeleteProduct,
                                        UploadProductImages, CreateCategory, UpdateCategory,
                                        CreateBrand, UpdateBrand, CreateDrugInteraction
      queries/                          ListProducts, GetProduct, GetProductByBarcode,
                                        SearchProducts, ListCategories, GetCategory,
                                        ListProductsByCategory, ListBrands, ListProductsByBrand,
                                        ListDrugInteractions, ListFBT, GetProductAvailability
      handlers/                         Handler por cada command/query (21 total)
      validators/                       Validadores de commands
      preprocessors/                    SanitizeInput
      postprocessors/                   LogAudit
    domain/
      entities/                         Product, Category, Brand, ProductImage,
                                        DrugInteraction, FrequentlyBoughtTogether
      events/                           Eventos de catalogo
      repositories/                     ProductRepository, CategoryRepository,
                                        BrandRepository, ProductImageRepository,
                                        DrugInteractionRepository, FBTRepository
      services/                         CacheService, EventPublisher, FileStorage,
                                        PharmacyClient
    infrastructure/
      persistence/postgres/             Repositorios GORM (6)
      cache/                            Redis cache service
      clients/                          HTTP clients (PharmacyClient)
      messaging/                        SQS event publisher
      storage/                          S3 file storage
      security/                         JWT service (validacion)
    presentation/
      dto/requests/                     CreateProductRequest, SearchRequest,
                                        CreateInteractionRequest, etc.
      dto/responses/                    ProductResponse, CategoryResponse,
                                        InteractionResponse, FBTResponse,
                                        AvailabilityResponse, etc.
      http/controllers/                 CatalogController
      http/middlewares/                 AuthMiddleware, RequireAdmin, CorrelationID
      http/routes/                      Configuracion de rutas Chi
    shared/
      common/                           ApiResponse[T], response factories
      constants/                        Codigos HTTP, message codes
      errors/                           Domain errors
  pkg/
    config/                             Carga de configuracion (Viper)
    mediator/                           Mediator CQRS generico
    logger/                             Setup de Zap logger
  docs/                                 Swagger generado
```
