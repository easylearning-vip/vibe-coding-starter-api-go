# 代码生成器使用指南

## 概述

本项目包含一个强大的代码生成器，可以快速生成符合项目规范的业务模块代码。

## 文件命名规范

### 统一命名原则

为了保持代码库的一致性和简洁性，生成的文件采用统一的命名规范：

#### 主要文件命名
- **Model**: `internal/model/{model_name}.go`
- **Repository**: `internal/repository/{model_name}.go`
- **Service**: `internal/service/{model_name}.go`
- **Handler**: `internal/handler/{model_name}.go`

#### 示例
对于 `Product` 模型，生成的文件结构如下：

```
internal/
├── model/
│   └── product.go              # Product 数据模型
├── repository/
│   └── product.go              # Product 数据访问层
├── service/
│   └── product.go              # Product 业务逻辑层
└── handler/
    └── product.go              # Product API 处理器（包含请求结构体）
```

#### 测试文件命名
- **Handler 测试**: `test/handler/{model_name}_handler_test.go`
- **Service 测试**: `test/service/{model_name}_service_test.go`
- **Repository 测试**: `test/repository/{model_name}_repository_test.go`

#### 数据库迁移文件
根据配置文件中的数据库类型自动选择目录：
- **MySQL**: `migrations/mysql/{timestamp}_{migration_name}.sql`
- **PostgreSQL**: `migrations/postgres/{timestamp}_{migration_name}.sql`
- **SQLite**: `migrations/sqlite/{timestamp}_{migration_name}.sql`

## 使用方法

### 🚀 生成所有组件（推荐）

**新功能！** 支持两种方式一键生成模型的所有组件：

#### 方式一：手动定义字段
```bash
go run cmd/generator/main.go all --name=Product --fields="name:string,description:string,price:float64,active:bool"
```

#### 方式二：从数据库表生成（最新功能）
```bash
# 从数据库表结构生成完整的业务模块
go run cmd/generator/main.go all --name=Product --table=products \
  --host=localhost --port=3306 --user=root --password=secret --database=mydb

# 使用k3d环境的数据库
go run cmd/generator/main.go all --name=Product --table=products \
  --host=127.0.0.1 --port=3306 --user=vibe_user --password=vibe_password --database=vibe_coding_starter
```

这个命令会按正确的依赖顺序生成：
1. **Model** - 数据模型
2. **Repository** - 数据访问层 + 接口定义
3. **Service** - 业务逻辑层 + MockService
4. **Handler** - API 处理器
5. **Migration** - 数据库迁移

生成的文件：
- `internal/model/product.go`
- `internal/repository/product.go`
- `internal/service/product.go`
- `internal/handler/product.go`
- `migrations/{db_type}/{timestamp}_create_products_table.sql`
- `test/handler/product_handler_test.go`
- `test/service/product_service_test.go`
- `test/repository/product_repository_test.go`
- MockProductService 自动添加到 `test/mocks/service_mocks.go`

### 生成完整模块（传统方式）

```bash
go run cmd/generator/main.go module --name=Product --fields="name:string,description:string,price:float64,active:bool"
```

### 生成单独组件

**统一使用 `--model` 参数，组件名称自动按命名规范生成：**

#### 生成模型
```bash
go run cmd/generator/main.go model --name=Product --fields="name:string,price:float64,active:bool"
```

#### 生成仓储（自动命名为 ProductRepository）
```bash
go run cmd/generator/main.go repository --model=Product
```

#### 生成服务（自动命名为 ProductService，自动生成 MockService）
```bash
go run cmd/generator/main.go service --model=Product
```

#### 生成处理器（自动命名为 ProductHandler）
```bash
go run cmd/generator/main.go handler --model=Product
```

#### 生成测试（自动为所有组件生成测试）
```bash
go run cmd/generator/main.go test --model=Product
```

#### 生成数据库迁移
```bash
# 使用模型名称自动生成迁移名称
go run cmd/generator/main.go migration --model=Product

# 或手动指定迁移名称
go run cmd/generator/main.go migration --name=create_products_table
```

**命名规范：**
- Model: `Product`
- Repository: `ProductRepository`
- Service: `ProductService`
- Handler: `ProductHandler`
- Mock: `MockProductService`
- Migration: `create_products_table` (从模型名自动生成)

### 📊 数据库表相关命令

#### 列出数据库中的所有表
```bash
go run cmd/generator/main.go list-tables --host=localhost --port=3306 --user=root --password=secret --database=mydb
```

#### 从单个数据库表生成模型
```bash
go run cmd/generator/main.go from-table --table=users --host=localhost --port=3306 --user=root --password=secret --database=mydb

# 可选参数
go run cmd/generator/main.go from-table --table=users --model=CustomUser \
  --host=localhost --port=3306 --user=root --password=secret --database=mydb \
  --timestamps=true --soft-delete=false
```

#### 从数据库中的所有表生成模型
```bash
go run cmd/generator/main.go from-db --host=localhost --port=3306 --user=root --password=secret --database=mydb
```

#### 数据库连接参数说明
| 参数 | 默认值 | 说明 |
|------|--------|------|
| `--host` | localhost | 数据库主机地址 |
| `--port` | 3306 | 数据库端口 |
| `--user` | root | 数据库用户名 |
| `--password` | (空) | 数据库密码 |
| `--database` | (必需) | 数据库名称 |
| `--table` | (必需) | 表名称（仅用于from-table命令） |
| `--model` | (可选) | 自定义模型名称，默认从表名生成 |
| `--timestamps` | true | 是否包含created_at和updated_at字段 |
| `--soft-delete` | false | 是否包含deleted_at字段用于软删除 |

#### 支持的数据类型映射
| MySQL类型 | Go类型 | 说明 |
|-----------|--------|------|
| `VARCHAR`, `CHAR`, `TEXT` | `string` / `sql.NullString` | 字符串类型，可空字段使用Null类型 |
| `TINYINT`, `SMALLINT`, `INT` | `int8`, `int16`, `int32` / `sql.NullInt32` | 整数类型 |
| `BIGINT` | `int64` / `sql.NullInt64` | 64位整数类型 |
| `DECIMAL`, `FLOAT`, `DOUBLE` | `float64` / `sql.NullFloat64` | 浮点数类型 |
| `BOOLEAN`, `TINYINT(1)` | `bool` / `sql.NullBool` | 布尔类型，自动识别tinyint(1) |
| `DATE`, `DATETIME`, `TIMESTAMP` | `time.Time` / `sql.NullTime` | 时间类型 |
| `ENUM` | `sql.NullString` | 枚举类型，映射为字符串 |
| `JSON` | `string` | JSON类型，映射为字符串 |
| `BLOB`, `BINARY` | `[]byte` | 二进制数据类型 |

#### 字段跳过规则
生成器会自动跳过以下字段（因为BaseModel已提供）：
- `id` (主键)
- `created_at` (创建时间)
- `updated_at` (更新时间)
- `deleted_at` (软删除时间)

## 文件内容结构

### Handler 文件特点

Handler 文件现在包含所有相关的结构体定义：

```go
// 主要的 Handler 结构体和方法
type ProductHandler struct { ... }
func (h *ProductHandler) Create(c *gin.Context) { ... }
func (h *ProductHandler) GetByID(c *gin.Context) { ... }
// ... 其他 CRUD 方法

// 请求结构体（在同一文件中）
type CreateProductRequest struct {
    Name        string `json:"name" validate:"required,min=1,max=255"`
    Description string `json:"description" validate:"max=1000"`
}

type UpdateProductRequest struct {
    Name        *string `json:"name,omitempty" validate:"omitempty,min=1,max=255"`
    Description *string `json:"description,omitempty" validate:"omitempty,max=1000"`
}
```

### 优势

1. **简洁性**: 每个层级只有一个主文件，避免文件名重复
2. **一致性**: 所有模块遵循相同的命名规范
3. **可维护性**: 相关代码集中在一个文件中，便于维护
4. **可读性**: 文件名直接反映模块名称，易于理解

## 配置文件支持

生成器会自动读取 `configs/config.yaml` 文件来确定：
- 数据库类型（MySQL/PostgreSQL/SQLite）
- 迁移文件生成目录
- 其他项目特定配置

## 编译验证

生成代码后，建议运行以下命令验证：

```bash
# 编译验证
go build ./...

# 运行测试
go test ./test/... -v
```

## 注意事项

1. 生成器会检查文件是否已存在，避免意外覆盖
2. 所有生成的代码都包含完整的错误处理和日志记录
3. 生成的测试文件包含基本的单元测试用例
4. 数据库迁移文件会根据配置自动选择正确的 SQL 语法

## 模板修复记录

### v1.0.1 更新 (2025-07-31)

**修复的模板问题**:
1. **Service 模板**: 修复了接口和实现类型命名不一致的问题
   - 接口名: `ProductService` (正确)
   - 实现类型: `productService` (小写开头，正确)
   - 构造函数: `NewProductService` (正确)

2. **Repository 模板**: 修复了 `applyFilters` 方法接收者类型错误
   - 从 `*{{.NameCamel}}Repository` 改为 `*{{.ModelCamel}}Repository`

3. **Handler 模板**: 统一请求结构体到主文件中
   - 不再生成单独的 `*_requests.go` 文件
   - 请求结构体直接包含在 handler 文件末尾

**变量使用规范**:
- `{{.Model}}`: 模型名称 (如 `Product`)
- `{{.ModelCamel}}`: 模型驼峰命名 (如 `product`)
- `{{.ModelSnake}}`: 模型蛇形命名 (如 `product`)
- `{{.Name}}`: 服务/仓储全名 (如 `ProductService`)
- `{{.NameCamel}}`: 服务/仓储驼峰命名 (如 `productService`)

## 命令参考

### 🚀 all - 生成所有组件

**推荐使用！** 一键生成模型的所有组件，按正确的依赖顺序执行。

```bash
go run cmd/generator/main.go all --name=<ModelName> --fields="<field_definitions>" [--auth] [--cache]
```

**参数：**
- `--name`: 模型名称（必需）
- `--fields`: 字段定义（可选）
- `--auth`: 包含认证中间件（可选）
- `--cache`: 包含缓存支持（可选）

**示例：**
```bash
# 基本用法
go run cmd/generator/main.go all --name=Product --fields="name:string,price:float64"

# 包含认证和缓存
go run cmd/generator/main.go all --name=Order --fields="total:float64,status:string" --auth --cache
```

**生成顺序：**
1. Model → 2. Repository → 3. Service (+ Mock) → 4. Handler → 5. Migration

### module - 生成完整模块

传统的模块生成方式，一次性生成所有文件。

```bash
go run cmd/generator/main.go module --name=<name> --fields="<field_definitions>"
```

### 单独组件生成

如果需要单独生成某个组件，**统一使用 `--model` 参数**，组件名称自动按命名规范生成：

```bash
# 生成模型
go run cmd/generator/main.go model --name=Product --fields="name:string,price:float64"

# 生成仓储（自动命名为 ProductRepository）
go run cmd/generator/main.go repository --model=Product

# 生成服务（自动命名为 ProductService，会自动生成 MockService）
go run cmd/generator/main.go service --model=Product

# 生成处理器（自动命名为 ProductHandler）
go run cmd/generator/main.go handler --model=Product

# 生成测试（自动为 ProductService、ProductRepository、ProductHandler 生成测试）
go run cmd/generator/main.go test --model=Product

# 生成迁移（使用模型名称自动生成）
go run cmd/generator/main.go migration --model=Product

# 或手动指定迁移名称
go run cmd/generator/main.go migration --name=create_products_table
```

**命名规范：**
- Model: `Product`
- Repository: `ProductRepository`
- Service: `ProductService`
- Handler: `ProductHandler`
- Mock: `MockProductService`

### 字段类型支持

支持的字段类型：
- `string` - 字符串类型
- `int`, `int32`, `int64` - 整数类型
- `uint`, `uint32`, `uint64` - 无符号整数
- `float32`, `float64` - 浮点数类型
- `bool` - 布尔类型
- `time.Time` - 时间类型

**字段定义格式：**
```
--fields="field1:type1,field2:type2,field3:type3"
```

## 版本信息

- 生成器版本: v1.1.0
- 支持的 Go 版本: 1.19+
- 支持的数据库: MySQL, PostgreSQL, SQLite

