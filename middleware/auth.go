package middleware

import (
	"glower/auth"
	"glower/middleware/internal"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateAuth(isStrict bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr, err := c.Cookie(auth.AccessTokenName)
		if err != nil {
			if isStrict {
				internal.RenderErrorResponse(c, http.StatusUnauthorized, "User is not logged in.")
				return
			}
			c.Next()
			return
		}

		claims, err := auth.VerifyToken(tokenStr)
		if err != nil {
			if isStrict {
				internal.RenderErrorResponse(c, http.StatusUnauthorized, "Invalid user credentials.")
				return
			}
			c.Next()
			return
		}

		userData, err := auth.GetUserClaims(claims)
		if err != nil {
			if isStrict {
				internal.RenderErrorResponse(c, http.StatusInternalServerError, "Error getting user data.")
				return
			}
			c.Next()
			return
		}

		c.Set("id", userData.Id)
		c.Set("user", userData.User)
		c.Set("role", userData.Role)
		c.Next()
	}
}
