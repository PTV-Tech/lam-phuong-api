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

## Configuration

### Recommended: Use `.env` file

The project uses **`.env` files** for configuration (recommended for development and simplicity).

**Environment Variables:**
- `SERVER_PORT` - Server port (default: `8080`)
- `SERVER_HOST` - Server host (default: `0.0.0.0`)
- `AIRTABLE_API_KEY` - Your Airtable API key (required)
- `AIRTABLE_BASE_ID` - Your Airtable base ID (required)

See `internal/config/README.md` for more details on configuration.

## Project Structure

```
lam-phuong-api/
├── cmd/
│   └── server/          # Application entry point
├── internal/
│   ├── airtable/        # Airtable client wrapper
│   ├── book/            # Book domain
│   ├── config/          # Configuration management
│   ├── location/        # Location domain
│   └── server/          # HTTP server setup
├── .env                 # Environment variables (gitignored)
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
```

## Security Notes

- Never commit `.env` files to git (already in `.gitignore`)
- Use environment variables in production/CI/CD
- Keep API keys secure and rotate them regularly

