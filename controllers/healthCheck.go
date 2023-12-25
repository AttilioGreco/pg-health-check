package controllers

import (
	"context"
	"net/http"
	"pgHealtchCheck/database"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func WriteHealthCheck(c *gin.Context) {
	ctx := context.Background()
	// Esegui la query per controllare se il database è in recovery
	var isRecovery bool
	err := database.DB.QueryRow(ctx, "SELECT pg_is_in_recovery()").Scan(&isRecovery)
	if err != nil {
		log.Error().Err(err).Msg("Error executing query")
		c.JSON(http.StatusInternalServerError, gin.H{"message": "503 Service Unavailable"})
	}

	switch isRecovery {
	case false:
		c.JSON(http.StatusOK, gin.H{"message": "200 OK - Primary"})
	case true:
		c.JSON(http.StatusForbidden, gin.H{"message": "403 Forbidden - Follower"})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"message": "503 DB Unknow status"})
	}
}

func ReadHealthCheck(c *gin.Context) {
	// Esegui un ping sul database
	err := database.DB.Ping(context.Background())
	log.Err(err).Msg("Error pinging database")
	if err != nil {
		log.Error().Err(err).Msg("Error pinging database")
		c.JSON(http.StatusInternalServerError, gin.H{"message": "503 Service Unavailable"})
	}
	// Restituisci la risposta in caso di successo
	c.JSON(http.StatusOK, gin.H{"message": "200 OK"})
}
