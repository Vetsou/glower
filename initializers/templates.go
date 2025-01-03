package initializers

import (
	"log"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func LoadHTMLTemplates(e *gin.Engine) {
	patterns := []string{"templates/pages/*", "templates/partials/*"}

	var files []string
	for _, pattern := range patterns {
		matches, err := filepath.Glob(pattern)
		if err != nil {
			log.Fatalf("Failed to parse glob pattern: %v", err)
		}
		files = append(files, matches...)
	}

	e.LoadHTMLFiles(files...)
}
