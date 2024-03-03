package main

import (
	"context"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"gitlab.teamdev.huds.su/bivi/backend/internal/pkg/app"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func readConfig() error {
	var configPath string
	pflag.StringVarP(&configPath, "config", "c", "", "Config file path")
	pflag.Parse()
	if configPath == "" {
		return errors.New("config file is not specified")
	}
	slog.Info(fmt.Sprintf("Config path: %s", configPath))

	viper.SetConfigFile(configPath)

	return viper.ReadInConfig()
}

func runApp(webApp app.WebApp, port string, logger *slog.Logger) {
	logger.Debug(fmt.Sprintf("web app starts at port %s with configuration: \n%v",
		port, viper.AllSettings()),
	)

	go func() {
		err := webApp.Start(port)
		if err != nil {
			panic(err)
		}
	}()
}

func shutdownApp(webApp app.WebApp, logger *slog.Logger) {
	quit := make(chan os.Signal)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be caught, so don't need to add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Debug("shutdown web app ...")

	err := webApp.Stop(context.Background())
	if err != nil {
		panic(errors.Wrap(err, "app shutdown"))
	}
	logger.Debug("web app exited")
}

//	@title			bivi API
//	@version		0.1.0
//	@description	This is bivi backend API.
//	@contact.name	API Support
//	@contact.email	andreysapozhkov535@gmail.com
//	@host			localhost:8080
//	@BasePath		/api/v1
//	@Schemes		http
func main() {
	err := readConfig()
	if err != nil {
		panic(err)
	}

	logLevel := new(slog.LevelVar)
	logLevel.Set(slog.LevelDebug)
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource:   true,
		Level:       logLevel.Level(),
		ReplaceAttr: nil,
	}))

	// setup dependencies

	settings := app.ApiSettings{
		Port:      viper.GetString("APP_PORT"),
		ApiPrefix: viper.GetString("API_PREFIX"),
	}

	webApp := app.NewFiberApp(settings, logger)

	runApp(webApp, settings.Port, logger)
	shutdownApp(webApp, logger)
}
