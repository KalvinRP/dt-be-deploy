package usersdto

type UserResponse struct {
	ID       int    `json:"id"`
	Name     string `json:"name" form:"name"`
	Email    string `json:"email" form:"email"`
	Password string `json:"password" form:"password"`
	Phone    string `json:"phone" form:"phone"`
	Address  string `json:"address" form:"address"`
	Gender   string `json:"gender" form:"gender"`
	Image    string `json:"image" form:"image"`
}
