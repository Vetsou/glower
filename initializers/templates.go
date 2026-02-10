package initializers

import (
	"log"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func InitHTMLTemplates(e *gin.Engine, path string) {
	patterns := []string{
		path + "assets/pages/index.html",
		path + "assets/pages/user/*",
		path + "assets/pages/shop/*",
		path + "assets/pages/error/*",

		path + "assets/partials/*",
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
