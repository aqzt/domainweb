package repository

import (
	"database/sql"
	"fmt"
	"time"

	"domainweb/internal/model"
)

// HistoryRepository 处理查询历史相关的数据库操作
type HistoryRepository struct {
	db *sql.DB
}

// NewHistoryRepository 创建一个新的HistoryRepository实例
func NewHistoryRepository(db *sql.DB) *HistoryRepository {
	return &HistoryRepository{db: db}
}

// SaveHistory 保存查询历史记录
func (r *HistoryRepository) SaveHistory(record *model.HistoryRecord) error {
	query := `INSERT INTO history_records (domain, grade, price, estimation_date)
			  VALUES (?, ?, ?, ?)`

	_, err := r.db.Exec(
		query,
		record.Domain,
		record.Grade,
		record.Price,
		time.Now(),
	)

	if err != nil {
		return fmt.Errorf("保存历史记录失败: %w", err)
	}

	return nil
}

// GetHistory 获取查询历史记录，可选择按域名筛选
func (r *HistoryRepository) GetHistory(domain string, limit int) ([]model.HistoryRecord, error) {
	var query string
	var args []interface{}

	if domain != "" {
		query = `SELECT id, domain, grade, price, estimation_date
				FROM history_records
				WHERE domain LIKE ?
				ORDER BY estimation_date DESC
				LIMIT ?`
		args = append(args, "%"+domain+"%", limit)
	} else {
		query = `SELECT id, domain, grade, price, estimation_date
				FROM history_records
				ORDER BY estimation_date DESC
				LIMIT ?`
		args = append(args, limit)
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("查询历史记录失败: %w", err)
	}
	defer rows.Close()

	var records []model.HistoryRecord
	for rows.Next() {
		var record model.HistoryRecord
		if err := rows.Scan(
			&record.ID,
			&record.Domain,
			&record.Grade,
			&record.Price,
			&record.EstimationDate,
		); err != nil {
			return nil, fmt.Errorf("扫描历史记录行失败: %w", err)
		}
		records = append(records, record)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("迭代历史记录行失败: %w", err)
	}

	return records, nil
}
