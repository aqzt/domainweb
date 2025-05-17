package api

import (
	"net/http"
	"strconv"

	"domainweb/internal/service"
	"github.com/gin-gonic/gin"
)

// Handler 处理HTTP请求
type Handler struct {
	domainService  *service.DomainService
	historyService *service.HistoryService
}

// NewHandler 创建一个新的Handler实例
func NewHandler(domainService *service.DomainService, historyService *service.HistoryService) *Handler {
	return &Handler{
		domainService:  domainService,
		historyService: historyService,
	}
}

// HomePage 处理首页请求
func (h *Handler) HomePage(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"title": "域名估价系统",
	})
}

// EstimateDomain 处理域名估价请求（Web界面）
func (h *Handler) EstimateDomain(c *gin.Context) {
	domain := c.PostForm("domain")
	if domain == "" {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{
			"error": "请输入有效的域名",
		})
		return
	}

	result, err := h.domainService.EstimateDomain(domain)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"error": "估价失败: " + err.Error(),
		})
		return
	}

	// 保存查询历史
	if err := h.historyService.SaveHistory(result); err != nil {
		// 仅记录错误，不影响用户体验
		c.Error(err)
	}

	c.HTML(http.StatusOK, "result.html", gin.H{
		"title":  "估价结果",
		"result": result,
	})
}

// GetHistory 处理查询历史请求（Web界面）
func (h *Handler) GetHistory(c *gin.Context) {
	domain := c.Query("domain")
	limitStr := c.DefaultQuery("limit", "50")

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 50
	}

	records, err := h.historyService.GetHistory(domain, limit)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"error": "获取历史记录失败: " + err.Error(),
		})
		return
	}

	c.HTML(http.StatusOK, "history.html", gin.H{
		"title":   "查询历史",
		"records": records,
		"domain":  domain,
	})
}

// APIEstimateDomain 处理域名估价请求（API）
func (h *Handler) APIEstimateDomain(c *gin.Context) {
	var request struct {
		Domain string `json:"domain" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数"})
		return
	}

	result, err := h.domainService.EstimateDomain(request.Domain)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 保存查询历史
	if err := h.historyService.SaveHistory(result); err != nil {
		// 仅记录错误，不影响用户体验
		c.Error(err)
	}

	c.JSON(http.StatusOK, result)
}

// APIGetHistory 处理查询历史请求（API）
func (h *Handler) APIGetHistory(c *gin.Context) {
	domain := c.Query("domain")
	limitStr := c.DefaultQuery("limit", "50")

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 50
	}

	records, err := h.historyService.GetHistory(domain, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, records)
}
