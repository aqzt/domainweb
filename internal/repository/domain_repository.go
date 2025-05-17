package repository

import (
	"database/sql"
	"fmt"
	"time"

	"domainweb/internal/model"
)

// DomainRepository 处理域名相关的数据库操作
type DomainRepository struct {
	db *sql.DB
}

// NewDomainRepository 创建一个新的DomainRepository实例
func NewDomainRepository(db *sql.DB) *DomainRepository {
	return &DomainRepository{db: db}
}

// GetDomainAttributes 获取所有域名属性规则
func (r *DomainRepository) GetDomainAttributes() ([]model.DomainAttribute, error) {
	query := `SELECT id, attribute_name, attribute_type, price_factor, grade_factor, attribute_value
			  FROM domain_attributes`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("查询域名属性失败: %w", err)
	}
	defer rows.Close()

	var attributes []model.DomainAttribute
	for rows.Next() {
		var attr model.DomainAttribute
		if err := rows.Scan(
			&attr.ID,
			&attr.AttributeName,
			&attr.AttributeType,
			&attr.PriceFactor,
			&attr.GradeFactor,
			&attr.AttributeValue,
		); err != nil {
			return nil, fmt.Errorf("扫描域名属性行失败: %w", err)
		}
		attributes = append(attributes, attr)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("迭代域名属性行失败: %w", err)
	}

	return attributes, nil
}

// GetAttributesByType 根据属性类型获取域名属性
func (r *DomainRepository) GetAttributesByType(attrType string) ([]model.DomainAttribute, error) {
	query := `SELECT id, attribute_name, attribute_type, price_factor, grade_factor, attribute_value
			  FROM domain_attributes
			  WHERE attribute_type = ?`

	rows, err := r.db.Query(query, attrType)
	if err != nil {
		return nil, fmt.Errorf("查询属性类型失败: %w", err)
	}
	defer rows.Close()

	var attributes []model.DomainAttribute
	for rows.Next() {
		var attr model.DomainAttribute
		if err := rows.Scan(
			&attr.ID,
			&attr.AttributeName,
			&attr.AttributeType,
			&attr.PriceFactor,
			&attr.GradeFactor,
			&attr.AttributeValue,
		); err != nil {
			return nil, fmt.Errorf("扫描属性类型行失败: %w", err)
		}
		attributes = append(attributes, attr)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("迭代属性类型行失败: %w", err)
	}

	return attributes, nil
}

// GetTLDAttributes 获取所有TLD属性
func (r *DomainRepository) GetTLDAttributes() (map[string]model.DomainAttribute, error) {
	query := `SELECT id, attribute_name, attribute_type, price_factor, grade_factor, attribute_value
			  FROM domain_attributes
			  WHERE attribute_name LIKE '%后缀'`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("查询TLD属性失败: %w", err)
	}
	defer rows.Close()

	tldAttrs := make(map[string]model.DomainAttribute)
	for rows.Next() {
		var attr model.DomainAttribute
		if err := rows.Scan(
			&attr.ID,
			&attr.AttributeName,
			&attr.AttributeType,
			&attr.PriceFactor,
			&attr.GradeFactor,
			&attr.AttributeValue,
		); err != nil {
			return nil, fmt.Errorf("扫描TLD属性行失败: %w", err)
		}
		tldAttrs[attr.AttributeValue] = attr
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("迭代TLD属性行失败: %w", err)
	}

	return tldAttrs, nil
}

// SaveDomainInfo 保存域名基本信息
func (r *DomainRepository) SaveDomainInfo(domain *model.Domain) error {
	query := `INSERT INTO domains (name, tld, length, structure, register_date, expire_date, created_at, updated_at)
			  VALUES (?, ?, ?, ?, ?, ?, ?, ?)
			  ON DUPLICATE KEY UPDATE
			  structure = VALUES(structure),
			  register_date = VALUES(register_date),
			  expire_date = VALUES(expire_date),
			  updated_at = VALUES(updated_at)`

	now := time.Now()
	_, err := r.db.Exec(
		query,
		domain.Name,
		domain.TLD,
		domain.Length,
		domain.Structure,
		domain.RegisterDate,
		domain.ExpireDate,
		now,
		now,
	)

	if err != nil {
		return fmt.Errorf("保存域名信息失败: %w", err)
	}

	return nil
}
