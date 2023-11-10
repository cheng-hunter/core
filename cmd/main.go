package main

import (
	"fmt"
	"github.com/cheng-hunter/core/actuator"
	"github.com/cheng-hunter/core/manager"
)

func main() {
	actuator := actuator.NewPluginActuator("/home/hunter/plugin")
	actuator.Start()
	fmt.Println(manager.DefaultPluginManager.Get("root"))
}
