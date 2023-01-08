package repositories

import (
	models "dewetour/1models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TransactionRepository interface {
	FindTrans() ([]models.Transaction, error)
	GetTrans(ID int) (models.Transaction, error)
	MakeTrans(transaction models.Transaction) (models.Transaction, error)
	EditTrans(transaction models.Transaction, ID int) (models.Transaction, error)
	DeleteTrans(transaction models.Transaction, ID int) (models.Transaction, error)
	GetOneTrans(ID string) (models.Transaction, error)
	UpdateTrans(status string, ID models.Transaction) error
	UserHistory(ID int) ([]models.Transaction, error)
}

func RepositoryTrans(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) FindTrans() ([]models.Transaction, error) {
	var transaction []models.Transaction
	err := r.db.Preload("Trips").Preload("Trips.Country").Preload("Users").Find(&transaction).Error

	return transaction, err
}

func (r *repository) GetTrans(ID int) (models.Transaction, error) {
	var trans models.Transaction
	err := r.db.Preload("Trips").Preload("Trips.Country").Preload("Users").First(&trans, ID).Error

	return trans, err
}

func (r *repository) MakeTrans(trans models.Transaction) (models.Transaction, error) {
	err := r.db.Create(&trans).Error

	return trans, err
}

func (r *repository) EditTrans(trans models.Transaction, ID int) (models.Transaction, error) {
	err := r.db.Model(&trans).Updates(trans).Error

	return trans, err
}

func (r *repository) DeleteTrans(trans models.Transaction, ID int) (models.Transaction, error) {
	err := r.db.Delete(&trans).Error

	return trans, err
}

func (r *repository) GetOneTrans(ID string) (models.Transaction, error) {
	var transaction models.Transaction
	err := r.db.Preload(clause.Associations).Preload("Trips.Country").First(&transaction, "id = ?", ID).Error

	return transaction, err
}

func (r *repository) UpdateTrans(status string, transaction models.Transaction) error {
	// var transaction models.Transaction
	// r.db.Preload(clause.Associations).Preload("Trips.Country").First(&transaction, ID)

	if status != transaction.Status && status == "success" {
		var trips models.Trips
		r.db.First(&trips, transaction.Trips.ID)
		trips.Quota = trips.Quota - transaction.Qty
		r.db.Save(&trips)
	}

	transaction.Status = status

	err := r.db.Save(&transaction).Error

	return err
}

func (r *repository) UserHistory(ID int) ([]models.Transaction, error) {
	var history []models.Transaction
	err := r.db.Preload(clause.Associations).Preload("Trips.Country").Where("users_id = ?", ID).Find(&history).Error
	// err := r.db.Raw("SELECT * FROM transactions WHERE users_id = ?", ID).Scan(&history).Error

	return history, err
}
