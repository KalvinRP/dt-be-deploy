package repositories

import (
	models "dewetour/1models"

	"gorm.io/gorm"
)

type CountryRepository interface {
	FindCountry() ([]models.Country, error)
	GetCountry(ID int) (models.Country, error)
	MakeCountry(country models.Country) (models.Country, error)
	EditCountry(country models.Country, ID int) (models.Country, error)
	DeleteCountry(country models.Country, ID int) (models.Country, error)
}

func RepositoryCountry(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) FindCountry() ([]models.Country, error) {
	var country []models.Country
	err := r.db.Order("name").Find(&country).Error

	return country, err
}

func (r *repository) GetCountry(ID int) (models.Country, error) {
	var user models.Country
	err := r.db.First(&user, ID).Error

	return user, err
}

func (r *repository) MakeCountry(user models.Country) (models.Country, error) {
	err := r.db.Create(&user).Error

	return user, err
}

func (r *repository) EditCountry(user models.Country, ID int) (models.Country, error) {
	err := r.db.Save(&user).Error

	return user, err
}

func (r *repository) DeleteCountry(user models.Country, ID int) (models.Country, error) {
	err := r.db.Delete(&user).Error

	return user, err
}
