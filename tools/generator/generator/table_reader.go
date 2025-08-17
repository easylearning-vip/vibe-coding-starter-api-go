package generator

import (
	"fmt"
	"strings"
)

// DatabaseField 数据库字段信息（转换为生成器字段格式）
type DatabaseField struct {
	Name     string      // Go字段名，如 CustomerName
	Type     string      // Go字段类型，如 string
	JSONName string      // JSON 名称，如 customer_name
	GormTag  string      // GORM 标签
	Comment  string      // 注释
	Column   *ColumnInfo // 原始列信息
}

// TableReader 表结构读取器
type TableReader struct {
	analyzer *MySQLAnalyzer
	mapper   *MySQLTypeMapper
}

// NewTableReader 创建表结构读取器
func NewTableReader(host string, port int, user, password, database string) (*TableReader, error) {
	analyzer, err := NewMySQLAnalyzer(host, port, user, password, database)
	if err != nil {
		return nil, fmt.Errorf("failed to create MySQL analyzer: %w", err)
	}

	return &TableReader{
		analyzer: analyzer,
		mapper:   NewMySQLTypeMapper(),
	}, nil
}

// Close 关闭连接
func (r *TableReader) Close() error {
	return r.analyzer.Close()
}

// ReadTableStructure 读取表结构并转换为生成器字段格式
func (r *TableReader) ReadTableStructure(tableName string) ([]*DatabaseField, error) {
	// 获取表信息
	tableInfo, err := r.analyzer.GetTableInfo(tableName)
	if err != nil {
		return nil, fmt.Errorf("failed to get table info: %w", err)
	}

	var fields []*DatabaseField
	for _, column := range tableInfo.Columns {
		// 跳过一些系统字段
		if r.shouldSkipColumn(column) {
			continue
		}

		field := r.convertColumnToField(column)
		fields = append(fields, field)
	}

	return fields, nil
}

// shouldSkipColumn 判断是否应该跳过某个列
func (r *TableReader) shouldSkipColumn(column *ColumnInfo) bool {
	// 跳过主键ID字段（通常由BaseModel提供）
	if column.IsPrimary && strings.ToLower(column.Name) == "id" {
		return true
	}

	// 跳过时间戳字段（通常由BaseModel提供）
	columnName := strings.ToLower(column.Name)
	if columnName == "created_at" || columnName == "updated_at" || columnName == "deleted_at" {
		return true
	}

	return false
}

// convertColumnToField 将数据库列转换为生成器字段
func (r *TableReader) convertColumnToField(column *ColumnInfo) *DatabaseField {
	// 生成Go字段名
	fieldName := ToPascalCase(column.Name)

	// 获取Go类型
	goType := r.mapper.MapToGoType(column)

	// 生成JSON名称
	jsonName := ToSnakeCase(fieldName)

	// 生成GORM标签
	gormTag := r.mapper.GenerateGormTag(column)

	// 生成注释
	comment := r.generateComment(column, goType)

	return &DatabaseField{
		Name:     fieldName,
		Type:     goType,
		JSONName: jsonName,
		GormTag:  gormTag,
		Comment:  comment,
		Column:   column,
	}
}

// generateComment 生成字段注释
func (r *TableReader) generateComment(column *ColumnInfo, goType string) string {
	var parts []string

	// 使用数据库注释作为主要描述
	if column.Comment != "" {
		parts = append(parts, column.Comment)
	} else {
		// 如果没有注释，使用字段名
		parts = append(parts, ToPascalCase(column.Name))
	}

	// 添加类型信息
	typeComment := r.getTypeComment(goType)
	if typeComment != "" {
		parts = append(parts, typeComment)
	}

	return strings.Join(parts, " ")
}

// getTypeComment 获取类型注释
func (r *TableReader) getTypeComment(goType string) string {
	switch {
	case strings.HasPrefix(goType, "sql.Null"):
		baseType := strings.TrimPrefix(goType, "sql.Null")
		return fmt.Sprintf("可空%s", r.getBaseTypeComment(baseType))
	case strings.HasPrefix(goType, "*"):
		baseType := strings.TrimPrefix(goType, "*")
		return fmt.Sprintf("可空%s", r.getBaseTypeComment(baseType))
	default:
		return r.getBaseTypeComment(goType)
	}
}

// getBaseTypeComment 获取基础类型注释
func (r *TableReader) getBaseTypeComment(goType string) string {
	switch goType {
	case "string", "String":
		return "字符串"
	case "int", "int32":
		return "32位整数"
	case "int64":
		return "64位整数"
	case "int8":
		return "8位整数"
	case "int16":
		return "16位整数"
	case "uint", "uint32":
		return "32位无符号整数"
	case "uint64":
		return "64位无符号整数"
	case "uint8":
		return "8位无符号整数"
	case "uint16":
		return "16位无符号整数"
	case "float32":
		return "32位浮点数"
	case "float64":
		return "64位浮点数"
	case "bool", "Bool":
		return "布尔值"
	case "time.Time", "Time":
		return "时间类型"
	case "[]byte":
		return "字节数组"
	default:
		return "自定义类型"
	}
}

// GetRequiredImports 获取字段类型需要的导入包
func (r *TableReader) GetRequiredImports(fields []*DatabaseField) []string {
	imports := make(map[string]bool)

	for _, field := range fields {
		// 检查是否需要time包
		if strings.Contains(field.Type, "time.Time") {
			imports["time"] = true
		}

		// 检查是否需要database/sql包
		if strings.HasPrefix(field.Type, "sql.Null") {
			imports["database/sql"] = true
		}
	}

	var result []string
	for imp := range imports {
		result = append(result, imp)
	}

	return result
}

// ListTables 列出数据库中的所有表
func (r *TableReader) ListTables() ([]string, error) {
	return r.analyzer.ListTables()
}

// GetTableComment 获取表注释
func (r *TableReader) GetTableComment(tableName string) (string, error) {
	tableInfo, err := r.analyzer.GetTableInfo(tableName)
	if err != nil {
		return "", err
	}
	return tableInfo.Comment, nil
}
