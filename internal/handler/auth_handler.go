package handler

import (
	"encoding/json"
	"net/http"

	"github.com/example/todo-api/internal/domain"
	"github.com/example/todo-api/internal/service"
	"github.com/example/todo-api/pkg/logger"
	"github.com/example/todo-api/pkg/response"
	"github.com/go-chi/chi/v5"
)

// AuthHandler handles HTTP requests for authentication
type AuthHandler struct {
	authService service.AuthService
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Routes returns the router with all auth routes registered
func (h *AuthHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Post("/register", h.Register)
	r.Post("/login", h.Login)
	r.Post("/logout", h.Logout)
	return r
}

// Register handles user registration
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger.Info("AUTH_HANDLER", "Register", "Request received")

	var input domain.RegisterInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		logger.Warn("AUTH_HANDLER", "Register", "Invalid JSON: "+err.Error())
		response.BadRequest(w, "INVALID_JSON", "Invalid JSON payload")
		return
	}
	defer r.Body.Close()

	user, err := h.authService.Register(ctx, input)
	if err != nil {
		logger.Error("AUTH_HANDLER", "Register", "Service error: "+err.Error())
		response.FromDomainError(w, err)
		return
	}

	logger.Info("AUTH_HANDLER", "Register", "Successfully registered: "+user.ID.String())
	response.Created(w, user.SafeUser())
}

// Login handles user authentication
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger.Info("AUTH_HANDLER", "Login", "Request received")

	var input domain.LoginInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		logger.Warn("AUTH_HANDLER", "Login", "Invalid JSON: "+err.Error())
		response.BadRequest(w, "INVALID_JSON", "Invalid JSON payload")
		return
	}
	defer r.Body.Close()

	token, user, err := h.authService.Login(ctx, input)
	if err != nil {
		logger.Error("AUTH_HANDLER", "Login", "Service error: "+err.Error())
		response.FromDomainError(w, err)
		return
	}

	logger.Info("AUTH_HANDLER", "Login", "Successful login: "+user.ID.String())
	response.OK(w, map[string]interface{}{
		"token": token,
		"user":  user.SafeUser(),
	})
}

// Logout handles user logout (client should delete the token)
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	logger.Info("AUTH_HANDLER", "Logout", "Logout request received")

	// Note: With JWT, actual logout happens client-side by deleting the token
	// Server-side we could implement token blacklisting here if needed

	response.NewSuccessWithMessage("Logout successful. Please delete your token on the client side.", nil).Write(w, http.StatusOK)
}
