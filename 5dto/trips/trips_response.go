package tripsdto

type Trips struct {
	ID             int    `json:"id"`
	Name           string `json:"name" form:"name" gorm:"type: varchar(255)"`
	Desc           string `json:"desc" form:"desc" gorm:"type: text"`
	Price          int    `json:"price" form:"price" gorm:"type: int"`
	Accomodation   string `json:"accomodation" form:"accomodation" gorm:"type: varchar(255)"`
	Transportation string `json:"transportation" form:"transportation" gorm:"type: varchar(255)"`
	Eat            string `json:"eat" form:"eat" gorm:"type: varchar(255)"`
	DateTrip       string `json:"datetrip" form:"datetrip" gorm:"type: varchar(255)"`
	Quota          int    `json:"quota" form:"quota" gorm:"type: int"`
	Day            int    `json:"day" form:"day" gorm:"type: int"`
	Night          int    `json:"night" form:"night" gorm:"type: int"`
	Image          string `json:"image" form:"image" gorm:"type: varchar(255)"`
	UserID         int    `json:"user_id" gorm:"type: int"`
	CountryID      int    `json:"country_id" form:"country_id" gorm:"type: int"`
}
