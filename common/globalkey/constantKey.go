package globalkey

/**
global constant key
*/

// 用户状态
var UserDel int64 = -2  //已删除
var UserFreeze int64 = -1 //已冻结
var audited int64 = 1 //待审核
var normal int64 = 2 //正常
var ThirdPartyRegistLogin int64 = 3 //第三方直接注册登入
var SuperAdmin int64 = 9 //超管

// 设备状态
//var Del int64 = -2  //删除
//var Disable int64 = -1 //禁用
//var audited int64 = 1 //待审核
var Offline int64 = 2	//离线
var Online int64 = 3	//在线

//数据状态
var Del int64 = -2  //删除
var Disable int64 = -1 //禁用、弃用
var Enable int64 = 2 //启用

//待办状态
//var Del int64 = -2	//删除
var Pause int64 = -1	//暂停
var Activated int64 = 1	//待启动
var Ongoing int64 = 2	//进行中
var Finish int64 = 3	//完成

//好友状态
var Defriend int64 = -2	//拉黑
var Refused int64 = -1	//拒绝
var Applying int64 = 1	//申请中
var Agree int64 = 2	//同意

//消息状态
var ReceiverDel int64 = -3	//接收者删除
var SenderDel int64 = -2	//发送者删除
var Withdraw int64 = -1	//撤回
var Untreated int64 = 1	//未处理
var Readed int64 = 2	//已读

//用户身份认证方式
var Oauth2Phone string = "phone"			//手机号
var Oauth2Email string = "email"			//邮箱
var Oauth2IdCard string = "idCard"  		//身份证
var Oauth2ThirdParty string = "thirdParty"  //第三方

//用户oauth2认证方式
var Oauth2AuthorizationCode string = "authorization_code"			//授权码式
var Oauth2Implicit string = "implicit"							//隐藏式
var Oauth2Password string = "password"  							//密码式
var Oauth2ClientCredentials string = "client_credentials"  		//客户端凭证式

// 群组类型
var Department int64 = 1	// 部门
var UserGroup int64 = 2		// 用户组
var Group int64 = 3			// 群组
var Circle int64 = 4		// 圈子
var Topics int64 = 5		// 话题

var GroupType = map[int64]string{
	Department : "部门",
	UserGroup : "用户组",
	Group : "群组",
	Circle : "圈子",
	Topics : "话题",
}

var UserStatus = map[int64]string{
	UserDel : "已删除",
	UserFreeze : "已冻结",
	audited : "待审核",
	normal : "正常",
	ThirdPartyRegistLogin : "第三方直接注册登入",
	SuperAdmin : "超管",
}
var DeviceStatus = map[int64]string{
	Offline : "离线",
	Online : "在线",
}
var DataStatus = map[int64]string{
	Del : "删除",
	Disable : "禁用/弃用",
}
var TodoStatus = map[int64]string{
	Pause : "暂停",
	Activated : "待启动",
	Ongoing : "进行中",
	Finish : "完成",
}
var FriendStatus = map[int64]string{
	Defriend : "拉黑",
	Refused : "拒绝",
	Applying : "申请中",
	Agree : "同意",
}
var MsgStatus = map[int64]string{
	ReceiverDel : "接收者删除",
	SenderDel : "发送者删除",
	Withdraw : "撤回",
	Untreated : "未处理",
	Readed : "已读",
}
