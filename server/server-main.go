package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"pgHealtchCheck/controllers"
	"pgHealtchCheck/database"
	"pgHealtchCheck/middleware"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type serverConfig struct {
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

func Run(cmd *cobra.Command, args []string) {

	var srvconf serverConfig
	serverMap := viper.Sub("server")

	if serverMap == nil {
		log.Error().Msg("server configuration not found in your config file, exiting...")
		os.Exit(1)
	}

	serverMap.Unmarshal(&srvconf)
	log.Debug().Str("component", "cloud/Collector").Msgf("collector configuration: %v", srvconf)

	switch srvconf.Mode {
	case "release":
		gin.SetMode(gin.ReleaseMode)
	case "debug":
		gin.SetMode(gin.DebugMode)
	default:
		gin.SetMode(gin.ReleaseMode)
	}

	// Configura il router Gin
	router := gin.New()
	router.Use(middleware.DefaultStructuredLogger())
	router.Use(gin.Recovery())
	// Heatch Check Routes

	router.GET("/write", controllers.WriteHealthCheck)
	router.GET("/read", controllers.ReadHealthCheck)

	log.Info().Msg("Postgresql health check started")

	HostPortConfigString := fmt.Sprintf("%s:%s", srvconf.Host, srvconf.Port)

	srv := &http.Server{
		Addr:    HostPortConfigString,
		Handler: router,
	}

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Info().Msgf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	log.Info().Msg("Starting postgresql health check")
	database.InitDB()
	go database.PgsqlCheckLoop(ctx)

	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be caught, so don't need to add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	log.Debug().Msg("Waiting for signal")
	<-quit
	log.Info().Msg("Shutdown Server ...")
	cancel()

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancelTerminate := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelTerminate()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal().Err(err).Msg("Server forced to shutdown")
	}

	log.Info().Msg("Server exiting")

	log.Info().Msg("Stopping postgresql health check")
}
