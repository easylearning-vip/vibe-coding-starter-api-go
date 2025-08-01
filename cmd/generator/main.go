package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"vibe-coding-starter/tools/generator/generator"
)

const (
	version = "1.0.0"
	banner  = `
██╗   ██╗██╗██████╗ ███████╗     ██████╗ ███████╗███╗   ██╗
██║   ██║██║██╔══██╗██╔════╝    ██╔════╝ ██╔════╝████╗  ██║
██║   ██║██║██████╔╝█████╗      ██║  ███╗█████╗  ██╔██╗ ██║
╚██╗ ██╔╝██║██╔══██╗██╔══╝      ██║   ██║██╔══╝  ██║╚██╗██║
 ╚████╔╝ ██║██████╔╝███████╗    ╚██████╔╝███████╗██║ ╚████║
  ╚═══╝  ╚═╝╚═════╝ ╚══════╝     ╚═════╝ ╚══════╝╚═╝  ╚═══╝

Vibe Coding Starter - Code Generator v%s
`
)

func main() {
	// 显示banner
	fmt.Printf(banner, version)

	// 解析命令行参数
	if len(os.Args) < 2 {
		showUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "all":
		handleAllCommand()
	case "module":
		handleModuleCommand()
	case "handler":
		handleHandlerCommand()
	case "test":
		handleTestCommand()
	case "service":
		handleServiceCommand()
	case "repository":
		handleRepositoryCommand()
	case "model":
		handleModelCommand()
	case "migration":
		handleMigrationCommand()
	case "from-table":
		handleFromTableCommand()
	case "from-db":
		handleFromDatabaseCommand()
	case "list-tables":
		handleListTablesCommand()
	case "frontend":
		handleFrontendCommand()
	case "help", "-h", "--help":
		showUsage()
	case "version", "-v", "--version":
		fmt.Printf("Vibe Code Generator v%s\n", version)
	default:
		fmt.Printf("Unknown command: %s\n\n", command)
		showUsage()
		os.Exit(1)
	}
}

func showUsage() {
	fmt.Println("Usage: go run cmd/generator/main.go <command> [options]")
	fmt.Println()
	fmt.Println("Available commands:")
	fmt.Println("  all         Generate all components for a model (model + repository + service + handler + migration)")
	fmt.Println("  module      Generate a complete business module")
	fmt.Println("  handler     Generate API handler")
	fmt.Println("  service     Generate service layer")
	fmt.Println("  repository  Generate repository layer")
	fmt.Println("  model       Generate data model")
	fmt.Println("  test        Generate test code")
	fmt.Println("  migration   Generate database migration")
	fmt.Println("  from-table  Generate model from database table")
	fmt.Println("  from-db     Generate models from all tables in database")
	fmt.Println("  list-tables List all tables in database")
	fmt.Println("  frontend    Generate frontend code (Antd/Vue)")
	fmt.Println("  help        Show this help message")
	fmt.Println("  version     Show version information")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  # Generate all components with manual field definition")
	fmt.Println("  go run cmd/generator/main.go all --name=Product --fields=\"name:string,price:float64\"")
	fmt.Println("  # Generate all components from database table")
	fmt.Println("  go run cmd/generator/main.go all --name=Product --table=products --host=localhost --port=3306 --user=root --password=secret --database=mydb")
	fmt.Println("  go run cmd/generator/main.go module --name=product")
	fmt.Println("  go run cmd/generator/main.go handler --model=Product")
	fmt.Println("  go run cmd/generator/main.go service --model=Product")
	fmt.Println("  go run cmd/generator/main.go repository --model=Product")
	fmt.Println("  go run cmd/generator/main.go test --model=Product")
	fmt.Println("  go run cmd/generator/main.go migration --model=Product")
	fmt.Println("  go run cmd/generator/main.go migration --name=create_products_table")
	fmt.Println("  go run cmd/generator/main.go from-table --table=users --host=localhost --port=3306 --user=root --password=secret --database=mydb")
	fmt.Println("  go run cmd/generator/main.go from-db --host=localhost --port=3306 --user=root --password=secret --database=mydb")
	fmt.Println("  go run cmd/generator/main.go list-tables --host=localhost --port=3306 --user=root --password=secret --database=mydb")
	fmt.Println("  # Generate frontend code")
	fmt.Println("  go run cmd/generator/main.go frontend --model=Product --framework=antd --output=../vibe-coding-starter-ui-antd")
	fmt.Println("  go run cmd/generator/main.go frontend --model=User --framework=antd --output=../vibe-coding-starter-ui-antd --with-auth --with-search")
	fmt.Println("  go run cmd/generator/main.go frontend --model=Product --framework=antd --module-type=admin --output=../vibe-coding-starter-ui-antd")
	fmt.Println("  go run cmd/generator/main.go frontend --model=Article --framework=antd --module-type=public --output=../vibe-coding-starter-ui-antd")
}

func handleAllCommand() {
	fs := flag.NewFlagSet("all", flag.ExitOnError)
	name := fs.String("name", "", "Model name (required)")
	fields := fs.String("fields", "", "Model fields (e.g., 'name:string,price:float64,active:bool')")
	withAuth := fs.Bool("auth", false, "Include authentication middleware")
	withCache := fs.Bool("cache", false, "Include cache support")

	// 数据库连接参数（可选）
	table := fs.String("table", "", "Database table name (optional, for reading structure from database)")
	host := fs.String("host", "localhost", "Database host (used with --table)")
	port := fs.Int("port", 3306, "Database port (used with --table)")
	user := fs.String("user", "root", "Database user (used with --table)")
	password := fs.String("password", "", "Database password (used with --table)")
	database := fs.String("database", "", "Database name (used with --table)")

	fs.Parse(os.Args[2:])

	if *name == "" {
		fmt.Println("Error: --name is required")
		fs.Usage()
		os.Exit(1)
	}

	// 如果指定了表名，则需要数据库连接信息
	if *table != "" && *database == "" {
		fmt.Println("Error: --database is required when using --table")
		fs.Usage()
		os.Exit(1)
	}

	var fieldsStr string
	var tableComment string

	// 如果指定了表名，从数据库读取字段结构
	if *table != "" {
		fmt.Printf("🔍 Reading table structure from database '%s.%s'...\n", *database, *table)

		reader, err := generator.NewTableReader(*host, *port, *user, *password, *database)
		if err != nil {
			log.Fatalf("Failed to connect to database: %v", err)
		}
		defer reader.Close()

		// 读取表结构
		dbFields, err := reader.ReadTableStructure(*table)
		if err != nil {
			log.Fatalf("Failed to read table structure: %v", err)
		}

		// 获取表注释
		tableComment, err = reader.GetTableComment(*table)
		if err != nil {
			log.Fatalf("Failed to get table comment: %v", err)
		}

		// 将数据库字段转换为字段字符串格式
		var fieldParts []string
		for _, field := range dbFields {
			// 转换Go类型为简化的字段类型字符串
			fieldType := convertGoTypeToFieldType(field.Type)
			fieldParts = append(fieldParts, fmt.Sprintf("%s:%s", generator.ToSnakeCase(field.Name), fieldType))
		}
		fieldsStr = strings.Join(fieldParts, ",")

		fmt.Printf("✅ Found %d fields in table '%s'\n", len(dbFields), *table)
		if tableComment != "" {
			fmt.Printf("💬 Table comment: %s\n", tableComment)
		}
	} else {
		fieldsStr = *fields
	}

	fmt.Printf("🚀 Generating all components for '%s'...\n\n", *name)

	// 步骤1: 生成 Model
	fmt.Println("📦 Step 1/5: Generating Model...")

	var modelGen generator.Generator
	var modelConfig interface{}

	if *table != "" {
		// 使用数据库表生成器
		modelGen = generator.NewDatabaseTableGenerator()
		modelConfig = &generator.DatabaseTableConfig{
			DatabaseHost:     *host,
			DatabasePort:     *port,
			DatabaseUser:     *user,
			DatabasePassword: *password,
			DatabaseName:     *database,
			TableName:        *table,
			ModelName:        *name,
			WithTimestamps:   true,
			WithSoftDelete:   false,
		}
	} else {
		// 使用传统的字段字符串生成器
		modelGen = generator.NewModelGenerator()
		modelConfig = &generator.ModelConfig{
			Name:   *name,
			Fields: fieldsStr,
		}
	}

	if err := modelGen.Generate(modelConfig); err != nil {
		log.Fatalf("Failed to generate model: %v", err)
	}
	fmt.Println("✅ Model generated successfully!")

	// 步骤2: 生成 Repository
	fmt.Println("\n🗄️  Step 2/5: Generating Repository...")
	repoGen := generator.NewRepositoryGenerator()
	repoConfig := &generator.RepositoryConfig{
		Name:  *name + "Repository",
		Model: *name,
	}
	if err := repoGen.Generate(repoConfig); err != nil {
		log.Fatalf("Failed to generate repository: %v", err)
	}
	fmt.Println("✅ Repository generated successfully!")

	// 步骤3: 生成 Service
	fmt.Println("\n⚙️  Step 3/5: Generating Service...")
	serviceGen := generator.NewServiceGenerator()
	serviceConfig := &generator.ServiceConfig{
		Name:      *name + "Service",
		Model:     *name,
		Fields:    fieldsStr, // 使用从数据库读取或手动指定的字段
		WithCache: *withCache,
	}
	if err := serviceGen.Generate(serviceConfig); err != nil {
		log.Fatalf("Failed to generate service: %v", err)
	}
	fmt.Println("✅ Service generated successfully!")

	// 步骤4: 生成 Handler
	fmt.Println("\n🌐 Step 4/5: Generating Handler...")
	handlerGen := generator.NewHandlerGenerator()
	handlerConfig := &generator.HandlerConfig{
		Model:    *name,
		WithAuth: *withAuth,
	}
	if err := handlerGen.Generate(handlerConfig); err != nil {
		log.Fatalf("Failed to generate handler: %v", err)
	}
	fmt.Println("✅ Handler generated successfully!")

	// 步骤5: 生成 Migration
	fmt.Println("\n🗃️  Step 5/5: Generating Migration...")
	migrationGen := generator.NewMigrationGenerator()
	migrationConfig := &generator.MigrationConfig{
		Name:   "create_" + generator.ToSnakeCase(*name) + "s_table",
		Table:  *name,
		Action: "create",
	}
	if err := migrationGen.Generate(migrationConfig); err != nil {
		log.Fatalf("Failed to generate migration: %v", err)
	}
	fmt.Println("✅ Migration generated successfully!")

	fmt.Printf("\n🎉 All components for '%s' generated successfully!\n", *name)
	fmt.Println("\nGenerated files:")
	fmt.Printf("  📦 Model:      internal/model/%s.go\n", generator.ToSnakeCase(*name))
	fmt.Printf("  🗄️  Repository: internal/repository/%s.go\n", generator.ToSnakeCase(*name))
	fmt.Printf("  ⚙️  Service:    internal/service/%s.go\n", generator.ToSnakeCase(*name))
	fmt.Printf("  🌐 Handler:    internal/handler/%s.go\n", generator.ToSnakeCase(*name))
	fmt.Printf("  🗃️  Migration:  migrations/{db_type}/{timestamp}_create_%ss_table.sql\n", generator.ToSnakeCase(*name))
	fmt.Println("\n💡 Next steps:")
	fmt.Println("  1. Run 'go build ./...' to verify compilation")
	fmt.Println("  2. Run tests with 'go test ./test/...'")
	fmt.Println("  3. Register routes in your main server file")
}

func handleModuleCommand() {
	fs := flag.NewFlagSet("module", flag.ExitOnError)
	name := fs.String("name", "", "Module name (required)")
	fields := fs.String("fields", "", "Model fields (e.g., 'name:string,price:float64,active:bool')")
	withAuth := fs.Bool("auth", false, "Include authentication middleware")
	withCache := fs.Bool("cache", false, "Include cache support")

	fs.Parse(os.Args[2:])

	if *name == "" {
		fmt.Println("Error: --name is required")
		fs.Usage()
		os.Exit(1)
	}

	gen := generator.NewModuleGenerator()
	config := &generator.ModuleConfig{
		Name:      *name,
		Fields:    *fields,
		WithAuth:  *withAuth,
		WithCache: *withCache,
	}

	if err := gen.Generate(config); err != nil {
		log.Fatalf("Failed to generate module: %v", err)
	}

	fmt.Printf("✅ Module '%s' generated successfully!\n", *name)
}

func handleHandlerCommand() {
	fs := flag.NewFlagSet("handler", flag.ExitOnError)
	model := fs.String("model", "", "Model name (required)")
	withAuth := fs.Bool("auth", false, "Include authentication middleware")
	withValidation := fs.Bool("validation", true, "Include request validation")

	fs.Parse(os.Args[2:])

	if *model == "" {
		fmt.Println("Error: --model is required")
		fs.Usage()
		os.Exit(1)
	}

	gen := generator.NewHandlerGenerator()
	config := &generator.HandlerConfig{
		Model:          *model,
		WithAuth:       *withAuth,
		WithValidation: *withValidation,
	}

	if err := gen.Generate(config); err != nil {
		log.Fatalf("Failed to generate handler: %v", err)
	}

	fmt.Printf("✅ Handler for '%s' generated successfully!\n", *model)
}

func handleTestCommand() {
	fs := flag.NewFlagSet("test", flag.ExitOnError)
	model := fs.String("model", "", "Model name (required)")
	testType := fs.String("type", "unit", "Test type (unit, integration, e2e)")

	fs.Parse(os.Args[2:])

	if *model == "" {
		fmt.Println("Error: --model is required")
		fs.Usage()
		os.Exit(1)
	}

	// 自动生成各组件名称
	serviceName := *model + "Service"
	handlerName := *model + "Handler"
	repositoryName := *model + "Repository"

	gen := generator.NewTestGenerator()
	config := &generator.TestConfig{
		Service:    serviceName,
		Handler:    handlerName,
		Repository: repositoryName,
		Type:       *testType,
	}

	if err := gen.Generate(config); err != nil {
		log.Fatalf("Failed to generate test: %v", err)
	}

	fmt.Printf("✅ Test code for '%s' generated successfully!\n", *model)
}

func handleServiceCommand() {
	fs := flag.NewFlagSet("service", flag.ExitOnError)
	model := fs.String("model", "", "Model name (required)")
	withCache := fs.Bool("cache", false, "Include cache support")

	fs.Parse(os.Args[2:])

	if *model == "" {
		fmt.Println("Error: --model is required")
		fs.Usage()
		os.Exit(1)
	}

	// 自动生成服务名称：Product -> ProductService
	serviceName := *model + "Service"

	gen := generator.NewServiceGenerator()
	config := &generator.ServiceConfig{
		Name:      serviceName,
		Model:     *model,
		Fields:    "", // 使用模型反射，不需要字段字符串
		WithCache: *withCache,
	}

	if err := gen.Generate(config); err != nil {
		log.Fatalf("Failed to generate service: %v", err)
	}

	fmt.Printf("✅ Service '%s' generated successfully!\n", serviceName)
}

func handleRepositoryCommand() {
	fs := flag.NewFlagSet("repository", flag.ExitOnError)
	model := fs.String("model", "", "Model name (required)")

	fs.Parse(os.Args[2:])

	if *model == "" {
		fmt.Println("Error: --model is required")
		fs.Usage()
		os.Exit(1)
	}

	// 自动生成仓储名称：Product -> ProductRepository
	repositoryName := *model + "Repository"

	gen := generator.NewRepositoryGenerator()
	config := &generator.RepositoryConfig{
		Name:  repositoryName,
		Model: *model,
	}

	if err := gen.Generate(config); err != nil {
		log.Fatalf("Failed to generate repository: %v", err)
	}

	fmt.Printf("✅ Repository '%s' generated successfully!\n", repositoryName)
}

func handleModelCommand() {
	fs := flag.NewFlagSet("model", flag.ExitOnError)
	name := fs.String("name", "", "Model name (required)")
	fields := fs.String("fields", "", "Model fields (e.g., 'name:string,price:float64,active:bool')")
	withTimestamps := fs.Bool("timestamps", true, "Include created_at and updated_at fields")

	fs.Parse(os.Args[2:])

	if *name == "" {
		fmt.Println("Error: --name is required")
		fs.Usage()
		os.Exit(1)
	}

	gen := generator.NewModelGenerator()
	config := &generator.ModelConfig{
		Name:           *name,
		Fields:         *fields,
		WithTimestamps: *withTimestamps,
	}

	if err := gen.Generate(config); err != nil {
		log.Fatalf("Failed to generate model: %v", err)
	}

	fmt.Printf("✅ Model '%s' generated successfully!\n", *name)
}

func handleMigrationCommand() {
	fs := flag.NewFlagSet("migration", flag.ExitOnError)
	name := fs.String("name", "", "Migration name (optional if --model is provided)")
	model := fs.String("model", "", "Model name (optional if --name is provided)")
	action := fs.String("action", "create", "Migration action (create, alter, drop)")

	fs.Parse(os.Args[2:])

	var migrationName, tableName string

	if *model != "" {
		// 使用模型名称自动生成迁移名称和表名
		tableName = generator.Pluralize(generator.ToSnakeCase(*model))
		migrationName = fmt.Sprintf("create_%s_table", tableName)
	} else if *name != "" {
		// 使用手动指定的名称
		migrationName = *name
		tableName = *model // 可能为空
	} else {
		fmt.Println("Error: Either --model or --name is required")
		fs.Usage()
		os.Exit(1)
	}

	gen := generator.NewMigrationGenerator()
	config := &generator.MigrationConfig{
		Name:   migrationName,
		Table:  tableName,
		Action: *action,
	}

	if err := gen.Generate(config); err != nil {
		log.Fatalf("Failed to generate migration: %v", err)
	}

	fmt.Printf("✅ Migration '%s' generated successfully!\n", migrationName)
}

func handleFromTableCommand() {
	fs := flag.NewFlagSet("from-table", flag.ExitOnError)
	table := fs.String("table", "", "Table name (required)")
	host := fs.String("host", "localhost", "Database host")
	port := fs.Int("port", 3306, "Database port")
	user := fs.String("user", "root", "Database user")
	password := fs.String("password", "", "Database password")
	database := fs.String("database", "", "Database name (required)")
	modelName := fs.String("model", "", "Model name (optional, auto-generated from table name if not provided)")
	withTimestamps := fs.Bool("timestamps", true, "Include created_at and updated_at fields")
	withSoftDelete := fs.Bool("soft-delete", false, "Include deleted_at field for soft delete")

	fs.Parse(os.Args[2:])

	if *table == "" {
		fmt.Println("Error: --table is required")
		fs.Usage()
		os.Exit(1)
	}

	if *database == "" {
		fmt.Println("Error: --database is required")
		fs.Usage()
		os.Exit(1)
	}

	fmt.Printf("🚀 Generating model from table '%s' in database '%s'...\n\n", *table, *database)

	gen := generator.NewDatabaseTableGenerator()
	config := &generator.DatabaseTableConfig{
		DatabaseHost:     *host,
		DatabasePort:     *port,
		DatabaseUser:     *user,
		DatabasePassword: *password,
		DatabaseName:     *database,
		TableName:        *table,
		ModelName:        *modelName,
		WithTimestamps:   *withTimestamps,
		WithSoftDelete:   *withSoftDelete,
	}

	if err := gen.Generate(config); err != nil {
		log.Fatalf("Failed to generate model from table: %v", err)
	}
}

func handleFromDatabaseCommand() {
	fs := flag.NewFlagSet("from-db", flag.ExitOnError)
	host := fs.String("host", "localhost", "Database host")
	port := fs.Int("port", 3306, "Database port")
	user := fs.String("user", "root", "Database user")
	password := fs.String("password", "", "Database password")
	database := fs.String("database", "", "Database name (required)")
	withTimestamps := fs.Bool("timestamps", true, "Include created_at and updated_at fields")
	withSoftDelete := fs.Bool("soft-delete", false, "Include deleted_at field for soft delete")

	fs.Parse(os.Args[2:])

	if *database == "" {
		fmt.Println("Error: --database is required")
		fs.Usage()
		os.Exit(1)
	}

	fmt.Printf("🚀 Generating models from all tables in database '%s'...\n\n", *database)

	gen := generator.NewDatabaseTableGenerator()
	config := &generator.DatabaseTableConfig{
		DatabaseHost:     *host,
		DatabasePort:     *port,
		DatabaseUser:     *user,
		DatabasePassword: *password,
		DatabaseName:     *database,
		WithTimestamps:   *withTimestamps,
		WithSoftDelete:   *withSoftDelete,
	}

	if err := gen.GenerateFromAllTables(config); err != nil {
		log.Fatalf("Failed to generate models from database: %v", err)
	}
}

func handleListTablesCommand() {
	fs := flag.NewFlagSet("list-tables", flag.ExitOnError)
	host := fs.String("host", "localhost", "Database host")
	port := fs.Int("port", 3306, "Database port")
	user := fs.String("user", "root", "Database user")
	password := fs.String("password", "", "Database password")
	database := fs.String("database", "", "Database name (required)")

	fs.Parse(os.Args[2:])

	if *database == "" {
		fmt.Println("Error: --database is required")
		fs.Usage()
		os.Exit(1)
	}

	// 创建表结构读取器
	reader, err := generator.NewTableReader(*host, *port, *user, *password, *database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer reader.Close()

	// 获取所有表
	tables, err := reader.ListTables()
	if err != nil {
		log.Fatalf("Failed to list tables: %v", err)
	}

	fmt.Printf("📊 Found %d tables in database '%s':\n\n", len(tables), *database)
	for i, table := range tables {
		fmt.Printf("  %d. %s\n", i+1, table)
	}
	fmt.Println()
}

// convertGoTypeToFieldType 将Go类型转换为字段类型字符串
func convertGoTypeToFieldType(goType string) string {
	switch {
	case strings.Contains(goType, "sql.NullString"), goType == "string":
		return "string"
	case strings.Contains(goType, "sql.NullInt32"), goType == "int32", goType == "int":
		return "int"
	case strings.Contains(goType, "sql.NullInt64"), goType == "int64":
		return "int64"
	case strings.Contains(goType, "sql.NullFloat64"), goType == "float64":
		return "float64"
	case strings.Contains(goType, "sql.NullBool"), goType == "bool":
		return "bool"
	case strings.Contains(goType, "sql.NullTime"), goType == "time.Time":
		return "time"
	case goType == "int8":
		return "int"
	case goType == "int16":
		return "int"
	case goType == "uint32", goType == "uint":
		return "int"
	case goType == "uint64":
		return "int64"
	case goType == "float32":
		return "float64"
	default:
		return "string"
	}
}

func handleFrontendCommand() {
	fs := flag.NewFlagSet("frontend", flag.ExitOnError)
	model := fs.String("model", "", "Model name (required)")
	framework := fs.String("framework", "antd", "Frontend framework (antd, vue)")
	output := fs.String("output", "", "Output directory (required)")
	moduleType := fs.String("module-type", "admin", "Module type (admin, public)")
	withAuth := fs.Bool("with-auth", false, "Include authentication")
	withSearch := fs.Bool("with-search", true, "Include search functionality")
	withExport := fs.Bool("with-export", false, "Include export functionality")
	withBatch := fs.Bool("with-batch", false, "Include batch operations")
	apiPrefix := fs.String("api-prefix", "/api/v1", "API prefix")
	moduleName := fs.String("module", "", "Module name (default: lowercase model name)")

	fs.Parse(os.Args[2:])

	if *model == "" {
		fmt.Println("Error: --model is required")
		fs.Usage()
		os.Exit(1)
	}

	if *output == "" {
		fmt.Println("Error: --output is required")
		fmt.Println("Please specify the frontend project directory, e.g., --output=../vibe-coding-starter-ui-antd")
		fs.Usage()
		os.Exit(1)
	}

	// 验证框架类型
	var fwType generator.FrontendFramework
	switch *framework {
	case "antd":
		fwType = generator.FrameworkAntd
	case "vue":
		fwType = generator.FrameworkVue
	default:
		fmt.Printf("Error: unsupported framework '%s'. Supported: antd, vue\n", *framework)
		os.Exit(1)
	}

	// 验证模块类型
	var modType generator.ModuleType
	switch *moduleType {
	case "admin":
		modType = generator.ModuleTypeAdmin
	case "public":
		modType = generator.ModuleTypePublic
	default:
		fmt.Printf("Error: unsupported module type '%s'. Supported: admin, public\n", *moduleType)
		os.Exit(1)
	}

	// 设置默认模块名
	if *moduleName == "" {
		*moduleName = strings.ToLower(*model)
	}

	// 根据模块类型设置默认 API 前缀
	if *apiPrefix == "/api/v1" && modType == generator.ModuleTypeAdmin {
		*apiPrefix = "/api/v1/admin"
	}

	fmt.Printf("🚀 Generating %s frontend code for '%s'...\n\n", *framework, *model)

	gen := generator.NewFrontendGenerator()
	config := &generator.FrontendConfig{
		Model:      *model,
		Framework:  fwType,
		OutputDir:  *output,
		ModuleType: modType,
		WithAuth:   *withAuth,
		WithSearch: *withSearch,
		WithExport: *withExport,
		WithBatch:  *withBatch,
		ApiPrefix:  *apiPrefix,
		ModuleName: *moduleName,
	}

	if err := gen.Generate(config); err != nil {
		log.Fatalf("Failed to generate frontend code: %v", err)
	}

	fmt.Printf("✅ Frontend code for '%s' generated successfully!\n", *model)
	fmt.Printf("📁 Output directory: %s\n", *output)
	fmt.Printf("🎨 Framework: %s\n", *framework)
	fmt.Printf("🏗️  Module type: %s\n", *moduleType)
	fmt.Printf("📦 Module: %s\n", *moduleName)
}
