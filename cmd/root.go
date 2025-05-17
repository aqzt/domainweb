package cmd

import (
	"context"
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"domainweb/internal/api"
	"domainweb/internal/repository"
	"domainweb/internal/service"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

// Execute 是应用程序的入口点
func Execute() error {
	// 初始化数据库连接
	db, err := initDB()
	if err != nil {
		return fmt.Errorf("数据库初始化失败: %w", err)
	}
	defer db.Close()

	// 初始化存储库
	domainRepo := repository.NewDomainRepository(db)
	historyRepo := repository.NewHistoryRepository(db)

	// 初始化服务
	domainService := service.NewDomainService(domainRepo)
	historyService := service.NewHistoryService(historyRepo)

	// 设置Gin路由
	router := setupRouter(domainService, historyService)

	// 创建HTTP服务器
	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	// 在goroutine中启动服务器
	go func() {
		fmt.Println("服务器启动在 http://localhost:8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("监听失败: %s\n", err)
		}
	}()

	// 等待中断信号以优雅地关闭服务器
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("关闭服务器...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("服务器强制关闭:", err)
	}

	log.Println("服务器优雅退出")
	return nil
}

// 初始化数据库连接
func initDB() (*sql.DB, error) {
	// 从环境变量或配置文件获取数据库连接信息
	// 这里使用硬编码的连接信息作为示例
	dsn := "sfs:d3v777aaa@tcp(127.0.0.1:3306)/domainweb?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	// 测试数据库连接
	if err := db.Ping(); err != nil {
		return nil, err
	}

	// 设置连接池参数
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(time.Hour)

	return db, nil
}

// 设置Gin路由
func setupRouter(domainService *service.DomainService, historyService *service.HistoryService) *gin.Engine {
	router := gin.Default()

	// 加载HTML模板并添加自定义函数
	router.SetFuncMap(template.FuncMap{
		"contains": strings.Contains,
	})
	router.LoadHTMLGlob("web/templates/*")

	// 提供静态文件
	router.Static("/static", "./web/static")

	// 设置API处理器
	handler := api.NewHandler(domainService, historyService)

	// 定义路由
	router.GET("/", handler.HomePage)
	router.POST("/estimate", handler.EstimateDomain)
	router.GET("/history", handler.GetHistory)

	// API路由
	apiGroup := router.Group("/api")
	{
		apiGroup.POST("/estimate", handler.APIEstimateDomain)
		apiGroup.GET("/history", handler.APIGetHistory)
	}

	return router
}
