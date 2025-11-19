package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/unsrat-it-community/back-end-e-voting-2025/config/pkg/errs"
)

// Auth middleware validates JWT token
func Auth() gin.HandlerFunc {
    return func(ctx *gin.Context) {
        authHeader := ctx.GetHeader("Authorization")
        if authHeader == "" {
            ctx.JSON(http.StatusUnauthorized, gin.H{
                "success": false,
                "error":   "Authorization header required",
                "code":    "MISSING_TOKEN",
            })
            ctx.Abort()
            return
        }

        // Extract token from "Bearer <token>"
        parts := strings.Split(authHeader, " ")
        if len(parts) != 2 || parts[0] != "Bearer" {
            ctx.JSON(http.StatusUnauthorized, gin.H{
                "success": false,
                "error":   "Invalid authorization format. Use: Bearer <token>",
                "code":    "INVALID_TOKEN_FORMAT",
            })
            ctx.Abort()
            return
        }

        tokenString := parts[1]

        // Parse and validate token
        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            // TODO: Use RSA public key from config
            // For now, use the same secret key used in service
            return []byte("temporary-secret-key-replace-with-rsa"), nil
        })

        if err != nil || !token.Valid {
            ctx.JSON(http.StatusUnauthorized, gin.H{
                "success": false,
                "error":   "Invalid or expired token",
                "code":    "INVALID_TOKEN",
            })
            ctx.Abort()
            return
        }

        // Extract claims and store in context
        if claims, ok := token.Claims.(jwt.MapClaims); ok {
            ctx.Set("user_id", int(claims["user_id"].(float64)))
            ctx.Set("user_email", claims["email"].(string))
            ctx.Set("user_role", claims["role"].(string))
        }

        ctx.Next()
    }
}

// AdminOnly middleware checks if user has admin role
func AdminOnly() gin.HandlerFunc {
    return func(ctx *gin.Context) {
        role, exists := ctx.Get("user_role")
        if !exists {
            ctx.AbortWithStatusJSON(http.StatusUnauthorized, errs.Unauthorized("user not authenticated"))
            return
        }

        if role != "admin" {
            ctx.AbortWithStatusJSON(http.StatusForbidden, errs.Forbidden("admin access required"))
            return
        }

        ctx.Next()
    }
}