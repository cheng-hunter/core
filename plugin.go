package core

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"core/types"
)

//PluginCore 插件核心部分
type PluginCore struct {

}

func (core *PluginCore)Load()(string,int32){
	return "root",0
}

func (core *PluginCore)UnLoad(){

}

func (core *PluginCore)AddFunction(name string, function interface{}) error{
	return nil
}

func (core *PluginCore)InvokeFunction(name string, args []interface{}) ([]interface{}, error){
	return nil,nil
}
//Open 启动
//通过执行./plugin open
func (core *PluginCore)Open(path string) (map[string]types.Plugin, error){
	cmd := exec.Command(path)
	cmd.Start()
	return nil,nil
}

//Open 启动
func (core *PluginCore)Lookup(pluginName string) (types.Plugin, error){
	return nil,nil
}


//PluginActuator 插件执行器
type PluginActuator struct {
	core      PluginCore
	pluginMap map[string]types.Plugin
}


func (p *PluginActuator)LoadAllPlugin(pluginPath string)  {
	root:=pluginPath
	if root==""{
		root=types.DefaultPluginPath
	}
	err:=filepath.WalkDir(root,p.readPlugin)
	if err != nil {
		fmt.Println(err)
	}
}

func (p *PluginActuator)readPlugin(path string, info os.DirEntry, err error) error {
	if err != nil {
		fmt.Println(err)
		return nil
	}
	if info.IsDir() {
		fmt.Printf("Found directory: %s\n", path)
	} else {
		// 加载插件
		plugins, e:= p.core.Open(path) // 假设插件文件名为pluginB.so，编译为二进制文件后缀为.so或.dll
		if e != nil {
			fmt.Println("Failed to load plugin B:", err)
			return e
		}
		fmt.Println("Plugin B loaded successfully.")
		for k, v := range plugins {
			fmt.Println("%s,Plugin B loaded successfully.%s",k,v)
		}
		// 获取PluginB的实例
		//pluginBInstance, err := p.core.Lookup("PluginB") // 获取插件B的实例，可以通过类型断言进行类型转换操作
		//if err != nil {
		//	fmt.Println("Failed to lookup Plugin B:", err)
		//	return
		//}
	}
	return nil
}