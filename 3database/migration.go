package database

import (
	models "dewetour/1models"
	mysql "dewetour/2pkg/mysql"
	"fmt"
)

func Migrate() {
	err := mysql.DB.AutoMigrate(&models.User{}, &models.Trips{}, &models.Country{}, &models.Transaction{})

	if err != nil {
		fmt.Println(err)
		panic("Migration failed")
	}

	fmt.Println("Migration success")
}
