package plugins

import (
	"fmt"
	"github.com/xtaci/kcp-go/v5"
	"time"
)

type CountLinkPlugin struct {
}

func (c *CountLinkPlugin) Invok(app Meta) {
	fmt.Println("链接数:",len(app.App))
}

type Heartbeat struct {

}


// Invok TODO 这段代码暂时不知道咋写，写在那一块!
func (h Heartbeat) Invok(app Meta) {
	udpSession := app.Conn.(*kcp.UDPSession)
	sessionId := udpSession.GetConv()
	fmt.Println("SessionId: ", sessionId)
	app.TimeOutMap.Store(sessionId,time.Now().UnixMilli())
}
