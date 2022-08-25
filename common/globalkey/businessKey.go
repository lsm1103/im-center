package globalkey

var (
	ObjSaveType = "obj"
	BlockSaveType = "block"
	//原始dicom的描述json文件id
	DescribeJsonFid = "describeJsonFid"
	//筛选后dicom的描述json文件id
	DescribeFilterJsonFid = "describeFilterJsonFid"

	EventTypeOnline = "online"
	EventTypeOffline = "offline"
	EventTypeRenewal = "renewal"

	// 用户连接超时时间/秒 600 * 60
	HeartbeatExpirationTime = 60
)
