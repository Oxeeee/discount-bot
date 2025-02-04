package services

import (
	"errors"
	"fmt"
	"log/slog"
	"math/rand/v2"
	"strconv"
	"strings"
	"time"

	"github.com/Oxeeee/discont-bot/internal/domain"
	"github.com/Oxeeee/discont-bot/internal/repo"
	csv "github.com/Oxeeee/discont-bot/pkg/CSV"
	"gorm.io/gorm"
)

type UserService interface {
	CheckRole(userId uint, role domain.UserRole) (bool, error)
	RegisterUser(user *domain.User) error
	GetUserByID(userID uint) (bool, *domain.User, error)
	VerifyCode(code string) (bool, error)
	GetUserRole(userID uint) (string, error)
	CheckWhitelist(userID uint) (bool, error)
	ManageWhitelist(username string, command string) error
	Userlist() (string, error)
	ChangeRole(userID uint, role string) error
	GetUserByUsername(username string) (bool, *domain.User, error)
	GetDiscountList() (string, error)
	SaveDiscountList(csv string) error
	GenerateCode(userID uint) (string, error)
}

type userService struct {
	repo repo.UsersRepo
	log  *slog.Logger
}

func NewUserService(repo repo.UsersRepo, log *slog.Logger) UserService {
	return &userService{
		repo: repo,
		log:  log,
	}
}

func (s *userService) CheckRole(userID uint, role domain.UserRole) (bool, error) {
	userRole, err := s.repo.GetRoleByID(userID)
	if err != nil {
		s.log.Error("error while checking role", "error", err)
		return false, err
	}

	if userRole != string(role) {
		s.log.Info("user role does not match the expected role", "expected_role", role, "user_role", userRole)
		return false, nil
	}

	return true, nil
}

func (s *userService) RegisterUser(user *domain.User) error {
	return s.repo.SaveUser(user)
}

func (s *userService) ChangeRole(userID uint, role string) error {
	err := s.repo.Update(userID, "role", role)
	if err != nil {
		return err
	}
	return nil
}

func (s *userService) GetUserByID(userID uint) (bool, *domain.User, error) {
	user, err := s.repo.GetUserByID(userID)
	if err != nil {
		return false, nil, err
	}

	return true, user, err
}

func (s *userService) VerifyCode(code string) (bool, error) {
	codeInfo, err := s.repo.GetCodeInfoByCode(code)
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	} else if err == gorm.ErrRecordNotFound {
		return false, nil
	}

	if codeInfo.ExpDate.Before(time.Now()) {
		return false, nil
	}

	return true, nil
}

func (s *userService) GetUserRole(userID uint) (string, error) {
	role, err := s.repo.GetRoleByID(userID)
	if err != nil {
		return "", err
	}

	return role, err
}

func (s *userService) CheckWhitelist(userID uint) (bool, error) {
	whitelisted, err := s.repo.IsWhitelisted(userID)
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	} else if err == gorm.ErrRecordNotFound {
		return false, nil
	}
	return whitelisted, nil
}

func (s *userService) ManageWhitelist(username string, command string) error {
	switch command {
	case "add":
		user, err := s.repo.GetUserByUsername(username)
		if err != nil {
			s.log.Error("error while getting user by username", "error", err)
			return err
		}

		err = s.repo.Update(user.ID, "whitelist", true)
		if err != nil {
			s.log.Error("error while add user in whitelist", "error", err)
			return err
		}

		return nil
	case "delete":
		user, err := s.repo.GetUserByUsername(username)
		if err != nil {
			s.log.Error("error while getting user by username", "error", err)
			return err
		}

		err = s.repo.Update(user.ID, "whitelist", false)
		if err != nil {
			s.log.Error("error while remove user from whitelist", "error", err)
			return err
		}

		return nil
	default:
		return errors.New("unexpected command")
	}
}

func (s *userService) Userlist() (string, error) {
	users, err := s.repo.GetUserlist()
	if err != nil {
		return "", err
	}

	if len(users) == 0 {
		return "Список пользователей пуст", nil
	}

	var sb strings.Builder
	sb.WriteString("Список пользователей:\n\n")

	for i, user := range users {
		username := user.Username
		if username == "" {
			username = "Неизвестно"
		} else if username[0] != '@' {
			username = "@" + username
		}
		userRole := user.Role
		switch userRole {
		case string(domain.UserRoleStaff):
			userRole = "Сотрудник"
		case string(domain.UserRoleAdmin):
			userRole = "Администратор"
		case string(domain.UserRoleUser):
			userRole = "Пользователь"
		}

		sb.WriteString(fmt.Sprintf("%d. Имя пользователя: %s\n", i+1, username))
		sb.WriteString(fmt.Sprintf("   Роль: %s\n", userRole))
		sb.WriteString(fmt.Sprintf("   Разрешение на использование: %t\n\n", user.Whitelist))
	}

	return sb.String(), nil
}

func (s *userService) GetUserByUsername(username string) (bool, *domain.User, error) {
	user, err := s.repo.GetUserByUsername(username)
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, nil, err
	} else if err == gorm.ErrRecordNotFound {
		return false, nil, nil
	}

	return true, user, err
}

func (s *userService) GetDiscountList() (string, error) {
	places, err := s.repo.GetPlaces()
	if err != nil {
		return "", err
	}

	csv, err := csv.ConvertToCSV(places)
	if err != nil {
		return "", err
	}

	return csv, nil
}

func (s *userService) SaveDiscountList(csvContent string) error {
	places, err := csv.ConvertFromCSV(csvContent)
	if err != nil {
		return err
	}

	return s.repo.SavePlaces(places)
}

func (s *userService) GenerateCode(userID uint) (string, error) {
	var code = domain.DiscountCode{
		Code:    strconv.Itoa(rand.IntN(900000) + 100000),
		UserID:  userID,
		ExpDate: time.Now().Add(time.Minute * 15),
	}

	err := s.repo.SaveCode(&code)
	if err != nil {
		return "", err
	}

	return code.Code, nil
}
