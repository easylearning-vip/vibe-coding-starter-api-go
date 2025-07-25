package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"vibe-coding-starter/tools/generator/cmd"
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
	fmt.Println("Usage: go run tools/generator/main.go <command> [options]")
	fmt.Println()
	fmt.Println("Available commands:")
	fmt.Println("  module      Generate a complete business module")
	fmt.Println("  handler     Generate API handler")
	fmt.Println("  service     Generate service layer")
	fmt.Println("  repository  Generate repository layer")
	fmt.Println("  model       Generate data model")
	fmt.Println("  test        Generate test code")
	fmt.Println("  migration   Generate database migration")
	fmt.Println("  help        Show this help message")
	fmt.Println("  version     Show version information")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  go run tools/generator/main.go module --name=product")
	fmt.Println("  go run tools/generator/main.go handler --model=Product")
	fmt.Println("  go run tools/generator/main.go test --service=ProductService")
	fmt.Println("  go run tools/generator/main.go migration --name=create_products_table")
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

	generator := cmd.NewModuleGenerator()
	config := &cmd.ModuleConfig{
		Name:      *name,
		Fields:    *fields,
		WithAuth:  *withAuth,
		WithCache: *withCache,
	}

	if err := generator.Generate(config); err != nil {
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

	generator := cmd.NewHandlerGenerator()
	config := &cmd.HandlerConfig{
		Model:          *model,
		WithAuth:       *withAuth,
		WithValidation: *withValidation,
	}

	if err := generator.Generate(config); err != nil {
		log.Fatalf("Failed to generate handler: %v", err)
	}

	fmt.Printf("✅ Handler for '%s' generated successfully!\n", *model)
}

func handleTestCommand() {
	fs := flag.NewFlagSet("test", flag.ExitOnError)
	service := fs.String("service", "", "Service name")
	handler := fs.String("handler", "", "Handler name")
	repository := fs.String("repository", "", "Repository name")
	testType := fs.String("type", "unit", "Test type (unit, integration, e2e)")

	fs.Parse(os.Args[2:])

	if *service == "" && *handler == "" && *repository == "" {
		fmt.Println("Error: One of --service, --handler, or --repository is required")
		fs.Usage()
		os.Exit(1)
	}

	generator := cmd.NewTestGenerator()
	config := &cmd.TestConfig{
		Service:    *service,
		Handler:    *handler,
		Repository: *repository,
		Type:       *testType,
	}

	if err := generator.Generate(config); err != nil {
		log.Fatalf("Failed to generate test: %v", err)
	}

	fmt.Println("✅ Test code generated successfully!")
}

func handleServiceCommand() {
	fs := flag.NewFlagSet("service", flag.ExitOnError)
	name := fs.String("name", "", "Service name (required)")
	model := fs.String("model", "", "Associated model name")
	withCache := fs.Bool("cache", false, "Include cache support")

	fs.Parse(os.Args[2:])

	if *name == "" {
		fmt.Println("Error: --name is required")
		fs.Usage()
		os.Exit(1)
	}

	generator := cmd.NewServiceGenerator()
	config := &cmd.ServiceConfig{
		Name:      *name,
		Model:     *model,
		WithCache: *withCache,
	}

	if err := generator.Generate(config); err != nil {
		log.Fatalf("Failed to generate service: %v", err)
	}

	fmt.Printf("✅ Service '%s' generated successfully!\n", *name)
}

func handleRepositoryCommand() {
	fs := flag.NewFlagSet("repository", flag.ExitOnError)
	name := fs.String("name", "", "Repository name (required)")
	model := fs.String("model", "", "Associated model name")

	fs.Parse(os.Args[2:])

	if *name == "" {
		fmt.Println("Error: --name is required")
		fs.Usage()
		os.Exit(1)
	}

	generator := cmd.NewRepositoryGenerator()
	config := &cmd.RepositoryConfig{
		Name:  *name,
		Model: *model,
	}

	if err := generator.Generate(config); err != nil {
		log.Fatalf("Failed to generate repository: %v", err)
	}

	fmt.Printf("✅ Repository '%s' generated successfully!\n", *name)
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

	generator := cmd.NewModelGenerator()
	config := &cmd.ModelConfig{
		Name:           *name,
		Fields:         *fields,
		WithTimestamps: *withTimestamps,
	}

	if err := generator.Generate(config); err != nil {
		log.Fatalf("Failed to generate model: %v", err)
	}

	fmt.Printf("✅ Model '%s' generated successfully!\n", *name)
}

func handleMigrationCommand() {
	fs := flag.NewFlagSet("migration", flag.ExitOnError)
	name := fs.String("name", "", "Migration name (required)")
	table := fs.String("table", "", "Table name")
	action := fs.String("action", "create", "Migration action (create, alter, drop)")

	fs.Parse(os.Args[2:])

	if *name == "" {
		fmt.Println("Error: --name is required")
		fs.Usage()
		os.Exit(1)
	}

	generator := cmd.NewMigrationGenerator()
	config := &cmd.MigrationConfig{
		Name:   *name,
		Table:  *table,
		Action: *action,
	}

	if err := generator.Generate(config); err != nil {
		log.Fatalf("Failed to generate migration: %v", err)
	}

	fmt.Printf("✅ Migration '%s' generated successfully!\n", *name)
}
