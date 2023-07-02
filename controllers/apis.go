package controllers

import (
	"github.com/gin-gonic/gin"
)

func Pods(c *gin.Context) {
	logger := getLogger(c)

	pods, err := listPods(c)
	if err != nil {
		logger.Error().Err(err).Msg("list pods fail")
		c.JSON(500, gin.H{"msg": "internal error"})
		return
	}
	c.JSON(200, gin.H{"data": pods})
}
