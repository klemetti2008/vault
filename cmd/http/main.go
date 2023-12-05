package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"gitag.ir/cookthepot/services/vault/database"
	"gitag.ir/cookthepot/services/vault/notification"
	"gitag.ir/cookthepot/services/vault/validity"

	"gitag.ir/cookthepot/services/vault/config"
	"gitag.ir/cookthepot/services/vault/services/account"
	"gitag.ir/cookthepot/services/vault/services/category"
	"gitag.ir/cookthepot/services/vault/services/healthcheck"
	"gitag.ir/cookthepot/services/vault/services/permission"
	"gitag.ir/cookthepot/services/vault/services/role"
	"gitag.ir/cookthepot/services/vault/services/user"
	"gitag.ir/cookthepot/services/vault/services/verify"
	"gitag.ir/cookthepot/services/vault/services/welcome"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echohandlers "github.com/mhosseintaher/kit/echo"
	"github.com/mhosseintaher/kit/log"
)

func main() {

	flag.Parse()

	config.Load()

	ctx := context.Background()

	validity.ApplyTranslations()

	var (
		// APP
		URL     = config.AppConfig.AppUrl
		PORT    = strconv.Itoa(config.AppConfig.Port)
		VERSION = config.AppConfig.Version
		// JWT
		AccessTokenSigningKey      = config.AppConfig.AccessTokenSigningKey
		AccessTokenTokenExpiration = config.AppConfig.AccessTokenTokenExpiration
		RefreshTokenSigningKey     = config.AppConfig.RefreshTokenSigningKey
		RefreshTokenExpiration     = config.AppConfig.RefreshTokenExpiration
		// Database
	)
	logger := log.New().With(ctx, "version", VERSION)

	db := database.Connect()

	e := echo.New()
	e.Use(middleware.CORS())
	e.Use(
		middleware.LoggerWithConfig(
			middleware.LoggerConfig{
				Format: `[${time_rfc3339}] [${latency_human}] [${remote_ip}:${method}] [${status}]   err=[${error}]` + "\n",
			},
		),
	)
	e.Use(middleware.Recover())
	e.Use(middleware.Gzip())
	e.Use(middleware.BodyLimit("2M"))

	e.HTTPErrorHandler = echohandlers.CustomHandler

	healthcheck.RegisterHandlers(e, config.AppConfig.Version)

	welcome.RegisterHandlers(e)

	notifier := notification.MakeNotifier()

	account.Register(e, db, logger, notifier, AccessTokenSigningKey, RefreshTokenSigningKey, AccessTokenTokenExpiration, RefreshTokenExpiration, "v1")

	user.Register(e, db, logger, "v1")

	verify.Register(e, db, notifier, logger, "v1")

	category.Register(e, db, logger, "/v1")

	role.Register(e, db, logger, "v1")

	permission.Register(e, db, logger, "v1")

	msg := make(chan error)

	go func() {
		_ = fmt.Sprintf(
			"listening on http://%s:%s ",
			URL, PORT,
		)
		msg <- e.Start(URL + ":" + PORT)
	}()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		msg <- fmt.Errorf("%s", <-c)
	}()
	logger.Error("chan err : ", <-msg)
	os.Exit(1)
}
