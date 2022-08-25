package fileCenter

import "sync"

var m = sync.Map{}
var wg sync.WaitGroup
var ScriptRootPath string

const existFileNodes = "existFileNodes-%s"

//const ScriptRootPath = "/Users/xm/Desktop/work_project/im-center/dataFileScript"

type (
	// 配置文件结构
	//ServiceConfig struct {
	//	Port     int
	//	Address  string
	//	StoreDir string
	//}
	// 文件分片结构
	//FilePart struct {
	//	Fid   string `form:"fid"`   // 操作文件ID，随机生成的UUID
	//	Index int    `form:"index"` // 文件切片序号
	//	Data  []byte `form:"data"`  // 分片数据
	//}
	// 客户端传来的文件元数据结构
	ClientFileMetadata struct {
		Fid          string            `json:"fid,omitempty,optional"`  // 操作文件ID，随机生成的UUID
		Filesize     int64             `json:"filesize,optional"`       // 文件大小（字节单位）
		Filename     string            `json:"filename,optional"`       // 文件名称
		SliceNum     int               `json:"slice_num,optional"`      // 切片数量
		Md5sum       string            `json:"md5_sum,optional"`        // 文件md5值
		ModifyTime   string            `json:"modify_time,optional"`    // 文件修改时间
		FileDataType string            `json:"file_data_type,optional"` // 文件数据类型
		Info         map[string]string `json:"info,optional"`           // 文件附属信息
		UserId       string            `json:"user_id,optional"`        // 用户id
	}
	// 服务端保存的文件元数据结构
	ServerFileMetadata struct {
		ClientFileMetadata        // 隐式嵌套
		State              string // 文件状态，目前有uploading、downloading和active
	}
	SliceSeq struct {
		Slices []int // 需要重传的分片号
	}
	// 文件列表单元结构
	FileInfo struct {
		Filename string // 文件名
		Filesize int64  // 文件大小
		Filetype string // 文件类型（目前有普通文件和切片文件两种）
	}
	// 文件列表结构
	ListFileInfos struct {
		Files []FileInfo
	}
	// 执行脚本的入参结构
	ExecScriptReq struct {
		Script   string            `json:"script"`
		Fid      string            `json:"fid"`
		FileType string            `json:"file_type"`
		Args     map[string]string `json:"args"`
	}
	CompressionReq struct {
		Fid      string `json:"fid"`
		FileType string `json:"file_type"`
	}
)
