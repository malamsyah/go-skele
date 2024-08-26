package middleware

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func ValidateID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		if id != "" {
			parsedID, err := strconv.Atoi(id)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "ID should be an integer",
				})
				c.Abort()
				return
			}

			c.Set("parsedID", parsedID)
		}

		c.Next()
	}
}
