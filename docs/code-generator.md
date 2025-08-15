# Vibe Coding Starter - 代码生成器使用指南

## 概述

本项目包含一个强大的代码生成器，支持多种生成模式，可以快速生成符合项目规范的完整业务模块代码。生成器具备以下核心功能：

- 🚀 **智能代码生成**: 支持手动字段定义和数据库表结构读取
- 🎨 **前端代码生成**: 自动生成 Antd/Vue 前端组件
- 🔧 **增强模块生成**: 自动路由注册、数据库迁移、国际化支持
- 📊 **数据库集成**: 支持从现有数据库表生成完整模型
- 🧪 **完整测试覆盖**: 自动生成所有层的单元测试

## 核心特性

### 1. 多种生成模式

| 命令 | 功能描述 | 适用场景 |
|------|----------|----------|
| `all` | 生成完整业务模块 | 快速开始新功能开发 |
| `enhanced` | 生成增强模块（含前端） | 全栈开发，需要前后端同时生成 |
| `module` | 传统模块生成 | 兼容旧版本，生成基础后端代码 |
| `frontend` | 仅生成前端代码 | 已有后端，需要前端界面 |

### 2. 智能字段推断

生成器能够根据字段名称和类型自动推断：
- 表单控件类型（input、textarea、switch、datetime等）
- 搜索字段配置
- 验证规则
- 国际化标签

### 3. 自动化集成

- **自动路由注册**: 更新 `server.go` 和 `main.go`
- **自动数据库迁移**: 生成并执行迁移脚本
- **自动国际化**: 生成中英文语言包
- **自动前端集成**: 更新路由配置和国际化文件

## 文件命名规范

### 统一命名原则

为了保持代码库的一致性和简洁性，生成的文件采用统一的命名规范：

#### 后端文件结构
```
internal/
├── model/
│   └── {model_name}.go              # 数据模型
├── repository/
│   ├── interfaces.go               # 仓储接口（追加）
│   └── {model_name}_repository.go   # 仓储实现
├── service/
│   └── {model_name}_service.go      # 业务逻辑层
└── handler/
    └── {model_name}.go              # API处理器（含请求结构体）

test/
├── handler/
│   └── {model_name}_handler_test.go
├── service/
│   └── {model_name}_service_test.go
└── repository/
    └── {model_name}_repository_test.go

migrations/{db_type}/
└── {timestamp}_create_{table_name}_table.sql
```

#### 前端文件结构（Antd）
```
src/
├── pages/
│   ├── admin/          # 管理后台模块
│   │   └── {module}/
│   │       └── index.tsx
│   └── {module}/       # 公共模块
│       └── index.tsx
├── services/
│   └── {module}/
│       ├── api.ts         # API服务
│       └── typings.d.ts   # 类型定义
└── locales/
    ├── zh-CN/
    │   └── {module}.ts    # 中文语言包
    └── en-US/
        └── {module}.ts    # 英文语言包
```

## 使用方法

### 🚀 基础生成 - `all` 命令

**推荐用于快速开始**，生成完整的后端业务模块。

#### 方式一：手动定义字段
```bash
# 基础用法
go run cmd/generator/main.go all --name=Product --fields="name:string,description:string,price:float64,active:bool"

# 包含认证和缓存
go run cmd/generator/main.go all --name=Order --fields="total:float64,status:string" --auth --cache
```

#### 方式二：从数据库表生成
```bash
# 从本地数据库生成
go run cmd/generator/main.go all --name=Product --table=products \
  --host=localhost --port=3306 --user=root --password=secret --database=mydb

# 使用 k3d 环境数据库
go run cmd/generator/main.go all --name=Product --table=products \
  --host=127.0.0.1 --port=3306 --user=vibe_user --password=vibe_password --database=vibe_coding_starter
```

**生成内容：**
- ✅ 数据模型（Model）
- ✅ 数据访问层（Repository + 接口）
- ✅ 业务逻辑层（Service + Mock）
- ✅ API处理器（Handler + 请求结构体）
- ✅ 数据库迁移文件
- ✅ 完整测试套件
- ✅ 自动路由注册

### 🎨 增强生成 - `enhanced` 命令

**推荐用于全栈开发**，生成后端 + 前端完整模块。

```bash
# 基础增强模块（仅后端）
go run cmd/generator/main.go enhanced --name=ProductStockHistory \
  --fields="product_id:uint,change_type:string,quantity_change:int,reason:string"

# 完整增强模块（后端 + 前端）
go run cmd/generator/main.go enhanced --name=ProductStockHistory \
  --fields="product_id:uint,change_type:string,quantity_change:int,reason:string" \
  --frontend-output=../vibe-coding-starter-ui-antd \
  --frontend-framework=antd \
  --frontend-module-type=admin

# 高级配置（启用所有功能）
go run cmd/generator/main.go enhanced --name=ProductStockHistory \
  --fields="product_id:uint,change_type:string,quantity_change:int,reason:string" \
  --frontend-output=../vibe-coding-starter-ui-antd \
  --auto-route=true \
  --auto-migration=true \
  --auto-i18n=true \
  --smart-search=true
```

**增强功能：**
- 🔗 **自动路由注册**: 更新后端路由配置
- 🗄️ **自动数据库迁移**: 生成并执行迁移脚本
- 🌍 **自动国际化**: 生成中英文语言包并更新配置
- 🔍 **智能搜索字段**: 根据字段名称自动配置搜索功能
- 🎨 **前端代码生成**: 生成完整的 Antd 管理界面

### 📱 前端生成 - `frontend` 命令

**用于已有后端模型**，仅生成前端代码。

```bash
# 生成管理后台前端
go run cmd/generator/main.go frontend --model=Product \
  --framework=antd \
  --output=../vibe-coding-starter-ui-antd \
  --module-type=admin \
  --with-auth \
  --with-search \
  --with-export

# 生成公共页面前端
go run cmd/generator/main.go frontend --model=Article \
  --framework=antd \
  --output=../vibe-coding-starter-ui-antd \
  --module-type=public \
  --api-prefix=/api/v1
```

**前端功能支持：**
- 📋 CRUD 操作界面
- 🔍 高级搜索和筛选
- 📊 数据表格展示
- 📤 数据导出功能
- 🔄 批量操作
- 🌐 多语言支持

### 📊 数据库集成命令

#### 列出数据库表
```bash
go run cmd/generator/main.go list-tables \
  --host=localhost --port=3306 --user=root --password=secret --database=mydb
```

#### 从单个表生成模型
```bash
# 基础用法
go run cmd/generator/main.go from-table --table=users \
  --host=localhost --port=3306 --user=root --password=secret --database=mydb

# 完整配置
go run cmd/generator/main.go from-table --table=users \
  --model=CustomUser \
  --host=localhost --port=3306 --user=root --password=secret --database=mydb \
  --timestamps=true --soft-delete=false
```

#### 从数据库生成所有模型
```bash
go run cmd/generator/main.go from-db \
  --host=localhost --port=3306 --user=root --password=secret --database=mydb
```

### 🔧 单独组件生成

所有组件生成统一使用 `--model` 参数：

```bash
# 生成模型
go run cmd/generator/main.go model --name=Product --fields="name:string,price:float64"

# 生成仓储
go run cmd/generator/main.go repository --model=Product

# 生成服务
go run cmd/generator/main.go service --model=Product

# 生成处理器
go run cmd/generator/main.go handler --model=Product

# 生成测试
go run cmd/generator/main.go test --model=Product

# 生成迁移
go run cmd/generator/main.go migration --model=Product
```

## 高级配置

### 字段类型支持

| Go 类型 | 说明 | 示例 |
|---------|------|------|
| `string` | 字符串 | `name:string` |
| `int`, `int32`, `int64` | 整数 | `age:int` |
| `uint`, `uint32`, `uint64` | 无符号整数 | `id:uint` |
| `float32`, `float64` | 浮点数 | `price:float64` |
| `bool` | 布尔值 | `active:bool` |
| `time.Time` | 时间 | `created_at:time.Time` |

### 数据库连接参数

| 参数 | 默认值 | 说明 |
|------|--------|------|
| `--host` | localhost | 数据库主机 |
| `--port` | 3306 | 数据库端口 |
| `--user` | root | 数据库用户名 |
| `--password` | - | 数据库密码 |
| `--database` | - | 数据库名称 |
| `--table` | - | 表名称 |
| `--timestamps` | true | 包含时间戳字段 |
| `--soft-delete` | false | 包含软删除字段 |

### 前端生成参数

| 参数 | 说明 | 默认值 |
|------|------|--------|
| `--framework` | 前端框架 | `antd` |
| `--output` | 输出目录 | 必需 |
| `--module-type` | 模块类型 | `admin` |
| `--with-auth` | 包含认证 | `false` |
| `--with-search` | 包含搜索 | `true` |
| `--with-export` | 包含导出 | `false` |
| `--with-batch` | 包含批量操作 | `false` |
| `--api-prefix` | API 前缀 | `/api/v1` |

## 增强功能详解

### 1. 智能字段推断

生成器会根据字段名称自动推断最适合的配置：

```go
// 字段名称推断示例
"email"    -> 类型: email,    搜索: true,  表单: email
"password" -> 类型: password, 搜索: false, 表单: password
"name"     -> 类型: input,    搜索: true,  表单: input
"active"   -> 类型: switch,   搜索: false, 表单: switch
"price"    -> 类型: number,   搜索: false, 表单: number
"content"  -> 类型: textarea, 搜索: true,  表单: textarea
```

### 2. 自动化集成

#### 后端集成
- **路由注册**: 自动更新 `internal/server/server.go` 和 `cmd/server/main.go`
- **依赖注入**: 自动添加到 Uber FX 容器
- **中间件配置**: 根据模块类型自动配置认证和权限中间件

#### 前端集成
- **路由配置**: 自动更新 `config/routes.ts`
- **国际化**: 自动更新 `src/locales/zh-CN.ts` 和 `src/locales/en-US.ts`
- **菜单集成**: 自动添加到管理后台菜单

### 3. 数据库迁移自动化

```bash
# 生成迁移文件
migrations/mysql/20240101_120000_create_products_table.sql

# 自动执行迁移（可选）
# 生成迁移脚本: cmd/automigrate_product/main.go
# 执行迁移: go run cmd/automigrate_product/main.go
```

### 4. 国际化支持

生成器会自动生成中英文语言包：

```typescript
// src/locales/zh-CN/product.ts
export default {
  productId: '产品ID',
  name: '产品名称',
  price: '产品价格',
  // ...
};

// src/locales/en-US/product.ts
export default {
  productId: 'Product ID',
  name: 'Product Name',
  price: 'Product Price',
  // ...
};
```

## 最佳实践

### 1. 项目初始化

```bash
# 1. 生成用户管理模块（完整功能）
go run cmd/generator/main.go enhanced --name=User \
  --fields="username:string,email:string,password:string,active:bool" \
  --frontend-output=../vibe-coding-starter-ui-antd \
  --auto-route=true \
  --auto-migration=true \
  --auto-i18n=true

# 2. 生成产品管理模块
go run cmd/generator/main.go enhanced --name=Product \
  --fields="name:string,description:string,price:float64,stock:int" \
  --frontend-output=../vibe-coding-starter-ui-antd
```

### 2. 数据库优先开发

```bash
# 1. 设计数据库表
# 2. 从数据库生成模型
go run cmd/generator/main.go from-table --table=orders \
  --host=localhost --port=3306 --user=root --password=secret --database=mydb

# 3. 生成完整模块
go run cmd/generator/main.go enhanced --name=Order \
  --fields="user_id:uint,total:float64,status:string" \
  --frontend-output=../vibe-coding-starter-ui-antd
```

### 3. 迭代开发

```bash
# 添加新功能模块
go run cmd/generator/main.go enhanced --name=OrderItem \
  --fields="order_id:uint,product_id:uint,quantity:int,price:float64" \
  --frontend-output=../vibe-coding-starter-ui-antd

# 生成报表模块
go run cmd/generator/main.go enhanced --name=SalesReport \
  --fields="date:time.Time,total_amount:float64,order_count:int" \
  --frontend-output=../vibe-coding-starter-ui-antd \
  --module-type=admin
```

## 故障排除

### 常见问题

1. **数据库连接失败**
   ```bash
   # 检查数据库连接参数
   go run cmd/generator/main.go list-tables \
     --host=localhost --port=3306 --user=root --password=secret --database=mydb
   ```

2. **前端输出目录无效**
   ```bash
   # 确保前端项目目录存在
   ls ../vibe-coding-starter-ui-antd/package.json
   ls ../vibe-coding-starter-ui-antd/src
   ```

3. **文件已存在冲突**
   ```bash
   # 生成器会提示覆盖确认，或者先删除现有文件
   rm internal/model/product.go
   rm internal/repository/product.go
   ```

### 验证生成结果

```bash
# 编译验证
go build ./...

# 运行测试
go test ./test/... -v

# 检查生成文件
ls internal/model/
ls internal/repository/
ls internal/service/
ls internal/handler/
ls test/
ls migrations/
```

## 版本信息

- **当前版本**: v2.0.0
- **支持的 Go 版本**: 1.19+
- **支持的数据库**: MySQL, PostgreSQL, SQLite
- **支持的前端框架**: Ant Design (Vue 开发中)
- **生成器类型**: 后端 + 前端全栈生成器

## 更新日志

### v2.0.0 (2025-01-15)
- ✨ 新增增强模块生成器 (`enhanced` 命令)
- 🎨 新增前端代码生成功能
- 🔗 新增自动路由注册
- 🗄️ 新增自动数据库迁移
- 🌍 新增自动国际化支持
- 🔍 新增智能搜索字段配置
- 📱 新增 Antd 管理界面生成

### v1.1.0 (2024-12-01)
- 📊 新增数据库表结构读取功能
- 🧪 改进测试代码生成
- 🔧 修复模板变量命名问题
- 📝 完善文档和示例

### v1.0.1 (2025-07-31)
- 🐛 修复 Service 模板命名不一致问题
- 🔧 修复 Repository 模板接收者类型错误
- 📁 统一 Handler 请求结构体到主文件

