package server

import (
	"fmt"
	"os"
	"pgHealtchCheck/controllers"
	"pgHealtchCheck/database"
	"pgHealtchCheck/middleware"

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

	database.InitDB()
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
	router.Run(HostPortConfigString)

	log.Info().Msg("Stopping postgresql health check")
}
