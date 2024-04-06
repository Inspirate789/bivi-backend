package app_test

import (
	"bufio"
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/ozontech/allure-go/pkg/allure"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"gitlab.teamdev.huds.su/bivi/backend/internal/pkg/app"
	"gitlab.teamdev.huds.su/bivi/backend/internal/stream/delivery"
	"gitlab.teamdev.huds.su/bivi/backend/internal/stream/repository"
	"gitlab.teamdev.huds.su/bivi/backend/internal/stream/usecase"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"
)

type StreamerSuite struct {
	suite.Suite
	app                *app.FiberApp
	appAddr            string
	streamsInfo        delivery.StreamsInfo
	nestedPlaylistPath string
	nestedMediaPath    string
}

const waitTimeout = 5 * time.Second

func (s *StreamerSuite) BeforeAll(t provider.T) {
	t.Epic("bivi streamer")

	t.WithNewStep("read config", func(sCtx provider.StepCtx) {
		var configPath string

		pflag.StringVarP(&configPath, "config", "c", "env/app.test.yaml", "Config file path")
		pflag.Parse()
		viper.SetConfigFile(configPath)
		sCtx.Require().NoError(viper.ReadInConfig())
	})

	t.WithNewStep("setup web app", func(_ provider.StepCtx) {
		logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			AddSource: viper.GetBool("debug"),
			Level:     slog.LevelDebug,
		}))
		contentRoute := viper.GetString("streams.content_route")
		streamRepository := repository.NewRepository(viper.GetString("streams.content_path"), logger)
		streamNameEncoder := base64.StdEncoding
		streamUseCase := usecase.NewUseCase(contentRoute, streamRepository, streamNameEncoder, logger)
		settings := app.APISettings{
			Prefix:              viper.GetString("api_prefix"),
			Port:                viper.GetString("port"),
			ContentRoute:        contentRoute,
			ContentPath:         viper.GetString("streams.content_path"),
			ClientLogPath:       viper.GetString("client_log_path"),
			UploadFilesizeLimit: viper.GetInt64("upload_filesize_limit"),
		}
		s.app = app.NewFiberApp(settings, streamUseCase, streamNameEncoder, logger)
		s.appAddr = "http://localhost:" + settings.Port
	})
}

func (s *StreamerSuite) BeforeEach(t provider.T) {
	t.WithNewStep("run web app", func(_ provider.StepCtx) {
		go func() {
			_ = s.app.Start(viper.GetString("port"))
		}()
	})
}

func (s *StreamerSuite) AfterEach(t provider.T) {
	t.WithNewStep("shutdown web app", func(sCtx provider.StepCtx) {
		ctx, cancel := context.WithTimeout(context.Background(), waitTimeout)
		err := s.app.Shutdown(ctx)
		sCtx.Require().NoError(err)
		cancel()
	})
}

func (s *StreamerSuite) TestGetStreamsInfo() func(provider.StepCtx) {
	return func(sCtx provider.StepCtx) {
		url := s.appAddr + viper.GetString("api_prefix") + "/streams"
		req := httptest.NewRequest(http.MethodGet, url, nil)

		resp, err := s.app.Test(req, viper.GetInt("test.request_timeout_ms"))
		sCtx.Require().NoError(err)
		defer resp.Body.Close()

		sCtx.Require().Equal(fiber.StatusOK, resp.StatusCode)

		body, err := io.ReadAll(resp.Body)
		sCtx.Require().NoError(err)

		err = json.Unmarshal(body, &s.streamsInfo)
		sCtx.Require().NoError(err)
		sCtx.Require().NotZero(len(s.streamsInfo.Streams))
		sCtx.Logf("Stream name: %s", s.streamsInfo.Streams[0].Name)
		sCtx.Logf("Stream preview path: %s", s.streamsInfo.Streams[0].PreviewPath)
		sCtx.Logf("Stream playlist path: %s", s.streamsInfo.Streams[0].PlaylistPath)
	}
}

func (s *StreamerSuite) TestGetStreamPreview() func(provider.StepCtx) {
	return func(sCtx provider.StepCtx) {
		url := s.appAddr + s.streamsInfo.Streams[0].PreviewPath
		req := httptest.NewRequest(http.MethodGet, url, nil)

		resp, err := s.app.Test(req, viper.GetInt("test.request_timeout_ms"))
		sCtx.Require().NoError(err)
		defer resp.Body.Close()

		sCtx.Require().Equal(fiber.StatusOK, resp.StatusCode)
		sCtx.Require().NotZero(resp.ContentLength)

		body, err := io.ReadAll(resp.Body)
		sCtx.Require().NoError(err)
		sCtx.Require().NotZero(len(body))
	}
}

func (s *StreamerSuite) readFilename(data []byte) string {
	scanner := bufio.NewScanner(bytes.NewBuffer(data))
	for scanner.Scan() {
		line := scanner.Text()

		ext := filepath.Ext(line)
		if ext == ".m3u8" || ext == ".ts" {
			return line
		}
	}

	return ""
}

func (s *StreamerSuite) TestGetMainPlaylist() func(provider.StepCtx) {
	return func(sCtx provider.StepCtx) {
		url := s.appAddr + s.streamsInfo.Streams[0].PlaylistPath
		req := httptest.NewRequest(http.MethodGet, url, nil)

		resp, err := s.app.Test(req, viper.GetInt("test.request_timeout_ms"))
		sCtx.Require().NoError(err)
		defer resp.Body.Close()

		sCtx.Require().Equal(fiber.StatusOK, resp.StatusCode)
		sCtx.Require().NotZero(resp.ContentLength)

		body, err := io.ReadAll(resp.Body)
		sCtx.Require().NoError(err)
		sCtx.Require().NotZero(len(body))

		s.nestedPlaylistPath = s.readFilename(body)
		sCtx.Require().NotEmpty(s.nestedPlaylistPath)
	}
}

func (s *StreamerSuite) TestGetNestedPlaylist() func(provider.StepCtx) {
	return func(sCtx provider.StepCtx) {
		pathDir := filepath.Dir(s.streamsInfo.Streams[0].PlaylistPath)
		url := s.appAddr + pathDir + "/" + s.nestedPlaylistPath
		req := httptest.NewRequest(http.MethodGet, url, nil)

		resp, err := s.app.Test(req, viper.GetInt("test.request_timeout_ms"))
		sCtx.Require().NoError(err)
		defer resp.Body.Close()

		sCtx.Require().Equal(fiber.StatusOK, resp.StatusCode)
		sCtx.Require().NotZero(resp.ContentLength)

		body, err := io.ReadAll(resp.Body)
		sCtx.Require().NoError(err)
		sCtx.Require().NotZero(len(body))

		s.nestedMediaPath = s.readFilename(body)
		sCtx.Require().NotEmpty(s.nestedMediaPath)
	}
}

func (s *StreamerSuite) TestGetVideoFragment() func(provider.StepCtx) {
	return func(sCtx provider.StepCtx) {
		pathDir := filepath.Dir(s.streamsInfo.Streams[0].PlaylistPath)
		url := s.appAddr + pathDir + "/" + filepath.Dir(s.nestedPlaylistPath) + "/" + s.nestedMediaPath
		req := httptest.NewRequest(http.MethodGet, url, nil)

		resp, err := s.app.Test(req, viper.GetInt("test.request_timeout_ms"))
		sCtx.Require().NoError(err)
		defer resp.Body.Close()

		sCtx.Require().Equal(fiber.StatusOK, resp.StatusCode)
		sCtx.Require().NotZero(resp.ContentLength)

		body, err := io.ReadAll(resp.Body)
		sCtx.Require().NoError(err)
		sCtx.Require().NotZero(len(body))
	}
}

func (s *StreamerSuite) TestStreamer(t provider.T) {
	t.WithNewStep("Get streams info", s.TestGetStreamsInfo(), allure.NewParameter("time", time.Now()))
	t.WithNewStep("Get stream preview", s.TestGetStreamPreview(), allure.NewParameter("time", time.Now()))
	t.WithNewStep("Get main playlist", s.TestGetMainPlaylist(), allure.NewParameter("time", time.Now()))
	t.WithNewStep("Get nested playlist", s.TestGetNestedPlaylist(), allure.NewParameter("time", time.Now()))
	t.WithNewStep("Get video fragment", s.TestGetVideoFragment(), allure.NewParameter("time", time.Now()))
}

func TestStreamer(t *testing.T) {
	t.Parallel()

	if testing.Short() {
		t.Skip()
	}

	suite.RunSuite(t, new(StreamerSuite))
}
