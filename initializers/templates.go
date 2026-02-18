package initializers

import (
	"glower/resources"
	"html/template"
	"log"

	"github.com/gin-gonic/gin"
)

func InitHTMLTemplates(e *gin.Engine) {
	tmpl, err := template.ParseFS(
		resources.AssetsFS,
		"assets/pages/**/*.html",
		"assets/partials/*.html",
		"assets/partials/**/*.html",
		"assets/templates/*.html",
	)

	if err != nil {
		log.Fatalf("Error creating HTML templates %s", err.Error())
	}

	e.SetHTMLTemplate(tmpl)
}
