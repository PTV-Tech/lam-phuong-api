# API Response Format

This document describes the consistent response format used across all API endpoints in the Lam Phuong API.

## Response Structure

All API responses follow a consistent structure:

### Success Response

```json
{
  "success": true,
  "data": { ... },
  "message": "Optional success message"
}
```

### Error Response

```json
{
  "success": false,
  "error": {
    "code": "ERROR_CODE",
    "message": "Human-readable error message",
    "details": {
      "additional": "context"
    }
  }
}
```

## Success Response Examples

### 1. Single Resource (GET /api/users/:id)

```json
{
  "success": true,
  "data": {
    "id": "123",
    "email": "user@example.com",
    "role": "User"
  },
  "message": "User retrieved successfully"
}
```

### 2. List Resources (GET /api/users)

```json
{
  "success": true,
  "data": [
    {
      "id": "123",
      "email": "user1@example.com",
      "role": "User"
    },
    {
      "id": "456",
      "email": "user2@example.com",
      "role": "Admin"
    }
  ],
  "message": "Users retrieved successfully"
}
```

### 3. Created Resource (POST /api/users)

```json
{
  "success": true,
  "data": {
    "id": "789",
    "email": "newuser@example.com",
    "role": "User"
  },
  "message": "User created successfully"
}
```

### 4. Authentication Response (POST /api/auth/login)

```json
{
  "success": true,
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "token_type": "Bearer",
    "expires_in": 86400,
    "user": {
      "id": "123",
      "email": "user@example.com",
      "role": "User"
    }
  },
  "message": "Login successful"
}
```

### 5. No Content Response (DELETE /api/users/:id)

```json
{
  "success": true,
  "message": "User deleted successfully"
}
```

### 6. Health Check (GET /health)

```json
{
  "success": true,
  "data": {
    "status": "ok",
    "version": "1.0.0-alpha.0"
  },
  "message": "Service is healthy"
}
```

## Error Response Examples

### 1. Validation Error (400 Bad Request)

```json
{
  "success": false,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Invalid request data",
    "details": {
      "validation_error": "Key: 'RegisterRequest.Email' Error:Field validation for 'Email' failed on the 'email' tag"
    }
  }
}
```

### 2. Unauthorized (401 Unauthorized)

```json
{
  "success": false,
  "error": {
    "code": "UNAUTHORIZED",
    "message": "Authorization header required"
  }
}
```

### 3. Invalid Token (401 Unauthorized)

```json
{
  "success": false,
  "error": {
    "code": "INVALID_TOKEN",
    "message": "Invalid or expired token"
  }
}
```

### 4. Expired Token (401 Unauthorized)

```json
{
  "success": false,
  "error": {
    "code": "EXPIRED_TOKEN",
    "message": "Token has expired"
  }
}
```

### 5. Invalid Authentication (401 Unauthorized)

```json
{
  "success": false,
  "error": {
    "code": "INVALID_AUTH",
    "message": "Invalid email or password"
  }
}
```

### 6. Forbidden (403 Forbidden)

```json
{
  "success": false,
  "error": {
    "code": "FORBIDDEN",
    "message": "Insufficient permissions. Required roles: Admin"
  }
}
```

### 7. Not Found (404 Not Found)

```json
{
  "success": false,
  "error": {
    "code": "NOT_FOUND",
    "message": "User not found"
  }
}
```

### 8. Conflict (409 Conflict)

```json
{
  "success": false,
  "error": {
    "code": "DUPLICATE_EMAIL",
    "message": "Email already registered"
  }
}
```

### 9. Internal Server Error (500 Internal Server Error)

```json
{
  "success": false,
  "error": {
    "code": "INTERNAL_ERROR",
    "message": "Failed to create user: database connection failed"
  }
}
```

## Error Codes

The following error codes are used throughout the API:

- `VALIDATION_ERROR` - Request validation failed
- `BAD_REQUEST` - Invalid request format or parameters
- `UNAUTHORIZED` - Authentication required or failed
- `INVALID_TOKEN` - Invalid JWT token format
- `EXPIRED_TOKEN` - JWT token has expired
- `INVALID_AUTH` - Invalid credentials
- `FORBIDDEN` - Insufficient permissions
- `NOT_FOUND` - Resource not found
- `CONFLICT` - Resource conflict (e.g., duplicate)
- `DUPLICATE_EMAIL` - Email already exists
- `INTERNAL_ERROR` - Internal server error

## HTTP Status Codes

The API uses standard HTTP status codes:

- `200 OK` - Successful GET, PUT, PATCH requests
- `201 Created` - Successful POST requests that create resources
- `400 Bad Request` - Invalid request data or validation errors
- `401 Unauthorized` - Authentication required or failed
- `403 Forbidden` - Insufficient permissions
- `404 Not Found` - Resource not found
- `409 Conflict` - Resource conflict
- `500 Internal Server Error` - Server-side errors

## Usage in Code

### Success Responses

```go
import "lam-phuong-api/internal/response"

// With data
response.Success(c, http.StatusOK, user, "User retrieved successfully")

// Without data
response.SuccessNoContent(c, "User deleted successfully")
```

### Error Responses

```go
// Validation error
response.ValidationError(c, "Invalid request data", map[string]interface{}{
    "validation_error": err.Error(),
})

// Unauthorized
response.Unauthorized(c, "Authorization header required")

// Not found
response.NotFound(c, "User")

// Conflict
response.DuplicateEmail(c)

// Internal error
response.InternalError(c, "Failed to process request")
```

## Benefits

1. **Consistency**: All endpoints follow the same response structure
2. **Predictability**: Clients can reliably parse responses
3. **Error Handling**: Structured error information with codes and details
4. **Debugging**: Clear error messages and codes help identify issues
5. **Type Safety**: Consistent structure enables better type checking in clients

