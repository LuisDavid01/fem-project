package api

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/LuisDavid01/femProject/internal/store"
	"github.com/LuisDavid01/femProject/internal/tokens"
	"github.com/LuisDavid01/femProject/internal/utils"
)

type TokenHandler struct {
	tokenStore store.TokenStore
	userStore  store.UserStore
	logger     *log.Logger
}
type createTokenRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func NewTokenHandler(tokenStore store.TokenStore, userStore store.UserStore, logger *log.Logger) *TokenHandler {
	return &TokenHandler{
		tokenStore: tokenStore,
		userStore:  userStore,
		logger:     logger,
	}
}

func (h *TokenHandler) HandleCreateToken(w http.ResponseWriter, r *http.Request) {
	var req createTokenRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		h.logger.Println("Error decoding request body:", err)
		utils.WriteJson(w, http.StatusBadRequest, utils.Envelope{"error": "Invalid request body"})
		return
	}

	user, err := h.userStore.GetUserByUsername(req.Username)
	if err != nil || user == nil {
		h.logger.Printf("User not found: %v", req.Username)
		utils.WriteJson(w, http.StatusUnauthorized, utils.Envelope{"error": "Invalid username or password"})
		return
	}
	passwordDoMatch, err := user.PasswordHash.Matches(req.Password)
	if err != nil {
		h.logger.Printf("Error matching the password %v:", err)
		utils.WriteJson(w, http.StatusInternalServerError, utils.Envelope{"error": "Internal server error"})
		return
	}
	if !passwordDoMatch {

		utils.WriteJson(w, http.StatusUnauthorized, utils.Envelope{"error": "Invalid username or password"})
		return
	}

	token, err := h.tokenStore.CreateNewToken(user.ID, 24*time.Hour, tokens.ScopeAuth)
	if err != nil {
		h.logger.Printf("ERROR creating the token: %v", err)
		utils.WriteJson(w, http.StatusInternalServerError, utils.Envelope{"error": "Internal server error"})
		return
	}
	utils.WriteJson(w, http.StatusOK, utils.Envelope{"auth_token": token})

}
