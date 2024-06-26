package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"gitlab.teamdev.huds.su/bivi/backend/internal/pkg/app"
	"gitlab.teamdev.huds.su/bivi/backend/internal/stream/repository"
	"gitlab.teamdev.huds.su/bivi/backend/internal/stream/usecase"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

type WebApp interface {
	Start(port string) error
	Shutdown(ctx context.Context) error
}

func readConfig() error {
	var configPath string

	pflag.StringVarP(&configPath, "config", "c", "env/app.local.yaml", "Config file path")
	pflag.Parse()

	if configPath == "" {
		return errors.New("config file is not specified")
	}

	slog.Info("Config path: " + configPath)
	viper.SetConfigFile(configPath)

	return errors.Wrap(viper.ReadInConfig(), "read configuration")
}

func runApp(webApp WebApp, port string, logger *slog.Logger) {
	logger.Debug(fmt.Sprintf("web app starts at port %s with configuration: \n%v", port, viper.AllSettings()))

	go func() {
		err := webApp.Start(port)
		if err != nil {
			panic(err)
		}
	}()
}

func shutdownApp(webApp WebApp, logger *slog.Logger) {
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be caught, so don't need to add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Debug("shutdown web app ...")

	err := webApp.Shutdown(context.Background())
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
//	@BasePath		/api/v1
//	@Schemes		http
func main() {
	err := readConfig()
	if err != nil {
		panic(err)
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: viper.GetBool("debug"),
		Level:     slog.LevelDebug,
	}))

	contentRoute := viper.GetString("streams.content_route")
	streamRepository := repository.NewRepository(viper.GetString("streams.content_path"), logger)
	streamNameEncoder := base64.URLEncoding
	streamUseCase := usecase.NewUseCase(contentRoute, streamRepository, streamNameEncoder, logger)

	settings := app.APISettings{
		Prefix:              viper.GetString("api_prefix"),
		Port:                viper.GetString("port"),
		ContentRoute:        contentRoute,
		ContentPath:         viper.GetString("streams.content_path"),
		ClientLogPath:       viper.GetString("client_log_path"),
		UploadFilesizeLimit: viper.GetInt64("upload_filesize_limit"),
	}
	webApp := app.NewFiberApp(settings, streamUseCase, streamNameEncoder, logger)

	runApp(webApp, settings.Port, logger)
	shutdownApp(webApp, logger)
}
