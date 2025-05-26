package app

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
	"time"
	_ "user-service/docs"
	"user-service/internal/entity"
	httpApp "user-service/internal/infrastructure/http"
	"user-service/internal/infrastructure/storage"
	"user-service/internal/usecase"
)

type App struct {
	httpApp *httpApp.Handler
	log     entity.Logger
}

// @title User Service API
// @version 1.0
// @description API для управления пользователями, заказами и балансом
// @host localhost:8080
// @BasePath /api
// @schemes http
func New(log entity.Logger, connStr string, tokenTTL time.Duration) *App {
	strg, err := storage.New(connStr, log)
	if err != nil {
		panic(err)
	}

	authService := usecase.NewUserService(strg, strg, log, tokenTTL)
	orderService := usecase.NewOrderService(strg, log)
	balanceService := usecase.NewBalanceService(strg, log)

	r := mux.NewRouter()
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	httpApp := httpApp.New(authService, orderService, balanceService, r, log)

	return &App{
		httpApp: httpApp,
		log:     log,
	}
}

func (app *App) Run(host, port string) {
	app.httpApp.SetupRoutes()
	conn := fmt.Sprintf("%s:%s", host, port)
	app.log.Info(context.Background(), "Server started at "+conn)
	if err := http.ListenAndServe(host+":"+port, app.httpApp.Router()); err != nil {
		panic(err)
	}
}
