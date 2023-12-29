package server

import (
	"log/syslog"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/journald"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

func SetupLogger() zerolog.Logger {
	/*
		+----------------------------------------------------------------------------------------+
		| LOGGER INIZIALIZATION
		+----------------------------------------------------------------------------------------+
	*/
	var logger zerolog.Logger

	logCollector := viper.Get("log.collector")
	if logCollector == nil {
		logCollector = "stdout"
	}

	log.Debug().Str("logCollector", logCollector.(string)).Msg("Log collector")

	switch logCollector {
	case "journald":
		logger = log.Output(journald.NewJournalDWriter())
		log.Debug().Msg("Using journald logger")
	case "console":
		logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
		log.Warn().Msg("Console logger is not recommended for production use, use journald or stdout or syslog instead")
	case "syslog":
		/*
			+----------------------------------------------------------------------------------------+
			| SYSLOG INIZIALIZATION
			+----------------------------------------------------------------------------------------+
		*/
		var priority syslog.Priority
		switch viper.Get("log.syslog.local") {
		case "local0":
			priority = syslog.LOG_LOCAL0
		case "local1":
			priority = syslog.LOG_LOCAL1
		case "local2":
			priority = syslog.LOG_LOCAL2
		case "local3":
			priority = syslog.LOG_LOCAL3
		case "local4":
			priority = syslog.LOG_LOCAL4
		case "local5":
			priority = syslog.LOG_LOCAL5
		case "local6":
			priority = syslog.LOG_LOCAL6
		case "local7":
			priority = syslog.LOG_LOCAL7
		default:
			priority = syslog.LOG_LOCAL7
		}

		switch viper.Get("log.syslog.protocol") {
		case "tcp":
			viper.Get("log.syslog.host")

			syslogHost := viper.Get("log.syslog.host").(string)
			zsyslog, err := syslog.Dial(
				"tcp",
				syslogHost,
				priority,
				"datacollector")
			if err != nil {
				log.Fatal().Err(err).Msg("Error during syslog configuration fallback to stdout")
				break
			}
			log.Output(zerolog.New(zsyslog).With().Caller().Logger())
		case "udp":
			syslogHost := viper.Get("log.syslog.host").(string)
			zsyslog, err := syslog.Dial(
				"udp",
				syslogHost,
				priority,
				"datacollector")
			if err != nil {
				log.Fatal().Err(err).Msg("Error during syslog configuration fallback to stdout")
				break
			}
			log.Output(zerolog.New(zsyslog).With().Caller().Logger())
		default:
			zsyslog, err := syslog.New(
				priority,
				"datacollector")
			if err != nil {
				log.Fatal().Err(err).Msg("Error during syslog configuration fallback to stdout")
				break
			}
			log.Logger = zerolog.New(zerolog.SyslogLevelWriter(zsyslog))
			log.Debug().Msg("Using syslog logger")
		}
		// Configure zerolog to use syslog
	case "stdout":
	default:
		logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
	}

	logLevel := viper.Get("log.level")
	if logLevel == nil {
		logLevel = "info"
	}

	switch logLevel {
	case "trace":
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case "fatal":
		zerolog.SetGlobalLevel(zerolog.FatalLevel)
	case "panic":
		zerolog.SetGlobalLevel(zerolog.PanicLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	return logger
}
