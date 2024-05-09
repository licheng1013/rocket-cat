package ctx

import (
	"fmt"
	"testing"
	"time"
)

func TestCtx(t *testing.T) {
	fmt.Println("main goroutine id:", localCtx.GetGID())
	Go(1)
	Go(2)
	Go(3)
	time.Sleep(1 * time.Second) // 等待1秒加入上下文
	fmt.Println(len(localCtx.Ctx))
	time.Sleep(5 * time.Second)
	fmt.Println(len(localCtx.Ctx))
}

// 定义一个上下文,测试是否线程安全
var localCtx = NewLocalContext()

func Go(num int) {
	go func() {
		for i := 0; i < 3; i++ {
			localCtx.SaveCtx(num)
			// 打印
			time.Sleep(time.Second)
			// 打印协程id
			fmt.Println("上下文数据:", localCtx.GetCtx(), "协程id:", localCtx.GetGID())
		}
		localCtx.ClearCtx()
	}()
}
