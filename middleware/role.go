package middleware

import (
	"glower/database/model"
	"glower/middleware/internal"
	"net/http"
	"slices"

	"github.com/gin-gonic/gin"
)

func CreateRolesAuth(allowedRoles ...model.Role) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleVal, exists := c.Get("role")
		if !exists {
			internal.RenderErrorResponse(c, http.StatusUnauthorized, "You cannot view this page.")
			return
		}

		userRole, ok := roleVal.(model.Role)
		if !ok {
			internal.RenderErrorResponse(c, http.StatusInternalServerError, "We couldn't verify your access rights.")
			return
		}

		if slices.Contains(allowedRoles, userRole) {
			c.Next()
			return
		}

		internal.RenderErrorResponse(c, http.StatusForbidden, "You don't have the right permissions.")
	}
}
