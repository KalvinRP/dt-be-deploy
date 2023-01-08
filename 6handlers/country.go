package handlers

import (
	models "dewetour/1models"
	repositories "dewetour/4repositories"
	countrydto "dewetour/5dto/country"
	dto "dewetour/5dto/result"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type handlerCountry struct {
	CountryRepository repositories.CountryRepository
}

func HandlerCountry(CountryRepository repositories.CountryRepository) *handlerCountry {
	return &handlerCountry{CountryRepository}
}

func (h *handlerCountry) FindCountry(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	country, err := h.CountryRepository.FindCountry()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
	}

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: country}
	json.NewEncoder(w).Encode(response)
}

func (h *handlerCountry) GetCountry(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	country, err := h.CountryRepository.GetCountry(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: convertCountryResponse(country)}
	json.NewEncoder(w).Encode(response)
}

func convertCountryResponse(u models.Country) countrydto.CountryResponse {
	return countrydto.CountryResponse{
		Name: u.Name,
	}
}

func (h *handlerCountry) MakeCountry(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")

	request := new(countrydto.CountryRequest)
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	validation := validator.New()
	err := validation.Struct(request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	country := models.Country{
		Name: request.Name,
	}

	data, err := h.CountryRepository.MakeCountry(country)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err.Error())
	}

	data, _ = h.CountryRepository.GetCountry(data.ID)

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: convertCountryResponse(data)}
	json.NewEncoder(w).Encode(response)
}

func (h *handlerCountry) EditCountry(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	request := new(countrydto.CountryRequest) //take pattern data submission
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	ID, _ := strconv.Atoi(mux.Vars(r)["id"])

	country, err := h.CountryRepository.GetCountry(ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
	}

	// country := models.Country{}

	if request.Name != "" {
		country.Name = request.Name
	}

	data, err := h.CountryRepository.EditCountry(country, ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	data, _ = h.CountryRepository.GetCountry(data.ID)

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: convertCountryResponse(data)}
	json.NewEncoder(w).Encode(response)
}

func (h *handlerCountry) DeleteCountry(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	country, err := h.CountryRepository.GetCountry(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	data, err := h.CountryRepository.DeleteCountry(country, id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: convertCountryResponse(data)}
	json.NewEncoder(w).Encode(response)
}
