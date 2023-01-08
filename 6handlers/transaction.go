package handlers

import (
	models "dewetour/1models"
	repositories "dewetour/4repositories"
	dto "dewetour/5dto/result"
	transactiondto "dewetour/5dto/transaction"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
	"github.com/midtrans/midtrans-go/snap"

	"gopkg.in/gomail.v2"
)

var c = coreapi.Client{
	ServerKey: os.Getenv("SERVER_KEY"),
	ClientKey: os.Getenv("CLIENT_KEY"),
}

type handlerTransaction struct {
	TransactionRepository repositories.TransactionRepository
}

func HandlerTransaction(TransactionRepository repositories.TransactionRepository) *handlerTransaction {
	return &handlerTransaction{TransactionRepository}
}

func convertTransResponse(u models.Transaction) transactiondto.TransactionResponse {
	return transactiondto.TransactionResponse{
		TripsID:   u.TripsID,
		Trips:     u.Trips,
		Users:     u.Users,
		Qty:       u.Qty,
		TotalPrc:  u.TotalPrc,
		Status:    u.Status,
		CreatedAt: u.CreatedAt,
	}
}

func (h *handlerTransaction) FindTrans(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	trans, err := h.TransactionRepository.FindTrans()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
	}

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: trans}
	json.NewEncoder(w).Encode(response)
}

func (h *handlerTransaction) GetTrans(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	trans, err := h.TransactionRepository.GetTrans(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: convertTransResponse(trans)}
	json.NewEncoder(w).Encode(response)
}

func (h *handlerTransaction) MakeTrans(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")

	userInfo := r.Context().Value("userInfo").(jwt.MapClaims)
	userId := int(userInfo["id"].(float64))

	request := new(transactiondto.TransactionRequest)
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

	var TransIdIsMatch = false
	var TransactionId int
	for !TransIdIsMatch {
		TransactionId = userId + int(time.Now().Unix()) + request.TripsID
		transactionData, _ := h.TransactionRepository.GetTrans(TransactionId)
		if transactionData.ID == 0 {
			TransIdIsMatch = true
		}
	}

	trans := models.Transaction{
		ID:       TransactionId,
		TripsID:  request.TripsID,
		UsersID:  userId,
		Qty:      request.Qty,
		TotalPrc: request.TotalPrc,
		Status:   request.Status,
	}

	transaction, err := h.TransactionRepository.MakeTrans(trans)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err.Error())
	}

	transaction, _ = h.TransactionRepository.GetTrans(transaction.ID)

	var s = snap.Client{}
	s.New("SB-Mid-server-sLvepaczAK9qycKAb61XruM3", midtrans.Sandbox)

	req := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  strconv.Itoa(transaction.ID),
			GrossAmt: int64(transaction.TotalPrc),
		},
		CreditCard: &snap.CreditCardDetails{
			Secure: true,
		},
		CustomerDetail: &midtrans.CustomerDetails{
			FName: transaction.Users.Name,
			Email: transaction.Users.Email,
		},
	}

	snapResp, _ := s.CreateTransaction(req)

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: snapResp}
	json.NewEncoder(w).Encode(response)

	// w.WriteHeader(http.StatusOK)
	// response := dto.SuccessResult{Code: http.StatusOK, Data: convertTransResponse(transaction)}
	// json.NewEncoder(w).Encode(response)
}

func (h *handlerTransaction) Notification(w http.ResponseWriter, r *http.Request) {
	var notificationPayload map[string]interface{}

	err := json.NewDecoder(r.Body).Decode(&notificationPayload)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	transactionStatus := notificationPayload["transaction_status"].(string)
	fraudStatus := notificationPayload["fraud_status"].(string)
	orderId := notificationPayload["order_id"].(string)

	transaction, _ := h.TransactionRepository.GetOneTrans(orderId)
	fmt.Println(transactionStatus, fraudStatus, orderId, transaction)

	if transactionStatus == "capture" {
		if fraudStatus == "challenge" {
			// KirimEmail("pending", transaction)
			h.TransactionRepository.UpdateTrans("pending", transaction)
		} else if fraudStatus == "accept" {
			KirimEmail("Success", transaction)
			h.TransactionRepository.UpdateTrans("success", transaction)
		}
	} else if transactionStatus == "settlement" {
		KirimEmail("Success", transaction)
		h.TransactionRepository.UpdateTrans("success", transaction)
	} else if transactionStatus == "deny" {
		KirimEmail("Denied by System", transaction)
		h.TransactionRepository.UpdateTrans("failed", transaction)
	} else if transactionStatus == "cancel" || transactionStatus == "expire" {
		KirimEmail("Cancelled", transaction)
		h.TransactionRepository.UpdateTrans("failed", transaction)
	} else if transactionStatus == "pending" {
		h.TransactionRepository.UpdateTrans("pending", transaction)
	}

	w.WriteHeader(http.StatusOK)

}

func (h *handlerTransaction) EditTrans(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userInfo := r.Context().Value("userInfo").(jwt.MapClaims)
	userId := int(userInfo["id"].(float64))

	request := new(transactiondto.UpdateTransactionRequest)
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	// ID, _ := strconv.Atoi(mux.Vars(r)["id"])

	trans, err := h.TransactionRepository.GetTrans(userId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
	}

	// trans := models.Transaction{}

	if request.TripsID != 0 {
		trans.TripsID = request.TripsID
	}

	if request.TotalPrc != 0 {
		trans.TotalPrc = request.TotalPrc
	}

	if request.Status != "" {
		trans.TotalPrc = request.TotalPrc
	}

	trans.UsersID = userId

	data, err := h.TransactionRepository.EditTrans(trans, userId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	data, _ = h.TransactionRepository.GetTrans(data.ID)

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: convertTransResponse(data)}
	json.NewEncoder(w).Encode(response)
}

func (h *handlerTransaction) DeleteTrans(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	trans, err := h.TransactionRepository.GetTrans(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	data, err := h.TransactionRepository.DeleteTrans(trans, id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: convertTransResponse(data)}
	json.NewEncoder(w).Encode(response)
}

func (h *handlerTransaction) UserHistory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userInfo := r.Context().Value("userInfo").(jwt.MapClaims)
	userId := int(userInfo["id"].(float64))

	trans, err := h.TransactionRepository.UserHistory(userId)
	for i, p := range trans {
		trans[i].Users.Image = os.Getenv("PATH_FILE") + p.Users.Image
	}

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: trans}
	json.NewEncoder(w).Encode(response)
}

func KirimEmail(status string, transaction models.Transaction) {
	var CONFIG_SMTP_HOST = "smtp.gmail.com"
	var CONFIG_SMTP_PORT = 587
	var CONFIG_SENDER_NAME = os.Getenv("SENDER_NAME")
	var CONFIG_AUTH_EMAIL = os.Getenv("EMAIL_SYSTEM")
	var CONFIG_AUTH_PASSWORD = os.Getenv("PASSWORD_SYSTEM")

	var productName = transaction.Trips.Name
	var price = strconv.Itoa(transaction.TotalPrc)
	var desc = "Transaction success. Enjoy your trip and please be safe!"
	if status == "Denied by System" {
		desc = "Transaction denied. Please make sure to do payment properly."
	} else if status == "Cancelled" {
		desc = "Transaction is cancelled."
	}

	mailer := gomail.NewMessage()
	mailer.SetHeader("From", CONFIG_SENDER_NAME)
	mailer.SetHeader("To", transaction.Users.Email)
	mailer.SetHeader("Subject", "Transaction Status Update")
	mailer.SetBody("text/html", fmt.Sprintf(`<!DOCTYPE html>
    <html lang="en">
      <head>
      <meta charset="UTF-8" />
      <meta http-equiv="X-UA-Compatible" content="IE=edge" />
      <meta name="viewport" content="width=device-width, initial-scale=1.0" />
      <title>Document</title>
      <style>
        h1 {
        color: orange;
        }
      </style>
      </head>
      <body>
      <h2>Product payment :</h2>
      <ul style="list-style-type:none;">
        <li>Name : %s</li>
        <li>Total payment: Rp.%s</li>
        <li>Status : <b>%s</b></li>
		<li>%s</li>
      </ul>
      </body>
    </html>`, productName, price, status, desc))

	dialer := gomail.NewDialer(
		CONFIG_SMTP_HOST,
		CONFIG_SMTP_PORT,
		CONFIG_AUTH_EMAIL,
		CONFIG_AUTH_PASSWORD,
	)

	err := dialer.DialAndSend(mailer)
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Println("Pesan terkirim!")
}
