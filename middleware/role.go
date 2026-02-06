package middleware

import (
	"glower/database/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateRolesAuth(allowedRoles ...model.Role) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleVal, exists := c.Get("role")
		if !exists {
			RenderResponse(c, http.StatusUnauthorized, "You cannot view this page.")
			return
		}

		userRole, ok := roleVal.(model.Role)
		if !ok {
			RenderResponse(c, http.StatusInternalServerError, "We couldn't verify your access rights.")
			return
		}

		for _, role := range allowedRoles {
			if userRole == role {
				c.Next()
				return
			}
		}

		RenderResponse(c, http.StatusForbidden, "You don't have the right permissions.")
	}
}
