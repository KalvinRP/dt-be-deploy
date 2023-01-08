package transactiondto

import (
	models "dewetour/1models"
	"time"
)

type TransactionResponse struct {
	TotalPrc  int                  `json:"totalprc"`
	Status    string               `json:"status"`
	TripsID   int                  `json:"trips_id"`
	Trips     models.TripsResponse `json:"trips" gorm:"foreignKey: TripsID"`
	Users     models.UsersResponse `json:"name" gorm:"foreignKey: UsersID"`
	Qty       int                  `json:"qty"`
	CreatedAt time.Time            `json:"bookdate"`
}
