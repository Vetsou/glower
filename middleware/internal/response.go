package internal

import "github.com/gin-gonic/gin"

func RenderErrorResponse(c *gin.Context, code int, msg string) {
	isHTMX := c.Request.Header.Get("HX-Request") == "true"

	if isHTMX {
		c.HTML(code, "error-alert.html", gin.H{"errorMessage": msg})
	} else {
		c.HTML(code, "error-page.html", gin.H{
			"code":    code,
			"message": msg,
		})
	}
	c.Abort()
}
