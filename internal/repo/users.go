package repo

import (
	"github.com/Oxeeee/discont-bot/internal/domain"
	"gorm.io/gorm"
)

type UsersRepo interface {
	GetRoleByID(userID uint) (string, error)
	SaveUser(user *domain.User) error
	Update(userID uint, col string, data any) error
	GetUserByID(userID uint) (*domain.User, error)
	GetUserByUsername(username string) (*domain.User, error)
	GetCodeInfoByCode(code string) (*domain.DiscountCode, error)
	IsWhitelisted(userID uint) (bool, error)
	GetUserlist() ([]struct {
		Username  string
		Role      string
		Whitelist bool
	}, error)
	GetPlaces() ([]domain.Place, error)
	SavePlaces(places []domain.Place) error
	SaveCode(code *domain.DiscountCode) error
	SaveCodeLog(log *domain.DiscountLog) error
	DeactivateCode(code *domain.DiscountCode) error
}

type userDatabase struct {
	DB *gorm.DB
}

func NewUsersRepo(DB *gorm.DB) UsersRepo {
	return &userDatabase{DB}
}

func (db *userDatabase) GetRoleByID(userID uint) (string, error) {
	var role string
	if err := db.DB.Model(&domain.User{}).Select("role").Where("id = ?", userID).Scan(&role).Error; err != nil {
		return "", err
	}
	return role, nil
}

func (db *userDatabase) SaveUser(user *domain.User) error {
	return db.DB.Model(&domain.User{}).Where("id = ?", user.ID).Save(user).Error
}

func (db *userDatabase) Update(userID uint, col string, data any) error {
	return db.DB.Model(&domain.User{}).Where("id = ?", userID).Update(col, data).Error
}

func (db *userDatabase) GetUserByID(userID uint) (*domain.User, error) {
	var user domain.User
	err := db.DB.Model(&domain.User{}).Where("id = ?", userID).First(&user).Error
	return &user, err
}

func (db *userDatabase) GetCodeInfoByCode(code string) (*domain.DiscountCode, error) {
	var codeInfo domain.DiscountCode
	err := db.DB.Model(&domain.DiscountCode{}).Where("code = ?", code).First(&codeInfo).Error
	return &codeInfo, err
}

func (db *userDatabase) GetColByID(userID uint, col string) (any, error) {
	var data any
	err := db.DB.Model(&domain.User{}).Select("email").Where("id = ?", userID).First(&data).Error
	return data, err
}

func (db *userDatabase) IsWhitelisted(userID uint) (bool, error) {
	var isWhitelisted bool
	err := db.DB.Model(&domain.User{}).Select("whitelist").Where("id = ?", userID).First(&isWhitelisted).Error
	return isWhitelisted, err
}

func (db *userDatabase) GetUserByUsername(username string) (*domain.User, error) {
	var user domain.User
	err := db.DB.Model(&domain.User{}).Where("username = ?", username).First(&user).Error
	return &user, err
}

func (db *userDatabase) GetUserlist() ([]struct {
	Username  string
	Role      string
	Whitelist bool
}, error) {
	var users []struct {
		Username  string
		Role      string
		Whitelist bool
	}

	err := db.DB.Model(&domain.User{}).Select("username, role, whitelist").Find(&users).Error
	return users, err
}

func (db *userDatabase) GetPlaces() ([]domain.Place, error) {
	places := make([]domain.Place, 5, 5)
	err := db.DB.Model(&domain.Place{}).Select("id, name, address, discount_factor").Find(&places).Error
	return places, err
}

func (db *userDatabase) SavePlaces(places []domain.Place) error {
	return db.DB.Model(&domain.Place{}).Save(places).Error
}

func (db *userDatabase) SaveCode(code *domain.DiscountCode) error {
	return db.DB.Model(&domain.DiscountCode{}).Save(code).Error
}

func (db *userDatabase) SaveCodeLog(log *domain.DiscountLog) error {
	return db.DB.Model(&domain.DiscountLog{}).Save(log).Error
}

func (db *userDatabase) DeactivateCode(code *domain.DiscountCode) error {
	return db.DB.Model(&domain.DiscountCode{}).Where("id = ?", code.ID).Update("exp_date", code.ExpDate).Error
}
