package initializers

import (
	"log"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func InitHTMLTemplates(e *gin.Engine) {
	patterns := []string{
		"templates/pages/index.html",
		"templates/pages/user/*",
		"templates/pages/shop/*",
		"templates/pages/error/*",

		"templates/partials/*",
	}

	var files []string
	for _, pattern := range patterns {
		matches, err := filepath.Glob(pattern)
		if err != nil {
			log.Fatalf("Failed to parse HTML filepath pattern: %v", err)
		}
		files = append(files, matches...)
	}

	e.LoadHTMLFiles(files...)
}
