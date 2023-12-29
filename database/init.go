package database

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type pgConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Dbname   string `mapstructure:"dbname"`
}

var DB *pgxpool.Pool
var pgconf pgConfig

func InitDB() {

	postgresqMap := viper.Sub("postgres")
	if postgresqMap == nil {
		log.Error().Msg("postgres configuration not found in your config file, exiting...")
		os.Exit(1)
	}
	postgresqMap.Unmarshal(&pgconf)

	log.Debug().Str("component", "cloud/Collector").Msgf("collector configuration: %v", pgconf)

	db_connection := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable pool_max_conns=10",
		pgconf.Host,
		pgconf.Port,
		pgconf.User,
		pgconf.Password,
		pgconf.Dbname,
	)
	log.Info().Str("component", "cloud/Collector").Msgf("db_connection: %s", db_connection)

	var err error
	DB, err = pgxpool.New(context.Background(), db_connection)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		os.Exit(1)
	}
}

func DBClose() {
	DB.Close()
}

func CheckDB() bool {
	var result bool

	ctx := context.Background()
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 1*time.Second)

	err := DB.QueryRow(ctxWithTimeout, "SELECT true").Scan(&result)
	defer cancel()

	select {
	case <-time.After(5 * time.Second):
		if err != nil {
			log.Error().Msgf("Error checking DB: %v", err)
			return false
		}
		return true

	case <-ctx.Done():
		log.Warn().Msg("Postgresql healch check Process timed out")
		return false
	}
}

func PgsqlCheckLoop(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-ctx.Done():
			log.Info().Msg("PgsqlCheckLoop stopped")
			ticker.Stop()

		case <-ticker.C:
			log.Trace().Msg("PgsqlCheckLoop running")
			if !CheckDB() {
				log.Warn().Msg("DB is not responding, retrying conection...")
				InitDB()
			}

		}
	}

}
