package http

import (
	"encoding/json"
	"net/http"
	"user-service/internal/entity"
	"user-service/pkg/errWrap"
)

// @Summary registration
// @Description registers the user
// @Tags auth
// @Accept json
// @Param request body userDTO true "Registration data"
// @Success 200 "user is successfully registered and authenticated"
// @Failure 400 {object} errWrap.ErrorResponse "wrong request format"
// @Failure 409 {object} errWrap.ErrorResponse "login is already occupied"
// @Failure 500 {object} errWrap.ErrorResponse "internal server error"
// @Router /api/user/register [post]
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req userDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Error(ctx, "json decode failed", "error", err)
		errWrap.HandleError(w, err)
		return
	}
	user := entity.User{
		Username: req.Username,
		PassHash: []byte(req.Password),
	}

	userID, err := h.userService.Registre(ctx, user)
	if err != nil {
		h.log.Error(ctx, "register failed", "error", err)
		errWrap.HandleError(w, err)
		return
	}

	h.log.Info(ctx, "register success", "userID", userID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

// @Summary login
// @Description login the user
// @Tags auth
// @Param request body userDTO true "Login data"
// @Success 200 "user successfully authenticated"
// @Failure 400 {object} errWrap.ErrorResponse "invalid request format"
// @Failure 401 {object} errWrap.ErrorResponse "invalid login/password pair"
// @Failure 500 {object} errWrap.ErrorResponse "internal server error"
// @Router /api/user/login [post]
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req *userDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Error(ctx, "json decode failed", "error", err)
		errWrap.HandleError(w, err)
		return
	}

	h.log.Info(ctx, "login attempt", "username", req.Username)

	token, err := h.userService.Login(ctx, req.Username, []byte(req.Password))
	if err != nil {
		h.log.Error(ctx, "login failed", "username", req.Username, "error", err)
		errWrap.HandleError(w, err)
		return
	}

	h.log.Info(ctx, "login success", "username", req.Username)

	w.Header().Set("Authorization", "Bearer "+token)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
