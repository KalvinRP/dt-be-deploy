package models

import "time"

type User struct {
	ID        int       `json:"id" gorm:"primary_key:auto_increment"`
	Name      string    `json:"name" form:"name" gorm:"type: varchar(255)"`
	Email     string    `json:"email" form:"email" gorm:"type: varchar(255)"`
	Password  string    `json:"password" form:"password" gorm:"type: varchar(255)"`
	Phone     string    `json:"phone" form:"phone" gorm:"type: varchar(20)"`
	Address   string    `json:"address" form:"address" gorm:"type: varchar(255)"`
	Gender    string    `json:"gender" form:"gender" gorm:"type: varchar(20)"`
	Role      string    `gorm:"type:varchar(5)"`
	Image     string    `json:"image" form:"image" gorm:"type: varchar(255)"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

type UsersResponse struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Gender  string `json:"gender"`
	Phone   string `json:"phone"`
	Email   string `json:"email"`
	Address string `json:"address"`
	Image   string `json:"image"`
}

type AuthResponse struct {
	Email string
	Role  string
}

func (UsersResponse) TableName() string {
	return "users"
}
