package engine

import (
	"github.com/gin-contrib/multitemplate"
	"path/filepath"
)

type GinTemplate = struct {
	Files  string `json:"files"`
	Layout string `json:"layout"`
}

func (e *GinEngine) UseMultiTemplate(templates []*GinTemplate) (err error) {
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
			r.AddFromFiles(filepath.Base(file), files...)
		}

	}

	e.HTMLRender = r
	return
}
