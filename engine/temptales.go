package engine

import (
	"encoding/json"
	"github.com/gin-contrib/multitemplate"
	"html/template"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type GinTemplate = struct {
	Files  string `json:"files"`
	Layout string `json:"layout"`
}

func (e *GinEngine) UseMultiTemplate(templates []*GinTemplate, templateBasePath string) (err error) {
	e.SetFuncMap(template.FuncMap{
		"include": func(name string, args map[string]interface{}, additionalArgs ...interface{}) template.HTML {
			cb, err1 := os.ReadFile(filepath.Join(templateBasePath, name))
			if nil != err1 {
				return ""
			}
			var w strings.Builder
			for i, arg := range additionalArgs {
				args["_props_arg_"+strconv.Itoa(i)] = arg
			}
			tmpl, err1 := template.New(name).
				Funcs(e.FuncMap).
				Parse(string(cb))
			if nil != err1 {
				return ""
			}
			err1 = tmpl.Execute(&w, args)
			if nil != err1 {
				return ""
			}
			return template.HTML(w.String())
		},
		"formatCurrentTime": func(layout string) string {
			return time.Now().Format(layout)
		},
		"toJson": func(v interface{}) string {
			s, _ := json.Marshal(v)
			return string(s)
		},
	})

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
			r.AddFromFilesFuncs(templateName, e.FuncMap, files...)
		}

	}

	e.HTMLRender = r

	return
}
func (e *GinEngine) ExtendTemplateFuncMap(name string, callback any) (err error) {

	if nil == e.FuncMap {
		e.FuncMap = template.FuncMap{}
	}
	e.FuncMap[name] = callback
	return
}
