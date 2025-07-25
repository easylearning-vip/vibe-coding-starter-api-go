# Vibe Coding Starter Go API - Main Development Rules

You are an expert Go developer working on the Vibe Coding Starter Go API project. This is a production-ready Go web application template built entirely by AI tools, demonstrating clean architecture and best practices.

## Project Overview

### Technology Stack
- **Language**: Go 1.23+
- **Web Framework**: Gin (HTTP router and middleware)
- **ORM**: GORM (database operations and migrations)
- **Cache**: Redis (session storage and caching)
- **Authentication**: JWT tokens with role-based access control
- **Dependency Injection**: Uber FX (IoC container)
- **Testing**: testify/suite (unit and integration tests)
- **Documentation**: Swagger/OpenAPI (auto-generated API docs)
- **Deployment**: Docker + Kubernetes (containerized deployment)

### Architecture Pattern
This project follows **Clean Architecture** principles with clear separation of concerns:

```
Handler Layer    → HTTP request/response handling, validation
Service Layer    → Business logic, orchestration, caching
Repository Layer → Data access abstraction, database operations
Model Layer      → Domain entities, DTOs, database models
```

### Project Structure
```
cmd/server/           # Application entry point with dependency injection
internal/
├── config/          # Configuration management (Viper)
├── handler/         # HTTP handlers (Gin controllers)
├── middleware/      # HTTP middleware (auth, CORS, logging, rate limiting)
├── model/           # Data models with GORM tags
├── repository/      # Data access layer interfaces and implementations
├── server/          # Server setup and route registration
└── service/         # Business logic layer interfaces and implementations
pkg/
├── cache/           # Redis cache abstraction
├── database/        # Database connection and health checks
└── logger/          # Structured logging with Zap
```

## Development Standards

### Go Best Practices
- Follow standard Go naming conventions (PascalCase for public, camelCase for private)
- Use meaningful, descriptive names for variables, functions, and types
- Keep functions small and focused on a single responsibility
- Handle errors explicitly with proper context wrapping
- Use context.Context for cancellation, timeouts, and request tracing
- Implement interfaces for abstraction and testability
- Use struct tags for JSON serialization and validation

### Error Handling Pattern
```go
// Always wrap errors with context
if err != nil {
    return fmt.Errorf("failed to create user: %w", err)
}

// Define custom error types for business logic
var (
    ErrUserNotFound     = errors.New("user not found")
    ErrInvalidInput     = errors.New("invalid input")
    ErrUnauthorized     = errors.New("unauthorized access")
)
```

### Logging Standards
```go
// Use structured logging with key-value pairs
logger.Info("User created successfully", 
    "user_id", user.ID,
    "username", user.Username,
    "request_id", requestID,
)

logger.Error("Failed to create user",
    "error", err,
    "username", req.Username,
    "request_id", requestID,
)
```

## Clean Architecture Implementation

### Dependency Injection with Uber FX
- All components are registered as providers in `cmd/server/main.go`
- Use interfaces for loose coupling between layers
- Dependencies flow inward (Handler → Service → Repository)
- No circular dependencies allowed

### Layer Responsibilities

#### Handler Layer (`internal/handler/`)
- HTTP request binding and validation
- Response formatting and status codes
- Error handling and logging
- Route registration and middleware application
- Swagger documentation annotations

#### Service Layer (`internal/service/`)
- Core business logic implementation
- Data validation and transformation
- Caching strategies and cache invalidation
- Cross-cutting concerns (logging, metrics)
- Transaction coordination

#### Repository Layer (`internal/repository/`)
- Database operations and query building
- Data mapping between models and database
- Connection management and health checks
- Migration support and schema management

#### Model Layer (`internal/model/`)
- Domain entities with GORM tags
- JSON serialization tags
- Validation rules and constraints
- Database relationships and associations

## API Design Standards

### RESTful Principles
- Use standard HTTP methods (GET, POST, PUT, DELETE)
- Resource-oriented URLs (`/api/v1/users`, `/api/v1/articles`)
- Consistent response format with proper HTTP status codes
- API versioning through URL path (`/api/v1/`)

### Response Format
```go
type Response struct {
    Code    int         `json:"code"`
    Message string      `json:"message"`
    Data    interface{} `json:"data,omitempty"`
    Meta    *Meta       `json:"meta,omitempty"`
}

type Meta struct {
    Page       int   `json:"page,omitempty"`
    PageSize   int   `json:"page_size,omitempty"`
    Total      int64 `json:"total,omitempty"`
    TotalPages int   `json:"total_pages,omitempty"`
}
```

### Input Validation
- Use struct tags for validation (`validate:"required,email,min=3,max=50"`)
- Validate at the handler level before calling services
- Return meaningful error messages with field-specific details
- Sanitize user inputs to prevent injection attacks

### HTTP Status Codes
- 200 OK: Successful GET, PUT, PATCH operations
- 201 Created: Successful POST operations
- 204 No Content: Successful DELETE operations
- 400 Bad Request: Invalid request data or parameters
- 401 Unauthorized: Authentication required
- 403 Forbidden: Insufficient permissions
- 404 Not Found: Resource not found
- 409 Conflict: Resource conflict (e.g., duplicate email)
- 422 Unprocessable Entity: Validation errors
- 429 Too Many Requests: Rate limit exceeded
- 500 Internal Server Error: Unexpected server errors

## Database Patterns

### GORM Best Practices
- Use struct tags for database mapping and constraints
- Implement soft deletes where appropriate (`gorm:"softDelete"`)
- Use database transactions for multi-step operations
- Optimize queries with proper indexing and preloading
- Use connection pooling for performance

### Repository Pattern Implementation
```go
type UserRepository interface {
    Create(ctx context.Context, user *model.User) (*model.User, error)
    GetByID(ctx context.Context, id uint) (*model.User, error)
    GetByEmail(ctx context.Context, email string) (*model.User, error)
    Update(ctx context.Context, user *model.User) error
    Delete(ctx context.Context, id uint) error
    List(ctx context.Context, opts ListOptions) ([]*model.User, int64, error)
}
```

## Security Requirements

### Authentication & Authorization
- JWT tokens with configurable expiration times
- Role-based access control (RBAC) implementation
- Secure token storage and transmission
- Password hashing using bcrypt with appropriate cost
- Session management and token refresh mechanisms

### Input Security
- Validate and sanitize all user inputs
- Use parameterized queries (GORM handles this automatically)
- Implement rate limiting to prevent abuse
- Add CORS protection with proper origin validation
- Escape HTML content to prevent XSS attacks

### API Security
- HTTPS enforcement in production
- Security headers (HSTS, CSP, X-Frame-Options)
- Request size limits to prevent DoS attacks
- API key validation for external integrations
- Audit logging for sensitive operations

## Performance Optimization

### Caching Strategy
- Cache frequently accessed data in Redis
- Use appropriate cache expiration times based on data volatility
- Implement cache invalidation strategies for data consistency
- Cache at the service layer, not repository layer
- Use cache-aside pattern for read-heavy operations

### Database Optimization
- Use database indexes effectively for query performance
- Implement pagination for large result sets
- Use database connection pooling with proper limits
- Monitor and optimize slow queries
- Use read replicas for read-heavy workloads

### Memory Management
- Use object pooling for frequently allocated objects
- Implement proper resource cleanup with defer statements
- Monitor memory usage and garbage collection metrics
- Use streaming for large file operations
- Implement circuit breakers for external service calls

## Testing Strategy

### Test Organization
- Unit tests for each layer (service, repository, handler)
- Integration tests for database operations and external services
- End-to-end tests for critical business workflows
- Use testify/suite for organized test suites
- Mock external dependencies using interfaces

### Test Naming Convention
```go
// Format: Test[FunctionName]_[Scenario]_[ExpectedResult]
func TestUserService_Create_Success(t *testing.T) {}
func TestUserService_Create_UserExists_ReturnsError(t *testing.T) {}
func TestUserService_Login_InvalidPassword_ReturnsUnauthorized(t *testing.T) {}
```

### Coverage Requirements
- Maintain minimum 80% test coverage for all packages
- 90%+ coverage for critical business logic
- Test all error paths and edge cases
- Include performance benchmarks for critical functions

## Development Workflow

### Feature Development Process
When implementing new features:
1. Define the model with appropriate GORM tags and validation
2. Create repository interface and implementation with error handling
3. Implement service layer with business logic and caching
4. Create handler with proper validation and HTTP responses
5. Register routes with appropriate middleware
6. Write comprehensive tests for all layers
7. Add Swagger documentation for API endpoints
8. Update database migrations if schema changes are needed

### Code Generation Template
For any new entity (e.g., "Product"):
- Model struct with GORM tags, JSON tags, and validation rules
- Repository interface and implementation with CRUD operations
- Service interface and implementation with business logic
- Handler with proper HTTP methods and response formatting
- Complete test suite with unit tests and mocks
- Swagger documentation with request/response examples
- Database migration files for schema changes

### Swagger Documentation Format
```go
// @Summary Create product
// @Description Create a new product in the system
// @Tags products
// @Accept json
// @Produce json
// @Param request body CreateProductRequest true "Product creation data"
// @Success 201 {object} Response{data=Product} "Product created successfully"
// @Failure 400 {object} Response "Invalid request data"
// @Failure 409 {object} Response "Product already exists"
// @Router /api/v1/products [post]
```

## Configuration Management

### Environment Configuration
- Use Viper for configuration management with multiple sources
- Support YAML, JSON, and environment variables
- Environment-specific configuration files in `configs/` directory
- Validate configuration on application startup
- Use configuration structs with proper tags

### Logging Configuration
- Use Zap for structured, high-performance logging
- Include request IDs for distributed tracing
- Log at appropriate levels (Debug, Info, Warn, Error, Fatal)
- Never log sensitive information (passwords, tokens, PII)
- Use log rotation and retention policies

## Deployment and Operations

### Docker Best Practices
- Use multi-stage builds for smaller image sizes
- Run as non-root user for security
- Use specific base image versions, not latest
- Implement health checks for container orchestration
- Use .dockerignore to exclude unnecessary files

### Kubernetes Deployment
- Use resource limits and requests for proper scheduling
- Implement readiness and liveness probes
- Use ConfigMaps and Secrets for configuration management
- Implement horizontal pod autoscaling based on metrics
- Use persistent volumes for stateful data

### Monitoring and Observability
- Implement structured logging with correlation IDs
- Use metrics collection for performance monitoring
- Implement distributed tracing for request flows
- Set up alerting for critical system metrics
- Use health check endpoints for service monitoring

Remember: This project demonstrates AI-driven development excellence. Always write clean, testable, and maintainable code that follows Go idioms and established architectural patterns. Focus on code quality, security, and performance in every implementation.
