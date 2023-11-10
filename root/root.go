package root

import (
	"errors"
	"fmt"
	"github.com/cheng-hunter/core/manager"
	"github.com/cheng-hunter/core/types"
	"github.com/cheng-hunter/core/util"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"sync/atomic"
)

//PluginCore 核心插件部分
type PluginCore struct {
	PluginPath  string
	PluginIndex int32
}

func (core *PluginCore) Activate() types.PluginInfo {
	return types.PluginInfo{Name: "root", Version: 0}
}

func (core *PluginCore) Destroy() {

}

func (core *PluginCore) InvokeFunction(name string, args ...interface{}) ([]interface{}, error) {
	switch name {
	case "startPlugin":
		a0, ok := args[0].(string)
		if !ok {
			return nil, errors.New("参数转换失败")
		}
		a1, ok := args[1].(types.PluginInfo)
		if !ok {
			return nil, errors.New("参数转换失败")
		}
		err := core.startPlugin(a0, a1)
		if err != nil {
			return nil, err
		}
	case "testAdd":
		a, ok := args[0].(string)
		if !ok {
			return nil, errors.New("参数转换失败")
		}
		rlt := core.testAdd(a)
		return []interface{}{rlt}, nil
	default:

	}
	return nil, nil
}

func (core *PluginCore) testAdd(test string) string {
	rlt := test + "999999"
	return rlt
}

func (core *PluginCore) startPlugin(pluginId string, info types.PluginInfo) error {
	return util.StartPlugin(pluginId, info.Name)
}

//Start 核心插件 激活其他插件，启动其他插件
func (core *PluginCore) Start() {
	root := core.PluginPath
	if root == "" {
		root = types.DefaultPluginPath
	}
	manager.DefaultPluginManager.AddLocal("root", core)
	err := filepath.WalkDir(root, core.activatePlugin)
	if err != nil {
		fmt.Println(err)
	}
}

func (core *PluginCore) activatePlugin(path string, info os.DirEntry, err error) error {
	if err != nil {
		fmt.Println(err)
		return nil
	}
	if info.IsDir() {
		fmt.Printf("Found directory: %s\n", path)
	} else {
		//激活插件
		atomic.AddInt32(&core.PluginIndex, 1)
		cmd := exec.Command(path)
		pluginId := strconv.Itoa(int(core.PluginIndex))
		cmd.Env = []string{
			fmt.Sprintf("%s=%s", types.PluginIdEnvKey, pluginId),
		}
		err = cmd.Run()
		if err != nil {
			return errors.New(fmt.Sprintf("Failed to Activate plugin:%v", err))
		}
		rlt, err := util.ActivatePlugin(pluginId)
		if err != nil {
			return errors.New(fmt.Sprintf("Failed to Activate plugin:%v", err))
		}
		fmt.Sprintf("Activate plugin successfully:%s,%s", pluginId, path)
		//插件添加
		for _, share := range rlt {
			//插件启动
			if e := core.startPlugin(pluginId, share); e == nil {
				manager.DefaultPluginManager.Add(share.Name, types.SharePluginInfo{
					PluginInfo: share,
					Id:         pluginId,
				})
			}
		}
	}
	return nil
}
