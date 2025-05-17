/*
 * @Author: 5553557@qq.com
 * @Date: 2025-05-10 14:07:23
 * @LastEditors: 5553557@qq.com
 * @LastEditTime: 2025-05-10 15:00:15
 * @FilePath: \domainweb\internal\service\history_service.go
 * @Description:
 *
 * Copyright (c) 2022 by 5553557@qq.com, All Rights Reserved.
 */
package service

import (
	"domainweb/internal/model"
	"domainweb/internal/repository"
)

// HistoryService 处理查询历史的业务逻辑
type HistoryService struct {
	repo *repository.HistoryRepository
}

// NewHistoryService 创建一个新的HistoryService实例
func NewHistoryService(repo *repository.HistoryRepository) *HistoryService {
	return &HistoryService{repo: repo}
}

// SaveHistory 保存查询历史记录
func (s *HistoryService) SaveHistory(result *model.EstimationResult) error {
	record := &model.HistoryRecord{
		Domain:         result.Domain,
		Grade:          result.Grade,
		Price:          result.Price,
		EstimationDate: result.EstimationDate,
	}

	return s.repo.SaveHistory(record)
}

// GetHistory 获取查询历史记录
func (s *HistoryService) GetHistory(domain string, limit int) ([]model.HistoryRecord, error) {
	if limit <= 0 {
		limit = 50 // 默认限制为50条记录
	}

	return s.repo.GetHistory(domain, limit)
}
