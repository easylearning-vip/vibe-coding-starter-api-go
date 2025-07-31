package model

import (

	"database/sql"

	"time"


	"gorm.io/gorm"
)

// Order Order模型
type Order struct {
	BaseModel

	OrderNumber string `json:"order_number" gorm:"column:order_number;type:varchar(50);not null"` // OrderNumber 字符串

	CustomerName string `json:"customer_name" gorm:"column:customer_name;type:varchar(100);not null"` // CustomerName 字符串

	TotalAmount float64 `json:"total_amount" gorm:"column:total_amount;type:decimal(12,2);not null"` // TotalAmount 64位浮点数

	DiscountAmount sql.NullFloat64 `json:"discount_amount" gorm:"column:discount_amount;type:decimal(10,2);default:0.00"` // DiscountAmount 可空自定义类型

	Status sql.NullString `json:"status" gorm:"column:status;default:pending"` // Status 可空字符串

	IsPaid bool `json:"is_paid" gorm:"column:is_paid;type:tinyint;not null;default:0"` // IsPaid 布尔值

	PaymentMethod sql.NullString `json:"payment_method" gorm:"column:payment_method;type:varchar(50)"` // PaymentMethod 可空字符串

	OrderDate time.Time `json:"order_date" gorm:"column:order_date;type:date;not null"` // OrderDate 时间类型

	DeliveryDate sql.NullTime `json:"delivery_date" gorm:"column:delivery_date;type:datetime"` // DeliveryDate 可空时间类型

	Notes sql.NullString `json:"notes" gorm:"column:notes;type:text"` // Notes 可空字符串

}

// TableName 获取表名
func (Order) TableName() string {
	return "orders"
}

// BeforeCreate GORM 钩子：创建前
func (order *Order) BeforeCreate(tx *gorm.DB) error {
	// 调用基础模型的钩子
	if err := order.BaseModel.BeforeCreate(tx); err != nil {
		return err
	}

	// 在这里添加特定的创建前逻辑
	return nil
}

// BeforeUpdate GORM 钩子：更新前
func (order *Order) BeforeUpdate(tx *gorm.DB) error {
	// 调用基础模型的钩子
	if err := order.BaseModel.BeforeUpdate(tx); err != nil {
		return err
	}

	// 在这里添加特定的更新前逻辑
	return nil
}
