package models

import "time"

type Transaction struct {
	ID        int           `json:"id" gorm:"primary_key:auto_increment"`
	TripsID   int           `json:"trips_id"`
	Trips     TripsResponse `json:"trips" gorm:"foreignKey:TripsID"`
	UsersID   int           `json:"users_id"`
	Users     UsersResponse `json:"name" gorm:"foreignKey: UsersID"`
	TotalPrc  int           `json:"totalprc"`
	Qty       int           `json:"qty"`
	Status    string        `json:"status" gorm:"type: varchar(20)"`
	CreatedAt time.Time     `json:"bookdate"`
	UpdatedAt time.Time     `json:"-"`
}

type TransactionResponse struct {
	ID        int           `json:"id"`
	TotalPrc  string        `json:"totalprc"`
	Qty       int           `json:"qty"`
	Status    string        `json:"status"`
	TripsID   int           `json:"trips_id"`
	Trips     Trips         `json:"accomodation" gorm:"foreignKey: TripsID"`
	UsersID   int           `json:"users_id"`
	Users     UsersResponse `json:"name" gorm:"foreignKey: UsersID"`
	CreatedAt time.Time     `json:"bookdate"`
}

func (TransactionResponse) TableName() string {
	return "profiles"
}
