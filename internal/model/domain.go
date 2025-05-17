package model

import (
	"time"
)

// Domain 表示域名及其属性
type Domain struct {
	ID           int64     `json:"id"`
	Name         string    `json:"name"`         // 域名名称，如 aqzt.com
	TLD          string    `json:"tld"`          // 顶级域名，如 com
	Length       int       `json:"length"`       // 域名长度（不含TLD）
	Structure    string    `json:"structure"`    // 域名结构，如 纯字母、数字字母混合等
	RegisterDate time.Time `json:"registerDate"` // 注册日期
	ExpireDate   time.Time `json:"expireDate"`   // 到期日期
	CreatedAt    time.Time `json:"createdAt"`    // 记录创建时间
	UpdatedAt    time.Time `json:"updatedAt"`    // 记录更新时间
}

// DomainAttribute 表示域名的各种属性及其对估价的影响
type DomainAttribute struct {
	ID             int64   `json:"id"`
	AttributeName  string  `json:"attributeName"`  // 属性名称，如 "com后缀"
	AttributeType  string  `json:"attributeType"`  // 属性类型，如 "基础属性" 或 "其他属性"
	PriceFactor    float64 `json:"priceFactor"`    // 估价倍数，如 9.55
	GradeFactor    float64 `json:"gradeFactor"`    // 等级增量，如 0.5
	AttributeValue string  `json:"attributeValue"` // 属性值，如 "com"
}

// EstimationResult 表示域名估价结果
type EstimationResult struct {
	Domain            string              `json:"domain"`            // 域名
	Grade             float64             `json:"grade"`             // 品相等级
	Price             float64             `json:"price"`             // 保守估价
	BaseAttributes    []AttributeDetail   `json:"baseAttributes"`    // 基础属性详情
	OtherAttributes   []AttributeDetail   `json:"otherAttributes"`   // 其他属性详情
	EstimationDate    time.Time           `json:"estimationDate"`    // 估价日期
}

// AttributeDetail 表示属性详情
type AttributeDetail struct {
	Name        string  `json:"name"`        // 属性名称
	Value       string  `json:"value"`       // 属性值
	Description string  `json:"description"` // 属性描述
	PriceFactor float64 `json:"priceFactor"` // 估价倍数
	GradeFactor float64 `json:"gradeFactor"` // 等级增量
}

// HistoryRecord 表示查询历史记录
type HistoryRecord struct {
	ID             int64     `json:"id"`
	Domain         string    `json:"domain"`         // 查询的域名
	Grade          float64   `json:"grade"`          // 品相等级
	Price          float64   `json:"price"`          // 估价结果
	EstimationDate time.Time `json:"estimationDate"` // 查询时间
}
