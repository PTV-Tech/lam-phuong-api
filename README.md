# Lam Phuong API

A Go API server with Airtable integration.

## Quick Start

### 1. Set up environment variables

Create a `.env` file in the root directory:

```bash
# Copy example (if available) or create manually
cp .env.example .env
```

Edit `.env` with your values:

```env
SERVER_PORT=8080
SERVER_HOST=0.0.0.0
AIRTABLE_API_KEY=your_api_key_here
AIRTABLE_BASE_ID=your_base_id_here
```

### 2. Get Airtable credentials

- **API Key**: Get from [Airtable Account](https://airtable.com/account)
- **Base ID**: Found in your base's API documentation or URL: `https://airtable.com/[BASE_ID]/...`

### 3. Run the server

```bash
go run cmd/server/main.go
```

The server will start on `http://0.0.0.0:8080` (or your configured port).

### 4. Access Swagger Documentation

Once the server is running, access the Swagger UI at:
- **Swagger UI**: `http://localhost:8080/swagger/index.html`

The Swagger UI provides interactive API documentation where you can:
- View all available endpoints
- See request/response schemas
- Test API endpoints directly from the browser

## Configuration

### Recommended: Use `.env` file

The project uses **`.env` files** for configuration (recommended for development and simplicity).

**Environment Variables:**
- `SERVER_PORT` - Server port (default: `8080`)
- `SERVER_HOST` - Server host (default: `0.0.0.0`)
- `AIRTABLE_API_KEY` - Your Airtable API key (required)
- `AIRTABLE_BASE_ID` - Your Airtable base ID (required)
- `AIRTABLE_LOCATIONS_TABLE_NAME` - Airtable table name for locations (default: `Địa điểm`)


## API Endpoints

### Locations

- **GET** `/api/locations` - List all locations
- **POST** `/api/locations` - Create a new location
  - Body: `{ "name": "string" (required), "slug": "string" (optional) }`
  - If slug is not provided, it will be auto-generated from the name
  - If slug already exists, a unique slug will be generated with a numeric suffix
- **DELETE** `/api/locations/:slug` - Delete a location by slug

### Health Check

- **GET** `/api/ping` - Health check endpoint

For detailed API documentation with request/response schemas, visit the [Swagger UI](#4-access-swagger-documentation).

## Project Structure

```
lam-phuong-api/
├── cmd/
│   └── server/          # Application entry point
├── docs/                # Generated Swagger documentation
│   ├── docs.go          # Swagger package
│   ├── swagger.json     # OpenAPI JSON spec
│   └── swagger.yaml     # OpenAPI YAML spec
├── internal/
│   ├── airtable/        # Airtable client wrapper
│   ├── config/          # Configuration management
│   ├── location/        # Location domain
│   │   ├── handler.go   # HTTP handlers
│   │   ├── model.go     # Location model
│   │   └── repository.go # Repository implementations
│   └── server/          # HTTP server setup
├── .env                 # Environment variables (gitignored)
├── .env.example         # Example environment variables
├── .air.toml            # Air live reload configuration
└── go.mod
```

## Configuration Priority

1. **Environment variables** (highest priority)
2. **`.env` file**
3. **Default values** (lowest priority)

## Development

```bash
# Install dependencies
go mod download

# Run server
go run cmd/server/main.go

# Build
go build -o bin/server cmd/server/main.go

# Live reload (with Air)
air
```

### Regenerating Swagger Documentation

If you modify Swagger annotations in the code, regenerate the documentation:

```bash
go run github.com/swaggo/swag/cmd/swag@latest init -g cmd/server/main.go -o docs
```

Or if you have `swag` installed globally:

```bash
swag init -g cmd/server/main.go -o docs
```

### Live Reload (Air)

For automatic rebuilds and restarts during development, install and use [Air](https://github.com/cosmtrek/air):

```bash
go install github.com/cosmtrek/air@latest
air
```

Air reads the `.air.toml` file in the project root to watch Go files and restart the server when changes are detected.

## Security Notes

- Never commit `.env` files to git (already in `.gitignore`)
- Use environment variables in production/CI/CD
- Keep API keys secure and rotate them regularly

