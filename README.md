# Repartners Home Task

A Go backend API for managing pack sizes and calculating optimal packaging using the chi framework with SQLite persistence.

## Features

- RESTful API using chi router
- SQLite database for persistence
- Packaging algorithm with GCD optimization and dynamic programming
- Web interface with template-based rendering and XSS protection
- Docker Compose for easy deployment
- Architecture with proper separation of concerns
- Health check endpoint
- Unit tests
- Makefile

## API Endpoints

### Pack Sizes

- `GET /api/pack-sizes` - Retrieve all pack sizes
- `PUT /api/pack-sizes` - Replace all pack sizes with new array

### Packaging

- `POST /api/package` - Calculate optimal packaging for a given number of items

### Health Check

- `GET /health` - Health check endpoint

### Web Interface (Forms)

- `GET /` - Interactive web interface for pack size management and packaging calculation
- `POST /pack-sizes` - Handle pack size form submissions
- `POST /package` - Handle packaging calculation form submissions

## Request/Response Formats

### Replace Pack Sizes

**Request:**
```json
{
  "sizes": [100, 250, 500, 1000, 2000, 5000]
}
```

**Response:**
```json
{
  "message": "Pack sizes replaced successfully",
  "count": 6
}
```

### Get All Pack Sizes

**Response:**
```json
[
  {
    "id": 1,
    "size": 100,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
]
```

### Calculate Optimal Packaging

**Request:**
```json
{
  "items": 1200
}
```

**Response:**
```json
{
  "items": 1200,
  "total_items": 1250,
  "packages": {
    "1000": 1,
    "250": 1
  },
  "details": [
    {
      "size": 1000,
      "count": 1
    },
    {
      "size": 250,
      "count": 1
    }
  ]
}
```

## Algorithm Description

The packaging algorithm uses dynamic programming:

1. **GCD reduction**: Divides all pack sizes by their greatest common divisor to reduce problem size
2. **Search limit**: Only searches up to `order_amount + largest_package_size - 1`
3. **Table building**: Builds a table of optimal package combinations for each possible total
4. **Two-step optimization**:
   - First minimizes total shipped items (prefers slight overage)
   - Then minimizes number of packages for that optimal total

## Web Interface

The application serves an HTML interface at `http://localhost:8080/`.

- Pack size management form (comma-separated values)
- Packaging calculator with result display
- Go templates with auto-escaping for XSS protection

## Running the Application

### Using Docker Compose

```bash
make docker-up
```

The API and web interface will be available at `http://localhost:8080`

### Local Development

```bash
make run
```

Or use Docker:
```bash
make docker-up
```

## Makefile Targets

```
build   - Build the application binary
run     - Run the application locally
test    - Run all tests
clean   - Clean build artifacts and test databases
docker-up    - Start with docker compose
docker-down  - Stop docker compose services
fmt     - Format Go code
deps    - Download dependencies
```

## Project Structure

```
.
- cmd/server/                # Application entry point
- internal/
-- config/                   # Configuration management
-- database/                 # Database layer
-- handlers/                 # HTTP handlers
-- models/                   # Data models
-- services/                 # Business logic
- pkg/utils/                 # Utility packages
- web/                       # Web interface files
- docker-compose.yml         # Docker Compose configuration
- Dockerfile                 # Docker image configuration
- Makefile                   # Development commands
- go.mod                     # Go module definition
- go.sum                     # Go module checksums
- .gitignore                 # Git ignore file
```

## Environment Variables

- `PORT` - Port to run the server on (default: 8080)
- `DATABASE_URL` - SQLite database file path (default: ./data/pack_sizes.db)

## Testing the API

### Replace pack sizes:
```bash
curl -X PUT http://localhost:8080/api/pack-sizes \
  -H "Content-Type: application/json" \
  -d '{"sizes": [100, 250, 500, 1000, 2000, 5000]}'
```

### Get all pack sizes:
```bash
curl http://localhost:8080/api/pack-sizes
```

### Calculate optimal packaging:
```bash
curl -X POST http://localhost:8080/api/package \
  -H "Content-Type: application/json" \
  -d '{"items": 1200}'
```

### Health check:
```bash
curl http://localhost:8080/health
```

## Technical Stack

### Backend
- Go 1.22
- Chi router
- SQLite with modernc.org/sqlite (pure-Go)

### Frontend
- HTML templates
- CSS

### Deployment
- Docker / Docker Compose
- Makefile
