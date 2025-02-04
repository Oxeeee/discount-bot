package domain

import "time"

type User struct {
	ID        uint `gorm:"primaryKey"`
	ChatID    int64
	Username  string
	Role      string `gorm:"default:'user'"`
	Whitelist bool   `gorm:"default:false"`
}

type DiscountCode struct {
	ID      uint `gorm:"primaryKey"`
	UserID  uint `gorm:"index"`
	Code    string
	ExpDate time.Time
}

type DiscountLog struct {
	ID      uint      `gorm:"primaryKey"`
	UserID  uint      `gorm:"index"`
	UseTime time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	PlaceID uint
}

type Place struct {
	ID             uint `gorm:"primaryKey"`
	Name           string
	Address        string
	DiscountFactor string
}
