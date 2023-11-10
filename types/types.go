package types

import (
	"encoding/json"
	"fmt"
	"github.com/valyala/fasthttp"
	"net/http"
)

const (
	DefaultPluginPath      = "/tmp/core.socket"
	PluginIdEnvKey         = "pluginId"
	PluginSocketPath       = "/tmp/%s.socket"
	PluginSocketServiceUrl = "unix:///tmp/%s.socket"

	PluginReqKey = "pluginReq"
)

type PluginInfo struct {
	Name    string `json:"name"`
	Version int32  `json:"version"`
}

type SharePluginInfo struct {
	PluginInfo
	Id string `json:"id"`
}

// PluginResult 插件结果
type PluginResultBase struct {
	//状态码(HTTP状态码)
	StatusCode int `json:"statusCode"`
	//业务是否成功
	Success bool `json:"success"`
	//业务消息(业务失败需要关注)
	Message string `json:"message,omitempty"`
}

// PluginResult 插件结果
type PluginResult struct {
	PluginResultBase
	//业务数据
	Detail interface{} `json:"detail,omitempty"`
}

// PluginResult 插件结果
type PluginInvokeResp struct {
	PluginResultBase
	//业务数据
	Detail []interface{} `json:"detail,omitempty"`
}

// ResSuccess 响应成功
func ResSuccess(ctx *fasthttp.RequestCtx, v interface{}) {
	ResJSON(ctx, http.StatusOK, PluginResult{
		PluginResultBase: PluginResultBase{
			StatusCode: 200,
			Success:    true,
		},
		Detail: v,
	})
}

// ResFail 业务响应失败
func ResFail(ctx *fasthttp.RequestCtx, errorCode int, errorMessage string) {
	ResJSON(ctx, http.StatusOK, PluginResult{
		PluginResultBase: PluginResultBase{
			StatusCode: 200,
			Success:    false,
			Message:    fmt.Sprintf("%s(%d)", errorMessage, errorCode),
		},
	})
}

// ResFail 业务响应失败
func ResFailMessage(ctx *fasthttp.RequestCtx, errorMessage string) {
	ResJSON(ctx, http.StatusOK, PluginResult{
		PluginResultBase: PluginResultBase{
			StatusCode: 200,
			Success:    false,
			Message:    errorMessage,
		},
	})
}

// ResHttpFail Http失败
func ResHttpFail(ctx *fasthttp.RequestCtx, statusCode int, v interface{}) {
	ResJSON(ctx, statusCode, PluginResult{
		PluginResultBase: PluginResultBase{
			StatusCode: statusCode,
			Success:    false,
		},
	})
}

// ResJSON 响应JSON数据
func ResJSON(ctx *fasthttp.RequestCtx, status int, v interface{}) {
	buf, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	ctx.SetStatusCode(status)
	ctx.Success("application/json; charset=utf-8", buf)
}

type Plugin interface {
	Activate() PluginInfo
	Start()
	Destroy()
	InvokeFunction(name string, args ...interface{}) ([]interface{}, error)
}
