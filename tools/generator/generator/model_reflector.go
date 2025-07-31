package generator

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"strings"
)

// ModelReflector 模型反射器
type ModelReflector struct {
	fieldParser *FieldParser
}

// NewModelReflector 创建模型反射器
func NewModelReflector() *ModelReflector {
	return &ModelReflector{
		fieldParser: NewFieldParser(),
	}
}

// ReflectModelFields 通过反射获取模型字段
func (r *ModelReflector) ReflectModelFields(modelName string) ([]*Field, error) {
	// 构建模型文件路径
	modelFile := filepath.Join("internal", "model", ToSnakeCase(modelName)+".go")
	
	// 解析 Go 文件
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, modelFile, nil, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("failed to parse model file %s: %w", modelFile, err)
	}

	// 查找模型结构体
	var fields []*Field
	ast.Inspect(node, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.TypeSpec:
			// 检查是否是目标模型
			if x.Name.Name == modelName {
				if structType, ok := x.Type.(*ast.StructType); ok {
					fields = r.extractFieldsFromStruct(structType)
					return false // 找到了，停止遍历
				}
			}
		}
		return true
	})

	if len(fields) == 0 {
		return nil, fmt.Errorf("model %s not found or has no fields", modelName)
	}

	return fields, nil
}

// extractFieldsFromStruct 从结构体中提取字段
func (r *ModelReflector) extractFieldsFromStruct(structType *ast.StructType) []*Field {
	var fields []*Field

	for _, field := range structType.Fields.List {
		// 跳过嵌入字段（如 BaseModel）
		if len(field.Names) == 0 {
			continue
		}

		for _, name := range field.Names {
			// 跳过私有字段
			if !name.IsExported() {
				continue
			}

			fieldName := name.Name
			fieldType := r.extractTypeFromExpr(field.Type)
			
			// 跳过某些系统字段
			if r.shouldSkipField(fieldName, fieldType) {
				continue
			}

			// 解析标签获取 JSON 名称
			jsonName := r.extractJSONNameFromTag(field.Tag)
			if jsonName == "" {
				jsonName = ToSnakeCase(fieldName)
			}

			// 创建字段对象
			f := &Field{
				Name:     fieldName,
				Type:     fieldType,
				JSONName: jsonName,
				GormTag:  r.fieldParser.generateGormTag(fieldName, fieldType),
				Comment:  fmt.Sprintf("%s %s", fieldName, r.fieldParser.getTypeComment(fieldType)),
			}

			fields = append(fields, f)
		}
	}

	return fields
}

// extractTypeFromExpr 从表达式中提取类型名称
func (r *ModelReflector) extractTypeFromExpr(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.SelectorExpr:
		// 处理 time.Time 等包类型
		if pkg, ok := t.X.(*ast.Ident); ok {
			return pkg.Name + "." + t.Sel.Name
		}
		return t.Sel.Name
	case *ast.StarExpr:
		// 处理指针类型
		return "*" + r.extractTypeFromExpr(t.X)
	case *ast.ArrayType:
		// 处理数组类型
		return "[]" + r.extractTypeFromExpr(t.Elt)
	case *ast.MapType:
		// 处理 map 类型
		keyType := r.extractTypeFromExpr(t.Key)
		valueType := r.extractTypeFromExpr(t.Value)
		return fmt.Sprintf("map[%s]%s", keyType, valueType)
	default:
		return "interface{}"
	}
}

// extractJSONNameFromTag 从标签中提取 JSON 名称
func (r *ModelReflector) extractJSONNameFromTag(tag *ast.BasicLit) string {
	if tag == nil {
		return ""
	}

	// 移除引号
	tagValue := strings.Trim(tag.Value, "`")
	
	// 查找 json 标签
	parts := strings.Fields(tagValue)
	for _, part := range parts {
		if strings.HasPrefix(part, "json:") {
			// 提取 json 标签值
			jsonTag := strings.TrimPrefix(part, "json:")
			jsonTag = strings.Trim(jsonTag, `"`)
			
			// 处理 json 标签选项（如 omitempty）
			if idx := strings.Index(jsonTag, ","); idx != -1 {
				jsonTag = jsonTag[:idx]
			}
			
			return jsonTag
		}
	}

	return ""
}

// shouldSkipField 判断是否应该跳过字段
func (r *ModelReflector) shouldSkipField(fieldName, fieldType string) bool {
	// 跳过时间字段，通常由系统自动管理
	if fieldType == "time.Time" {
		return true
	}

	// 跳过某些系统字段
	systemFields := map[string]bool{
		"ID":        true,
		"CreatedAt": true,
		"UpdatedAt": true,
		"DeletedAt": true,
	}

	return systemFields[fieldName]
}

// GetFieldsFromModelOrFallback 优先从模型反射获取字段，失败则使用字段字符串
func (r *ModelReflector) GetFieldsFromModelOrFallback(modelName, fieldsStr string) ([]*Field, error) {
	// 首先尝试从模型文件反射获取字段
	if modelName != "" {
		fields, err := r.ReflectModelFields(modelName)
		if err == nil && len(fields) > 0 {
			return fields, nil
		}
		// 如果反射失败，记录但不报错，继续使用字段字符串
		fmt.Printf("Warning: Failed to reflect model %s, falling back to fields string: %v\n", modelName, err)
	}

	// 回退到使用字段字符串解析
	if fieldsStr != "" {
		return r.fieldParser.ParseFields(fieldsStr)
	}

	// 如果都没有，返回空字段列表
	return []*Field{}, nil
}
