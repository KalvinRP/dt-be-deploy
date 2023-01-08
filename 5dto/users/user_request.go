package usersdto

type CreateUserRequest struct {
	Name     string `json:"name" form:"name" validate:"required"`
	Email    string `json:"email" form:"email" validate:"required"`
	Password string `json:"password" form:"password" validate:"required"`
	Phone    string `json:"phone" form:"phone" validate:"required"`
	Address  string `json:"address" form:"address" validate:"required"`
	Image    string `json:"image" form:"image"`
	// Gender   string `json:"gender" form:"gender" validate:"required"`
}

type UpdateUserRequest struct {
	Name     string `json:"name" form:"name"`
	Email    string `json:"email" form:"email"`
	Password string `json:"password" form:"password"`
	Phone    string `json:"phone" form:"phone"`
	Address  string `json:"address" form:"address"`
	Gender   string `json:"gender" form:"gender"`
	Image    string `json:"image" form:"image"`
}
