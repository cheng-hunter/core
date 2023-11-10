package manager

import (
	"github.com/cheng-hunter/core/types"
	"github.com/cheng-hunter/core/util"
)

type PluginProxy struct {
	types.SharePluginInfo
}

func (p *PluginProxy) Activate() types.PluginInfo {
	return p.PluginInfo
}
func (p *PluginProxy) Start() {

}
func (p *PluginProxy) Destroy() {

}

func (p *PluginProxy) InvokeFunction(name string, args ...interface{}) ([]interface{}, error) {
	return util.InvokeFunction(p.Id, p.Name, name, args)
}
