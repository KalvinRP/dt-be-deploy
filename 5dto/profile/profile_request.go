package profiledto

type ProfileRequest struct {
	Phone   string `json:"phone" gorm:"type: varchar(255)"`
	Address string `json:"address" gorm:"type: text"`
	Gender  string `json:"gender" gorm:"type: enum('Laki-laki', 'Perempuan', 'Lainnya')"`
	UserID  int    `json:"user_id" gorm:"type: int"`
}

type UpdateProfileRequest struct {
	Phone   string `json:"phone" form:"phone"`
	Address string `json:"address" form:"address"`
	Gender  string `json:"gender" form:"gender"`
}
