package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"golang.org/x/sync/errgroup"
	"net/http"
	"os/signal"
	"syscall"
	"time"
	"user-service/internal/entity"
	httpApp "user-service/internal/infrastructure/http"
	"user-service/internal/infrastructure/storage"
	"user-service/internal/usecase"
	"user-service/pkg/mock"
)

type App struct {
	httpApp *httpApp.Handler
	storage *storage.Storage
	log     entity.Logger
}

func New(log entity.Logger, connStr string, tokenTTL time.Duration) *App {
	strg, err := storage.New(connStr, log)
	if err != nil {
		panic(err)
	}

	accrualService := mock.NewAccrualService()
	authService := usecase.NewUserService(strg, strg, log, tokenTTL)
	orderService := usecase.NewOrderService(strg, accrualService, log)
	balanceService := usecase.NewBalanceService(strg, log)

	r := mux.NewRouter()

	httpService := httpApp.New(authService, orderService, balanceService, r, log)

	return &App{
		httpApp: httpService,
		storage: strg,
		log:     log,
	}
}

func (app *App) Run(host, port string) {
	app.httpApp.SetupRoutes()
	conn := fmt.Sprintf("%s:%s", host, port)
	app.log.Info(context.Background(), "Server started at "+conn)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	g, ctx := errgroup.WithContext(ctx)

	server := &http.Server{
		Addr:    conn,
		Handler: app.httpApp.Router(),
	}

	g.Go(func() error {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}
		return nil
	})

	g.Go(func() error {
		<-ctx.Done()
		app.log.Info(context.Background(), "Shutting down server...")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			return err
		}

		if err := app.storage.Pool().Close; err != nil {
			app.log.Error(context.Background(), "Failed to close storage connection", "error", err)
		}

		return nil
	})

	if err := g.Wait(); err != nil {
		app.log.Error(context.Background(), "Server shutdown with error", "error", err)
		return
	}

	app.log.Info(context.Background(), "Server shutdown gracefully")
}
