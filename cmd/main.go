package main

import (
	"context"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/kk7453603/avito_2024_summer/docs"
	"github.com/kk7453603/avito_2024_summer/internal/delivery"
	"github.com/kk7453603/avito_2024_summer/internal/repository"
	"github.com/kk7453603/avito_2024_summer/internal/service"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	echoSwagger "github.com/swaggo/echo-swagger"
)

//	@title			Merch Store API
//	@version		1.0
//	@description	This is a merch store server.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	@lettons
//	@contact.url	https://t.me/lettons
//	@contact.email	kk7453603@gmail.com

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

// @host		0.0.0.0
// @BasePath	/
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("godotenv error: %v", err)
	}
	e := echo.New()
	e.Logger.Info("Переменные среды загружены")
	if os.Getenv("DEBUG") == "on" {
		e.Debug = true
		e.Logger.SetLevel(log.DEBUG)
		e.Logger.Info("DEBUG режим включен")
	}

	sql_handler := repository.New(context.Background(), e.Logger)
	sql_handler.Migrate()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Swagger documentation route
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	g := e.Group("")
	serv := service.New(sql_handler, e.Logger)
	deliv := delivery.New(serv, e.Logger)
	deliv.InitRoutes(g)

	e.Logger.Infof("service starts on: %s", e.Server.Addr)

	if err = e.Start(os.Getenv("Service_Url")); err != nil {
		e.Logger.Fatalf("Ошибка запуска сервера: %v", err)
	}
	e.Logger.Info("service start")
}
