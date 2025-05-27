package http

import (
	"encoding/json"
	"net/http"
	"user-service/pkg/errWrap"
)

// @Summary Get Balance
// @Tags balance
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Authorization token (Bearer)" default(Bearer <token>)
// @Success 200 {object} BalanceResponse "successful processing of the request"
// @Failure 401 {object} errWrap.ErrorResponse "user is not authorized"
// @Failure 500 {object} errWrap.ErrorResponse "internal server error"
// @Router /api/user/balance [get]
func (h *Handler) GetBalance(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, ok := ctx.Value("user_id").(int64)
	if !ok {
		h.log.Error(ctx, "unauthorized", "method", "GetBalance")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	h.log.Info(ctx, "balance requested")

	balance, err := h.balanceService.Balance(ctx, userID)
	if err != nil {
		h.log.Error(ctx, "failed to get balance", "error", err)
		errWrap.HandleError(w, err)
		return
	}

	h.log.Info(ctx, "balance retrieved", "balance", balance)

	resp := BalanceResponse{
		Current:  CentsToRubles(balance.Current),
		Withdraw: CentsToRubles(balance.Withdraw),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(resp); err != nil {
		h.log.Error(ctx, "failed to encode balance", "error", err)
		errWrap.HandleError(w, err)
		return
	}
}

// @Summary Write-off request
// @Tags balance
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Authorization token (Bearer)" default(Bearer <token>)
// @Param request body WithdrawRequest true "Withdrawal data"
// @Success 200 "successful processing of the request"
// @Failure 401 {object} errWrap.ErrorResponse "user is not authorized"
// @Failure 402 {object} errWrap.ErrorResponse "insufficient funds on the account"
// @Failure 422 {object} errWrap.ErrorResponse "incorrect order number"
// @Failure 500 {object} errWrap.ErrorResponse "internal server error"
// @Router /api/user/balance/withdraw [post]
func (h *Handler) WithdrawBalance(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, ok := ctx.Value("user_id").(int64)
	if !ok {
		h.log.Error(ctx, "unauthorized", "method", "WithdrawBalance")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req WithdrawRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Error(ctx, "failed to decode withdraw request", "error", err)
		errWrap.HandleError(w, err)
		return
	}

	h.log.Info(ctx, "withdraw request accepted", "order", req.Order)

	sumInCent := RublesToCents(req.Sum)

	if err := h.balanceService.Withdraw(ctx, userID, sumInCent, req.Order); err != nil {
		h.log.Error(ctx, "failed to withdraw", "order", req.Order, "sum", req.Sum, "error", err)
		errWrap.HandleError(w, err)
		return
	}
	h.log.Info(ctx, "withdraw successful", "order", req.Order, "sum", req.Sum)
	w.WriteHeader(http.StatusOK)
}

// @Summary information on withdrawal of funds
// @Tags balance
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Authorization token (Bearer)" default(Bearer <token>)
// @Success 200 {object} []WithdrawalResponse "successful request processing"
// @Failure 204 {object} errWrap.ErrorResponse "there are no write-offs"
// @Failure 401 {object} errWrap.ErrorResponse "user is not authorized"
// @Failure 500 {object} errWrap.ErrorResponse "internal server error"
// @Router /api/user/withdrawals [get]
func (h *Handler) Withdrawals(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, ok := ctx.Value("user_id").(int64)
	if !ok {
		h.log.Error(ctx, "unauthorized", "method", "Withdrawals")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	h.log.Info(ctx, "withdraw request accepted")

	withdrawals, err := h.balanceService.Withdrawals(ctx, userID)
	if err != nil {
		h.log.Error(ctx, "failed to get withdrawals", "error", err)
		errWrap.HandleError(w, err)
		return
	}

	h.log.Info(ctx, "withdrawals retrieved", "count", len(withdrawals))

	resp := make([]WithdrawalResponse, 0, len(withdrawals))
	for _, wd := range withdrawals {
		resp = append(resp, WithdrawalResponse{
			OrderNumber: wd.OrderNumber,
			Sum:         CentsToRubles(wd.Sum),
			ProcessedAt: wd.ProcessedAt,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(resp); err != nil {
		h.log.Error(ctx, "failed to encode withdrawals", "error", err)
		errWrap.HandleError(w, err)
		return
	}
}
