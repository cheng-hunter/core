package run

import (
	"encoding/json"
	"fmt"
	"github.com/cheng-hunter/core/manager"
	"github.com/cheng-hunter/core/types"
	"github.com/valyala/fasthttp"
	"net"
	"os"
)

var (
	registryMap map[string]func(ctx *fasthttp.RequestCtx)
)

func Start(plugin []types.Plugin) error {
	registryMap["/Plugin.Activate"] = PluginActivate
	registryMap["/Plugin.Start"] = PluginStart
	registryMap["/Plugin.Invoke"] = PluginInvoke
	pluginId, _ := os.LookupEnv(types.PluginIdEnvKey)
	go func() {
		//启动插件
		listener, err := net.Listen("unix", fmt.Sprintf(types.PluginSocketPath, pluginId))
		if err != nil {
			panic(err)
		}
		fasthttp.Serve(listener, func(ctx *fasthttp.RequestCtx) {
			if f, ok := registryMap[string(ctx.Path())]; ok {
				ctx.SetUserValue(types.PluginReqKey, plugin)
				ctx.SetUserValue(types.PluginIdEnvKey, pluginId)
				f(ctx)
			}
		})
	}()
	return nil
}

func PluginActivate(ctx *fasthttp.RequestCtx) {
	plugins, ok := ctx.UserValue(types.PluginReqKey).([]types.Plugin)
	if !ok {
		types.ResFailMessage(ctx, "插件激活失败")
		return
	}
	var info []types.PluginInfo
	for _, t := range plugins {
		p := t.Activate()
		//
		manager.DefaultPluginManager.AddLocal(p.Name, t)
		info = append(info, p)
	}
	types.ResSuccess(ctx, info)
}

func PluginStart(ctx *fasthttp.RequestCtx) {
	by := ctx.Request.PostArgs().Peek("pluginName")
	if by == nil {
		types.ResFailMessage(ctx, "插件启动失败")
		return
	}
	manager.DefaultPluginManager.Get(string(by)).Start()
	types.ResSuccess(ctx, nil)
}

func PluginInvoke(ctx *fasthttp.RequestCtx) {
	pluginName := ctx.Request.PostArgs().Peek("pluginName")
	if pluginName == nil {
		types.ResFail(ctx, 1000, "插件调用失败")
		return
	}
	funcName := ctx.Request.PostArgs().Peek("funcName")
	if funcName == nil {
		types.ResFail(ctx, 1001, "插件调用失败")
		return
	}
	arg := ctx.Request.PostArgs().Peek("args")
	if arg == nil {
		types.ResFail(ctx, 1002, "插件调用失败")
		return
	}
	var args []interface{}
	err := json.Unmarshal(arg, &args)
	if err != nil {
		types.ResFail(ctx, 1003, "插件调用失败")
		return
	}
	rlt, err := manager.DefaultPluginManager.InvokeFunction(string(pluginName), string(funcName), args)
	if err != nil {
		types.ResFail(ctx, 1004, "插件调用失败")
		return
	}
	types.ResSuccess(ctx, rlt)
}
