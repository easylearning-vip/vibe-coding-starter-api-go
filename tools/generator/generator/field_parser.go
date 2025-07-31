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
			Type:     fieldType,
			JSONName: ToSnakeCase(fieldName),
			GormTag:  p.generateGormTag(fieldName, fieldType),
			Comment:  fmt.Sprintf("%s %s", ToPascalCase(fieldName), p.getTypeComment(fieldType)),
		}

		fields = append(fields, field)
	}

	return fields, nil
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
