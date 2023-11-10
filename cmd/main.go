package main

import (
	"fmt"
	"github.com/cheng-hunter/core/actuator"
	"github.com/cheng-hunter/core/manager"
)

func main() {
	actuator := actuator.NewPluginActuator("/home/hunter/plugin")
	actuator.Start()
	rlt, err := manager.DefaultPluginManager.Get("root").InvokeFunction("testAdd", "sssssss")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(rlt)
	}

}
