package util

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cheng-hunter/core/types"
	"github.com/valyala/fasthttp"
	"time"
)

var client *fasthttp.Client

func init() {
	// You may read the timeouts from some config
	readTimeout, _ := time.ParseDuration("500ms")
	writeTimeout, _ := time.ParseDuration("500ms")
	maxIdleConnDuration, _ := time.ParseDuration("1h")
	client = &fasthttp.Client{
		ReadTimeout:                   readTimeout,
		WriteTimeout:                  writeTimeout,
		MaxIdleConnDuration:           maxIdleConnDuration,
		NoDefaultUserAgentHeader:      true, // Don't send: User-Agent: fasthttp
		DisableHeaderNamesNormalizing: true, // If you set the case on your headers correctly you can enable this
		DisablePathNormalizing:        true,
		// increase DNS cache time to an hour instead of default minute
		Dial: (&fasthttp.TCPDialer{
			Concurrency:      4096,
			DNSCacheDuration: time.Hour,
		}).Dial,
	}
}

func ActivatePlugin(pluginId string) ([]types.PluginInfo, error) {
	req := fasthttp.AcquireRequest()
	req.SetRequestURI(fmt.Sprintf(types.PluginSocketServiceUrl, pluginId))
	req.Header.SetMethod(fasthttp.MethodGet)
	resp := fasthttp.AcquireResponse()
	err := client.Do(req, resp)
	fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)
	if err != nil {
		return nil, err
	}
	var rlt []types.PluginInfo
	err = json.Unmarshal(resp.Body(), &rlt)
	if err != nil {
		return nil, err
	}
	return rlt, nil
}

func StartPlugin(pluginId string, pluginName string) error {
	req := fasthttp.AcquireRequest()
	req.SetRequestURI(fmt.Sprintf(types.PluginSocketServiceUrl, pluginId))
	req.Header.SetMethod(fasthttp.MethodPost)
	req.PostArgs().Set("pluginName", pluginName)
	resp := fasthttp.AcquireResponse()
	err := client.Do(req, resp)
	fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)
	if err != nil {
		return err
	}
	var rlt types.PluginResult
	err = json.Unmarshal(resp.Body(), &rlt)
	if err != nil {
		return err
	}
	if !rlt.Success {
		return errors.New(rlt.Message)
	}
	return nil
}

func InvokeFunction(pluginId string, pluginName string, funcName string, args []interface{}) ([]interface{}, error) {
	req := fasthttp.AcquireRequest()
	req.SetRequestURI(fmt.Sprintf(types.PluginSocketServiceUrl, pluginId))
	req.Header.SetMethod(fasthttp.MethodPost)
	req.PostArgs().Set("pluginName", pluginName)
	req.PostArgs().Set("funcName", funcName)
	a, err := json.Marshal(args)
	//var arrays []string
	//for _, arg := range args {
	//	a,err :=json.Marshal(arg)
	//	if err != nil {
	//		return nil,err
	//	}
	//	arrays=append(arrays,string(a))
	//}
	req.PostArgs().Set("args", string(a))
	resp := fasthttp.AcquireResponse()
	err = client.Do(req, resp)
	fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)
	if err != nil {
		return nil, err
	}
	var rlt types.PluginInvokeResp
	err = json.Unmarshal(resp.Body(), &rlt)
	if err != nil {
		return nil, err
	}
	if !rlt.Success {
		return nil, errors.New(rlt.Message)
	}
	return rlt.Detail, nil
}
