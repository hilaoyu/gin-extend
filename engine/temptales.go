package engine

import (
	"github.com/gin-contrib/multitemplate"
	"html/template"
	"path/filepath"
	"strings"
	"time"
)

type GinTemplate = struct {
	Files  string `json:"files"`
	Layout string `json:"layout"`
}

func (e *GinEngine) UseMultiTemplate(templates []*GinTemplate, templateBasePath string) (err error) {
	r := multitemplate.NewRenderer()

	for _, t := range templates {
		var layouts []string
		var templateFiles []string
		if "" != t.Layout {
			layouts, err = filepath.Glob(t.Layout)
			if nil != err {
				return
			}
		}

		templateFiles, err = filepath.Glob(t.Files)
		if err != nil {
			return
		}

		// Generate our templates map from our articleLayouts/ and articles/ directories
		for _, file := range templateFiles {
			layoutCopy := make([]string, len(layouts))
			copy(layoutCopy, layouts)
			files := append(layoutCopy, file)
			templateName := filepath.Base(file)
			if "" != templateBasePath {
				templateName = strings.TrimPrefix(file, templateBasePath)
			}
			templateName = strings.TrimLeft(templateName, "/")
			r.AddFromFiles(templateName, files...)
		}

	}

	e.HTMLRender = r

	e.SetFuncMap(template.FuncMap{
		"formatCurrentTime": func(layout string) string {
			return time.Now().Format(layout)
		},
	})
	return
}
