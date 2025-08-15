# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is the Vibe Coding Starter API Go - a production-ready web API template built entirely by AI tools using clean architecture principles. It serves as both a learning project for AI-assisted development and a foundation for building production-grade business systems.

### Technology Stack
- **Framework**: Gin Web Framework
- **ORM**: GORM with MySQL, PostgreSQL, SQLite support
- **Cache**: Redis
- **Authentication**: JWT with role-based access control
- **Dependency Injection**: Uber FX
- **Logging**: Zap structured logging
- **Testing**: Testify with comprehensive test coverage
- **Documentation**: Swagger/OpenAPI
- **Deployment**: Docker + Kubernetes ready

### Architecture
- **Clean Architecture**: Handler → Service → Repository → Database
- **Dependency Injection**: Using Uber FX for container management
- **Modular Design**: Each business entity has dedicated handler, service, repository layers
- **Comprehensive Testing**: Unit tests with mocks for all layers

## Common Development Commands

### Development Server
```bash
# Run with different configurations
go run cmd/server/main.go -c configs/config.yaml                    # Local development
go run cmd/server/main.go -c configs/config-docker.yaml            # Docker environment
go run cmd/server/main.go -c configs/config-k3d.yaml              # k3d/Kubernetes
```

### Build and Test
```bash
# Build application
make build                           # Build for Linux
make build-local                     # Build for local platform

# Testing
make test                            # Run all tests
make test-coverage                   # Generate coverage report

# Code quality
make fmt                             # Format code
make lint                            # Run golangci-lint
```

### Database Operations
```bash
# Database migrations
make migrate-up                      # Run migrations
make migrate-down                    # Rollback migrations
make migrate-version                 # Check migration version

# Auto-migration
go run cmd/automigrate/main.go -c configs/config.yaml
```

### Docker and Kubernetes
```bash
# Development environments
make dev-docker                      # Start with Docker Compose
make dev-k3d                         # Start with k3d cluster

# Container operations
make docker-build                    # Build Docker image
make docker-push                     # Push to registry
make k8s-deploy                      # Deploy to Kubernetes
```

### Code Generation
```bash
# Generate complete business module (recommended)
go run cmd/generator/main.go all --name=Product --fields="name:string,price:float64,active:bool"

# Generate from database table
go run cmd/generator/main.go all --name=Product --table=products \
  --host=localhost --port=3306 --user=root --password=secret --database=mydb

# Individual component generation
go run cmd/generator/main.go model --name=Product --fields="name:string,price:float64"
go run cmd/generator/main.go service --model=Product
go run cmd/generator/main.go repository --model=Product
go run cmd/generator/main.go handler --model=Product
go run cmd/generator/main.go test --model=Product
```

## Architecture Patterns

### Clean Architecture Layers

1. **Handler Layer** (`internal/handler/`): HTTP request/response handling, validation, Swagger documentation
2. **Service Layer** (`internal/service/`): Business logic, orchestration, caching
3. **Repository Layer** (`internal/repository/`): Data access abstraction, GORM operations
4. **Model Layer** (`internal/model/`): Domain entities, GORM models, validation

### Dependency Injection with Uber FX
All components are registered in `cmd/server/main.go` using Uber FX. The application uses a layered approach:

```go
app := fx.New(
    // Configuration
    fx.Provide(configProvider),
    
    // Infrastructure
    fx.Provide(logger.New, database.New, cache.New),
    
    // Business layers
    fx.Provide(repository.NewUserRepository, service.NewUserService, handler.NewUserHandler),
    
    // Server setup
    fx.Provide(server.New),
)
```

### Code Generation System
The project includes a sophisticated code generator that can create complete business modules:

- **Database-driven**: Read table structure and generate corresponding Go code
- **Field-driven**: Define fields manually with type system
- **Complete modules**: Generate model, repository, service, handler, tests, migrations
- **Frontend integration**: Generate Ant Design/Vue frontend components

## Configuration Management

### Environment-Specific Configs
- `configs/config.yaml` - Local development
- `configs/config-docker.yaml` - Docker container environment
- `configs/config-k3d.yaml` - Kubernetes development environment
- `configs/config.sqlite.yaml` - SQLite for testing

### Configuration Structure
The configuration system supports:
- YAML files with environment variable overrides
- Database connections (MySQL, PostgreSQL, SQLite)
- Redis caching
- JWT authentication settings
- CORS and security policies
- Logging configuration

## Database Schema

### Core Entities
- **Users**: Authentication and role-based access control
- **Articles**: Content management with CRUD operations
- **Files**: File upload and management
- **Dictionary**: System configuration and data dictionary

### Migration System
- Database-agnostic migrations in `migrations/` directory
- Support for MySQL, PostgreSQL, SQLite
- Auto-migration capabilities
- Version-controlled schema changes

## API Structure

### RESTful Endpoints
```
/api/v1/
├── health/           # Health check endpoints
├── users/           # User management
│   ├── register
│   ├── login
│   └── profile
├── articles/        # Article management
│   ├── (public)     # Public access
│   └── user/        # User-owned content
└── admin/           # Admin-only endpoints
```

### Authentication & Authorization
- JWT-based authentication
- Role-based access control (User, Admin)
- Middleware for protected routes
- Public vs protected endpoint separation

## Testing Strategy

### Test Organization
```
test/
├── handler/         # HTTP handler tests
├── service/         # Business logic tests
├── repository/      # Data access tests
├── mocks/           # Auto-generated mocks
└── testutil/        # Test utilities
```

### Testing Best Practices
- Use testify/suite for organized test suites
- Mock external dependencies
- Test all layers independently
- Integration tests for database operations
- API endpoint testing with HTTP requests

## Development Workflow

### Adding New Features
1. **Use Code Generator**: `go run cmd/generator/main.go all --name=Entity --fields="..."`
2. **Register Dependencies**: Add to `cmd/server/main.go`
3. **Update Routes**: Add routes in `internal/server/routes.go`
4. **Write Tests**: Comprehensive test coverage for all layers
5. **Update Documentation**: Add Swagger annotations

### Code Generation Workflow
1. **Define Model**: Either manually or from database table
2. **Generate Components**: Use the `all` command for complete modules
3. **Customize**: Modify generated code as needed
4. **Test**: Verify with comprehensive test suite
5. **Deploy**: Use Docker/Kubernetes deployment

## Security Considerations

### Authentication & Authorization
- JWT tokens with configurable expiration
- Role-based access control (RBAC)
- Password hashing with bcrypt
- Secure session management

### API Security
- Input validation and sanitization
- Rate limiting middleware
- CORS protection
- Request size limits
- Security headers

### Data Protection
- Never log sensitive information
- Use environment variables for secrets
- Database connection encryption
- File upload validation

## Performance Optimization

### Caching Strategy
- Redis for frequently accessed data
- Cache invalidation on data changes
- Configurable TTL settings

### Database Optimization
- Connection pooling
- Proper indexing
- Query optimization
- Pagination for large datasets

### HTTP Performance
- Request/response compression
- Keep-alive connections
- Proper timeout configuration
- Efficient middleware chain

## File Structure Conventions

### Go Project Structure
```
cmd/                    # Application entry points
├── server/            # Main server application
├── generator/         # Code generator
├── migrate/           # Database migration tool
└── automigrate/       # Auto-migration utility

internal/              # Private application code
├── config/           # Configuration management
├── handler/          # HTTP handlers (controllers)
├── middleware/       # HTTP middleware
├── model/            # Data models and DTOs
├── repository/       # Data access layer
├── service/          # Business logic layer
└── server/           # Server setup and routing

pkg/                  # Public packages
├── cache/            # Cache abstraction
├── database/         # Database abstraction
├── logger/           # Logging abstraction
└── migration/       # Migration utilities

tools/                # Development tools
└── generator/        # Code generation engine

test/                 # Test files
├── handler/          # Handler tests
├── service/          # Service tests
├── repository/       # Repository tests
├── mocks/            # Mock implementations
└── testutil/         # Test utilities
```

### Naming Conventions
- **Files**: snake_case (user_service.go)
- **Packages**: lowercase, short (user, article)
- **Structs**: PascalCase (UserService)
- **Interfaces**: PascalCase with -er suffix (UserRepository)
- **Functions**: PascalCase for exported, camelCase for private
- **Constants**: UPPER_SNAKE_CASE

## Common Issues and Solutions

### Database Connection Issues
- Check configuration in configs/ directory
- Verify database service is running
- Use appropriate config file for environment (dev/docker/k3d)

### Code Generation Problems
- Ensure model names are in PascalCase
- Verify database connection parameters for table-based generation
- Check field type compatibility
- Run `go mod tidy` after generating new code

### Testing Failures
- Mock all external dependencies
- Use testutil for database setup in tests
- Verify test database configuration
- Check for missing interface implementations

### Build Issues
- Run `go mod tidy` to clean dependencies
- Check for missing imports in generated code
- Verify all components are registered in main.go
- Use proper build tags if needed

## Code Quality Standards

### Go Best Practices
- Follow standard Go conventions and idioms
- Use meaningful names with proper casing
- Handle errors explicitly with context
- Keep functions small and focused
- Use interfaces for abstraction and testability

### Error Handling
```go
// Wrap errors with context
if err != nil {
    return fmt.Errorf("failed to create user: %w", err)
}

// Define custom error types
var ErrUserNotFound = errors.New("user not found")
```

### Logging
```go
// Use structured logging
logger.Info("User created", "user_id", user.ID, "username", user.Username)
logger.Error("Failed to create user", "error", err, "username", req.Username)
```

### Documentation
- Swagger annotations for all API endpoints
- Godoc comments for exported functions
- Clear README for complex features
- Inline comments for business logic

## Cursor Rules Integration

The project includes comprehensive Cursor rules in `.cursorrules` that provide detailed guidance for:
- Code standards and best practices
- Architecture patterns and implementation
- Testing strategies and mocking
- Security considerations
- Performance optimization
- AI-assisted development workflows

When working with this project, refer to the Cursor rules for specific implementation patterns and coding standards.