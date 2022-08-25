package all

type (
	UserItem struct {
		Id          int64  `json:"id"`          // 用户id
		Nickname    string `json:"nickname"`    // 昵称
		RealName    string `json:"realName"`    // 真实姓名
		Sex         int64  `json:"sex"`         // 性别，0:未知；1:男；2:女
		Ico         string `json:"ico"`         // 用户头像
		Status      int64  `json:"status"`      // 用户状态
		CreateTime string `json:"create_time"` // 创建时间
		UpdateTime  string `json:"update_time"`  // 更新时间
	}

	//PredictItem Predict
	//DictionaryItem Dictionary
	//FileItem File
)
