// Package server contains the logic for setting up and running an HTTP server.
// Includes route handling, middleware setup, and server configuration.
package server

import (
	"github.com/kk7453603/avito_2024_summer/internal/server/middlewares"
)

// configureRouter sets up the HTTP route handlers.
func (as *APIServer) configureRouter() {
	api := as.router.Group("/api")
	{
		api.POST("/auth", as.usrHandlers.AuthHandler)

		meddlers := middlewares.NewMiddlewares(as.tknMng)
		authorized := api.Group("/", meddlers.JWTMiddleware())
		{
			authorized.GET("/info", as.usrHandlers.InfoHandler)
			authorized.POST("/sendCoin", as.usrHandlers.SendCoinsHandler)
			authorized.GET("/buy/:item", as.usrHandlers.BuyItemHandler)
		}
	}
}
