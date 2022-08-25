package utils

import (
	"im-center/service/connect/internal/types"
	"testing"
)

func TestBuildMsg(t *testing.T) {
	info := BuildMsg(&types.SendOneMsgReq{
		SenderType:       "1",
		SenderId:         "小民",
		SenderDeviceId:   "web",
		ReceiverId:       "415028502295743063",
		ReceiverDeviceId: "web",
		ParentId:         "",
		SendTime:         "111111",
		MsgType:          "text",
		MsgContent: 	  types.Text{
			Content: "hello",
		},
		//TextMsg:          types.Text{
		//	Content: "hello",
		//},
	})
	t.Logf("%+v",info)
}
