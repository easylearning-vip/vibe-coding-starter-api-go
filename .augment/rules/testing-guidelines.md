# Testing Guidelines for Vibe Coding Starter Go API

When working on tests for this Go API project, follow these comprehensive testing guidelines that align with clean architecture principles and Go best practices.

## Testing Philosophy

### Test Pyramid Strategy
- **Unit Tests (80%)**: Fast, isolated tests for individual functions and methods
- **Integration Tests (15%)**: Test interactions between components, especially database operations
- **End-to-End Tests (5%)**: Test complete user workflows and business processes

### Testing Principles
- Tests should be fast, reliable, and independent
- Each test should focus on a single behavior or scenario
- Tests should be easy to read and understand
- Mock external dependencies to ensure test isolation
- Use descriptive test names that explain the scenario and expected outcome

## Test Organization Structure

### Directory Layout
```
test/
├── config/                 # Test configuration and setup
├── handler/                # HTTP handler tests
├── service/                # Business logic tests
├── repository/             # Data access layer tests
├── mocks/                  # Mock implementations
├── testutil/               # Test utilities and helpers
└── scripts/                # Test execution scripts
```

### Test File Naming
- Test files should end with `_test.go`
- Place tests in the same package as the code being tested
- Use descriptive names: `user_service_test.go`, `article_handler_test.go`

## Test Naming Conventions

### Test Function Naming
Use the format: `Test[FunctionName]_[Scenario]_[ExpectedResult]`

```go
func TestUserService_Create_Success(t *testing.T) {}
func TestUserService_Create_UserExists_ReturnsError(t *testing.T) {}
func TestUserService_Login_InvalidPassword_ReturnsUnauthorized(t *testing.T) {}
func TestUserHandler_Register_InvalidEmail_ReturnsBadRequest(t *testing.T) {}
```

### Test Suite Naming
```go
type UserServiceTestSuite struct {
    suite.Suite
    service     UserService
    mockRepo    *mocks.MockUserRepository
    mockCache   *mocks.MockCache
    mockLogger  *mocks.MockLogger
}
```

## Unit Testing Patterns

### Service Layer Testing
```go
func (s *UserServiceTestSuite) SetupTest() {
    s.mockRepo = &mocks.MockUserRepository{}
    s.mockCache = &mocks.MockCache{}
    s.mockLogger = &mocks.MockLogger{}
    
    s.service = NewUserService(s.mockRepo, s.mockCache, s.mockLogger)
}

func (s *UserServiceTestSuite) TestRegister_Success() {
    // Given - Arrange
    req := &RegisterRequest{
        Username: "testuser",
        Email:    "test@example.com",
        Password: "password123",
        Nickname: "Test User",
    }
    
    expectedUser := &model.User{
        ID:       1,
        Username: "testuser",
        Email:    "test@example.com",
        Nickname: "Test User",
        Role:     "user",
        Status:   "active",
    }
    
    // Setup mocks
    s.mockRepo.On("GetByUsername", mock.Anything, "testuser").Return(nil, repository.ErrUserNotFound)
    s.mockRepo.On("GetByEmail", mock.Anything, "test@example.com").Return(nil, repository.ErrUserNotFound)
    s.mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*model.User")).Return(expectedUser, nil)
    
    // When - Act
    user, err := s.service.Register(context.Background(), req)
    
    // Then - Assert
    s.NoError(err)
    s.NotNil(user)
    s.Equal("testuser", user.Username)
    s.Equal("test@example.com", user.Email)
    s.Equal("Test User", user.Nickname)
    s.Equal("user", user.Role)
    s.Equal("active", user.Status)
    
    // Verify all mock expectations were met
    s.mockRepo.AssertExpectations(s.T())
}

func (s *UserServiceTestSuite) TestRegister_UserExists_ReturnsError() {
    // Given
    req := &RegisterRequest{
        Username: "existinguser",
        Email:    "existing@example.com",
        Password: "password123",
    }
    
    existingUser := &model.User{
        ID:       1,
        Username: "existinguser",
        Email:    "existing@example.com",
    }
    
    s.mockRepo.On("GetByUsername", mock.Anything, "existinguser").Return(existingUser, nil)
    
    // When
    user, err := s.service.Register(context.Background(), req)
    
    // Then
    s.Error(err)
    s.Nil(user)
    s.Equal(ErrUserExists, err)
    s.mockRepo.AssertExpectations(s.T())
}
```

### Handler Layer Testing
```go
func (s *UserHandlerTestSuite) SetupTest() {
    gin.SetMode(gin.TestMode)
    s.router = gin.New()
    
    s.mockService = &mocks.MockUserService{}
    s.handler = handler.NewUserHandler(s.mockService)
    
    // Register routes
    api := s.router.Group("/api/v1")
    s.handler.RegisterRoutes(api)
}

func (s *UserHandlerTestSuite) TestRegister_Success() {
    // Given
    reqBody := map[string]interface{}{
        "username": "testuser",
        "email":    "test@example.com",
        "password": "password123",
        "nickname": "Test User",
    }
    
    expectedUser := &model.User{
        ID:       1,
        Username: "testuser",
        Email:    "test@example.com",
        Nickname: "Test User",
    }
    
    s.mockService.On("Register", mock.Anything, mock.AnythingOfType("*service.RegisterRequest")).Return(expectedUser, nil)
    
    jsonBody, _ := json.Marshal(reqBody)
    req := httptest.NewRequest(http.MethodPost, "/api/v1/users/register", bytes.NewBuffer(jsonBody))
    req.Header.Set("Content-Type", "application/json")
    
    // When
    w := httptest.NewRecorder()
    s.router.ServeHTTP(w, req)
    
    // Then
    s.Equal(http.StatusCreated, w.Code)
    
    var response map[string]interface{}
    err := json.Unmarshal(w.Body.Bytes(), &response)
    s.NoError(err)
    s.Equal(float64(201), response["code"])
    s.Equal("user registered successfully", response["message"])
    
    data := response["data"].(map[string]interface{})
    s.Equal("testuser", data["username"])
    s.Equal("test@example.com", data["email"])
    
    s.mockService.AssertExpectations(s.T())
}
```

## Integration Testing Patterns

### Repository Integration Tests
```go
type UserRepositoryIntegrationTestSuite struct {
    suite.Suite
    db   database.Database
    repo repository.UserRepository
}

func (s *UserRepositoryIntegrationTestSuite) SetupSuite() {
    // Setup test database
    db, err := testutil.SetupTestDatabase()
    s.Require().NoError(err)
    s.db = db
    
    // Run migrations
    err = s.db.AutoMigrate(&model.User{})
    s.Require().NoError(err)
    
    // Create repository
    logger := testutil.NewTestLogger()
    s.repo = repository.NewUserRepository(s.db, logger)
}

func (s *UserRepositoryIntegrationTestSuite) TearDownSuite() {
    s.db.Close()
}

func (s *UserRepositoryIntegrationTestSuite) SetupTest() {
    // Clean test data before each test
    s.db.GetDB().Exec("DELETE FROM users")
}

func (s *UserRepositoryIntegrationTestSuite) TestCreate_Success() {
    // Given
    user := &model.User{
        Username: "testuser",
        Email:    "test@example.com",
        Password: "hashedpassword",
        Nickname: "Test User",
        Role:     "user",
        Status:   "active",
    }
    
    // When
    createdUser, err := s.repo.Create(context.Background(), user)
    
    // Then
    s.NoError(err)
    s.NotNil(createdUser)
    s.NotZero(createdUser.ID)
    s.Equal("testuser", createdUser.Username)
    s.Equal("test@example.com", createdUser.Email)
    s.Equal("Test User", createdUser.Nickname)
    s.NotZero(createdUser.CreatedAt)
    s.NotZero(createdUser.UpdatedAt)
}

func (s *UserRepositoryIntegrationTestSuite) TestGetByID_NotFound_ReturnsError() {
    // When
    user, err := s.repo.GetByID(context.Background(), 999)
    
    // Then
    s.Error(err)
    s.Nil(user)
    s.Equal(repository.ErrUserNotFound, err)
}
```

## Mock Implementation Guidelines

### Repository Mock Example
```go
type MockUserRepository struct {
    mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *model.User) (*model.User, error) {
    args := m.Called(ctx, user)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) GetByID(ctx context.Context, id uint) (*model.User, error) {
    args := m.Called(ctx, id)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
    args := m.Called(ctx, email)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*model.User), args.Error(1)
}
```

### Service Mock Example
```go
type MockUserService struct {
    mock.Mock
}

func (m *MockUserService) Register(ctx context.Context, req *RegisterRequest) (*model.User, error) {
    args := m.Called(ctx, req)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserService) Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error) {
    args := m.Called(ctx, req)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*LoginResponse), args.Error(1)
}
```

## Test Utilities and Helpers

### Test Database Setup
```go
func SetupTestDatabase() (database.Database, error) {
    cfg := &config.Config{
        Database: config.DatabaseConfig{
            Driver:   "sqlite",
            Database: ":memory:",
        },
    }
    
    logger := NewTestLogger()
    return database.New(cfg, logger)
}
```

### Test Cache Implementation
```go
type TestCache struct {
    data map[string]string
    mu   sync.RWMutex
}

func NewTestCache() cache.Cache {
    return &TestCache{
        data: make(map[string]string),
    }
}

func (c *TestCache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.data[key] = value.(string)
    return nil
}

func (c *TestCache) Get(ctx context.Context, key string) (string, error) {
    c.mu.RLock()
    defer c.mu.RUnlock()
    if val, exists := c.data[key]; exists {
        return val, nil
    }
    return "", cache.ErrKeyNotFound
}
```

### Test Data Factories
```go
func CreateTestUser() *model.User {
    return &model.User{
        Username: "testuser",
        Email:    "test@example.com",
        Password: "hashedpassword",
        Nickname: "Test User",
        Role:     "user",
        Status:   "active",
    }
}

func CreateTestUserWithID(id uint) *model.User {
    user := CreateTestUser()
    user.ID = id
    user.CreatedAt = time.Now()
    user.UpdatedAt = time.Now()
    return user
}

func CreateTestArticle(authorID uint) *model.Article {
    return &model.Article{
        Title:     "Test Article",
        Slug:      "test-article",
        Content:   "This is a test article content.",
        Summary:   "Test article summary",
        AuthorID:  authorID,
        Status:    "published",
        ViewCount: 0,
    }
}
```

## Table-Driven Tests

### Example Implementation
```go
func TestValidateEmail(t *testing.T) {
    tests := []struct {
        name     string
        email    string
        expected bool
    }{
        {"valid email", "test@example.com", true},
        {"valid email with subdomain", "user@mail.example.com", true},
        {"invalid email without @", "testexample.com", false},
        {"invalid email without domain", "test@", false},
        {"invalid email without local part", "@example.com", false},
        {"empty email", "", false},
        {"email with spaces", "test @example.com", false},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := validateEmail(tt.email)
            assert.Equal(t, tt.expected, result, "Email validation failed for: %s", tt.email)
        })
    }
}
```

## Test Coverage and Quality

### Coverage Requirements
- Maintain minimum 80% test coverage for all packages
- Critical business logic should have 90%+ coverage
- All error paths and edge cases must be tested
- Include performance benchmarks for critical functions

### Running Tests
```bash
# Run all tests
go test -v ./...

# Run tests with coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html

# Run specific test suite
go test -v ./internal/service/...

# Run tests with race detection
go test -race ./...

# Run benchmarks
go test -bench=. ./...
```

### Test Quality Checklist
- [ ] Tests are independent and can run in any order
- [ ] Each test focuses on a single behavior
- [ ] Test names clearly describe the scenario and expected outcome
- [ ] Mocks are properly configured and expectations are verified
- [ ] Test data is realistic and covers edge cases
- [ ] Error conditions are thoroughly tested
- [ ] Tests run quickly (unit tests < 100ms each)
- [ ] Integration tests clean up after themselves

## Best Practices

### Test Organization
- Group related tests in test suites using testify/suite
- Use SetupTest/TearDownTest for test isolation
- Keep tests independent and idempotent
- Use meaningful test data that reflects real-world scenarios

### Mock Best Practices
- Only mock external dependencies and interfaces
- Verify mock expectations to ensure correct interactions
- Keep mocks simple and focused on the interface contract
- Use dependency injection to make code more testable

### Test Maintenance
- Keep tests up to date with code changes
- Refactor tests when refactoring production code
- Remove obsolete tests that no longer provide value
- Document complex test scenarios and edge cases

### Performance Testing
- Use benchmarks for performance-critical code paths
- Test with realistic data volumes and load patterns
- Monitor test execution time and optimize slow tests
- Include stress tests for high-load scenarios

Remember: Tests are living documentation of your code's behavior. Write tests that are clear, maintainable, and provide confidence in your system's correctness and reliability.
