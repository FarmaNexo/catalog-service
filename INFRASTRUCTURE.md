# INFRAESTRUCTURA - Catalog Service

## Resumen

Este documento describe la infraestructura especifica utilizada por Catalog Service (puerto 4003).

---

## SERVICIOS REQUERIDOS

### PostgreSQL
- **Host:** localhost:5432 (local) / RDS endpoint (cloud)
- **Database:** `catalog_db`
- **User:** admin (local) / `${DB_USER}` (cloud)
- **Password:** admin (local) / `${DB_PASSWORD}` (cloud)
- **Schema:** `catalog`
- **SSL:** disable (local) / require (produccion)

### Redis
- **Host:** localhost:6379 (local) / ElastiCache (cloud)
- **Password:** farmanexo2026 (local) / `${REDIS_PASSWORD}` (cloud)
- **DB:** 0 (compartido)
- **Pool size:** 10 (local) / 20 (produccion)
- **Max retries:** 3
- **Uso en este servicio:**
  - Cache de detalle de productos (1hr TTL)
  - Cache de lista de categorias (6hr TTL)
  - Invalidacion automatica en operaciones de escritura

### LocalStack (Local) / AWS (Cloud)
- **Endpoint:** http://localhost:4566 (local)
- **Region:** us-east-1
- **Credenciales (local):** test/test (fake)

---

## RECURSOS AWS UTILIZADOS

### S3 Buckets

**Bucket:** `farmanexo-products`
**Uso:** Almacenamiento de imagenes de productos

**Estructura:**
```
farmanexo-products/
  products/
    {product_id}/
      {uuid}.jpg     (imagen principal y adicionales)
      {uuid}.png
      {uuid}.webp
```

**Operaciones:**
- **Upload:** Recibe imagenes multipart, redimensiona a max 800x800, sube a S3
- **Download:** URL directa al objeto S3 (almacenada en `product_images.image_url`)
- **Delete:** Elimina imagenes anteriores al reemplazar (PUT /images)

**URLs generadas:**
- **Local:** `http://localhost:4566/farmanexo-products/products/{product_id}/{uuid}.{ext}`
- **Cloud:** `https://farmanexo-products.s3.{region}.amazonaws.com/products/{product_id}/{uuid}.{ext}`

**Restricciones:**
- Tamano maximo total: 50 MB
- Tamano maximo por imagen: 10 MB
- Tipos permitidos: `image/jpeg`, `image/png`, `image/gif`, `image/webp`
- Dimensiones output: max 800x800 px

### SQS Queues

**Cola que PUBLICA:**

**Cola:** `farmanexo-catalog-events`
- **URL (local):** `http://sqs.us-east-1.localhost.localstack.cloud:4566/000000000000/farmanexo-catalog-events`
- **URL (cloud):** `${SQS_CATALOG_EVENTS_QUEUE_URL}`

**Eventos que genera:**

1. **PRODUCT_CREATED** - Cuando se crea un producto nuevo
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

2. **PRODUCT_UPDATED** - Cuando se actualiza un producto
```json
{
  "event_type": "PRODUCT_UPDATED",
  "product_id": "uuid",
  "timestamp": "2026-02-22T12:00:00Z",
  "metadata": {
    "source": "catalog-service",
    "version": "1.0"
  }
}
```

3. **PRODUCT_DELETED** - Cuando se elimina un producto (soft delete)
```json
{
  "event_type": "PRODUCT_DELETED",
  "product_id": "uuid",
  "timestamp": "2026-02-22T12:00:00Z",
  "metadata": {
    "source": "catalog-service",
    "version": "1.0"
  }
}
```

**Patron de publicacion:** Fire-and-forget en goroutines.

**Catalog Service no consume eventos de ninguna cola.**

---

## ESQUEMA DE BASE DE DATOS

### Tabla: `catalog.categories`

```sql
CREATE TABLE catalog.categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(255) UNIQUE,
    description TEXT,
    parent_id UUID REFERENCES catalog.categories(id),
    image_url VARCHAR(500),
    is_active BOOLEAN DEFAULT true,
    display_order INTEGER DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_categories_parent_id ON catalog.categories(parent_id);
CREATE UNIQUE INDEX idx_categories_slug ON catalog.categories(slug);
```

**Proposito:** Categorias jerarquicas de productos (parent_id permite arbol).

### Tabla: `catalog.brands`

```sql
CREATE TABLE catalog.brands (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) UNIQUE NOT NULL,
    slug VARCHAR(255) UNIQUE,
    description TEXT,
    logo_url VARCHAR(500),
    website VARCHAR(500),
    country VARCHAR(100),
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX idx_brands_slug ON catalog.brands(slug);
```

**Proposito:** Marcas farmaceuticas. Cada producto pertenece a una marca.

### Tabla: `catalog.products`

```sql
CREATE TABLE catalog.products (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(500) NOT NULL,
    slug VARCHAR(500) UNIQUE,
    description TEXT,
    active_ingredient VARCHAR(500),
    presentation VARCHAR(255),
    concentration VARCHAR(100),
    requires_prescription BOOLEAN DEFAULT false,
    category_id UUID REFERENCES catalog.categories(id),
    brand_id UUID REFERENCES catalog.brands(id),
    sku VARCHAR(100) UNIQUE,
    barcode VARCHAR(100),
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_products_category_id ON catalog.products(category_id);
CREATE INDEX idx_products_brand_id ON catalog.products(brand_id);
CREATE UNIQUE INDEX idx_products_slug ON catalog.products(slug);
CREATE INDEX idx_products_active_ingredient ON catalog.products(active_ingredient);
CREATE INDEX idx_products_is_active ON catalog.products(is_active);
```

**Proposito:** Productos farmaceuticos. Soft delete via `deleted_at`.

**Slugs:** Auto-generados desde el nombre, con remocion de acentos (golang.org/x/text). Si hay duplicado, se agrega sufijo UUID.

**Campos farmaceuticos:**
- `active_ingredient` - Principio activo (ej: "Acido acetilsalicilico")
- `presentation` - Forma farmaceutica (ej: "Tabletas", "Jarabe", "Capsulas")
- `concentration` - Dosis (ej: "500mg", "10mg/5ml")
- `requires_prescription` - Si requiere receta medica

### Tabla: `catalog.product_images`

```sql
CREATE TABLE catalog.product_images (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id UUID NOT NULL REFERENCES catalog.products(id) ON DELETE CASCADE,
    image_url VARCHAR(500) NOT NULL,
    is_primary BOOLEAN DEFAULT false,
    display_order INTEGER DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_product_images_product_id ON catalog.product_images(product_id);
```

**Proposito:** Imagenes de productos. Una puede ser marcada como primaria. Cascade delete con producto.

### Tabla: `catalog.drug_interactions`

```sql
CREATE TABLE catalog.drug_interactions (
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
```

**Proposito:** Interacciones medicamentosas entre productos. Severidades: leve, moderada, grave.

### Tabla: `catalog.frequently_bought_together`

```sql
CREATE TABLE catalog.frequently_bought_together (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id UUID NOT NULL REFERENCES catalog.products(id) ON DELETE CASCADE,
    related_product_id UUID NOT NULL REFERENCES catalog.products(id) ON DELETE CASCADE,
    score DECIMAL(5,2) NOT NULL DEFAULT 0.00,
    co_purchase_count INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT uq_fbt_pair UNIQUE (product_id, related_product_id)
);

CREATE INDEX idx_fbt_product_id ON catalog.frequently_bought_together(product_id);
CREATE INDEX idx_fbt_related_product_id ON catalog.frequently_bought_together(related_product_id);
CREATE INDEX idx_fbt_score ON catalog.frequently_bought_together(score DESC);
```

**Proposito:** Relaciones de productos frecuentemente comprados juntos. Score indica relevancia (0-100).

---

## COMUNICACION CON OTROS SERVICIOS

### Pharmacy Service (HTTP)
- **URL:** `http://localhost:4004` (local) / `${PHARMACY_SERVICE_URL}` (cloud)
- **Endpoint consumido:** `GET /api/v1/pharmacies/inventory/product/{productId}`
- **Timeout:** 3 segundos
- **Uso:** Consulta de disponibilidad en tiempo real de productos en farmacias
- **Degradacion graceful:** Si Pharmacy Service no responde, retorna lista vacia

---

## CONFIGURACION POR AMBIENTE

### Local (config.local.yaml)
```yaml
environment: local
server:
  port: 4003
  read_timeout: 15s
  write_timeout: 15s
database:
  host: localhost
  user: admin
  password: admin
  db_name: catalog_db
  schema: catalog
  sslmode: disable
  max_open_conns: 25
jwt:
  secret: "dev-super-secret-key-change-in-production-min-32-chars"
redis:
  host: localhost
  password: farmanexo2026
  pool_size: 10
aws:
  endpoint: "http://localhost:4566"
sqs:
  catalog_events_queue_url: "http://sqs.us-east-1.localhost.localstack.cloud:4566/000000000000/farmanexo-catalog-events"
s3:
  products_bucket: "farmanexo-products"
log:
  level: debug
  encoding: console
```

### Development (config.development.yaml)
```yaml
environment: development
database:
  host: ${DB_HOST}
  user: ${DB_USER}
  password: ${DB_PASSWORD}
  sslmode: require
jwt:
  secret: ${JWT_SECRET}
aws:
  endpoint: ""  # AWS real
log:
  level: info
  encoding: json
```

### Production (config.production.yaml)
```yaml
environment: production
server:
  read_timeout: 10s
  write_timeout: 10s
  idle_timeout: 120s
database:
  max_open_conns: 50
  max_idle_conns: 10
  sslmode: require
redis:
  pool_size: 20
log:
  level: warn
  encoding: json
```

---

## SECRETS Y CREDENCIALES

### Secrets Manager (Produccion)
- `farmanexo/auth/jwt-secret` - JWT validation key (mismo que Auth Service)
- `farmanexo/database/password` - Database password

### Variables de Entorno (Ambientes desplegados)

| Variable | Descripcion |
|---|---|
| `ENV` | local, development, qa, uat, production |
| `DB_HOST` | Host PostgreSQL |
| `DB_USER` | Usuario PostgreSQL |
| `DB_PASSWORD` | Password PostgreSQL |
| `JWT_SECRET` | JWT secret (mismo que Auth Service) |
| `REDIS_HOST` | Host Redis |
| `REDIS_PASSWORD` | Password Redis |
| `AWS_REGION` | Region AWS |
| `SQS_CATALOG_EVENTS_QUEUE_URL` | URL cola SQS |

---

## CACHE REDIS - PATRONES

### Keys Utilizadas

| Pattern | TTL | Uso |
|---------|-----|-----|
| `cache:product:{id}` | 1 hora | Detalle de producto individual |
| `cache:product:barcode:{barcode}` | 1 hora | Producto por codigo de barras |
| `cache:catalog:categories:all` | 6 horas | Lista completa de categorias |
| `cache:product:{id}:interactions` | 24 horas | Interacciones medicamentosas del producto |
| `cache:product:{id}:fbt` | 6 horas | Productos frecuentemente comprados juntos |
| `cache:product:{id}:availability` | 5 minutos | Disponibilidad en farmacias (tiempo real) |
| `cache:search:*` | Variable | Resultados de busqueda (pattern delete) |

### Invalidacion de Cache
- **Crear producto:** Invalida `cache:search:*` (pattern delete)
- **Actualizar producto:** Invalida `cache:catalog:product:{id}` + `cache:search:*`
- **Eliminar producto:** Invalida `cache:catalog:product:{id}` + `cache:search:*`
- **Crear/actualizar categoria:** Invalida `cache:catalog:categories:all`

### Implementacion
- Servicio: `RedisCacheService`
- Metodos: `Get()`, `Set()`, `Delete()`, `DeleteByPattern()`
- Serializacion: JSON marshal/unmarshal
- Context-aware con structured logging

---

## EVENTOS SQS

### Flujo de Eventos

```
[Catalog Service] --PRODUCT_CREATED--> [farmanexo-catalog-events] --> [Pharmacy Service] (prepara inventario)
[Catalog Service] --PRODUCT_UPDATED--> [farmanexo-catalog-events] --> [Pharmacy Service]
[Catalog Service] --PRODUCT_DELETED--> [farmanexo-catalog-events] --> [Pharmacy Service]
```

### Procesamiento
- Los eventos se publican de forma asincrona en goroutines
- Fire-and-forget: no bloquea el flujo principal
- Los errores de publicacion se loguean con `logger.Warn()`

---

## DESPLIEGUE

### Checklist Pre-Deploy
- [ ] Migraciones aplicadas
- [ ] Variables de entorno configuradas
- [ ] JWT_SECRET identico al de Auth Service
- [ ] S3 bucket `farmanexo-products` creado
- [ ] SQS queue `farmanexo-catalog-events` creada
- [ ] Redis accesible
- [ ] PostgreSQL accesible con database `catalog_db`

### Comandos de Deploy
```bash
# Build
make build

# Migraciones
make migrate-up ENV=production

# Docker
make docker-build
make docker-run
```

---

## TESTING LOCAL

### 1. Levantar Infraestructura
```bash
cd FarmaNexo/Helpers
./start-local.sh --full
./init-localstack-resources.sh
```

### 2. Verificar Servicios
```bash
# PostgreSQL
docker exec -it farmanexo-postgres psql -U admin -d catalog_db

# Redis
docker exec -it farmanexo-redis redis-cli -a farmanexo2026

# LocalStack S3
aws --endpoint-url=http://localhost:4566 s3 ls s3://farmanexo-products/

# LocalStack SQS
aws --endpoint-url=http://localhost:4566 sqs list-queues
```

### 3. Crear Base de Datos
```bash
docker exec -it farmanexo-postgres psql -U admin -c "CREATE DATABASE catalog_db;"
```

### 4. Ejecutar Migraciones
```bash
cd services/catalog-service
make migrate-up
```

### 5. Ejecutar Servicio
```bash
make dev
```

---

## MONITOREO

### Metricas Importantes
- Tasa de busquedas (search queries/min)
- Cache hit rate (Redis)
- Tamano del catalogo (total productos activos)
- Latencia de busqueda avanzada
- Tamano promedio de imagenes subidas

### Logs
- **Formato:** Console (local/dev), JSON (produccion)
- **Logger:** Zap structured logging
- **Campos contextuales:** product_id, category_id, brand_id, correlation_id, sku

### Alertas Recomendadas
- Rate de errores 5xx > 1%
- Cache miss rate > 80%
- Latencia de busqueda p99 > 3s
- S3 upload failures
- Redis no disponible

---

## TROUBLESHOOTING

### Problema: "Product not found" despues de crear
**Sintoma:** GET /products/{id} retorna 404 justo despues de POST
**Causa:** Cache desactualizada (poco probable) o ID incorrecto
**Solucion:** Verificar que el ID retornado en la creacion sea correcto. La cache se invalida automaticamente al crear.

### Problema: "Image upload failed"
**Sintoma:** Error al subir imagenes de producto
**Causa:** S3 bucket no existe o LocalStack no esta corriendo
**Solucion:**
```bash
aws --endpoint-url=http://localhost:4566 s3 mb s3://farmanexo-products
```

### Problema: "Duplicate SKU"
**Sintoma:** Error 409 al crear producto
**Causa:** Ya existe un producto con ese SKU
**Solucion:** Usar un SKU unico. SKUs son unicos en toda la tabla de productos.

### Problema: "Categories cache stale"
**Sintoma:** Categorias nuevas no aparecen en listado
**Causa:** Cache de 6 horas no se invalido
**Solucion:** La cache se invalida automaticamente al crear/actualizar categorias. Si persiste, verificar conexion a Redis.

### Problema: "Slug conflict"
**Sintoma:** Error al crear producto con nombre similar
**Causa:** Slug generado ya existe
**Solucion:** El sistema agrega sufijo UUID automaticamente en caso de conflicto.

---

## Referencias

- [PostgreSQL Documentation](https://www.postgresql.org/docs/)
- [Redis Documentation](https://redis.io/documentation)
- [AWS S3 Documentation](https://docs.aws.amazon.com/s3/)
- [AWS SQS Documentation](https://docs.aws.amazon.com/sqs/)
- [LocalStack Documentation](https://docs.localstack.cloud/)
- [GORM Documentation](https://gorm.io/docs/)
- [disintegration/imaging](https://github.com/disintegration/imaging)

---

Ultima actualizacion: 2026-02-22
