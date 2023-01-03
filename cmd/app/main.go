package main

import (
	"context"
	defaultLog "log"
	"os/signal"
	"syscall"

	_ "github.com/e1leet/simple-auth-service/docs"
	"github.com/e1leet/simple-auth-service/internal/app"
	"github.com/e1leet/simple-auth-service/internal/config"
	"github.com/e1leet/simple-auth-service/internal/logging"
	"github.com/e1leet/simple-auth-service/internal/utils"
	"github.com/rs/zerolog/log"
)

//	@title						Simple auth service
//	@version					1.0.0
//	@description				Simple auth service API documentation
//	@contact.name				Damir Mirasov
//	@contact.url				https://github.com/e1leet
//	@contact.email				damirmirasovmain@gmail.com
//	@license.name				MIT
//	@license.url				https://opensource.org/licenses/MIT
//	@BasePath					/api
//	@securityDefinitions.apikey	ApiKeyAuth
//	@in							header
//	@name						Authorization

func main() {
	cfgPath := utils.ConfigPath()

	cfg, err := config.LoadConfig(cfgPath)
	if err != nil {
		defaultLog.Fatal(err)
	}

	logging.ConfigureLogging(cfg.Log.Level)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	a := app.New(cfg)

	if err := a.Run(ctx); err != nil {
		stop()
		log.Fatal().Err(err).Send() //nolint:gocritic
	}
}
