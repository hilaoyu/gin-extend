package engine

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hilaoyu/go-utils/utilHttp"
	"strings"
)

type Response struct {
	gc         *gin.Context
	Status     bool
	StatusCode int
	Message    string

	Debugs    []string
	Errors    []string
	Variables map[string]interface{}
	Data      interface{}
}

const (
	ContextVariablesKeyResponse = "_gin_extend_context_variables_key_response"
)

var ErrorPageTemplates map[string]string

func GetResponse(gc *gin.Context) (response *Response) {
	tmp, exist := gc.Get(ContextVariablesKeyResponse)
	if exist {
		if tmpRes, ok := tmp.(*Response); ok {
			response = tmpRes
		}
	}

	if nil == response {
		response = &Response{
			gc:        gc,
			Variables: map[string]interface{}{},
		}
		gc.Set(ContextVariablesKeyResponse, response)
	}

	return
}

func (res *Response) Success(msg string, code ...int) *Response {
	res.Status = true
	res.Message = msg
	if len(code) > 0 {
		res.StatusCode = code[0]
	} else {
		res.StatusCode = 200
	}
	return res
}
func (res *Response) Failed(msg string, code ...int) *Response {
	res.Status = false
	res.Message = msg
	if len(code) > 0 {
		res.StatusCode = code[0]
	} else {
		res.StatusCode = 501
	}
	return res
}

func (res *Response) WithDebug(v string) *Response {
	res.Debugs = append(res.Debugs, v)
	return res
}
func (res *Response) WithError(v string) *Response {
	res.Errors = append(res.Errors, v)
	return res
}
func (res *Response) WithVariables(v interface{}, k string) *Response {
	k = strings.TrimSpace(k)
	if "" == k {
		return res
	}
	res.Variables[k] = v
	return res
}
func (res *Response) SetData(v interface{}) *Response {
	res.Data = v
	return res
}

func (res *Response) RenderJson(c *gin.Context) {
	gc := c
	if nil == gc {
		gc = res.gc
	}
	if nil == gc {
		return
	}
	if nil != res.Data {
		gc.JSON(res.StatusCode, res.Data)
		return
	}
	data := res.Variables
	if len(res.Debugs) > 0 {
		data["_debug"] = res.Debugs
	}
	if len(res.Errors) > 0 {
		if _, ok := data["errors"]; !ok {
			data["errors"] = res.Errors
		}
	}
	gc.JSON(res.StatusCode, data)
	return
}
func (res *Response) RenderApiJson(c *gin.Context) {
	gc := c
	if nil == gc {
		gc = res.gc
	}
	if nil == gc {
		return
	}

	data := utilHttp.ApiReturnJson{
		Status:  res.Status,
		Code:    res.StatusCode,
		Message: res.Message,
		Errors:  res.Errors,
	}

	if nil != res.Data {
		data.Data = res.Data
	}
	if len(res.Debugs) > 0 {
		data.Debug = res.Debugs
	}
	if len(res.Errors) > 0 {
		data.Errors = res.Errors
	}
	gc.JSON(200, data)
	return
}
func (res *Response) RenderHtml(templateName string, c *gin.Context) {
	gc := c
	if nil == gc {
		gc = res.gc
	}
	if nil == gc {
		return
	}

	data := res.Variables
	if len(res.Debugs) > 0 {
		data["_debug"] = res.Debugs
	}
	if len(res.Errors) > 0 {
		if _, ok := data["errors"]; !ok {
			data["errors"] = res.Errors
		}
	}

	gc.HTML(res.StatusCode, templateName, data)
	return
}

func (res *Response) RenderErrorPage(errorType string, c *gin.Context) {
	gc := c
	if nil == gc {
		gc = res.gc
	}
	if nil == gc {
		return
	}

	templateName, ok := ErrorPageTemplates[errorType]
	if !ok {
		gc.AbortWithError(res.StatusCode, fmt.Errorf(res.Message))
		return
	}

	data := res.Variables
	if len(res.Debugs) > 0 {
		data["_debug"] = res.Debugs
	}
	if len(res.Errors) > 0 {
		if _, ok := data["errors"]; !ok {
			data["errors"] = res.Errors
		}
	}
	data["message"] = res.Message
	data["_prev_url"] = gc.Request.Referer()

	gc.HTML(res.StatusCode, templateName, data)
	gc.Abort()
	return
}

func SetErrorPageTemplates(templates map[string]string) {
	ErrorPageTemplates = templates
}

func AddErrorPageTemplates(errorType string, template string) {
	if nil == ErrorPageTemplates {
		ErrorPageTemplates = map[string]string{}
	}
	ErrorPageTemplates[errorType] = template
}
