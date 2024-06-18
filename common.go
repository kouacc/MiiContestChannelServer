package main

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

func writeResult(c *gin.Context, result int) {
	c.Header("X-RESULT", strconv.Itoa(result))
}
