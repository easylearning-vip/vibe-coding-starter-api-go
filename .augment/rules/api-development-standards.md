# API Development Standards for Vibe Coding Starter Go API

When developing HTTP APIs, handlers, routes, and middleware for this Go project, follow these comprehensive standards that ensure consistency, security, and maintainability.

## RESTful API Design Principles

### HTTP Method Usage
- **GET**: Retrieve resources (safe and idempotent)
- **POST**: Create new resources or non-idempotent operations
- **PUT**: Update entire resources (idempotent)
- **PATCH**: Partial resource updates
- **DELETE**: Remove resources (idempotent)

### URL Structure Standards
```
/api/v1/users                    # Collection operations
/api/v1/users/{id}              # Individual resource operations
/api/v1/users/{id}/articles     # Nested resource relationships
/api/v1/articles?author_id={id} # Query parameters for filtering
/api/v1/articles?page=1&size=10 # Pagination parameters
```

### HTTP Status Code Guidelines
- **200 OK**: Successful GET, PUT, PATCH operations
- **201 Created**: Successful POST operations (resource creation)
- **204 No Content**: Successful DELETE operations
- **400 Bad Request**: Invalid request data or malformed syntax
- **401 Unauthorized**: Authentication required or invalid credentials
- **403 Forbidden**: Authenticated but insufficient permissions
- **404 Not Found**: Resource does not exist
- **409 Conflict**: Resource conflict (duplicate email, username)
- **422 Unprocessable Entity**: Validation errors with valid syntax
- **429 Too Many Requests**: Rate limit exceeded
- **500 Internal Server Error**: Unexpected server errors

## Response Format Standards

### Unified Response Structure
```go
type Response struct {
    Code    int         `json:"code" example:"200"`
    Message string      `json:"message" example:"success"`
    Data    interface{} `json:"data,omitempty"`
    Meta    *Meta       `json:"meta,omitempty"`
}

type Meta struct {
    Page       int   `json:"page,omitempty" example:"1"`
    PageSize   int   `json:"page_size,omitempty" example:"10"`
    Total      int64 `json:"total,omitempty" example:"100"`
    TotalPages int   `json:"total_pages,omitempty" example:"10"`
}
```

### Success Response Examples
```json
// Single resource response
{
    "code": 200,
    "message": "success",
    "data": {
        "id": 1,
        "username": "john_doe",
        "email": "john@example.com",
        "created_at": "2024-01-01T00:00:00Z"
    }
}

// Collection response with pagination
{
    "code": 200,
    "message": "success",
    "data": [
        {"id": 1, "title": "Article 1"},
        {"id": 2, "title": "Article 2"}
    ],
    "meta": {
        "page": 1,
        "page_size": 10,
        "total": 25,
        "total_pages": 3
    }
}
```

### Error Response Format
```json
{
    "code": 422,
    "message": "validation failed",
    "data": {
        "errors": [
            {
                "field": "email",
                "message": "invalid email format"
            },
            {
                "field": "password",
                "message": "password must be at least 6 characters"
            }
        ]
    }
}
```

## Handler Implementation Patterns

### Standard Handler Structure
```go
type UserHandler struct {
    service service.UserService
    logger  logger.Logger
}

func NewUserHandler(service service.UserService, logger logger.Logger) *UserHandler {
    return &UserHandler{
        service: service,
        logger:  logger,
    }
}

func (h *UserHandler) RegisterRoutes(router *gin.RouterGroup) {
    users := router.Group("/users")
    {
        // Public routes
        users.POST("/register", h.Register)
        users.POST("/login", h.Login)
        
        // Protected routes (require authentication)
        protected := users.Group("")
        protected.Use(middleware.RequireAuth())
        {
            protected.GET("/profile", h.GetProfile)
            protected.PUT("/profile", h.UpdateProfile)
            protected.PUT("/password", h.ChangePassword)
        }
    }
}
```

### Handler Method Implementation Pattern
```go
// @Summary Register user
// @Description Create a new user account
// @Tags users
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "User registration data"
// @Success 201 {object} Response{data=User} "User registered successfully"
// @Failure 400 {object} Response "Invalid request data"
// @Failure 409 {object} Response "User already exists"
// @Failure 422 {object} Response "Validation failed"
// @Router /api/v1/users/register [post]
func (h *UserHandler) Register(c *gin.Context) {
    // 1. Extract request ID for tracing
    requestID := c.GetString("request_id")
    
    // 2. Bind and validate request
    var req RegisterRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        h.logger.Error("Failed to bind request", 
            "error", err, 
            "request_id", requestID,
            "path", c.Request.URL.Path,
        )
        c.JSON(http.StatusBadRequest, Response{
            Code:    http.StatusBadRequest,
            Message: "invalid request format",
            Data:    map[string]string{"error": err.Error()},
        })
        return
    }

    // 3. Call service layer
    user, err := h.service.Register(c.Request.Context(), &req)
    if err != nil {
        h.handleError(c, err, requestID)
        return
    }

    // 4. Log success and return response
    h.logger.Info("User registered successfully",
        "user_id", user.ID,
        "username", user.Username,
        "request_id", requestID,
    )

    c.JSON(http.StatusCreated, Response{
        Code:    http.StatusCreated,
        Message: "user registered successfully",
        Data:    user,
    })
}
```

### Centralized Error Handling
```go
func (h *UserHandler) handleError(c *gin.Context, err error, requestID string) {
    h.logger.Error("Handler error", 
        "error", err, 
        "path", c.Request.URL.Path,
        "method", c.Request.Method,
        "request_id", requestID,
    )

    switch {
    case errors.Is(err, service.ErrUserNotFound):
        c.JSON(http.StatusNotFound, Response{
            Code:    http.StatusNotFound,
            Message: "user not found",
        })
    case errors.Is(err, service.ErrUserExists):
        c.JSON(http.StatusConflict, Response{
            Code:    http.StatusConflict,
            Message: "user already exists",
        })
    case errors.Is(err, service.ErrInvalidCredentials):
        c.JSON(http.StatusUnauthorized, Response{
            Code:    http.StatusUnauthorized,
            Message: "invalid credentials",
        })
    case errors.Is(err, service.ErrValidationFailed):
        c.JSON(http.StatusUnprocessableEntity, Response{
            Code:    http.StatusUnprocessableEntity,
            Message: "validation failed",
            Data:    map[string]string{"error": err.Error()},
        })
    case errors.Is(err, service.ErrRateLimitExceeded):
        c.JSON(http.StatusTooManyRequests, Response{
            Code:    http.StatusTooManyRequests,
            Message: "rate limit exceeded",
            Data:    map[string]string{"retry_after": "60"},
        })
    default:
        c.JSON(http.StatusInternalServerError, Response{
            Code:    http.StatusInternalServerError,
            Message: "internal server error",
        })
    }
}
```

## Request and Response Models

### Request Validation Models
```go
type RegisterRequest struct {
    Username string `json:"username" validate:"required,min=3,max=50,alphanum" example:"john_doe"`
    Email    string `json:"email" validate:"required,email,max=100" example:"john@example.com"`
    Password string `json:"password" validate:"required,min=6,max=100" example:"password123"`
    Nickname string `json:"nickname" validate:"max=50" example:"John Doe"`
}

type LoginRequest struct {
    Username string `json:"username" validate:"required" example:"john_doe"`
    Password string `json:"password" validate:"required" example:"password123"`
}

type UpdateProfileRequest struct {
    Nickname string `json:"nickname" validate:"max=50" example:"John Doe"`
    Avatar   string `json:"avatar" validate:"omitempty,url" example:"https://example.com/avatar.jpg"`
}

type ChangePasswordRequest struct {
    CurrentPassword string `json:"current_password" validate:"required" example:"oldpassword"`
    NewPassword     string `json:"new_password" validate:"required,min=6,max=100" example:"newpassword123"`
}
```

### Response Models
```go
type User struct {
    ID        uint      `json:"id" example:"1"`
    Username  string    `json:"username" example:"john_doe"`
    Email     string    `json:"email" example:"john@example.com"`
    Nickname  string    `json:"nickname" example:"John Doe"`
    Avatar    string    `json:"avatar" example:"https://example.com/avatar.jpg"`
    Role      string    `json:"role" example:"user"`
    Status    string    `json:"status" example:"active"`
    CreatedAt time.Time `json:"created_at" example:"2024-01-01T00:00:00Z"`
    UpdatedAt time.Time `json:"updated_at" example:"2024-01-01T00:00:00Z"`
}

type LoginResponse struct {
    Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
    User  *User  `json:"user"`
}
```

## Middleware Implementation Patterns

### Authentication Middleware
```go
func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
    return func(c *gin.Context) {
        token := m.extractToken(c)
        if token == "" {
            c.JSON(http.StatusUnauthorized, Response{
                Code:    http.StatusUnauthorized,
                Message: "authentication required",
            })
            c.Abort()
            return
        }

        claims, err := m.jwtService.ValidateToken(token)
        if err != nil {
            m.logger.Error("Invalid token", "error", err, "token", token[:10]+"...")
            c.JSON(http.StatusUnauthorized, Response{
                Code:    http.StatusUnauthorized,
                Message: "invalid or expired token",
            })
            c.Abort()
            return
        }

        // Set user context
        c.Set("user_id", claims.UserID)
        c.Set("user_role", claims.Role)
        c.Set("username", claims.Username)
        
        c.Next()
    }
}

func (m *AuthMiddleware) extractToken(c *gin.Context) string {
    // Try Authorization header first
    authHeader := c.GetHeader("Authorization")
    if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
        return strings.TrimPrefix(authHeader, "Bearer ")
    }
    
    // Fallback to query parameter (not recommended for production)
    return c.Query("token")
}
```

### Role-Based Authorization Middleware
```go
func (m *AuthMiddleware) RequireRole(roles ...string) gin.HandlerFunc {
    return func(c *gin.Context) {
        userRole := c.GetString("user_role")
        if userRole == "" {
            c.JSON(http.StatusUnauthorized, Response{
                Code:    http.StatusUnauthorized,
                Message: "authentication required",
            })
            c.Abort()
            return
        }

        hasRole := false
        for _, role := range roles {
            if userRole == role {
                hasRole = true
                break
            }
        }

        if !hasRole {
            c.JSON(http.StatusForbidden, Response{
                Code:    http.StatusForbidden,
                Message: "insufficient permissions",
            })
            c.Abort()
            return
        }

        c.Next()
    }
}
```

### Rate Limiting Middleware
```go
func (m *RateLimitMiddleware) IPRateLimit(requests int, window time.Duration) gin.HandlerFunc {
    limiter := rate.NewLimiter(rate.Every(window/time.Duration(requests)), requests)
    
    return func(c *gin.Context) {
        clientIP := c.ClientIP()
        
        if !limiter.Allow() {
            m.logger.Warn("Rate limit exceeded", "client_ip", clientIP)
            c.JSON(http.StatusTooManyRequests, Response{
                Code:    http.StatusTooManyRequests,
                Message: "rate limit exceeded",
                Data: map[string]interface{}{
                    "retry_after": int(window.Seconds()),
                    "limit":       requests,
                },
            })
            c.Abort()
            return
        }
        
        c.Next()
    }
}
```

### CORS Middleware
```go
func (m *CORSMiddleware) CORS() gin.HandlerFunc {
    return func(c *gin.Context) {
        origin := c.Request.Header.Get("Origin")
        
        // Check if origin is allowed
        if m.isAllowedOrigin(origin) {
            c.Header("Access-Control-Allow-Origin", origin)
        }
        
        c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
        c.Header("Access-Control-Allow-Credentials", "true")
        c.Header("Access-Control-Max-Age", "86400")

        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(http.StatusNoContent)
            return
        }

        c.Next()
    }
}
```

## Pagination Implementation

### Pagination Parameters
```go
type PaginationParams struct {
    Page     int `form:"page" validate:"min=1" example:"1"`
    PageSize int `form:"page_size" validate:"min=1,max=100" example:"10"`
}

func (p *PaginationParams) Normalize() {
    if p.Page <= 0 {
        p.Page = 1
    }
    if p.PageSize <= 0 {
        p.PageSize = 10
    }
    if p.PageSize > 100 {
        p.PageSize = 100
    }
}

func (p *PaginationParams) GetOffset() int {
    return (p.Page - 1) * p.PageSize
}

func (p *PaginationParams) GetLimit() int {
    return p.PageSize
}
```

### Paginated List Handler
```go
func (h *ArticleHandler) List(c *gin.Context) {
    var params PaginationParams
    if err := c.ShouldBindQuery(&params); err != nil {
        c.JSON(http.StatusBadRequest, Response{
            Code:    http.StatusBadRequest,
            Message: "invalid pagination parameters",
            Data:    map[string]string{"error": err.Error()},
        })
        return
    }

    params.Normalize()

    // Optional filtering
    var filters service.ArticleFilters
    if status := c.Query("status"); status != "" {
        filters.Status = status
    }
    if authorID := c.Query("author_id"); authorID != "" {
        if id, err := strconv.ParseUint(authorID, 10, 32); err == nil {
            filters.AuthorID = uint(id)
        }
    }

    articles, total, err := h.service.List(c.Request.Context(), service.ListOptions{
        Offset:  params.GetOffset(),
        Limit:   params.GetLimit(),
        Filters: filters,
    })
    if err != nil {
        h.handleError(c, err, c.GetString("request_id"))
        return
    }

    totalPages := int(math.Ceil(float64(total) / float64(params.PageSize)))

    c.JSON(http.StatusOK, Response{
        Code:    http.StatusOK,
        Message: "success",
        Data:    articles,
        Meta: &Meta{
            Page:       params.Page,
            PageSize:   params.PageSize,
            Total:      total,
            TotalPages: totalPages,
        },
    })
}
```

## Swagger Documentation Standards

### Complete Swagger Example
```go
// @Summary List articles
// @Description Get a paginated list of articles with optional filtering
// @Tags articles
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1) minimum(1)
// @Param page_size query int false "Page size" default(10) minimum(1) maximum(100)
// @Param status query string false "Article status filter" Enums(draft, published, archived)
// @Param author_id query int false "Filter by author ID"
// @Success 200 {object} Response{data=[]Article,meta=Meta} "Articles retrieved successfully"
// @Failure 400 {object} Response "Invalid parameters"
// @Failure 500 {object} Response "Internal server error"
// @Router /api/v1/articles [get]
func (h *ArticleHandler) List(c *gin.Context) {
    // Implementation
}

// @Summary Create article
// @Description Create a new article (requires authentication)
// @Tags articles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateArticleRequest true "Article creation data"
// @Success 201 {object} Response{data=Article} "Article created successfully"
// @Failure 400 {object} Response "Invalid request data"
// @Failure 401 {object} Response "Authentication required"
// @Failure 422 {object} Response "Validation failed"
// @Router /api/v1/user/articles [post]
func (h *ArticleHandler) Create(c *gin.Context) {
    // Implementation
}
```

## Security Best Practices

### Input Sanitization and Validation
- Always validate and sanitize user inputs using struct tags
- Use custom validators for complex business rules
- Escape HTML content to prevent XSS attacks
- Validate file uploads (type, size, content)

### Authentication and Authorization
- Use secure JWT token storage and transmission
- Implement proper token expiration and refresh mechanisms
- Use role-based access control for fine-grained permissions
- Log authentication and authorization events

### API Security Headers
```go
func (m *SecurityMiddleware) SecurityHeaders() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Header("X-Content-Type-Options", "nosniff")
        c.Header("X-Frame-Options", "DENY")
        c.Header("X-XSS-Protection", "1; mode=block")
        c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
        c.Header("Content-Security-Policy", "default-src 'self'")
        c.Next()
    }
}
```

### Request Size Limiting
```go
func (m *SecurityMiddleware) RequestSizeLimit(maxSize int64) gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxSize)
        c.Next()
    }
}
```

Remember: Always prioritize security, performance, and user experience when developing API endpoints. Follow RESTful principles, maintain consistent response formats, and implement comprehensive error handling and logging.
