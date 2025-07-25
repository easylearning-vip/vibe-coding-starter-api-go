# API 设计文档

## API 设计原则

### RESTful 设计
- 使用标准 HTTP 方法 (GET, POST, PUT, DELETE)
- 资源导向的 URL 设计
- 统一的响应格式
- 合理的 HTTP 状态码使用

### 版本管理
- URL 路径版本控制: `/api/v1/`
- 向后兼容性保证
- 废弃 API 的优雅处理

### 安全设计
- JWT Token 认证
- 角色权限控制
- 请求限流
- 输入验证

## API 路由结构

```
/api/v1/
├── /users                    # 用户管理
│   ├── POST /register       # 用户注册
│   ├── POST /login          # 用户登录
│   ├── GET /profile         # 获取用户信息
│   ├── PUT /profile         # 更新用户信息
│   └── PUT /password        # 修改密码
├── /articles                # 文章管理
│   ├── GET /                # 获取文章列表
│   ├── GET /search          # 搜索文章
│   ├── GET /:id             # 获取文章详情
│   └── /user/articles       # 用户文章管理
│       ├── GET /            # 获取用户文章列表
│       ├── POST /           # 创建文章
│       ├── PUT /:id         # 更新文章
│       └── DELETE /:id      # 删除文章
├── /files                   # 文件管理
│   ├── POST /upload         # 文件上传
│   ├── GET /:id             # 获取文件信息
│   ├── GET /:id/download    # 下载文件
│   └── DELETE /:id          # 删除文件
└── /health                  # 健康检查
    ├── GET /                # 基础健康检查
    ├── GET /ready           # 就绪检查
    └── GET /live            # 存活检查
```

## 中间件架构

### 全局中间件
```go
// 应用于所有路由
engine.Use(gin.Recovery())                    // 恢复中间件
engine.Use(m.logging.StructuredLogging())     // 结构化日志
engine.Use(m.security.SecurityHeaders())      // 安全头
engine.Use(m.cors.CORS())                     // CORS
engine.Use(m.rateLimit.IPRateLimit(100, 200)) // IP限流
engine.Use(m.security.RequestSizeLimit(10MB)) // 请求大小限制
```

### API 中间件
```go
// 应用于 /api 路由组
api.Use(m.rateLimit.UserRateLimit(60, 120))  // 用户限流
api.Use(m.logging.ErrorLogging())            // 错误日志
```

### 认证中间件
```go
// 应用于需要认证的路由
protected.Use(m.auth.RequireAuth())          // JWT认证
protected.Use(m.auth.RequireRole("admin"))   // 角色验证
```

## 数据模型设计

### 用户模型 (User)
```go
type User struct {
    ID        uint      `json:"id" gorm:"primaryKey"`
    Username  string    `json:"username" gorm:"uniqueIndex;size:50"`
    Email     string    `json:"email" gorm:"uniqueIndex;size:100"`
    Password  string    `json:"-" gorm:"size:255"`
    Nickname  string    `json:"nickname" gorm:"size:50"`
    Avatar    string    `json:"avatar" gorm:"size:255"`
    Role      string    `json:"role" gorm:"size:20;default:user"`
    Status    string    `json:"status" gorm:"size:20;default:active"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
```

### 文章模型 (Article)
```go
type Article struct {
    ID          uint      `json:"id" gorm:"primaryKey"`
    Title       string    `json:"title" gorm:"size:200"`
    Slug        string    `json:"slug" gorm:"uniqueIndex;size:200"`
    Content     string    `json:"content" gorm:"type:text"`
    Summary     string    `json:"summary" gorm:"size:500"`
    AuthorID    uint      `json:"author_id" gorm:"index"`
    Author      User      `json:"author" gorm:"foreignKey:AuthorID"`
    Status      string    `json:"status" gorm:"size:20;default:draft"`
    ViewCount   int       `json:"view_count" gorm:"default:0"`
    PublishedAt *time.Time `json:"published_at"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}
```

### 文件模型 (File)
```go
type File struct {
    ID        uint      `json:"id" gorm:"primaryKey"`
    Filename  string    `json:"filename" gorm:"size:255"`
    OriginalName string `json:"original_name" gorm:"size:255"`
    MimeType  string    `json:"mime_type" gorm:"size:100"`
    Size      int64     `json:"size"`
    Path      string    `json:"path" gorm:"size:500"`
    OwnerID   uint      `json:"owner_id" gorm:"index"`
    Owner     User      `json:"owner" gorm:"foreignKey:OwnerID"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
```

## 请求响应格式

### 统一响应格式
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

### 成功响应示例
```json
{
    "code": 200,
    "message": "success",
    "data": {
        "id": 1,
        "username": "john_doe",
        "email": "john@example.com"
    }
}
```

### 分页响应示例
```json
{
    "code": 200,
    "message": "success",
    "data": [
        {
            "id": 1,
            "title": "Article 1"
        }
    ],
    "meta": {
        "page": 1,
        "page_size": 10,
        "total": 100,
        "total_pages": 10
    }
}
```

### 错误响应示例
```json
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

## API 端点详细设计

### 用户认证 API

#### 用户注册
```http
POST /api/v1/users/register
Content-Type: application/json

{
    "username": "john_doe",
    "email": "john@example.com",
    "password": "password123",
    "nickname": "John"
}
```

#### 用户登录
```http
POST /api/v1/users/login
Content-Type: application/json

{
    "username": "john_doe",
    "password": "password123"
}

Response:
{
    "code": 200,
    "message": "login successful",
    "data": {
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
        "user": {
            "id": 1,
            "username": "john_doe",
            "email": "john@example.com"
        }
    }
}
```

### 文章管理 API

#### 获取文章列表
```http
GET /api/v1/articles?page=1&page_size=10&status=published

Response:
{
    "code": 200,
    "message": "success",
    "data": [
        {
            "id": 1,
            "title": "Sample Article",
            "slug": "sample-article",
            "summary": "This is a sample article",
            "author": {
                "id": 1,
                "username": "john_doe"
            },
            "view_count": 100,
            "published_at": "2024-01-01T00:00:00Z"
        }
    ],
    "meta": {
        "page": 1,
        "page_size": 10,
        "total": 50,
        "total_pages": 5
    }
}
```

#### 创建文章
```http
POST /api/v1/user/articles
Authorization: Bearer <token>
Content-Type: application/json

{
    "title": "New Article",
    "content": "Article content here...",
    "summary": "Article summary",
    "status": "draft"
}
```

### 文件管理 API

#### 文件上传
```http
POST /api/v1/files/upload
Authorization: Bearer <token>
Content-Type: multipart/form-data

file: <binary data>

Response:
{
    "code": 200,
    "message": "file uploaded successfully",
    "data": {
        "id": 1,
        "filename": "uuid-filename.jpg",
        "original_name": "photo.jpg",
        "mime_type": "image/jpeg",
        "size": 1024000,
        "path": "/uploads/2024/01/uuid-filename.jpg"
    }
}
```

## 错误处理

### HTTP 状态码使用
- `200 OK`: 请求成功
- `201 Created`: 资源创建成功
- `400 Bad Request`: 请求参数错误
- `401 Unauthorized`: 未认证
- `403 Forbidden`: 无权限
- `404 Not Found`: 资源不存在
- `409 Conflict`: 资源冲突
- `422 Unprocessable Entity`: 验证失败
- `429 Too Many Requests`: 请求过于频繁
- `500 Internal Server Error`: 服务器内部错误

### 错误码定义
```go
const (
    // 通用错误码
    CodeSuccess           = 200
    CodeBadRequest        = 400
    CodeUnauthorized      = 401
    CodeForbidden         = 403
    CodeNotFound          = 404
    CodeConflict          = 409
    CodeValidationFailed  = 422
    CodeTooManyRequests   = 429
    CodeInternalError     = 500
    
    // 业务错误码
    CodeUserNotFound      = 1001
    CodeUserExists        = 1002
    CodeInvalidPassword   = 1003
    CodeArticleNotFound   = 2001
    CodeFileNotFound      = 3001
    CodeFileUploadFailed  = 3002
)
```

## 认证授权

### JWT Token 结构
```json
{
    "header": {
        "alg": "HS256",
        "typ": "JWT"
    },
    "payload": {
        "user_id": 1,
        "username": "john_doe",
        "role": "user",
        "exp": 1640995200,
        "iat": 1640908800,
        "iss": "vibe-coding-starter"
    }
}
```

### 权限控制
```go
// 角色定义
const (
    RoleUser  = "user"
    RoleAdmin = "admin"
)

// 权限检查
func RequireRole(role string) gin.HandlerFunc {
    return func(c *gin.Context) {
        userRole := c.GetString("user_role")
        if userRole != role {
            c.JSON(403, gin.H{"error": "insufficient permissions"})
            c.Abort()
            return
        }
        c.Next()
    }
}
```

## 请求限流

### 限流策略
```go
// IP 限流: 每分钟 100 次请求
engine.Use(rateLimit.IPRateLimit(100, 200))

// 用户限流: 每分钟 60 次请求
api.Use(rateLimit.UserRateLimit(60, 120))

// 管理员限流: 每分钟 200 次请求
admin.Use(rateLimit.AdminRateLimit())
```

### 限流响应
```json
{
    "code": 429,
    "message": "too many requests",
    "data": {
        "retry_after": 60
    }
}
```

## 输入验证

### 验证规则
```go
type RegisterRequest struct {
    Username string `json:"username" validate:"required,min=3,max=50"`
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=6"`
    Nickname string `json:"nickname" validate:"max=50"`
}
```

### 验证错误响应
```json
{
    "code": 422,
    "message": "validation failed",
    "data": {
        "errors": [
            {
                "field": "username",
                "message": "username is required"
            },
            {
                "field": "email",
                "message": "invalid email format"
            }
        ]
    }
}
```

## API 文档

### Swagger 注释示例
```go
// Register 用户注册
// @Summary 用户注册
// @Description 创建新用户账户
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "注册信息"
// @Success 201 {object} Response{data=User} "注册成功"
// @Failure 400 {object} Response "请求参数错误"
// @Failure 409 {object} Response "用户已存在"
// @Router /api/v1/users/register [post]
func (h *UserHandler) Register(c *gin.Context) {
    // 实现逻辑
}
```

### 文档访问
- Swagger UI: `http://localhost:8080/swagger/index.html`
- API 文档: 自动生成，实时更新

## 性能优化

### 缓存策略
```go
// 文章列表缓存
key := fmt.Sprintf("articles:list:%d:%d", page, pageSize)
if cached, err := cache.Get(ctx, key); err == nil {
    return cached, nil
}

// 查询数据库
articles, err := repo.List(ctx, opts)
if err != nil {
    return nil, err
}

// 缓存结果
cache.Set(ctx, key, articles, 5*time.Minute)
```

### 分页优化
```go
// 使用游标分页替代偏移分页
type CursorPagination struct {
    Cursor   string `json:"cursor,omitempty"`
    PageSize int    `json:"page_size"`
    HasNext  bool   `json:"has_next"`
}
```

### 查询优化
```go
// 预加载关联数据
db.Preload("Author").Find(&articles)

// 选择特定字段
db.Select("id", "title", "summary").Find(&articles)

// 使用索引
db.Where("status = ? AND author_id = ?", "published", userID).Find(&articles)
```

这个 API 设计确保了系统的可扩展性、安全性和性能，为前端开发和第三方集成提供了清晰的接口规范。
