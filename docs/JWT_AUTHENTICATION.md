# JWT Authentication Implementation

## Overview

This document describes the JWT (JSON Web Token) authentication implementation for the classroom analysis system. The implementation provides secure, stateless authentication with automatic token refresh and blacklist support.

## Architecture

### Components

1. **JWT Service** (`internal/service/jwt_service.go`)
   - Token generation for both Patient and User claims
   - Token validation
   - Token refresh logic
   - Configurable expiration and refresh thresholds

2. **JWT Blacklist Service** (`internal/service/jwt_blacklist.go`)
   - Redis-based token blacklist for logout functionality
   - Automatic expiration based on token TTL

3. **JWT Middleware** (`internal/web/middleware/jwt.go`)
   - Request authentication
   - Path matching (exact, prefix, wildcard)
   - Automatic token refresh
   - Blacklist checking

4. **Auth Handler** (`internal/web/auth_handler.go`)
   - Logout endpoint

## Configuration

JWT settings are configured in `config/conf.yaml`:

```yaml
jwt:
  secret: "your-secret-key"              # JWT signing secret
  expiration_hours: 24                   # Token expiration time (hours)
  refresh_threshold_hours: 1             # Auto-refresh threshold (hours)
  ignore_paths:                          # Paths that bypass JWT auth
    - /login
    - /register
    - /api/user/login
    - /api/user/register
    - /health
    - /swagger/*
    - /metrics
    - /debug/*
```

### Configuration Parameters

- **secret**: Secret key used for signing JWT tokens
- **expiration_hours**: How long tokens are valid (default: 24 hours)
- **refresh_threshold_hours**: When to auto-refresh tokens (default: 1 hour before expiry)
- **ignore_paths**: List of paths that don't require authentication

## Path Matching

The middleware supports three types of path patterns:

1. **Exact Match**: `/login` matches only `/login`
2. **Prefix Match**: `/api/*` matches `/api/user`, `/api/user/login`, etc.
3. **Wildcard**: `/swagger/*` matches `/swagger/index.html`, `/swagger/v2/doc.json`, etc.

## Token Structure

### PatientClaims

```go
type PatientClaims struct {
    PatientId   string             // Patient identifier
    Id          primitive.ObjectID // MongoDB ObjectID
    Name        string             // Patient name
    RefreshedAt int64              // Last refresh timestamp
    TokenType   string             // "patient"
    jwt.RegisteredClaims
}
```

### UserClaims

```go
type UserClaims struct {
    UserID      primitive.ObjectID // User MongoDB ObjectID
    Phone       string             // User phone number
    Name        string             // User name
    RefreshedAt int64              // Last refresh timestamp
    TokenType   string             // "user"
    jwt.RegisteredClaims
}
```

## Token Lifecycle

### 1. Login

When a user logs in:
- Server generates a JWT token with claims
- Token is returned in the response
- Token expiration is set based on `expiration_hours`

Example login response:
```json
{
  "code": 200,
  "msg": "登录成功",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "data": { ... }
}
```

### 2. Authentication

For protected endpoints:
- Client includes token in `Authorization` header
- Middleware validates token
- Token is checked against blacklist
- Claims are stored in request context

Example request:
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

or simply:
```
Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

### 3. Automatic Token Refresh

When a token is within the refresh threshold:
- Middleware generates a new token with extended expiry
- New token is returned in `X-New-Token` response header
- Client should update stored token

Example response headers:
```
X-New-Token: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**Client Implementation:**
```javascript
// Check for token refresh in response
const newToken = response.headers.get('X-New-Token');
if (newToken) {
  // Update stored token
  localStorage.setItem('token', newToken);
}
```

### 4. Logout

To logout and invalidate a token:

**Endpoint:** `POST /api/auth/logout`

**Headers:**
```
Authorization: Bearer <token>
```

**Response:**
```json
{
  "code": 200,
  "msg": "登出成功",
  "success": true
}
```

The token is added to the Redis blacklist and cannot be used until it naturally expires.

## Error Responses

All JWT errors return a consistent format:

### 401 Unauthorized

**Missing Token:**
```json
{
  "code": 401,
  "msg": "未提供认证token"
}
```

**Invalid Token:**
```json
{
  "code": 401,
  "msg": "token无效或已过期"
}
```

**Blacklisted Token:**
```json
{
  "code": 401,
  "msg": "token已失效，请重新登录"
}
```

**Invalid Format:**
```json
{
  "code": 401,
  "msg": "token格式应为 Bearer <Token> 或 <Token>"
}
```

## Usage Examples

### Backend - Getting Claims from Context

```go
import "classroom-analysis/internal/web/middleware"

func MyHandler(c *gin.Context) {
    // Get PatientClaims
    patientClaims, exists := middleware.GetPatientClaimsFromContext(c)
    if exists {
        patientID := patientClaims.Id
        // Use patient information
    }
    
    // Get UserClaims
    userClaims, exists := middleware.GetUserClaimsFromContext(c)
    if exists {
        userID := userClaims.UserID
        // Use user information
    }
}
```

### Frontend - Token Management

```javascript
// Login
async function login(credentials) {
  const response = await fetch('/api/user/login', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(credentials)
  });
  
  const data = await response.json();
  if (data.code === 200) {
    localStorage.setItem('token', data.token);
  }
}

// Authenticated Request
async function fetchData() {
  const token = localStorage.getItem('token');
  const response = await fetch('/api/customer/list', {
    headers: {
      'Authorization': `Bearer ${token}`
    }
  });
  
  // Check for token refresh
  const newToken = response.headers.get('X-New-Token');
  if (newToken) {
    localStorage.setItem('token', newToken);
  }
  
  return response.json();
}

// Logout
async function logout() {
  const token = localStorage.getItem('token');
  await fetch('/api/auth/logout', {
    method: 'POST',
    headers: {
      'Authorization': `Bearer ${token}`
    }
  });
  
  localStorage.removeItem('token');
}
```

## Security Considerations

1. **Secret Key**: Use a strong, randomly generated secret in production
2. **HTTPS**: Always use HTTPS in production to prevent token interception
3. **Token Storage**: 
   - Frontend: Use `localStorage` or `sessionStorage`
   - Never store tokens in cookies without proper security flags
4. **Blacklist Cleanup**: Redis automatically removes expired entries
5. **Token Expiration**: Balance security (shorter) vs UX (longer)

## Testing

### Unit Tests

Run JWT-related tests:
```bash
go test ./internal/service/jwt_service_test.go
go test ./internal/service/jwt_blacklist_test.go
go test ./internal/web/middleware/jwt_test.go
```

### Manual Testing

#### Test Login
```bash
curl -X POST http://localhost:8081/api/user/login \
  -H "Content-Type: application/json" \
  -d '{"phone":"13800138000","password":"password123"}'
```

#### Test Protected Endpoint
```bash
curl -X GET http://localhost:8081/api/customer/list \
  -H "Authorization: Bearer <your-token>"
```

#### Test Logout
```bash
curl -X POST http://localhost:8081/api/auth/logout \
  -H "Authorization: Bearer <your-token>"
```

#### Test Blacklisted Token
```bash
# After logout, try to use the same token
curl -X GET http://localhost:8081/api/customer/list \
  -H "Authorization: Bearer <blacklisted-token>"
# Should return: {"code":401,"msg":"token已失效，请重新登录"}
```

## Migration Guide

### From Old JWT Implementation

1. **Update Claims**: Claims now include `RefreshedAt` and `TokenType`
2. **Check Response Headers**: Watch for `X-New-Token` header
3. **Handle Logout**: Use `/api/auth/logout` endpoint
4. **Update Config**: Add JWT section to `conf.yaml`

### Breaking Changes

- Old middleware signature changed (now requires services)
- Claims structure updated with new fields
- Path matching improved (supports wildcards)

## Troubleshooting

### Token Not Refreshing

- Check `refresh_threshold_hours` in config
- Verify token is within refresh window
- Check client is handling `X-New-Token` header

### Authentication Failing

1. Verify token format: `Bearer <token>` or just `<token>`
2. Check token hasn't been blacklisted (logout)
3. Verify token hasn't expired
4. Check path is not in `ignore_paths`

### Redis Connection Issues

- Verify Redis is running
- Check Redis connection settings in `conf.yaml`
- Test Redis connection: `redis-cli ping`

## Future Enhancements

- [ ] Refresh token support (separate from access token)
- [ ] Token revocation by user ID (logout all sessions)
- [ ] Rate limiting for login attempts
- [ ] JWT token introspection endpoint
- [ ] Support for OAuth2/OIDC flows
