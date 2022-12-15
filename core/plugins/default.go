package plugins

import (
	"core/core"
	"fmt"
	"github.com/xtaci/kcp-go/v5"
	"log"
	"time"
)

type CountLinkPlugin struct {
}

func (c *CountLinkPlugin) Invok(app *core.App) {
	fmt.Println("连接数: " + fmt.Sprint(len(app.Conns)))
}

type Heartbeat struct {

}

var timeOutMap = map[uint32]int64{}

// Invok TODO 这段代码暂时不知道咋写，写在那一块!
func (h Heartbeat) Invok(app *core.App) {
	conn := app.Conns[len(app.Conns)-1]
	udpSession := conn.(*kcp.UDPSession)
	sessionId := udpSession.GetConv()
	fmt.Println("SessionId: ", sessionId)
	if app.Result.GetHeartbeat() {
		timeOutMap[sessionId] = time.Now().UnixMilli()
		return
	}
	if timeOutMap[sessionId] == 0 {
		timeOutMap[sessionId] = time.Now().UnixMilli()
	}
	// 最大超时3秒
	if time.Now().UnixMilli()-timeOutMap[sessionId] > 3000 {
		err := conn.Close()
		if err != nil {
			log.Println(err)
		}
	}
}
