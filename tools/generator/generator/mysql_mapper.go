package generator

import (
	"fmt"
	"strings"
)

// MySQLTypeMapper MySQL类型映射器
type MySQLTypeMapper struct{}

// NewMySQLTypeMapper 创建MySQL类型映射器
func NewMySQLTypeMapper() *MySQLTypeMapper {
	return &MySQLTypeMapper{}
}

// MapToGoType 将MySQL类型映射为Go类型
func (m *MySQLTypeMapper) MapToGoType(column *ColumnInfo) string {
	mysqlType := strings.ToLower(column.Type)

	// 处理可空类型
	goType := m.getBaseGoType(mysqlType, column)

	// 如果字段可为空且不是指针类型，则使用指针类型
	if column.IsNullable && !column.IsPrimary && !strings.HasPrefix(goType, "*") {
		// 对于某些类型，使用sql包的Null类型更合适
		switch goType {
		case "string":
			return "sql.NullString"
		case "int", "int32":
			return "sql.NullInt32"
		case "int64":
			return "sql.NullInt64"
		case "float64":
			return "sql.NullFloat64"
		case "bool":
			return "sql.NullBool"
		case "time.Time":
			return "sql.NullTime"
		default:
			return "*" + goType
		}
	}

	return goType
}

// getBaseGoType 获取基础Go类型
func (m *MySQLTypeMapper) getBaseGoType(mysqlType string, column *ColumnInfo) string {
	switch {
	// 字符串类型
	case strings.Contains(mysqlType, "char"), strings.Contains(mysqlType, "text"):
		return "string"
	case mysqlType == "json":
		return "string" // 或者可以使用 json.RawMessage

	// 整数类型
	case mysqlType == "tinyint":
		// 检查是否为tinyint(1)，通常用作布尔类型
		if strings.Contains(strings.ToLower(column.FullType), "tinyint(1)") {
			return "bool"
		}
		return "int8"
	case mysqlType == "smallint":
		return "int16"
	case mysqlType == "mediumint", mysqlType == "int":
		return "int32"
	case mysqlType == "bigint":
		return "int64"

	// 无符号整数类型
	case strings.Contains(mysqlType, "unsigned"):
		if strings.Contains(mysqlType, "tinyint") {
			return "uint8"
		} else if strings.Contains(mysqlType, "smallint") {
			return "uint16"
		} else if strings.Contains(mysqlType, "mediumint") || strings.Contains(mysqlType, "int") {
			return "uint32"
		} else if strings.Contains(mysqlType, "bigint") {
			return "uint64"
		}
		return "uint32"

	// 浮点数类型
	case mysqlType == "float":
		return "float32"
	case mysqlType == "double", strings.Contains(mysqlType, "decimal"):
		return "float64"

	// 布尔类型
	case mysqlType == "boolean", mysqlType == "bool":
		return "bool"

	// 时间类型
	case mysqlType == "date", mysqlType == "datetime", mysqlType == "timestamp":
		return "time.Time"
	case mysqlType == "time":
		return "string" // 时间类型可以用string表示

	// 二进制类型
	case strings.Contains(mysqlType, "binary"), strings.Contains(mysqlType, "blob"):
		return "[]byte"

	default:
		return "string" // 默认使用string类型
	}
}

// GenerateGormTag 生成GORM标签
func (m *MySQLTypeMapper) GenerateGormTag(column *ColumnInfo) string {
	var tags []string

	// 列名
	tags = append(tags, fmt.Sprintf("column:%s", column.Name))

	// 类型
	typeTag := m.generateTypeTag(column)
	if typeTag != "" {
		tags = append(tags, fmt.Sprintf("type:%s", typeTag))
	}

	// 主键
	if column.IsPrimary {
		tags = append(tags, "primaryKey")
	}

	// 自增
	if column.IsAutoIncr {
		tags = append(tags, "autoIncrement")
	}

	// 非空
	if !column.IsNullable && !column.IsPrimary {
		tags = append(tags, "not null")
	}

	// 默认值
	if column.DefaultValue != nil && *column.DefaultValue != "" {
		defaultVal := *column.DefaultValue
		// 对于字符串默认值，需要加引号
		if m.isStringType(column.Type) && !strings.HasPrefix(defaultVal, "'") {
			defaultVal = fmt.Sprintf("'%s'", defaultVal)
		}
		tags = append(tags, fmt.Sprintf("default:%s", defaultVal))
	}

	// 注释
	if column.Comment != "" {
		tags = append(tags, fmt.Sprintf("comment:%s", column.Comment))
	}

	return fmt.Sprintf("gorm:\"%s\"", strings.Join(tags, ";"))
}

// generateTypeTag 生成类型标签
func (m *MySQLTypeMapper) generateTypeTag(column *ColumnInfo) string {
	mysqlType := strings.ToLower(column.Type)

	switch {
	case strings.Contains(mysqlType, "varchar"):
		if column.MaxLength != nil {
			return fmt.Sprintf("varchar(%d)", *column.MaxLength)
		}
		return "varchar(255)"

	case strings.Contains(mysqlType, "char"):
		if column.MaxLength != nil {
			return fmt.Sprintf("char(%d)", *column.MaxLength)
		}
		return "char(255)"

	case mysqlType == "text":
		return "text"

	case mysqlType == "longtext":
		return "longtext"

	case mysqlType == "mediumtext":
		return "mediumtext"

	case strings.Contains(mysqlType, "decimal"):
		if column.MaxLength != nil && column.NumericScale != nil {
			return fmt.Sprintf("decimal(%d,%d)", *column.MaxLength, *column.NumericScale)
		}
		return "decimal(10,2)"

	case mysqlType == "int":
		return "int"

	case mysqlType == "bigint":
		return "bigint"

	case mysqlType == "tinyint":
		return "tinyint"

	case mysqlType == "smallint":
		return "smallint"

	case mysqlType == "float":
		return "float"

	case mysqlType == "double":
		return "double"

	case mysqlType == "datetime":
		return "datetime"

	case mysqlType == "timestamp":
		return "timestamp"

	case mysqlType == "date":
		return "date"

	case mysqlType == "time":
		return "time"

	case mysqlType == "json":
		return "json"

	case mysqlType == "boolean", mysqlType == "bool":
		return "boolean"

	default:
		return ""
	}
}

// isStringType 判断是否为字符串类型
func (m *MySQLTypeMapper) isStringType(mysqlType string) bool {
	mysqlType = strings.ToLower(mysqlType)
	return strings.Contains(mysqlType, "char") ||
		strings.Contains(mysqlType, "text") ||
		mysqlType == "json"
}

// GetRequiredImports 获取需要的导入包
func (m *MySQLTypeMapper) GetRequiredImports(columns []*ColumnInfo) []string {
	imports := make(map[string]bool)

	for _, column := range columns {
		goType := m.MapToGoType(column)

		// 检查是否需要time包
		if strings.Contains(goType, "time.Time") {
			imports["time"] = true
		}

		// 检查是否需要database/sql包
		if strings.HasPrefix(goType, "sql.Null") {
			imports["database/sql"] = true
		}
	}

	var result []string
	for imp := range imports {
		result = append(result, imp)
	}

	return result
}
