package engine

import (
	"encoding/gob"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/hilaoyu/go-utils/utilHttp"
	"github.com/hilaoyu/go-utils/utilRandom"
)

const (
	ContextVariablesKeyResponse    = "_gin_extend_context_variables_key_response"
	ContextSessionKeyErrors        = "_gin_extend_res_session_errors"
	ContextSessionKeyAlertMessages = "_gin_extend_res_session_alert_messages"
	ContextSessionKeyVariables     = "_gin_extend_res_session_variables"
	ContextSessionKeyData          = "_gin_extend_res_session_data"
)

type Response struct {
	gc            *gin.Context
	status        bool
	statusCode    int
	message       string
	alertMessages []*ResponseAlertMessage

	debugs    []string
	errors    map[string]string
	variables ResponseVariables
	data      interface{}
}
type ResponseVariables map[string]interface{}

type ResponseAlertMessage struct {
	AlertType    string
	AlertMessage string
}

var (
	ErrorPageTemplates map[string]string
)

func init() {
	gob.Register(&ResponseAlertMessage{})
	gob.Register([]*ResponseAlertMessage{})
	gob.Register(ResponseVariables{})
}

func GetResponse(gc *gin.Context) (response *Response) {
	tmp, exist := gc.Get(ContextVariablesKeyResponse)
	if exist {
		if tmpRes, ok := tmp.(*Response); ok {
			response = tmpRes
		}
	}

	if nil == response {
		response = &Response{
			gc:         gc,
			statusCode: 200,
			variables:  map[string]interface{}{},
		}
		gc.Set(ContextVariablesKeyResponse, response)
	}

	session, _ := GetSession(gc)
	if nil != session {
		sessionErrors := session.Get(ContextSessionKeyErrors)
		if nil != sessionErrors {
			if errors, ok := sessionErrors.(map[string]string); ok {
				response.errors = errors
			}
		}
		session.Delete(ContextSessionKeyErrors)

		sessionAlertMessages := session.Get(ContextSessionKeyAlertMessages)
		if nil != sessionAlertMessages {
			if alertMessages, ok := sessionAlertMessages.([]*ResponseAlertMessage); ok {
				response.alertMessages = alertMessages
			}
		}
		session.Delete(ContextSessionKeyAlertMessages)

		sessionVariables := session.Get(ContextSessionKeyVariables)
		if nil != sessionVariables {
			if variables, ok := sessionVariables.(map[string]interface{}); ok {
				response.variables = variables
			}
		}
		session.Delete(ContextSessionKeyVariables)

		sessionData := session.Get(ContextSessionKeyData)
		if nil != sessionVariables {
			response.data = sessionData
		}
		session.Delete(ContextSessionKeyData)

		_ = session.Save()
	}

	return
}

func (res *Response) Success(msg string, code ...int) *Response {
	res.status = true
	res.message = msg
	if len(code) > 0 {
		res.statusCode = code[0]
	} else {
		res.statusCode = 200
	}
	return res
}
func (res *Response) Failed(msg string, code ...int) *Response {
	res.status = false
	res.message = msg
	if len(code) > 0 {
		res.statusCode = code[0]
	} else {
		res.statusCode = 501
	}
	return res
}

func (res *Response) AlertSuccess(msg string) *Response {
	res.alertMessages = append(res.alertMessages, &ResponseAlertMessage{
		AlertType:    "success",
		AlertMessage: msg,
	})
	return res
}
func (res *Response) AlertError(msg string) *Response {
	res.alertMessages = append(res.alertMessages, &ResponseAlertMessage{
		AlertType:    "error",
		AlertMessage: msg,
	})
	return res
}
func (res *Response) AlertWarning(msg string) *Response {
	res.alertMessages = append(res.alertMessages, &ResponseAlertMessage{
		AlertType:    "warning",
		AlertMessage: msg,
	})
	return res
}
func (res *Response) AlertInfo(msg string) *Response {
	res.alertMessages = append(res.alertMessages, &ResponseAlertMessage{
		AlertType:    "info",
		AlertMessage: msg,
	})
	return res
}

func (res *Response) WithDebug(v string) *Response {
	res.debugs = append(res.debugs, v)
	return res
}
func (res *Response) WithError(v string, k string) *Response {
	if nil == res.errors {
		res.errors = map[string]string{}
	}
	k = strings.TrimSpace(k)
	if "" == k {
		k = utilRandom.UniqId("")
	}
	res.errors[k] = v
	return res
}
func (res *Response) WithVariables(v interface{}, k string) *Response {
	k = strings.TrimSpace(k)
	if "" == k {
		return res
	}
	res.variables[k] = v
	return res
}
func (res *Response) SetData(v interface{}) *Response {
	res.data = v
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
	if nil != res.data {
		gc.JSON(res.statusCode, res.data)
		return
	}
	data := res.variables
	if len(res.debugs) > 0 {
		data["_debug"] = res.debugs
	}
	if len(res.errors) > 0 {
		if _, ok := data["errors"]; !ok {
			data["errors"] = res.errors
		}
	}
	gc.JSON(res.statusCode, data)
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

	data := utilHttp.ApiDataJson{
		Status:  res.status,
		Code:    res.statusCode,
		Message: res.message,
		Errors:  res.errors,
	}

	if nil != res.data {
		data.Data = res.data
	}
	if len(res.debugs) > 0 {
		data.Debug = res.debugs
	}
	if len(res.errors) > 0 {
		data.Errors = res.errors
	}
	gc.JSON(200, data)
	return
}
func (res *Response) RenderHtml(templateName string, c ...*gin.Context) {
	gc := res.gc
	if len(c) > 0 {
		gc = c[0]
	}
	if nil == gc {
		return
	}

	data := res.variables
	if len(res.message) > 0 {
		data["_message"] = res.message
	}
	if len(res.debugs) > 0 {
		data["_debug"] = res.debugs
	}
	if len(res.errors) > 0 {
		data["_errors"] = res.errors
	}
	if len(res.alertMessages) > 0 {
		data["_alert_messages"] = res.alertMessages
	}

	gc.HTML(res.statusCode, templateName, data)
	return
}

func (res *Response) RenderErrorPage(errorType string, c ...*gin.Context) {
	gc := res.gc
	if len(c) > 0 {
		gc = c[0]
	}
	if nil == gc {
		return
	}

	templateName, ok := ErrorPageTemplates[errorType]
	if !ok {
		gc.AbortWithError(res.statusCode, fmt.Errorf(res.message))
		return
	}

	data := res.variables
	if len(res.debugs) > 0 {
		data["_debug"] = res.debugs
	}
	if len(res.errors) > 0 {
		data["_errors"] = res.errors
	}

	if len(res.alertMessages) > 0 {
		data["_alert_messages"] = res.alertMessages
	}
	data["_message"] = res.message
	data["_prev_url"] = gc.Request.Referer()

	gc.HTML(res.statusCode, templateName, data)
	gc.Abort()
	return
}

func (res *Response) SendFileBytes(filename string, contentType string, data []byte, c ...*gin.Context) {
	gc := res.gc
	if len(c) > 0 {
		gc = c[0]
	}
	if nil == gc {
		return
	}

	gc.Header("Content-Description", "File Transfer")
	gc.Header("Content-Disposition", "attachment; filename="+filename)
	gc.Data(200, contentType, data)
}

func (res *Response) Redirect(code int, location string, c ...*gin.Context) {
	gc := res.gc
	if len(c) > 0 {
		gc = c[0]
	}
	if nil == gc {
		return
	}

	session, _ := GetSession(gc)
	if nil != session {

		if len(res.errors) > 0 {
			session.Set(ContextSessionKeyErrors, res.errors)
		}
		if len(res.alertMessages) > 0 {
			session.Set(ContextSessionKeyAlertMessages, res.alertMessages)
		}
		if len(res.variables) > 0 {
			session.Set(ContextSessionKeyVariables, res.variables)
		}
		if nil != res.data {
			session.Set(ContextSessionKeyData, res.data)
		}
		_ = session.Save()
	}

	gc.Redirect(code, location)
	return
}
func (res *Response) Abort(c ...*gin.Context) {
	gc := res.gc
	if len(c) > 0 {
		gc = c[0]
	}
	if nil == gc {
		return
	}

	if !res.status {
		_ = gc.AbortWithError(res.statusCode, fmt.Errorf(res.message))
	} else {
		gc.Abort()
	}

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
