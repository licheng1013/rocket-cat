package plugins

import (
	"core/core"
	"fmt"
)



type CountLinkPlugin struct {
	Count int64
}

func (c *CountLinkPlugin) Invok(app *core.App) {
	c.Count++
	fmt.Println("连接数: "+fmt.Sprint(c.Count))
}
