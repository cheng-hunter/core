package manager

import (
	"errors"
	"fmt"
	"github.com/cheng-hunter/core/types"
)

var DefaultPluginManager *PluginManager

type PluginManager struct {
	pluginInfo  map[string]types.SharePluginInfo
	localPlugin map[string]types.Plugin
}

func init() {
	DefaultPluginManager = &PluginManager{
		pluginInfo:  map[string]types.SharePluginInfo{},
		localPlugin: map[string]types.Plugin{},
	}
}

func (m *PluginManager) Add(pluginName string, pluginInfo types.SharePluginInfo) {
	m.pluginInfo[pluginName] = pluginInfo
}

func (m *PluginManager) AddLocal(pluginName string, plugin types.Plugin) {
	if _, ok := m.pluginInfo[pluginName]; !ok {
		m.localPlugin[pluginName] = plugin
	}
}

func (m *PluginManager) Get(pluginName string) types.Plugin {
	local, ok := m.localPlugin[pluginName]
	if ok {
		return local
	}
	share, ok := m.pluginInfo[pluginName]
	if !ok {
		return nil
	}
	return &PluginProxy{
		share,
	}
}

func (m *PluginManager) InvokeFunction(pluginName string, funcName string, args []interface{}) ([]interface{}, error) {
	plugin := m.Get(pluginName)
	if plugin == nil {
		return nil, errors.New(fmt.Sprintf("%s插件不存在", pluginName))
	}
	return plugin.InvokeFunction(funcName, args)
}
