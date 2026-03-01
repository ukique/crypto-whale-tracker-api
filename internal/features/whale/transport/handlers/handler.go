package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/ukique/crypto-whale-tracker-api/internal/features/whale/models"
)

func WhaleHandler(whaleChan <-chan models.Whale) gin.HandlerFunc {
	return func(c *gin.Context) {
		whale := <-whaleChan
		c.JSON(200, gin.H{
			"whale": whale,
		})
	}
}
