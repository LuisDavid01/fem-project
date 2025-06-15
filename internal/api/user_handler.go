package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"regexp"

	"github.com/LuisDavid01/femProject/internal/store"
	"github.com/LuisDavid01/femProject/internal/utils"
)

type registerUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Bio      string `json:"bio"`
}
type UserHandler struct {
	userStore store.UserStore
	logger    *log.Logger
}

func NewUserHandler(userStore store.UserStore, logger *log.Logger) *UserHandler {
	return &UserHandler{
		userStore: userStore,
		logger:    logger,
	}
}

func (h *UserHandler) validateRegisterUserRequest(req *registerUserRequest) error {
	if req.Username == "" {
		return errors.New("Err Invalid Username")
	}
	if req.Email == "" {
		return errors.New("Err Invalid Email")
	}
	if req.Password == "" {
		return errors.New("Err Invalid Password")
	}

	if len(req.Password) < 8 {
		return errors.New("Err Password must be at least 8 characters long")
	}
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(req.Email) {
		return errors.New("Err Invalid Email format")
	}
	return nil
}

func (h *UserHandler) HandlerRegisterUser(w http.ResponseWriter, r *http.Request) {
	var req registerUserRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		h.logger.Printf("Error decoding the req body: %v", err)
		utils.WriteJson(w, http.StatusBadRequest, utils.Envelope{"error": "Invalid request body"})
		return
	}

	err = h.validateRegisterUserRequest(&req)
	if err != nil {
		h.logger.Printf("Error validating the req: %v", err)
		utils.WriteJson(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}
	user := &store.User{
		Username: req.Username,
		Email:    req.Email,
	}
	if req.Bio != "" {
		user.Bio = req.Bio
	}
	//lets handle user passwords
	err = user.PasswordHash.Set(req.Password)
	if err != nil {
		h.logger.Printf("Error setting user password: %v", err)
		utils.WriteJson(w, http.StatusInternalServerError, utils.Envelope{"error": "Error setting user password"})
		return
	}
	err = h.userStore.CreateUser(user)
	if err != nil {
		h.logger.Printf("Error creating the user: %v", err)
		utils.WriteJson(w, http.StatusInternalServerError, utils.Envelope{"error": "Error creating the user"})

		return
	}
	utils.WriteJson(w, http.StatusCreated, utils.Envelope{"user": user})
}
