package model

import "time"

type AccountCreationData struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AccountLogingData struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type User struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string    `gorm:"type:varchar(255);not null" json:"name"`
	Email     string    `gorm:"type:varchar(255);unique;not null" json:"email"`
	Password  string    `gorm:"type:text;not null" json:"-"` // hide password
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

type URL struct {
	ID         int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	LongURL    string     `gorm:"type:text;not null" json:"long_url"`
	ShortURLID string     `gorm:"type:text;not null;unique" json:"short_url_id"`
	UserID     int64      `gorm:"constraint:OnDelete:CASCADE;foreignKey:UserID;references:ID" json:"user_id"`
	TTL        *time.Time `gorm:"type:timestamp" json:"ttl"`
	CreatedAt  time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
}

type URLClicks struct {
	ID        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	URLId     int64     `gorm:"constraint:OnDelete:CASCADE;foreignKey:URLId;references:ID" json:"url_id"`
	ClickedAt time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"clicked_at"`
	IPAddress string    `gorm:"type:varchar(50)" json:"ip_address"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

type LongUrl struct {
	LongURL string `json:"long_url" binding:"required,url"`
}
