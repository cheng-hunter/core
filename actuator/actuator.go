package actuator

import (
	"github.com/cheng-hunter/core/root"
)

//PluginActuator 插件执行器
type PluginActuator struct {
	core *root.PluginCore
}

func NewPluginActuator(pluginPath string) *PluginActuator {
	return &PluginActuator{
		core: &root.PluginCore{
			PluginIndex: 0,
			PluginPath:  pluginPath,
		},
	}
}

func (p *PluginActuator) Start() {
	//启动核心插件
	p.core.Start()
}
