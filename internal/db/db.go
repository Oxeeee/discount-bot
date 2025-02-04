package db

import (
	"log"

	"github.com/Oxeeee/discont-bot/internal/config"
	"github.com/Oxeeee/discont-bot/internal/domain"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Database interface {
	AutoMigrate(dst ...any) error
	GetDB() *gorm.DB
}

type gormDatabase struct {
	Conn *gorm.DB
}

func (g *gormDatabase) AutoMigrate(dst ...any) error {
	return g.Conn.AutoMigrate(dst...)
}

func (g *gormDatabase) GetDB() *gorm.DB {
	return g.Conn
}

func ConntectDatabase(cfg *config.Config) Database {
	conn, err := gorm.Open(sqlite.Open(cfg.DatabaseRoute), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error while connecting to database: %v", err)
	}

	db := &gormDatabase{Conn: conn}
	err = db.AutoMigrate(&domain.User{}, &domain.Place{}, &domain.DiscountCode{}, &domain.DiscountLog{})
	if err != nil {
		log.Fatalf("Error while migrating tables: %v", err)
	}

	CreateDefaultAdmin(db.GetDB(), cfg.DefaultAdmin.Username, cfg.UserID)

	return db
}

func CreateDefaultAdmin(DB *gorm.DB, username string, userID uint) {
	admin := domain.User{
		Username:  username,
		ChatID:    int64(userID),
		ID:        userID,
		Role:      "admin",
		Whitelist: true,
	}

	var exitingAdmin domain.User
	res := DB.Where("username = ?", admin.Username).First(&exitingAdmin)
	if res.Error != nil && res.Error != gorm.ErrRecordNotFound {
		log.Fatalf("Error while searching administrator: %v", res.Error)
	}

	if res.RowsAffected == 0 {
		DB.Create(&admin)
		log.Println("Default administrator created")
	} else {
		log.Println("Administrator already exists")
	}
}
