package main

import (
	"fmt"
	"log"
	"net/http"
	"research-ability-assessment/internal/agent"
	"research-ability-assessment/internal/config"
	"research-ability-assessment/internal/handler"
	"research-ability-assessment/internal/llm"
	"research-ability-assessment/internal/middleware"
	"research-ability-assessment/internal/models"
	repoNeo4j "research-ability-assessment/internal/repository/neo4j"
	repoPostgres "research-ability-assessment/internal/repository/postgres"
	"research-ability-assessment/internal/service"

	"github.com/gin-gonic/gin"
	neo4jdriver "github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// 加载配置
	log.Println("开始加载配置...")
	cfg, err := config.GetConfig()
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}
	log.Println("配置加载成功")

	// 连接数据库
	log.Println("开始连接数据库...")
	db, err := connectDatabase(cfg)
	if err != nil {
		log.Fatalf("连接数据库失败: %v", err)
	}
	log.Println("数据库连接成功")

	// 设置数据库字符集
	db.Exec("SET NAMES utf8mb4")
	db.Exec("SET CHARACTER SET utf8mb4")
	db.Exec("SET character_set_connection=utf8mb4")
	log.Println("数据库字符集设置成功")

	// 连接Neo4j
	log.Println("开始连接Neo4j...")
	driver, err := connectNeo4j(cfg)
	if err != nil {
		log.Fatalf("连接Neo4j失败: %v", err)
	}
	defer driver.Close()
	log.Println("Neo4j连接成功")

	// 自动迁移数据库表结构
	log.Println("开始迁移数据库表结构...")
	migrateDatabase(db)
	log.Println("数据库表结构迁移成功")

	// 初始化仓库
	log.Println("开始初始化仓库...")
	userRepo := repoPostgres.NewUserRepo(db)
	taskRepo := repoPostgres.NewTaskRepo(db)
	resultRepo := repoPostgres.NewResultRepo(db)
	graphRepo := repoNeo4j.NewGraphRepo(driver)
	log.Println("仓库初始化成功")

	// 初始化LLM客户端
	log.Println("开始初始化LLM客户端...")
	llmClient := llm.NewClient(&cfg.LLM)
	log.Println("LLM客户端初始化成功")

	// 初始化服务
	log.Println("开始初始化服务...")
	evidenceService := service.NewEvidenceService(db, llmClient)
	authService := service.NewAuthService(userRepo)
	taskService := service.NewTaskService(taskRepo, userRepo)
	inferenceService := service.NewInferenceService(resultRepo, evidenceService)
	reportService := service.NewReportService(inferenceService)
	log.Println("服务初始化成功")

	// 初始化Agent
	log.Println("开始初始化Agent...")
	ioUnit := agent.NewIOUnit(taskService, evidenceService)
	evidenceAgent := agent.NewEvidenceAgent(evidenceService, ioUnit)
	logicUnit := agent.NewLogicUnit()
	storageUnit := agent.NewStorageUnit(resultRepo, graphRepo)
	inferenceAgent := agent.NewInferenceAgent(evidenceAgent, logicUnit, inferenceService)
	controlUnit := agent.NewControlUnit(taskService, evidenceService, inferenceService, inferenceAgent, storageUnit)
	_ = controlUnit
	log.Println("Agent初始化成功")

	// 初始化处理器
	log.Println("开始初始化处理器...")
	authHandler := handler.NewAuthHandler(authService)
	taskHandler := handler.NewTaskHandler(taskService)
	evidenceHandler := handler.NewEvidenceHandler(evidenceService)
	resultHandler := handler.NewResultHandler(inferenceService, reportService)
	log.Println("处理器初始化成功")

	// 初始化路由
	log.Println("开始初始化路由...")
	r := setupRouter(authService, authHandler, taskHandler, evidenceHandler, resultHandler, llmClient)
	log.Println("路由初始化成功")

	// 启动服务器
	serverAddr := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Printf("服务器启动在 %s", serverAddr)
	if err := r.Run(serverAddr); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}

func connectDatabase(cfg *config.Config) (*gorm.DB, error) {
	var dialector gorm.Dialector

	switch cfg.Database.Type {
	case "mysql":
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local&collation=utf8mb4_unicode_ci",
			cfg.Database.User, cfg.Database.Password, cfg.Database.Host, cfg.Database.Port, cfg.Database.DBName)
		dialector = mysql.Open(dsn)
	case "postgres":
		fallthrough
	default:
		dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
			cfg.Database.Host, cfg.Database.Port, cfg.Database.User, cfg.Database.Password, cfg.Database.DBName, cfg.Database.SSLMode)
		dialector = postgres.Open(dsn)
	}

	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// 确保连接使用utf8mb4字符集
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.Exec("SET NAMES utf8mb4 COLLATE utf8mb4_unicode_ci")

	return db, nil
}

func connectNeo4j(cfg *config.Config) (neo4jdriver.Driver, error) {
	driver, err := neo4jdriver.NewDriver(cfg.Neo4j.URI, neo4jdriver.BasicAuth(cfg.Neo4j.Username, cfg.Neo4j.Password, ""))
	if err != nil {
		return nil, err
	}

	// 验证连接
	if err := driver.VerifyConnectivity(); err != nil {
		return nil, err
	}

	return driver, nil
}

func migrateDatabase(db *gorm.DB) {
	db.AutoMigrate(
		&models.User{},
		&models.Teacher{},
		&models.Student{},
		&models.Task{},
		&models.StudentTask{},
		&models.Evidence{},
		&models.InferenceResult{},
	)
}

func setupRouter(authService *service.AuthService, authHandler *handler.AuthHandler, taskHandler *handler.TaskHandler, evidenceHandler *handler.EvidenceHandler, resultHandler *handler.ResultHandler, llmClient *llm.Client) *gin.Engine {
	r := gin.Default()

	// 中间件
	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.LoggingMiddleware())

	// 公共路由
	public := r.Group("/api/v1")
	{
		// 认证路由
		auth := public.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}

		// 测试LLM API的路由
		public.POST("/test-llm", func(c *gin.Context) {
			var req struct {
				Content string `json:"content" binding:"required"`
			}

			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"code":    400,
					"message": "请求参数错误",
					"data":    nil,
				})
				return
			}

			// 调用LLM API
			messages := []llm.Message{
				{
					Role:    "system",
					Content: "你是一个研究能力评估专家，负责分析学生提交的证据并评估其质量。",
				},
				{
					Role:    "user",
					Content: req.Content,
				},
			}

			response, err := llmClient.Chat(c.Request.Context(), messages)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"code":    500,
					"message": "LLM API调用失败: " + err.Error(),
					"data":    nil,
				})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"code":    200,
				"message": "LLM API调用成功",
				"data":    map[string]string{"response": response},
			})
		})
	}

	// 受保护路由
	protected := r.Group("/api/v1")
	protected.Use(middleware.AuthMiddleware(authService))
	{
		// 用户路由
		user := protected.Group("/user")
		{
			user.GET("/info", authHandler.GetUserInfo)
		}

		// 任务路由
		task := protected.Group("/tasks")
		{
			task.POST("", taskHandler.CreateTask)
			task.GET("", taskHandler.GetTasksByTeacherID)
			task.GET("/:task_id", taskHandler.GetTaskByID)
			task.PUT("/:task_id/status", taskHandler.UpdateTaskStatus)
			task.POST("/:task_id/assign", taskHandler.AssignTask)
			task.GET("/:task_id/students", taskHandler.GetStudentTasks)
			task.GET("/students/list", taskHandler.GetStudents)
			task.GET("/students/assigned", taskHandler.GetAssignedTasks)
		}

		// 证据路由
		evidence := protected.Group("/evidences")
		{
			evidence.POST("", evidenceHandler.CreateEvidence)
			evidence.GET("", evidenceHandler.GetEvidences)
			evidence.GET("/:evidence_id", evidenceHandler.GetEvidenceByID)
			evidence.POST("/:evidence_id/analyze", evidenceHandler.AnalyzeEvidence)
			evidence.GET("/student-task/:student_task_id", evidenceHandler.GetEvidencesByStudentTaskID)
			evidence.GET("/student-task", evidenceHandler.GetEvidencesByStudentAndTask)
		}

		// 结果路由
		result := protected.Group("/results")
		{
			result.GET("", resultHandler.GetResults)
			result.GET("/student", resultHandler.GetStudentResults)
			result.GET("/:result_id", resultHandler.GetInferenceResultByID)
			result.GET("/task/:task_id", resultHandler.GetInferenceResultsByTaskID)
			result.GET("/student-task", resultHandler.GetInferenceResultByStudentAndTask)
			result.GET("/report/student", resultHandler.GenerateStudentReport)
			result.GET("/report/task/:task_id", resultHandler.GenerateTaskReport)
		}
	}

	return r
}
