package main

import (
	"fmt"
	"time"

	"github.com/Miuzarte/biligo"
	"github.com/gorilla/websocket"
	"github.com/tidwall/gjson"
)

func main() {
	const roomId = 25788785

RETRY:
	lms := biligo.NewLiveMsgStream(roomId)
	for body, err := range lms.RunIter() {
		if err != nil {
			if websocket.IsCloseError(err,
				websocket.CloseNormalClosure,
				websocket.CloseAbnormalClosure) {
				break
			}
			fmt.Printf("failed to listen live msg: %v\n", err)
			<-time.After(time.Second * 10)
			goto RETRY
		}
		fmt.Println(body)

		// 解析消息
		pkt := gjson.Parse(body)
		switch pkt.Get("cmd").String() {
		case biligo.LIVE_MSG_STREAM_LIVE: // 直播开始
		case biligo.LIVE_MSG_STREAM_PREPARING: // 直播准备中 (结束)
		case biligo.LIVE_MSG_STREAM_CHANGE: // 房间信息变更
		case biligo.LIVE_MSG_STREAM_WARNING: // 警告
		case biligo.LIVE_MSG_STREAM_CUT_OFF: // 切断
		}
	}
}
