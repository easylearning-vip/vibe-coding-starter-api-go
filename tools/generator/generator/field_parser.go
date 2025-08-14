package generator

import (
	"fmt"
	"strings"
)

// Field 表示模型字段
type Field struct {
	Name     string // 字段名，如 CustomerName
	Type     string // 字段类型，如 string
	JSONName string // JSON 名称，如 customer_name
	GormTag  string // GORM 标签
	Comment  string // 注释
}

// FieldParser 字段解析器
type FieldParser struct{}

// NewFieldParser 创建字段解析器
func NewFieldParser() *FieldParser {
	return &FieldParser{}
}

// ParseFields 解析字段定义字符串
// 输入格式: "name:string,description:string,price:float64,active:bool"
func (p *FieldParser) ParseFields(fieldsStr string) ([]*Field, error) {
	if fieldsStr == "" {
		return []*Field{}, nil
	}

	var fields []*Field
	fieldPairs := strings.Split(fieldsStr, ",")

	for _, pair := range fieldPairs {
		pair = strings.TrimSpace(pair)
		if pair == "" {
			continue
		}

		parts := strings.Split(pair, ":")
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid field format: %s, expected format: name:type", pair)
		}

		fieldName := strings.TrimSpace(parts[0])
		fieldType := strings.TrimSpace(parts[1])

		field := &Field{
			Name:     ToPascalCase(fieldName),
			Type:     p.mapToGoType(fieldType),
			JSONName: ToSnakeCase(fieldName),
			GormTag:  p.generateGormTag(fieldName, fieldType),
			Comment:  fmt.Sprintf("%s %s", ToPascalCase(fieldName), p.getTypeComment(fieldType)),
		}

		fields = append(fields, field)
	}

	return fields, nil
}

// mapToGoType 将字段类型映射为 Go 类型
func (p *FieldParser) mapToGoType(fieldType string) string {
	switch fieldType {
	case "text":
		return "string"
	case "varchar":
		return "string"
	case "char":
		return "string"
	case "int":
		return "int"
	case "integer":
		return "int"
	case "bigint":
		return "int64"
	case "smallint":
		return "int16"
	case "tinyint":
		return "int8"
	case "float":
		return "float32"
	case "double":
		return "float64"
	case "decimal":
		return "float64"
	case "bool":
		return "bool"
	case "boolean":
		return "bool"
	case "datetime":
		return "time.Time"
	case "timestamp":
		return "time.Time"
	case "date":
		return "time.Time"
	case "time":
		return "time.Time"
	default:
		return fieldType // 保持原类型
	}
}

// generateGormTag 生成 GORM 标签
func (p *FieldParser) generateGormTag(fieldName, fieldType string) string {
	columnName := ToSnakeCase(fieldName)

	var gormType string
	switch fieldType {
	case "string":
		gormType = "varchar(255)"
	case "int", "int32":
		gormType = "int"
	case "int64":
		gormType = "bigint"
	case "uint", "uint32":
		gormType = "int unsigned"
	case "uint64":
		gormType = "bigint unsigned"
	case "float32":
		gormType = "float"
	case "float64":
		gormType = "decimal(10,2)"
	case "bool":
		gormType = "boolean;default:false"
	case "time.Time":
		gormType = "datetime"
	default:
		gormType = "text"
	}

	// 为Name字段添加唯一约束
	if strings.ToLower(fieldName) == "name" && fieldType == "string" {
		return fmt.Sprintf("gorm:\"column:%s;type:%s;uniqueIndex\"", columnName, gormType)
	}

	return fmt.Sprintf("gorm:\"column:%s;type:%s\"", columnName, gormType)
}

// getTypeComment 获取类型注释
func (p *FieldParser) getTypeComment(fieldType string) string {
	switch fieldType {
	case "string":
		return "字符串"
	case "int", "int32":
		return "32位整数"
	case "int64":
		return "64位整数"
	case "uint", "uint32":
		return "32位无符号整数"
	case "uint64":
		return "64位无符号整数"
	case "float32":
		return "32位浮点数"
	case "float64":
		return "64位浮点数"
	case "bool":
		return "布尔值"
	case "time.Time":
		return "时间类型"
	default:
		return "自定义类型"
	}
}

// GetRequiredImports 获取字段类型需要的导入包
func (p *FieldParser) GetRequiredImports(fields []*Field) []string {
	imports := make(map[string]bool)

	for _, field := range fields {
		switch field.Type {
		case "time.Time":
			imports["time"] = true
		}
	}

	var result []string
	for imp := range imports {
		result = append(result, imp)
	}

	return result
}

// GenerateSQLColumn 生成SQL列定义
func (p *FieldParser) GenerateSQLColumn(field *Field, databaseType string) string {
	columnName := ToSnakeCase(field.Name)

	var sqlType string
	var defaultValue string

	switch field.Type {
	case "string":
		sqlType = "VARCHAR(255)"
		defaultValue = ""
	case "int", "int32":
		if databaseType == "postgres" {
			sqlType = "INTEGER"
		} else {
			sqlType = "INT"
		}
		defaultValue = ""
	case "int64":
		if databaseType == "postgres" {
			sqlType = "BIGINT"
		} else {
			sqlType = "BIGINT"
		}
		defaultValue = ""
	case "uint", "uint32":
		if databaseType == "postgres" {
			sqlType = "INTEGER"
		} else {
			sqlType = "INT UNSIGNED"
		}
		defaultValue = ""
	case "uint64":
		if databaseType == "postgres" {
			sqlType = "BIGINT"
		} else {
			sqlType = "BIGINT UNSIGNED"
		}
		defaultValue = ""
	case "float32":
		sqlType = "FLOAT"
		defaultValue = ""
	case "float64":
		sqlType = "DECIMAL(10,2)"
		defaultValue = ""
	case "bool":
		if databaseType == "postgres" {
			sqlType = "BOOLEAN"
		} else {
			sqlType = "BOOLEAN"
		}
		defaultValue = " DEFAULT FALSE"
	case "time.Time":
		if databaseType == "postgres" {
			sqlType = "TIMESTAMP WITH TIME ZONE"
		} else {
			sqlType = "DATETIME"
		}
		defaultValue = ""
	default:
		sqlType = "TEXT"
		defaultValue = ""
	}

	// 构建完整的列定义
	columnDef := fmt.Sprintf("    %s %s", columnName, sqlType)

	// 添加NOT NULL约束（除了某些特殊情况）
	if field.Type != "bool" && !strings.Contains(strings.ToLower(field.Name), "optional") {
		columnDef += " NOT NULL"
	}

	// 添加默认值
	if defaultValue != "" {
		columnDef += defaultValue
	}

	return columnDef
}

// GenerateCreateRequestFields 生成创建请求的字段
func (p *FieldParser) GenerateCreateRequestFields(fields []*Field) string {
	var lines []string

	for _, field := range fields {
		// 跳过时间字段，通常由系统自动设置
		if field.Type == "time.Time" {
			continue
		}

		validateTag := p.generateValidateTag(field)
		line := fmt.Sprintf("\t%s %s `json:\"%s\"%s`",
			field.Name, field.Type, field.JSONName, validateTag)
		lines = append(lines, line)
	}

	return strings.Join(lines, "\n")
}

// GenerateUpdateRequestFields 生成更新请求的字段
func (p *FieldParser) GenerateUpdateRequestFields(fields []*Field) string {
	var lines []string

	for _, field := range fields {
		// 跳过时间字段，通常由系统自动设置
		if field.Type == "time.Time" {
			continue
		}

		validateTag := p.generateValidateTag(field)
		if validateTag != "" {
			validateTag = strings.Replace(validateTag, "required,", "omitempty,", 1)
			validateTag = strings.Replace(validateTag, "required", "omitempty", 1)
		} else {
			validateTag = ` validate:"omitempty"`
		}

		line := fmt.Sprintf("\t%s *%s `json:\"%s,omitempty\"%s`",
			field.Name, field.Type, field.JSONName, validateTag)
		lines = append(lines, line)
	}

	return strings.Join(lines, "\n")
}

// generateValidateTag 生成验证标签
func (p *FieldParser) generateValidateTag(field *Field) string {
	switch field.Type {
	case "string":
		return ` validate:"required,min=1,max=255"`
	case "int", "int32", "int64", "uint", "uint32", "uint64":
		return ` validate:"required,min=0"`
	case "float32", "float64":
		return ` validate:"required,min=0"`
	case "bool":
		return ""
	default:
		return ` validate:"required"`
	}
}

// GenerateModelAssignment 生成模型赋值代码
func (p *FieldParser) GenerateModelAssignment(fields []*Field, varName string) string {
	var lines []string

	for _, field := range fields {
		// 跳过时间字段，通常由系统自动设置
		if field.Type == "time.Time" {
			continue
		}

		line := fmt.Sprintf("\t\t%s: %s.%s,", field.Name, varName, field.Name)
		lines = append(lines, line)
	}

	return strings.Join(lines, "\n")
}

// GenerateUpdateAssignment 生成更新赋值代码
func (p *FieldParser) GenerateUpdateAssignment(fields []*Field, entityVar, reqVar string) string {
	var lines []string

	for _, field := range fields {
		// 跳过时间字段，通常由系统自动设置
		if field.Type == "time.Time" {
			continue
		}

		lines = append(lines, fmt.Sprintf("\tif %s.%s != nil {", reqVar, field.Name))
		lines = append(lines, fmt.Sprintf("\t\t%s.%s = *%s.%s", entityVar, field.Name, reqVar, field.Name))
		lines = append(lines, "\t}")
	}

	return strings.Join(lines, "\n")
}
