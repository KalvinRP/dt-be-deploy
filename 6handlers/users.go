package handlers

import (
	models "dewetour/1models"
	repositories "dewetour/4repositories"
	dto "dewetour/5dto/result"
	usersdto "dewetour/5dto/users"
	"encoding/json"
	"net/http"
	"os"
	"strconv"

	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
)

type handlerUser struct {
	UserRepository repositories.UserRepository
}

func convertResponse(u models.User) usersdto.UserResponse {
	return usersdto.UserResponse{
		ID:       u.ID,
		Name:     u.Name,
		Email:    u.Email,
		Password: u.Password,
		Phone:    u.Phone,
		Address:  u.Address,
		Gender:   u.Gender,
		Image:    os.Getenv("PATH_FILE") + u.Image,
	}
}

func HandlerUser(UserRepository repositories.UserRepository) *handlerUser {
	return &handlerUser{UserRepository}
}

func (h *handlerUser) FindAcc(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	users, err := h.UserRepository.FindAcc()
	for i, p := range users {
		users[i].Image = os.Getenv("PATH_FILE") + p.Image
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
	}

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: users}
	json.NewEncoder(w).Encode(response)
}

func (h *handlerUser) GetAcc(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userInfo := r.Context().Value("userInfo").(jwt.MapClaims)
	userId := int(userInfo["id"].(float64))

	user, err := h.UserRepository.GetAcc(userId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: convertResponse(user)}
	json.NewEncoder(w).Encode(response)
}

func (h *handlerUser) EditAcc(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userInfo := r.Context().Value("userInfo").(jwt.MapClaims)
	userId := int(userInfo["id"].(float64))

	images := r.Context().Value("dataFile")
	filename := images.(string)

	// request := new(usersdto.UpdateUserRequest)
	// if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
	// 	json.NewEncoder(w).Encode(response)
	// 	return
	// }

	request := usersdto.UpdateUserRequest{
		Name:     r.FormValue("name"),
		Email:    r.FormValue("email"),
		Phone:    r.FormValue("phone"),
		Address:  r.FormValue("address"),
		Gender:   r.FormValue("gender"),
		Password: r.FormValue("password"),
		Image:    filename,
	}

	user, err := h.UserRepository.GetAcc(userId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
	}

	// user := models.User{}

	if request.Name != "" {
		user.Name = request.Name
	}

	if request.Email != "" {
		user.Email = request.Email
	}

	if request.Password != "" {
		user.Password = request.Password
	}

	if request.Phone != "" {
		user.Phone = request.Phone
	}

	if request.Address != "" {
		user.Address = request.Address
	}

	if request.Gender != "" {
		user.Gender = request.Gender
	}

	if request.Image != "" {
		user.Image = filename
	}

	data, err := h.UserRepository.EditAcc(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: convertResponse(data)}
	json.NewEncoder(w).Encode(response)
}

func (h *handlerUser) DeleteAcc(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	user, err := h.UserRepository.GetAcc(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	data, err := h.UserRepository.DeleteAcc(user, id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: convertResponse(data)}
	json.NewEncoder(w).Encode(response)
}
