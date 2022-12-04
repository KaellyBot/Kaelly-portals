package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/kaellybot/kaelly-portals/application"
	"github.com/kaellybot/kaelly-portals/models"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

func init() {
	initConfig()
	initLog()
}

func initConfig() {
	viper.SetConfigFile(models.ConfigFileName)

	for configName, defaultValue := range models.DefaultConfigValues {
		viper.SetDefault(configName, defaultValue)
	}

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal().Err(err).Str(models.LogFileName, models.ConfigFileName).Msgf("Failed to read config, shutting down.")
	}
}

func initLog() {
	zerolog.SetGlobalLevel(models.LogLevelFallback)
	zerolog.CallerMarshalFunc = func(pc uintptr, file string, line int) string {
		short := file
		for i := len(file) - 1; i > 0; i-- {
			if file[i] == '/' {
				short = file[i+1:]
				break
			}
		}
		return fmt.Sprintf("%s:%d", short, line)
	}
	log.Logger = log.With().Caller().Logger()

	logLevel, err := zerolog.ParseLevel(viper.GetString(models.LogLevel))
	if err != nil {
		log.Warn().Err(err).Msgf("Log level not set, continue with %s...", models.LogLevelFallback)
	} else {
		zerolog.SetGlobalLevel(logLevel)
		log.Debug().Msgf("Logger level set to '%s'", logLevel)
	}
}

func main() {
	app, err := application.New(
		models.RabbitMQClientId,
		viper.GetString(models.RabbitMqAddress),
		viper.GetInt(models.HttpTimeout),
	)
	if err != nil {
		log.Fatal().Err(err).Msgf("Shutting down after failing to instantiate application")
	}

	err = app.Run()
	if err != nil {
		log.Fatal().Err(err).Msgf("Shutting down after failing to run application.")
	}

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	log.Info().Msgf("%s v%s is now running. Press CTRL-C to exit.", models.InternalName, models.Version)
	<-sc

	log.Info().Msgf("Gracefully shutting down %s...", models.InternalName)
	app.Shutdown()
}
