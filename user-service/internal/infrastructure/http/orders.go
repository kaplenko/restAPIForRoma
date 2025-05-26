package http

import (
	"encoding/json"
	"net/http"
	"user-service/pkg/errWrap"
)

type OrderRequest struct {
	OrderNumber string `json:"order_number"`
}

// @Summary Add order
// @Description Creates a new order for an authorized user
// @Tags orders
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Authorization token (Bearer)" default(Bearer <token>)
// @Param input body OrderRequest true "Order data"
// @Success 202 "new order number has been accepted for processing"
// @Success 200 {object} errWrap.ErrorResponse "the order number has already been uploaded by this user"
// @Failure 400 {object} errWrap.ErrorResponse "wrong request format"
// @Failure 401 {object} errWrap.ErrorResponse "user is not authenticated"
// @Failure 409 {object} errWrap.ErrorResponse "the order number has already been uploaded by another user"
// @Failure 422 {object} errWrap.ErrorResponse "incorrect order number format"
// @Failure 500 {object} errWrap.ErrorResponse "internal server error"
// @Router /api/user/orders [post]
func (h *Handler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, ok := ctx.Value("user_id").(int64)
	if !ok {
		h.log.Error(ctx, "unauthorized", "method", "CreateOrder")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req OrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.OrderNumber == "" {
		h.log.Error(ctx, "json decode failed", "err", err)
		errWrap.HandleError(w, err)
		return
	}
	err := h.orderService.CreateOrder(ctx, userID, req.OrderNumber)
	if err != nil {
		h.log.Error(ctx, "create order failed", "err", err)
		errWrap.HandleError(w, err)
		return
	}

	h.log.Info(ctx, "create order succeed", "orderNumber", req.OrderNumber)

	w.WriteHeader(http.StatusAccepted)
}

// @Summary Get order
// @Description Creates a new order for an authorized user
// @Tags orders
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Authorization token (Bearer)" default(Bearer <token>)
// @Success 200 {array} entity.Order "List of user orders"
// @Success 204 {object} errWrap.ErrorResponse "no data to answer"
// @Failure 401 {object} errWrap.ErrorResponse "user is not authorized"
// @Failure 500 {object} errWrap.ErrorResponse "internal server error"
// @Router /api/user/orders [get]
func (h *Handler) GetOrders(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, ok := ctx.Value("user_id").(int64)
	if !ok {
		h.log.Error(ctx, "unauthorized", "method", "GetOrders")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	orders, err := h.orderService.OrdersByUser(ctx, userID)
	if err != nil {
		h.log.Error(ctx, "get orders failed", "err", err)
		errWrap.HandleError(w, err)
		return
	}

	h.log.Info(ctx, "get orders succeed", "orders", orders)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(orders); err != nil {
		h.log.Error(ctx, "json encode failed", "err", err)
		errWrap.HandleError(w, err)
	}
}
