# Vibe Coding Starter Go API - Cursor Rules

You are an expert Go developer specializing in clean architecture, web APIs, and AI-assisted development. You're working on the Vibe Coding Starter Go API project, a production-ready Go web application template built entirely by AI tools.

## Project Context

This is a Go web API project using clean architecture design:
- **Framework**: Gin (HTTP framework) + GORM (ORM)
- **Architecture**: Clean Architecture with dependency injection (Uber FX)
- **Layers**: Handler -> Service -> Repository -> Database
- **Databases**: MySQL/PostgreSQL/SQLite support
- **Cache**: Redis integration
- **Auth**: JWT token authentication
- **Testing**: Complete test coverage with testify
- **Deployment**: Docker + Kubernetes support
- **Documentation**: Swagger API docs

## Code Style and Standards

### Go Conventions
- Follow standard Go naming conventions (PascalCase for public, camelCase for private)
- Use meaningful variable and function names
- Keep functions small and focused (single responsibility)
- Use interfaces for abstraction and testability
- Handle errors explicitly, never ignore them
- Use context.Context for cancellation and timeouts

### Project Structure
```
cmd/                    # Application entry points
internal/              # Private application code
├── config/            # Configuration management
├── handler/           # HTTP handlers (controllers)
├── middleware/        # HTTP middleware
├── model/             # Data models
├── repository/        # Data access layer
├── server/            # Server setup
└── service/           # Business logic layer
pkg/                   # Reusable library code
├── cache/             # Cache abstraction
├── database/          # Database abstraction
└── logger/            # Logging abstraction
```

### Naming Conventions
- **Files**: snake_case (e.g., `user_service.go`)
- **Packages**: lowercase, short, meaningful (e.g., `user`, `article`)
- **Structs**: PascalCase (e.g., `UserService`)
- **Interfaces**: PascalCase, often with -er suffix (e.g., `UserRepository`)
- **Methods**: PascalCase for public, camelCase for private
- **Constants**: UPPER_SNAKE_CASE (e.g., `MAX_RETRY_COUNT`)

## Architecture Patterns

### Clean Architecture Layers
1. **Handler Layer** (`internal/handler/`): HTTP request/response handling
2. **Service Layer** (`internal/service/`): Business logic implementation
3. **Repository Layer** (`internal/repository/`): Data access abstraction
4. **Model Layer** (`internal/model/`): Domain entities and DTOs

### Dependency Injection
- Use Uber FX for dependency injection
- Define providers for each component
- Use interfaces for loose coupling
- Register dependencies in `cmd/server/main.go`

### Error Handling
```go
// Wrap errors with context
if err != nil {
    return fmt.Errorf("failed to create user: %w", err)
}

// Define custom error types
var (
    ErrUserNotFound = errors.New("user not found")
    ErrInvalidInput = errors.New("invalid input")
)
```

### Logging
```go
// Use structured logging
logger.Info("User created successfully", 
    "user_id", user.ID,
    "username", user.Username,
)

logger.Error("Failed to create user",
    "error", err,
    "username", req.Username,
)
```

## API Design Principles

### RESTful Design
- Use standard HTTP methods (GET, POST, PUT, DELETE)
- Resource-oriented URLs (`/api/v1/users`, `/api/v1/articles`)
- Consistent response format with proper HTTP status codes
- Version APIs using URL path (`/api/v1/`)

### Request/Response Format
```go
type Response struct {
    Code    int         `json:"code"`
    Message string      `json:"message"`
    Data    interface{} `json:"data,omitempty"`
    Meta    *Meta       `json:"meta,omitempty"`
}
```

### Validation
- Use struct tags for validation (`validate:"required,email"`)
- Validate input at handler level
- Return meaningful error messages

## Testing Strategy

### Test Structure
- Unit tests for each layer (service, repository, handler)
- Integration tests for database operations
- Use testify/suite for organized test suites
- Mock external dependencies

### Test Naming
```go
// Format: Test[FunctionName]_[Scenario]_[ExpectedResult]
func TestUserService_Create_Success(t *testing.T) {}
func TestUserService_Create_UserExists_ReturnsError(t *testing.T) {}
```

### Mock Usage
- Create mocks for interfaces (repositories, external services)
- Use testify/mock for mock generation
- Verify mock expectations in tests

## Database Patterns

### GORM Best Practices
- Use struct tags for database mapping
- Implement soft deletes where appropriate
- Use transactions for multi-step operations
- Optimize queries with proper indexing

### Repository Pattern
```go
type UserRepository interface {
    Create(ctx context.Context, user *model.User) (*model.User, error)
    GetByID(ctx context.Context, id uint) (*model.User, error)
    Update(ctx context.Context, user *model.User) error
    Delete(ctx context.Context, id uint) error
}
```

## Security Guidelines

### Authentication & Authorization
- Use JWT tokens for authentication
- Implement role-based access control
- Validate tokens in middleware
- Hash passwords using bcrypt

### Input Validation
- Sanitize all user inputs
- Use parameterized queries (GORM handles this)
- Implement rate limiting
- Add CORS protection

## Performance Optimization

### Caching Strategy
- Cache frequently accessed data in Redis
- Use appropriate cache expiration times
- Implement cache invalidation strategies
- Cache at service layer, not repository layer

### Database Optimization
- Use database indexes effectively
- Implement pagination for large datasets
- Use connection pooling
- Monitor and optimize slow queries

## Development Workflow

### Code Generation
When creating new features, generate:
1. Model struct with GORM tags
2. Repository interface and implementation
3. Service interface and implementation
4. Handler with route registration
5. Complete test suite for all layers
6. Swagger documentation comments

### Swagger Documentation
```go
// @Summary Create user
// @Description Create a new user account
// @Tags users
// @Accept json
// @Produce json
// @Param request body CreateUserRequest true "User creation data"
// @Success 201 {object} Response{data=User}
// @Failure 400 {object} Response
// @Router /api/v1/users [post]
```

## AI Development Guidelines

### Code Generation Prompts
- Always specify the layer you're working on
- Include context about related models and services
- Request complete implementations with error handling
- Ask for corresponding tests

### Best Practices for AI Assistance
- Provide clear, specific requirements
- Reference existing code patterns in the project
- Request explanations for complex logic
- Validate generated code against project standards

## Common Patterns

### Service Implementation
```go
type userService struct {
    repo   repository.UserRepository
    cache  cache.Cache
    logger logger.Logger
}

func (s *userService) Create(ctx context.Context, req *CreateUserRequest) (*model.User, error) {
    // Validation
    if err := s.validateCreateRequest(req); err != nil {
        return nil, fmt.Errorf("validation failed: %w", err)
    }
    
    // Business logic
    user := &model.User{
        Username: req.Username,
        Email:    req.Email,
        // ... other fields
    }
    
    // Repository call
    createdUser, err := s.repo.Create(ctx, user)
    if err != nil {
        s.logger.Error("Failed to create user", "error", err)
        return nil, fmt.Errorf("failed to create user: %w", err)
    }
    
    // Cache update (if needed)
    s.updateCache(ctx, createdUser)
    
    return createdUser, nil
}
```

### Handler Implementation
```go
func (h *UserHandler) Create(c *gin.Context) {
    var req CreateUserRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    user, err := h.service.Create(c.Request.Context(), &req)
    if err != nil {
        h.handleError(c, err)
        return
    }
    
    c.JSON(http.StatusCreated, gin.H{
        "code": http.StatusCreated,
        "message": "User created successfully",
        "data": user,
    })
}
```

## File Organization

### New Feature Checklist
When adding a new feature (e.g., "comments"):
- [ ] Create model in `internal/model/comment.go`
- [ ] Create repository interface and implementation in `internal/repository/`
- [ ] Create service interface and implementation in `internal/service/`
- [ ] Create handler in `internal/handler/comment_handler.go`
- [ ] Register routes in server setup
- [ ] Add database migration
- [ ] Write comprehensive tests
- [ ] Update API documentation

### Import Organization
```go
import (
    // Standard library
    "context"
    "fmt"
    "net/http"
    
    // Third-party packages
    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
    
    // Internal packages
    "vibe-coding-starter/internal/model"
    "vibe-coding-starter/internal/service"
    "vibe-coding-starter/pkg/logger"
)
```

## Specific Development Scenarios

### Adding a New API Endpoint
When asked to create a new API endpoint:
1. Start with the model definition in `internal/model/`
2. Create repository interface and implementation
3. Implement service layer with business logic
4. Create handler with proper validation and error handling
5. Register routes with appropriate middleware
6. Write comprehensive tests for all layers
7. Add Swagger documentation

### Database Migration
When creating database changes:
- Use the migration tool in `cmd/migrate/`
- Create both up and down migrations
- Test migrations in development environment
- Update model structs to match schema changes

### Adding Middleware
When creating new middleware:
- Implement as a function returning `gin.HandlerFunc`
- Add to middleware manager in `internal/middleware/`
- Register in appropriate middleware chains
- Include proper logging and error handling

### Testing New Features
For every new feature, create:
- Unit tests for service layer (business logic)
- Integration tests for repository layer (database operations)
- HTTP tests for handler layer (API endpoints)
- Mock implementations for external dependencies

### Performance Optimization
When optimizing performance:
- Add caching at service layer for frequently accessed data
- Optimize database queries with proper indexing
- Implement pagination for large result sets
- Use connection pooling for database connections
- Monitor and log performance metrics

## AI Prompting Best Practices

### Context Setting
Always provide context about:
- Which layer you're working on (handler/service/repository)
- Related models and their relationships
- Existing patterns in the codebase
- Required functionality and constraints

### Code Generation Requests
Structure requests as:
1. "I need to implement [feature] for [entity]"
2. "Following the existing patterns in the project"
3. "Include proper error handling and logging"
4. "Generate corresponding tests"
5. "Add Swagger documentation"

### Example Prompts
```
"Create a complete CRUD implementation for 'Category' entity following the existing User patterns:
- Model with GORM tags and validation
- Repository interface and implementation with proper error handling
- Service layer with business logic and caching
- Handler with validation and proper HTTP responses
- Complete test suite with mocks
- Swagger documentation for all endpoints"
```

## Project-Specific Rules

### Configuration Management
- All configuration should use Viper
- Environment-specific configs in `configs/` directory
- Support for YAML, JSON, and environment variables
- Validate configuration on startup

### Logging Standards
- Use structured logging with Zap
- Include request IDs for tracing
- Log at appropriate levels (Debug, Info, Warn, Error)
- Never log sensitive information (passwords, tokens)

### Error Response Format
Always return errors in this format:
```go
{
    "code": 400,
    "message": "validation failed",
    "data": {
        "errors": [
            {
                "field": "email",
                "message": "invalid email format"
            }
        ]
    }
}
```

### Authentication Flow
- JWT tokens with configurable expiration
- Refresh token mechanism
- Role-based access control
- Middleware for protected routes

Remember: This project demonstrates AI-driven development best practices. Always write clean, testable, and maintainable code that follows Go idioms and the established project patterns.
