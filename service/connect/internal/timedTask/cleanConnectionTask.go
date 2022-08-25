/**
* Created by GoLand.
* User: link1st
* Date: 2019-07-31
* Time: 15:17
 */

package timedTask

import (
	"fmt"
	"im-center/service/connect/core/connectManager"
	"runtime/debug"
)

// 清理超时连接
func cleanConnection(param interface{}) (result bool) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("ClearTimeoutConnections stop", r, string(debug.Stack()))
		}
	}()

	fmt.Println("定时任务-清理超时连接...")
	result = true
	connectManager.GetCM().ClearTimeoutConnections()
	//param.(*core.ConnectServer).cm.ClearTimeoutConnections()
	return
}
