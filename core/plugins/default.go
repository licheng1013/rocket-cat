package plugins

import (
	"core/core"
	"fmt"
)



type CountLinkPlugin struct {
}

func (c *CountLinkPlugin) Invok(app *core.App) {
	fmt.Println("连接数: "+fmt.Sprint(len(app.Conns)))
}

type Heartbeat struct {

}