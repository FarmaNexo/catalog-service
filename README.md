# FarmaNexo - Catalog Service

Microservicio de catálogo de productos farmacéuticos para FarmaNexo.

## Descripción

Gestiona el catálogo completo de productos farmacéuticos, categorías y marcas. Incluye búsqueda avanzada, caching con Redis, almacenamiento de imágenes en S3 y publicación de eventos vía SQS.

## Stack Tecnológico

- **Go 1.23+** con Chi Router
- **PostgreSQL** con GORM
- **Redis** para caching
- **AWS S3** (LocalStack local) para imágenes
- **AWS SQS** (LocalStack local) para eventos
- **JWT** para autenticación (validación únicamente)
- **Swagger/OpenAPI** para documentación

## Requisitos Previos

- Go 1.23+
- PostgreSQL (base de datos `catalog_db`)
- Redis
- LocalStack (para S3 y SQS en desarrollo local)

## Instalación

```bash
# Instalar dependencias
make install

# Aplicar migraciones
make migrate-up

# Ejecutar en modo desarrollo
make dev
```

## Endpoints

### Productos (Público)
```bash
# Listar productos
curl http://localhost:4003/api/v1/products?page=1&limit=20

# Obtener producto por ID
curl http://localhost:4003/api/v1/products/{id}

# Búsqueda avanzada
curl -X POST http://localhost:4003/api/v1/products/search \
  -H "Content-Type: application/json" \
  -d '{"query": "aspirina", "page": 1, "limit": 20}'
```

### Productos (Admin)
```bash
# Crear producto
curl -X POST http://localhost:4003/api/v1/products \
  -H "Authorization: Bearer {token}" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Aspirina 500mg",
    "description": "Analgésico y antipirético",
    "active_ingredient": "Ácido acetilsalicílico",
    "presentation": "Tabletas",
    "concentration": "500mg",
    "requires_prescription": false,
    "category_id": "uuid",
    "brand_id": "uuid",
    "sku": "ASP-500-TAB"
  }'

# Actualizar producto
curl -X PUT http://localhost:4003/api/v1/products/{id} \
  -H "Authorization: Bearer {token}" \
  -H "Content-Type: application/json" \
  -d '{"name": "Aspirina 500mg Actualizada"}'

# Eliminar producto (soft delete)
curl -X DELETE http://localhost:4003/api/v1/products/{id} \
  -H "Authorization: Bearer {token}"

# Subir imágenes
curl -X PUT http://localhost:4003/api/v1/products/{id}/images \
  -H "Authorization: Bearer {token}" \
  -F "images=@image1.jpg" \
  -F "images=@image2.jpg" \
  -F "primary=0"
```

### Categorías
```bash
# Listar categorías
curl http://localhost:4003/api/v1/categories

# Productos por categoría
curl http://localhost:4003/api/v1/categories/{id}/products

# Crear categoría (Admin)
curl -X POST http://localhost:4003/api/v1/categories \
  -H "Authorization: Bearer {token}" \
  -H "Content-Type: application/json" \
  -d '{"name": "Analgésicos", "description": "Medicamentos para el dolor"}'
```

### Marcas
```bash
# Listar marcas
curl http://localhost:4003/api/v1/brands

# Productos por marca
curl http://localhost:4003/api/v1/brands/{id}/products

# Crear marca (Admin)
curl -X POST http://localhost:4003/api/v1/brands \
  -H "Authorization: Bearer {token}" \
  -H "Content-Type: application/json" \
  -d '{"name": "Bayer", "country": "Alemania", "website": "https://bayer.com"}'
```

## Documentación API

Swagger UI disponible en: `http://localhost:4003/swagger/index.html`

## Comandos Disponibles

```bash
make help           # Ver todos los comandos
make build          # Compilar binario
make dev            # Ejecutar en modo local
make test           # Ejecutar tests
make lint           # Ejecutar linter
make swagger        # Generar docs Swagger
make migrate-up     # Aplicar migraciones
make migrate-down   # Revertir última migración
make docker-build   # Build imagen Docker
make docker-run     # Ejecutar en Docker
```

## Arquitectura

```
cmd/server/main.go          → Entry point + DI wiring
internal/
├── application/            → Commands, Queries, Handlers, Validators
├── domain/                 → Entities, Repository interfaces, Service interfaces
├── infrastructure/         → PostgreSQL, Redis, S3, SQS, JWT implementations
├── presentation/           → HTTP controllers, routes, middlewares, DTOs
└── shared/                 → ApiResponse, constants, domain errors
pkg/
├── config/                 → Configuration loading
└── mediator/               → CQRS Mediator implementation
```

## Variables de Entorno

| Variable | Descripción | Default |
|---|---|---|
| `ENV` | Ambiente (local/development/qa/uat/production) | `local` |
| `DATABASE_HOST` | Host de PostgreSQL | `localhost` |
| `DATABASE_PORT` | Puerto de PostgreSQL | `5432` |
| `REDIS_HOST` | Host de Redis | `localhost` |
| `REDIS_PORT` | Puerto de Redis | `6379` |

## Puerto

**4003** (desarrollo local)
