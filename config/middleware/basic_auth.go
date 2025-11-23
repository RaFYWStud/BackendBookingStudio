package middleware

import (
	"crypto/subtle"
	"net/http"

	"github.com/gin-gonic/gin"
)

func BasicAuthForSwagger() gin.HandlerFunc {
	const (
		SWAGGER_USERNAME = "admin"
		SWAGGER_PASSWORD = "kelompok3"
	)

	return func(c *gin.Context) {
		username, password, hasAuth := c.Request.BasicAuth()

		if !hasAuth {
			c.Header("WWW-Authenticate", `Basic realm="Swagger Documentation"`)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized - Swagger access requires authentication",
			})
			return
		}

		usernameMatch := subtle.ConstantTimeCompare([]byte(username), []byte(SWAGGER_USERNAME)) == 1
		passwordMatch := subtle.ConstantTimeCompare([]byte(password), []byte(SWAGGER_PASSWORD)) == 1

		if !usernameMatch || !passwordMatch {
			c.Header("WWW-Authenticate", `Basic realm="Swagger Documentation"`)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid username or password",
			})
			return
		}

		c.Next()
	}
}