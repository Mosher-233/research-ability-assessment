package main

import (
	"fmt"
	"log"
	"research-ability-assessment/internal/config"
	"research-ability-assessment/internal/models"
	"research-ability-assessment/pkg/utils"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	log.Println("开始初始化数据库...")

	// 加载配置
	cfg, err := config.GetConfig()
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 连接数据库
	db, err := connectDatabase(cfg)
	if err != nil {
		log.Fatalf("连接数据库失败: %v", err)
	}
	log.Println("数据库连接成功")

	// 清空所有表
	log.Println("清空数据库表...")
	clearDatabase(db)

	// 创建测试数据
	log.Println("创建测试数据...")
	createTestData(db)

	log.Println("数据库初始化完成！")
}

func connectDatabase(cfg *config.Config) (*gorm.DB, error) {
	var dsn string

	switch cfg.Database.Type {
	case "mysql":
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			cfg.Database.User, cfg.Database.Password, cfg.Database.Host, cfg.Database.Port, cfg.Database.DBName)
		return gorm.Open(mysql.Open(dsn), &gorm.Config{})
	default:
		return nil, fmt.Errorf("不支持的数据库类型: %s", cfg.Database.Type)
	}
}

func clearDatabase(db *gorm.DB) {
	// 清空表（注意顺序，避免外键约束问题）
	db.Exec("SET FOREIGN_KEY_CHECKS = 0")
	db.Exec("TRUNCATE TABLE feedbacks")
	db.Exec("TRUNCATE TABLE evidences")
	db.Exec("TRUNCATE TABLE student_tasks")
	db.Exec("TRUNCATE TABLE tasks")
	db.Exec("TRUNCATE TABLE users")
	db.Exec("SET FOREIGN_KEY_CHECKS = 1")
	log.Println("数据库表清空完成")
}

func createTestData(db *gorm.DB) {
	// 1. 创建3个教师账号
	teachers := []struct {
		name  string
		email string
	}{
		{"张三老师", "1@tea.com"},
		{"李四老师", "2@tea.com"},
		{"John Smith", "3@tea.com"},
	}

	var teacherIDs []string
	for _, t := range teachers {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
		teacher := &models.User{
			ID:       utils.GenerateUserID("teacher"),
			Name:     t.name,
			Email:    t.email,
			Password: string(hashedPassword),
			Role:     "teacher",
		}
		if err := db.Create(teacher).Error; err != nil {
			log.Printf("创建教师 %s 失败: %v", t.name, err)
		} else {
			teacherIDs = append(teacherIDs, teacher.ID)
			log.Printf("创建教师成功: %s (%s), 密码: 123456", t.name, teacher.ID)
		}
	}

	// 2. 创建5个学生账号
	students := []struct {
		name      string
		email     string
		studentID string
		major     string
		grade     string
	}{
		{"王五学生", "1@stu.com", "2024001", "计算机科学", "大三"},
		{"赵六学生", "2@stu.com", "2024002", "数据科学", "大三"},
		{"Alice Wang", "3@stu.com", "2024003", "人工智能", "大二"},
		{"钱七学生", "4@stu.com", "2024004", "软件工程", "大四"},
		{"Bob Chen", "5@stu.com", "2024005", "网络安全", "大三"},
	}

	var studentIDs []string
	for _, s := range students {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
		student := &models.User{
			ID:       utils.GenerateUserID("student"),
			Name:     s.name,
			Email:    s.email,
			Password: string(hashedPassword),
			Role:     "student",
		}
		if err := db.Create(student).Error; err != nil {
			log.Printf("创建学生 %s 失败: %v", s.name, err)
		} else {
			studentIDs = append(studentIDs, student.ID)
			log.Printf("创建学生成功: %s (%s), 密码: 123456", s.name, student.ID)
		}
	}

	// 3. 创建一个代码编写项目任务
	if len(teacherIDs) > 0 {
		task := &models.Task{
			ID:          utils.GenerateTaskID(),
			Name:        "Python数据分析项目",
			Description: "完成一个基于Python的数据分析项目，包括数据清洗、可视化和统计分析。",
			CourseID:    "DS2024",
			TeacherID:   teacherIDs[0],
			StartDate:   time.Now(),
			EndDate:     time.Now().AddDate(0, 1, 0),
			Status:      "active",
		}
		if err := db.Create(task).Error; err != nil {
			log.Printf("创建任务失败: %v", err)
		} else {
			log.Printf("创建任务成功: %s (%s)", task.Name, task.ID)

			// 4. 将任务分配给所有学生
			for _, studentID := range studentIDs {
				studentTask := &models.StudentTask{
					ID:        utils.GenerateStudentTaskID(),
					TaskID:    task.ID,
					StudentID: studentID,
					Status:    "pending",
					Progress:  0,
				}
				if err := db.Create(studentTask).Error; err != nil {
					log.Printf("分配任务给学生 %s 失败: %v", studentID, err)
				} else {
					log.Printf("分配任务给学生成功: %s -> %s", studentID, task.ID)
				}
			}
		}
	}

	log.Println("测试数据创建完成！")
	log.Println("\n=== 账号信息 ===")
	log.Println("教师账号（密码均为 123456）:")
	for i, t := range teachers {
		log.Printf("  %d. %s - %s", i+1, t.name, t.email)
	}
	log.Println("\n学生账号（密码均为 123456）:")
	for i, s := range students {
		log.Printf("  %d. %s - %s", i+1, s.name, s.email)
	}
}
