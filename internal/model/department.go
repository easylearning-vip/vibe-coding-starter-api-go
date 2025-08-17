package model

import (
	"fmt"
	"gorm.io/gorm"
)

// Department Department模型
type Department struct {
	BaseModel

	Name string `json:"name" gorm:"column:name;type:varchar(255);uniqueIndex"` // Name 字符串

	Code string `json:"code" gorm:"column:code;type:varchar(255)"` // Code 字符串

	Description string `json:"description" gorm:"column:description;type:varchar(255)"` // Description 字符串

	ParentId uint `json:"parent_id" gorm:"column:parent_id;type:int unsigned;default:0"` // ParentId 32位无符号整数

	Sort int `json:"sort" gorm:"column:sort;type:int;default:0"` // Sort 32位整数

	Status string `json:"status" gorm:"column:status;type:varchar(255);default:'active'"` // Status 字符串

	ManagerId uint `json:"manager_id" gorm:"column:manager_id;type:int unsigned"` // ManagerId 32位无符号整数

	// 树形结构相关字段
	Parent   *Department  `json:"parent,omitempty" gorm:"foreignKey:ParentId;references:ID"` // 父级部门
	Children []Department `json:"children,omitempty" gorm:"foreignKey:ParentId"`             // 子部门
	Path     string       `json:"path" gorm:"column:path;type:varchar(1000)"`                // 部门路径，如 "1,2,3"
	Level    int          `json:"level" gorm:"column:level;type:int;default:1"`              // 部门层级

}

// TableName 获取表名
func (Department) TableName() string {
	return "departments"
}

// BeforeCreate GORM 钩子：创建前
func (department *Department) BeforeCreate(tx *gorm.DB) error {
	// 调用基础模型的钩子
	if err := department.BaseModel.BeforeCreate(tx); err != nil {
		return err
	}

	// 设置路径和层级
	if department.ParentId == 0 {
		department.Path = ""
		department.Level = 1
	} else {
		var parent Department
		if err := tx.First(&parent, department.ParentId).Error; err != nil {
			return err
		}
		department.Level = parent.Level + 1
		if parent.Path == "" {
			department.Path = fmt.Sprintf("%d", parent.ID)
		} else {
			department.Path = fmt.Sprintf("%s,%d", parent.Path, parent.ID)
		}
	}

	return nil
}

// BeforeUpdate GORM 钩子：更新前
func (department *Department) BeforeUpdate(tx *gorm.DB) error {
	// 调用基础模型的钩子
	if err := department.BaseModel.BeforeUpdate(tx); err != nil {
		return err
	}

	// 在这里添加特定的更新前逻辑
	return nil
}
