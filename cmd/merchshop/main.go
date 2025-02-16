// Package main = entry point.
package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kk7453603/avito_2024_summer/internal/config"
	"github.com/kk7453603/avito_2024_summer/internal/db"
	"github.com/kk7453603/avito_2024_summer/internal/hasher"
	"github.com/kk7453603/avito_2024_summer/internal/logger"
	"github.com/kk7453603/avito_2024_summer/internal/modules/authentication"
	"github.com/kk7453603/avito_2024_summer/internal/modules/buy_item"
	"github.com/kk7453603/avito_2024_summer/internal/modules/jwt_token_manager"
	"github.com/kk7453603/avito_2024_summer/internal/modules/transaction"
	"github.com/kk7453603/avito_2024_summer/internal/modules/user_info"
	"github.com/kk7453603/avito_2024_summer/internal/server"
	"github.com/kk7453603/avito_2024_summer/internal/server/handlers"
)

var version = "v1.2.2"

func main() {
	cfg := config.MustLoad()
	logg := logger.Init(cfg.Log)
	logg.Info("Application loading...")

	ctx, ctxCancel := context.WithCancel(context.Background())

	storage, err := db.NewPostgresPool(ctx, cfg.DB)
	if err != nil {
		logg.Error("db.NewPostgresPool", "err", err.Error())
		os.Exit(1)
	}

	passwdHasher := hasher.New()
	tknMng, err := jwt_token_manager.New(cfg.JWT)
	if err != nil {
		logg.Error("jwt_token_manager.New", "err", err.Error())
		os.Exit(1)
	}
	authSrv := authentication.New(storage, passwdHasher)
	usrInfSrv := user_info.New(storage) // creating a user information module
	txSrv := transaction.New(storage)   // transaction module creation
	buyItmSrv := buy_item.New(storage)  // creating an item purchase module

	// creating the main request handler
	usrHandlers := handlers.NewUserHandlers(ctx, authSrv, tknMng, usrInfSrv, txSrv, buyItmSrv)
	// server creation
	serv := server.New(ctx, cfg.APIServer, usrHandlers, tknMng)

	// server startup
	go func() {
		logg.Info("Application Started! " + version)
		if err = serv.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logg.Error("serv.Start", "err", err)
			os.Exit(1)
		}
	}()

	// graceful shutdown
	sigCtx, sigStop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	<-sigCtx.Done()
	sigStop()
	logg.Info("Shutting down gracefully...")

	ctxTimeOut, ctxTimeOutCancel := context.WithTimeout(ctx, 5*time.Second)
	defer ctxTimeOutCancel()

	if err = serv.Shutdown(ctxTimeOut); err != nil {
		logg.Error("serv.Shutdown", "err", err.Error())
	}

	if storage != nil {
		storage.Close()
	}

	ctxCancel()
	logg.Info("Application Stopped!")
}
