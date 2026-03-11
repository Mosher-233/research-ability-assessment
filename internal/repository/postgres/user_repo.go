package postgres

import (
	"context"
	"research-ability-assessment/internal/models"

	"gorm.io/gorm"
)

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) CreateUser(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *UserRepo) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepo) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).First(&user, "email = ?", email).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepo) GetStudents(ctx context.Context) ([]models.Student, error) {
	var students []models.Student
	if err := r.db.WithContext(ctx).Where("role = ?", "student").Find(&students).Error; err != nil {
		return nil, err
	}
	return students, nil
}

func (r *UserRepo) GetTeachers(ctx context.Context) ([]models.Teacher, error) {
	var teachers []models.Teacher
	if err := r.db.WithContext(ctx).Where("role = ?", "teacher").Find(&teachers).Error; err != nil {
		return nil, err
	}
	return teachers, nil
}