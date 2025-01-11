package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetCookies(c *gin.Context, refreshToken *string, accessToken *string) {
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     RefreshTokenName,
		Value:    *refreshToken,
		Path:     "/",
		Domain:   DomainName,
		MaxAge:   3 * 24 * 60 * 60,
		Secure:   false,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     AccessTokenName,
		Value:    *accessToken,
		Path:     "/",
		Domain:   DomainName,
		MaxAge:   24 * 60 * 60,
		Secure:   false,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})
}

func CleanCookies(c *gin.Context) {
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     RefreshTokenName,
		Value:    "",
		Path:     "/",
		Domain:   DomainName,
		MaxAge:   -1,
		Secure:   false,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     AccessTokenName,
		Value:    "",
		Path:     "/",
		Domain:   DomainName,
		MaxAge:   -1,
		Secure:   false,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})
}
