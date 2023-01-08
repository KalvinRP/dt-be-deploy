package handlers

import (
	models "dewetour/1models"
	repositories "dewetour/4repositories"
	dto "dewetour/5dto/result"
	tripsdto "dewetour/5dto/trips"
	"encoding/json"
	"net/http"
	"os"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v4"

	"github.com/gorilla/mux"
)

// var path_file = "http://localhost:5000/uploads/"

type handlerTrips struct {
	TripsRepository repositories.TripsRepository
}

func HandlerTrips(TripsRepository repositories.TripsRepository) *handlerTrips {
	return &handlerTrips{TripsRepository}
}

func convertResponseTrips(u models.Trips) models.Trips {
	return models.Trips{
		ID:             u.ID,
		Name:           u.Name,
		Desc:           u.Desc,
		Price:          u.Price,
		Accomodation:   u.Accomodation,
		Transportation: u.Transportation,
		Eat:            u.Eat,
		DateTrip:       u.DateTrip,
		Quota:          u.Quota,
		Day:            u.Day,
		Night:          u.Night,
		Image:          u.Image,
		CountryID:      u.CountryID,
		Country:        u.Country,
	}
}

func (h *handlerTrips) MakeTrips(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userInfo := r.Context().Value("userInfo").(jwt.MapClaims)
	userId := int(userInfo["id"].(float64))

	images := r.Context().Value("dataFile")
	filename := images.(string)

	// request := new(tripsdto.TripsRequest)
	// if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
	// 	json.NewEncoder(w).Encode(response)
	// 	return
	// }

	price, _ := strconv.Atoi(r.FormValue("price"))
	quota, _ := strconv.Atoi(r.FormValue("quota"))
	day, _ := strconv.Atoi(r.FormValue("day"))
	night, _ := strconv.Atoi(r.FormValue("night"))
	countryid, _ := strconv.Atoi(r.FormValue("country_id"))
	request := tripsdto.TripsRequest{
		Name:           r.FormValue("name"),
		Desc:           r.FormValue("desc"),
		Accomodation:   r.FormValue("accomodation"),
		Transportation: r.FormValue("transport"),
		Eat:            r.FormValue("eat"),
		DateTrip:       r.FormValue("datetrip"),
		Price:          price,
		Quota:          quota,
		Day:            day,
		Night:          night,
		CountryID:      countryid,
	}

	validation := validator.New()
	err := validation.Struct(request)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	trips := models.Trips{
		Name:           request.Name,
		Desc:           request.Desc,
		Price:          request.Price,
		Accomodation:   request.Accomodation,
		Transportation: request.Transportation,
		Eat:            request.Eat,
		DateTrip:       request.DateTrip,
		Quota:          request.Quota,
		Day:            request.Day,
		Night:          request.Night,
		Image:          filename,
		CountryID:      request.CountryID,
		UserID:         userId,
	}

	trips, err = h.TripsRepository.MakeTrips(trips)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	// trips, _ = h.TripsRepository.GetTrips(trips.ID)

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: convertResponseTrips(trips)}
	json.NewEncoder(w).Encode(response)
}

func (h *handlerTrips) FindTrips(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	trips, err := h.TripsRepository.FindTrips()
	for i, p := range trips {
		trips[i].Image = os.Getenv("PATH_FILE") + p.Image
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
	}

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: trips}
	json.NewEncoder(w).Encode(response)
}

func (h *handlerTrips) GetTrips(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	trips, err := h.TripsRepository.GetTrips(id)
	trips.Image = os.Getenv("PATH_FILE") + trips.Image
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: convertResponseTrips(trips)}
	json.NewEncoder(w).Encode(response)
}

func (h *handlerTrips) EditTrips(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userInfo := r.Context().Value("userInfo").(jwt.MapClaims)
	userId := int(userInfo["id"].(float64))

	images := r.Context().Value("dataFile")
	filename := images.(string)

	// request := new(tripsdto.TripsRequest)
	// if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
	// 	json.NewEncoder(w).Encode(response)
	// 	return
	// }

	price, _ := strconv.Atoi(r.FormValue("price"))
	quota, _ := strconv.Atoi(r.FormValue("quota"))
	day, _ := strconv.Atoi(r.FormValue("day"))
	night, _ := strconv.Atoi(r.FormValue("night"))
	countryid, _ := strconv.Atoi(r.FormValue("country_id"))
	request := tripsdto.TripsRequest{
		Name:           r.FormValue("name"),
		Desc:           r.FormValue("desc"),
		Accomodation:   r.FormValue("accomodation"),
		Transportation: r.FormValue("transport"),
		Eat:            r.FormValue("eat"),
		DateTrip:       r.FormValue("datetrip"),
		Price:          price,
		Quota:          quota,
		Day:            day,
		Night:          night,
		CountryID:      countryid,
		Image:          filename,
	}

	ID, _ := strconv.Atoi(mux.Vars(r)["id"])

	trips, err := h.TripsRepository.GetTrips(ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
	}

	// trips := models.Trips{}

	if request.Name != "" {
		trips.Name = request.Name
	}

	if request.Desc != "" {
		trips.Desc = request.Desc
	}

	if request.Accomodation != "" {
		trips.Accomodation = request.Accomodation
	}

	if request.Transportation != "" {
		trips.Transportation = request.Transportation
	}

	if request.Eat != "" {
		trips.Eat = request.Eat
	}

	if request.DateTrip != "" {
		trips.DateTrip = request.DateTrip
	}

	if request.Quota != 0 {
		trips.Quota = request.Quota
	}

	if request.Image != "" {
		trips.Image = request.Image
	}

	if request.Price != 0 {
		trips.Price = request.Price
	}

	if request.CountryID != 0 {
		trips.CountryID = request.CountryID
	}

	trips.UserID = userId

	data, err := h.TripsRepository.EditTrips(trips, ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	data, _ = h.TripsRepository.GetTrips(data.ID)

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: convertResponseTrips(data)}
	json.NewEncoder(w).Encode(response)
}

func (h *handlerTrips) DeleteTrips(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userInfo := r.Context().Value("userInfo").(jwt.MapClaims)
	userId := int(userInfo["id"].(float64))

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	user, err := h.TripsRepository.GetTrips(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	user.UserID = userId

	data, err := h.TripsRepository.DeleteTrips(user, id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: convertResponseTrips(data)}
	json.NewEncoder(w).Encode(response)
}
