package transactiondto

type TransactionRequest struct {
	ID       int    `json:"id" gorm:"primary_key:auto_increment"`
	TripsID  int    `json:"trips_id"`
	UsersID  int    `json:"Users_id"`
	TotalPrc int    `json:"totalprc"`
	Qty      int    `json:"qty"`
	Status   string `json:"status" gorm:"type: varchar(20)"`
}

type UpdateTransactionRequest struct {
	TripsID  int
	TotalPrc int
	Qty      int
	Status   string
}
