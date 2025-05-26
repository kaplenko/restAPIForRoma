package http

import (
	"github.com/gorilla/mux"
	"net/http"
	"strings"
	"user-service/internal/entity"
	"user-service/internal/usecase"
	"user-service/pkg/jwt"
)

type Handler struct {
	userService    *usecase.UserService
	orderService   *usecase.OrderService
	balanceService *usecase.BalanceService
	router         *mux.Router
	log            entity.Logger
}

type userDTO struct {
	Username string `json:"login"`
	Password string `json:"password"`
}

func New(usecase *usecase.UserService, orderService *usecase.OrderService, balanceService *usecase.BalanceService, r *mux.Router, log entity.Logger) *Handler {
	return &Handler{
		userService:    usecase,
		orderService:   orderService,
		balanceService: balanceService,
		router:         r,
		log:            log,
	}
}

func (h *Handler) Router() *mux.Router {
	return h.router
}

func (h *Handler) SetupRoutes() {
	h.router.Use(Cors)
	public := h.router.PathPrefix("/api/user").Subrouter()

	public.HandleFunc("/register", h.Register).Methods(http.MethodPost)
	public.HandleFunc("/login", h.Login).Methods(http.MethodPost)

	private := h.router.PathPrefix("/api/user").Subrouter()
	private.Use(jwt.JWTMiddleware)

	private.HandleFunc("/orders", h.CreateOrder).Methods(http.MethodPost)
	private.HandleFunc("/orders", h.GetOrders).Methods(http.MethodGet)

	private.HandleFunc("/balance", h.GetBalance).Methods(http.MethodGet)
	private.HandleFunc("/balance/withdraw", h.WithdrawBalance).Methods(http.MethodPost)
	private.HandleFunc("/withdrawals", h.Withdrawals).Methods(http.MethodGet)
}

func Cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/api/user/balance") {
			next.ServeHTTP(w, r)
			return
		}
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}
