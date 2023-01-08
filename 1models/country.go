package models

import "time"

type Country struct {
	ID        int       `json:"id" gorm:"primary_key:auto_increment"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

type CountryResponse struct {
	Name string `json:"name"`
}

func (CountryResponse) TableName() string {
	return "countries"
}
