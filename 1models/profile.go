package models

import "time"

type Profile struct {
	ID        int           `json:"id" gorm:"primary_key:auto_increment"`
	UserID    int           `json:"user_id" gorm:"type: varchar(255)"`
	User      UsersResponse `json:"user"`
	Phone     string        `json:"phone" gorm:"type: varchar(255)"`
	Address   string        `json:"address" gorm:"type: text"`
	Gender    string        `json:"gender" gorm:"type: enum('Laki-laki', 'Perempuan', 'Lainnya')"`
	CreatedAt time.Time     `json:"-"`
	UpdatedAt time.Time     `json:"-"`
}

type ProfileResponse struct {
	UserID  int    `json:"-"`
	Phone   string `json:"phone"`
	Address string `json:"address"`
	Gender  string `json:"gender"`
}

func (ProfileResponse) TableName() string {
	return "profiles"
}
