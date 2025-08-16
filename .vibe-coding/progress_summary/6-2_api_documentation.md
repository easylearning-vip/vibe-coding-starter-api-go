# Vibe Coding Starter API - API 文档

## 概述
Vibe Coding Starter API 是一个基于 Go 语言开发的 RESTful API 服务，提供完整的产品管理功能。本文档详细描述了所有可用的 API 端点。

## 基础信息

- **基础 URL**: `http://localhost:8080/api/v1`
- **API 版本**: v1
- **认证方式**: JWT Bearer Token
- **数据格式**: JSON

## 认证

### JWT 认证
所有需要认证的 API 端点都需要在请求头中包含 JWT token：

```http
Authorization: Bearer <your-jwt-token>
```

### 获取 Token
```http
POST /api/v1/users/login
Content-Type: application/json

{
  "username": "admin",
  "password": "vibecoding"
}
```

## API 端点

### 1. 健康检查

#### GET /api/v1/health
检查服务健康状态

**请求**: 无

**响应**:
```json
{
  "status": "healthy",
  "timestamp": "2025-08-15T09:35:00Z",
  "version": "1.0.0",
  "services": {
    "database": {
      "status": "healthy",
      "message": "Database connection successful"
    },
    "redis": {
      "status": "healthy",
      "message": "Redis connection successful"
    }
  }
}
```

### 2. 用户管理

#### POST /api/v1/users/register
注册新用户

**请求**:
```json
{
  "username": "testuser",
  "email": "test@example.com",
  "password": "password123",
  "nickname": "Test User"
}
```

**响应**:
```json
{
  "code": 200,
  "message": "User registered successfully",
  "data": {
    "id": 1,
    "username": "testuser",
    "email": "test@example.com",
    "nickname": "Test User",
    "created_at": "2025-08-15T09:35:00Z"
  }
}
```

#### POST /api/v1/users/login
用户登录

**请求**:
```json
{
  "username": "testuser",
  "password": "password123"
}
```

**响应**:
```json
{
  "code": 200,
  "message": "Login successful",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "username": "testuser",
      "email": "test@example.com",
      "nickname": "Test User"
    }
  }
}
```

#### GET /api/v1/users/profile
获取用户信息（需要认证）

**响应**:
```json
{
  "code": 200,
  "message": "Profile retrieved successfully",
  "data": {
    "id": 1,
    "username": "testuser",
    "email": "test@example.com",
    "nickname": "Test User",
    "role": "user",
    "created_at": "2025-08-15T09:35:00Z",
    "updated_at": "2025-08-15T09:35:00Z"
  }
}
```

### 3. 产品分类管理

#### GET /api/v1/admin/productcategories
获取产品分类列表（需要管理员权限）

**参数**:
- `page`: 页码（默认：1）
- `page_size`: 每页数量（默认：10）
- `search`: 搜索关键词
- `sort`: 排序字段
- `order`: 排序方向（asc/desc）

**响应**:
```json
{
  "code": 200,
  "message": "Product categories retrieved successfully",
  "data": {
    "items": [
      {
        "id": 1,
        "name": "Electronics",
        "description": "Electronic devices and accessories",
        "parent_id": 0,
        "sort_order": 1,
        "is_active": true,
        "created_at": "2025-08-15T09:35:00Z",
        "updated_at": "2025-08-15T09:35:00Z"
      }
    ],
    "total": 1,
    "page": 1,
    "page_size": 10
  }
}
```

#### POST /api/v1/admin/productcategories
创建产品分类（需要管理员权限）

**请求**:
```json
{
  "name": "Electronics",
  "description": "Electronic devices and accessories",
  "parent_id": 0,
  "sort_order": 1,
  "is_active": true
}
```

**响应**:
```json
{
  "code": 200,
  "message": "Product category created successfully",
  "data": {
    "id": 1,
    "name": "Electronics",
    "description": "Electronic devices and accessories",
    "parent_id": 0,
    "sort_order": 1,
    "is_active": true,
    "created_at": "2025-08-15T09:35:00Z",
    "updated_at": "2025-08-15T09:35:00Z"
  }
}
```

#### GET /api/v1/admin/productcategories/{id}
获取产品分类详情（需要管理员权限）

**响应**:
```json
{
  "code": 200,
  "message": "Product category retrieved successfully",
  "data": {
    "id": 1,
    "name": "Electronics",
    "description": "Electronic devices and accessories",
    "parent_id": 0,
    "sort_order": 1,
    "is_active": true,
    "created_at": "2025-08-15T09:35:00Z",
    "updated_at": "2025-08-15T09:35:00Z"
  }
}
```

#### PUT /api/v1/admin/productcategories/{id}
更新产品分类（需要管理员权限）

**请求**:
```json
{
  "name": "Electronics Updated",
  "description": "Updated description",
  "sort_order": 2,
  "is_active": true
}
```

**响应**:
```json
{
  "code": 200,
  "message": "Product category updated successfully",
  "data": {
    "id": 1,
    "name": "Electronics Updated",
    "description": "Updated description",
    "parent_id": 0,
    "sort_order": 2,
    "is_active": true,
    "created_at": "2025-08-15T09:35:00Z",
    "updated_at": "2025-08-15T09:35:00Z"
  }
}
```

#### DELETE /api/v1/admin/productcategories/{id}
删除产品分类（需要管理员权限）

**响应**:
```json
{
  "code": 200,
  "message": "Product category deleted successfully"
}
```

### 4. 产品管理

#### GET /api/v1/admin/products
获取产品列表（需要管理员权限）

**参数**:
- `page`: 页码（默认：1）
- `page_size`: 每页数量（默认：10）
- `search`: 搜索关键词
- `category_id`: 分类ID筛选
- `sort`: 排序字段
- `order`: 排序方向（asc/desc）

**响应**:
```json
{
  "code": 200,
  "message": "Products retrieved successfully",
  "data": {
    "items": [
      {
        "id": 1,
        "name": "Smartphone",
        "description": "Latest smartphone model",
        "category_id": 1,
        "sku": "PHONE-001",
        "price": 599.99,
        "cost_price": 450.00,
        "stock_quantity": 100,
        "min_stock": 10,
        "is_active": true,
        "weight": 0.2,
        "dimensions": "15x7x0.8 cm",
        "created_at": "2025-08-15T09:35:00Z",
        "updated_at": "2025-08-15T09:35:00Z"
      }
    ],
    "total": 1,
    "page": 1,
    "page_size": 10
  }
}
```

#### POST /api/v1/admin/products
创建产品（需要管理员权限）

**请求**:
```json
{
  "name": "Smartphone",
  "description": "Latest smartphone model",
  "category_id": 1,
  "sku": "PHONE-001",
  "price": 599.99,
  "cost_price": 450.00,
  "stock_quantity": 100,
  "min_stock": 10,
  "is_active": true,
  "weight": 0.2,
  "dimensions": "15x7x0.8 cm"
}
```

**响应**:
```json
{
  "code": 200,
  "message": "Product created successfully",
  "data": {
    "id": 1,
    "name": "Smartphone",
    "description": "Latest smartphone model",
    "category_id": 1,
    "sku": "PHONE-001",
    "price": 599.99,
    "cost_price": 450.00,
    "stock_quantity": 100,
    "min_stock": 10,
    "is_active": true,
    "weight": 0.2,
    "dimensions": "15x7x0.8 cm",
    "created_at": "2025-08-15T09:35:00Z",
    "updated_at": "2025-08-15T09:35:00Z"
  }
}
```

#### GET /api/v1/admin/products/{id}
获取产品详情（需要管理员权限）

**响应**:
```json
{
  "code": 200,
  "message": "Product retrieved successfully",
  "data": {
    "id": 1,
    "name": "Smartphone",
    "description": "Latest smartphone model",
    "category_id": 1,
    "sku": "PHONE-001",
    "price": 599.99,
    "cost_price": 450.00,
    "stock_quantity": 100,
    "min_stock": 10,
    "is_active": true,
    "weight": 0.2,
    "dimensions": "15x7x0.8 cm",
    "created_at": "2025-08-15T09:35:00Z",
    "updated_at": "2025-08-15T09:35:00Z"
  }
}
```

#### PUT /api/v1/admin/products/{id}
更新产品（需要管理员权限）

**请求**:
```json
{
  "name": "Smartphone Pro",
  "price": 699.99,
  "stock_quantity": 150
}
```

**响应**:
```json
{
  "code": 200,
  "message": "Product updated successfully",
  "data": {
    "id": 1,
    "name": "Smartphone Pro",
    "description": "Latest smartphone model",
    "category_id": 1,
    "sku": "PHONE-001",
    "price": 699.99,
    "cost_price": 450.00,
    "stock_quantity": 150,
    "min_stock": 10,
    "is_active": true,
    "weight": 0.2,
    "dimensions": "15x7x0.8 cm",
    "created_at": "2025-08-15T09:35:00Z",
    "updated_at": "2025-08-15T09:35:00Z"
  }
}
```

#### DELETE /api/v1/admin/products/{id}
删除产品（需要管理员权限）

**响应**:
```json
{
  "code": 200,
  "message": "Product deleted successfully"
}
```

### 5. 批量操作

#### POST /api/v1/admin/products/batch-update-prices
批量更新产品价格（需要管理员权限）

**请求**:
```json
{
  "updates": [
    {
      "product_id": 1,
      "price": 699.99,
      "cost_price": 500.00
    },
    {
      "product_id": 2,
      "price": 299.99,
      "cost_price": 200.00
    }
  ]
}
```

**响应**:
```json
{
  "code": 200,
  "message": "Product prices updated successfully"
}
```

#### POST /api/v1/admin/products/batch-update-status
批量更新产品状态（需要管理员权限）

**请求**:
```json
{
  "updates": [
    {
      "product_id": 1,
      "is_active": true
    },
    {
      "product_id": 2,
      "is_active": false
    }
  ]
}
```

**响应**:
```json
{
  "code": 200,
  "message": "Product status updated successfully"
}
```

## 错误响应

所有 API 端点都可能返回以下错误响应：

### 400 Bad Request
```json
{
  "code": 400,
  "message": "Invalid request parameters",
  "error": "Validation failed"
}
```

### 401 Unauthorized
```json
{
  "code": 401,
  "message": "Unauthorized",
  "error": "Invalid or expired token"
}
```

### 403 Forbidden
```json
{
  "code": 403,
  "message": "Forbidden",
  "error": "Insufficient permissions"
}
```

### 404 Not Found
```json
{
  "code": 404,
  "message": "Resource not found",
  "error": "Product not found"
}
```

### 500 Internal Server Error
```json
{
  "code": 500,
  "message": "Internal server error",
  "error": "Database connection failed"
}
```

## 数据模型

### ProductCategory
```json
{
  "id": 1,
  "name": "Electronics",
  "description": "Electronic devices",
  "parent_id": 0,
  "sort_order": 1,
  "is_active": true,
  "created_at": "2025-08-15T09:35:00Z",
  "updated_at": "2025-08-15T09:35:00Z"
}
```

### Product
```json
{
  "id": 1,
  "name": "Smartphone",
  "description": "Latest smartphone",
  "category_id": 1,
  "sku": "PHONE-001",
  "price": 599.99,
  "cost_price": 450.00,
  "stock_quantity": 100,
  "min_stock": 10,
  "is_active": true,
  "weight": 0.2,
  "dimensions": "15x7x0.8 cm",
  "created_at": "2025-08-15T09:35:00Z",
  "updated_at": "2025-08-15T09:35:00Z"
}
```

## 使用示例

### JavaScript (Axios)
```javascript
// 登录
const login = async () => {
  const response = await axios.post('/api/v1/users/login', {
    username: 'admin',
    password: 'vibecoding'
  });
  const token = response.data.data.token;
  return token;
};

// 获取产品列表
const getProducts = async (token) => {
  const response = await axios.get('/api/v1/admin/products', {
    headers: {
      'Authorization': `Bearer ${token}`
    }
  });
  return response.data.data;
};

// 创建产品
const createProduct = async (token, productData) => {
  const response = await axios.post('/api/v1/admin/products', productData, {
    headers: {
      'Authorization': `Bearer ${token}`,
      'Content-Type': 'application/json'
    }
  });
  return response.data.data;
};
```

### Python (requests)
```python
import requests

# 登录
def login():
    response = requests.post('http://localhost:8080/api/v1/users/login', json={
        'username': 'admin',
        'password': 'vibecoding'
    })
    return response.json()['data']['token']

# 获取产品列表
def get_products(token):
    headers = {'Authorization': f'Bearer {token}'}
    response = requests.get('http://localhost:8080/api/v1/admin/products', headers=headers)
    return response.json()['data']

# 创建产品
def create_product(token, product_data):
    headers = {
        'Authorization': f'Bearer {token}',
        'Content-Type': 'application/json'
    }
    response = requests.post('http://localhost:8080/api/v1/admin/products', json=product_data, headers=headers)
    return response.json()['data']
```

### cURL
```bash
# 登录
curl -X POST http://localhost:8080/api/v1/users/login \
  -H "Content-Type: application/json" \
  -d '{"username": "admin", "password": "vibecoding"}'

# 获取产品列表
curl -X GET http://localhost:8080/api/v1/admin/products \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"

# 创建产品
curl -X POST http://localhost:8080/api/v1/admin/products \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Smartphone",
    "description": "Latest smartphone",
    "category_id": 1,
    "sku": "PHONE-001",
    "price": 599.99,
    "cost_price": 450.00,
    "stock_quantity": 100,
    "min_stock": 10,
    "is_active": true,
    "weight": 0.2,
    "dimensions": "15x7x0.8 cm"
  }'
```

## 限制和配额

- **请求频率限制**: 1000 请求/分钟
- **分页限制**: 最大每页 100 条记录
- **文件上传限制**: 10MB
- **请求超时**: 30 秒

## 版本历史

### v1.0 (2025-08-15)
- 初始版本发布
- 支持用户管理
- 支持产品分类管理
- 支持产品管理
- 支持批量操作

## 支持

如有问题或需要技术支持，请联系：
- **API 文档**: [Swagger UI](http://localhost:8080/swagger/index.html)
- **技术支持**: support@example.com
- **开发团队**: dev@example.com

---

*API 文档版本: 1.0*  
*最后更新: 2025-08-15*  
*维护者: Vibe Coding Starter Team*