package service

import (
	"blog-backend/internal/app/model"
	"blog-backend/internal/app/repository"
	"errors"
	"gorm.io/gorm"
)

// UserService 用户服务接口
type UserService struct {
	userRepo *repository.UserRepository
	DB       *gorm.DB // 添加数据库实例
}

// NewUserService 创建用户服务实例
func NewUserService(userRepo *repository.UserRepository, db *gorm.DB) *UserService {
	return &UserService{userRepo: userRepo, DB: db}
}

// GetUserByID 根据ID获取用户
func (s *UserService) GetUserByID(id uint) (*model.User, error) {
	return s.userRepo.GetByID(id)
}

// Login 用户登录
func (s *UserService) Login(username, password string) (*model.User, error) {
	// 根据用户名查找用户
	user, err := s.userRepo.GetByUsername(username)
	if err != nil {
		return nil, errors.New("用户名或密码错误")
	}

	// 检查用户状态
	if user.Status != 1 {
		return nil, errors.New("账号已被禁用")
	}

	// 验证密码
	if !user.CheckPassword(password) {
		return nil, errors.New("用户名或密码错误")
	}

	return user, nil
}

// Register 用户注册
func (s *UserService) Register(user *model.User) error {
	// 检查用户名是否已存在
	_, err := s.userRepo.GetByUsername(user.Username)
	if err == nil {
		return errors.New("用户名已存在")
	}

	// 检查邮箱是否已存在
	_, err = s.userRepo.GetByEmail(user.Email)
	if err == nil {
		return errors.New("邮箱已被注册")
	}

	// 设置默认角色和状态
	if user.Role == "" {
		user.Role = "user"
	}
	user.Status = 1

	return s.userRepo.Create(user)
}

// GetAllUsers 获取所有用户（管理员功能）
func (s *UserService) GetAllUsers() ([]*model.User, error) {
	var users []*model.User
	err := s.userRepo.Db.Find(&users).Error
	return users, err
}

// UpdateUserStatus 更新用户状态（管理员功能）
func (s *UserService) UpdateUserStatus(id uint, status int) error {
	return s.userRepo.Db.Model(&model.User{}).Where("id = ?", id).Update("status", status).Error
}

// DeleteUser 删除用户（管理员功能）
func (s *UserService) DeleteUser(id uint) error {
	// 不允许删除管理员账户
	var user model.User
	if err := s.userRepo.Db.First(&user, id).Error; err != nil {
		return err
	}
	if user.Role == "admin" {
		return errors.New("不能删除管理员账户")
	}
	return s.userRepo.Db.Delete(&model.User{}, id).Error
}
