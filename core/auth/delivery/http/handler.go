package authHttp

import (
	authInterfaces "GOSecretProject/core/auth/interfaces"
	"GOSecretProject/core/model/base"
	"GOSecretProject/core/utils/empty_status_json"
	"GOSecretProject/core/utils/sms"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/kataras/golog"
	"net/http"
)

type Handler struct {
	repo      authInterfaces.AuthRepository
	smsSender *sms.SMS
}

func NewHandler(repo authInterfaces.AuthRepository) *Handler {
	return &Handler{
		repo:      repo,
		smsSender: sms.NewSMS(),
	}
}

func (h *Handler) Test(w http.ResponseWriter, r *http.Request) {
	golog.Infof("Test ")

	w.Write([]byte("hello"))
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var user baseModels.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		golog.Errorf("Register error: ", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(empty_status_json.JsonWithStatusCode(http.StatusBadRequest))
		return
	}
	golog.Infof("Register: %s", user.Login)

	err = h.repo.Register(user)

	if err == nil {
		w.WriteHeader(http.StatusCreated)
		w.Write(empty_status_json.JsonWithStatusCode(http.StatusCreated))
	} else {
		golog.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(empty_status_json.JsonWithStatusCode(http.StatusInternalServerError))
	}
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var user baseModels.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		golog.Errorf("Login error: ", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(empty_status_json.JsonWithStatusCode(http.StatusBadRequest))
		return
	}
	golog.Infof("Login: %s", user)

	var code int
	user.ID, user.Session, code, err = h.repo.Login(user)

	switch code {
	case 201:
		golog.Error("Logined ", user.ID)
		w.WriteHeader(http.StatusCreated)
		json, _ := json.Marshal(user)
		w.Write(json)
	case 401:
		golog.Error("401")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write(empty_status_json.JsonWithStatusCode(http.StatusUnauthorized))
	default:
		golog.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(empty_status_json.JsonWithStatusCode(http.StatusInternalServerError))
	}
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	var user baseModels.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		golog.Errorf("Logout error: ", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(empty_status_json.JsonWithStatusCode(http.StatusBadRequest))
		return
	}
	golog.Infof("Logout: %s", user)

	err = h.repo.Logout(user.Session)

	if err == nil {
		w.WriteHeader(http.StatusOK)
		w.Write(empty_status_json.JsonWithStatusCode(http.StatusOK))
	} else {
		golog.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(empty_status_json.JsonWithStatusCode(http.StatusInternalServerError))
	}
}

func (h *Handler) Confirm(w http.ResponseWriter, r *http.Request) {
	phone, _ := mux.Vars(r)["phone"]
	if !h.smsSender.ValidatePhoneNumber(phone) {
		golog.Error("Not valid phone: ", phone)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(empty_status_json.JsonWithStatusCode(http.StatusBadRequest))
	}

	//code := base64.StdEncoding.EncodeToString([]byte(time.Now().String()))[:4]
	//err := h.smsSender.SendSMS("", phone)
	code := "2222"
	var err error = nil

	if err == nil {
		w.WriteHeader(http.StatusOK)
		json, _ := json.Marshal(baseModels.CodeConfirmation{Code: code})
		w.Write(json)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(empty_status_json.JsonWithStatusCode(http.StatusInternalServerError))
	}
}

func (h *Handler) CheckSession(w http.ResponseWriter, r *http.Request) {
	var session baseModels.SessionConfirmation
	err := json.NewDecoder(r.Body).Decode(&session)
	if err != nil {
		golog.Errorf("Session error: ", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(empty_status_json.JsonWithStatusCode(http.StatusBadRequest))
		return
	}
	golog.Infof("Session: %s", session)

	var user baseModels.User
	user, err = h.repo.CheckSession(session.Session)

	if err == nil {
		w.WriteHeader(http.StatusOK)
		json, _ := json.Marshal(user)
		w.Write(json)
	} else {
		golog.Error(err)
		w.WriteHeader(http.StatusUnauthorized)
		w.Write(empty_status_json.JsonWithStatusCode(http.StatusUnauthorized))
	}
}
