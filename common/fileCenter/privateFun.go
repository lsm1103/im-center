package fileCenter

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path"
	"im-center/common/globalkey"
	"strings"
	"time"

	"im-center/common/dynamicConfig/cfgByWatchFile"
	"im-center/common/tool"
)

// 获取obj存储分片文件元数据文件路径
func (fc *FileCenter) getMetadataFilepath(fid string) string {
	info := strings.Split(fid, "-")
	if len(info) < 3 {
		return ""
	}
	return path.Join(fc.NodeCfg.OssNodes[info[0]].DiskMountedOn, "slice", info[2]+".slice")
	//return fmt.Sprintf("%s/slice/%s.slice", fc.NodeCfg.OssNodes[info[0]].DiskMountedOn, info[2])
}

// 获取obj存储分片文件路径
func (fc *FileCenter) getSaveFilepath(fid string) string {
	info := strings.Split(fid, "-")
	if len(info) < 3 {
		return ""
	}
	return path.Join(fc.NodeCfg.OssNodes[info[0]].DiskMountedOn, "sliceFile", info[1][:4], info[1][4:6], info[1][6:8], info[2])
}

// 获取obj存储一般文件元数据文件路径
func (fc *FileCenter) getNormalFilepath(fid string) string {
	info := strings.Split(fid, "-")
	if len(info) < 2 {
		return ""
	}
	name_ := info[1]
	if len(info) == 3 {
		name_ = info[2]
	}
	return path.Join(fc.NodeCfg.OssNodes[info[0]].DiskMountedOn, "normalFile", name_)
}

//通过fid获取block存储相应文件的路径
func (fc *FileCenter) getBlockFilepath(fid string) string {
	info := strings.SplitN(fid, "-", 2)
	if info == nil {
		return ""
	}
	node, ok := fc.NodeCfg.BlockNodes[info[0]]
	if !ok {
		return ""
	}
	return path.Join(node.DiskMountedOn, info[1])
}

//通过元数据确认文件是否存在
func (fc *FileCenter) checkFileExistByMetadata(fid string) (bool, error) {
	sfPath := fc.getSaveFilepath(fid)
	if ok := IsFile(sfPath); !ok {
		return false, errors.New("该文件不存在")
	}
	mPath := fc.getMetadataFilepath(fid)
	// 加载元数据信息
	metadata, err := LoadMetadata(mPath)
	if err != nil {
		return false, err
	}
	// 校验fid和filename是匹配的
	if metadata.Fid != fid {
		return false, errors.New("文件名和fid对不上，请确认后重试")
	}
	return true, nil
}

//获取所有节点（配置的文件夹）
func (fc *FileCenter) objGetAllNode() *cfgByWatchFile.Config {
	//获取所有节点（配置的文件夹）
	//DictionaryModel := model.NewDictionaryModel(sqlx.NewMysql(c.DB.DataSource), c.Cache)
	//var val = []model.DictionaryItem{}
	//err := cacheHandle.NewCecheHandle(c.Cache).Get(&val, globalkey.CacheOAuth2CodeKey, func(val interface{}) error {
	//	val, err := DictionaryModel.FindAll(&tool.GetsReq{
	//		Query:    nil,
	//		OrderBy:  "",
	//		Sort:     "",
	//		Current:  0,
	//		PageSize: 0,
	//	})
	//	return err
	//})
	//if err != nil {
	//	return
	//}
	return fc.NodeCfg
}

//检查所有节点该文件是否存在
func (fc *FileCenter) objCheckFileExistInAllNode(fid string) []*cfgByWatchFile.Node {
	if len(fc.NodeCfg.OssNodes) == 1 {
		for _, v := range fc.NodeCfg.OssNodes {
			return []*cfgByWatchFile.Node{&v}
		}
	}

	//并发检查所有节点该文件是否存在
	for _, node := range fc.NodeCfg.OssNodes {
		wg.Add(1)
		//fmt.Println(k)
		go CheckFileExist(&node, fid)
	}
	wg.Wait()

	val, _ := m.Load(fmt.Sprintf(existFileNodes, fid))
	if val != nil {
		return val.([]*cfgByWatchFile.Node)
	}
	return nil
}

//获取一个可以存储的节点，需考虑磁盘容量、请求流量
func (fc *FileCenter) objGetStorageNode() *cfgByWatchFile.Node {
	node_num := len(fc.NodeCfg.OssNodes)
	nodes := []cfgByWatchFile.Node{}
	for _, v := range fc.NodeCfg.OssNodes {
		nodes = append(nodes, v)
	}
	if node_num == 1 {
		return &nodes[0]
	}

	// todo 先用随机数做一个简单的，后面可以用hash啊、加权啊实现更严谨的
	//rand.Seed(time.Now().UnixNano())	//默认加锁
	r := rand.New(rand.NewSource(time.Now().UnixNano())) //不加锁
	return &nodes[r.Intn(node_num)]
}

//解析该文件唯一id，确定存储在哪个节点；
func (fc *FileCenter) objGetNodeByFid(fid string) *cfgByWatchFile.Node {
	nodeId := strings.Split(fid, "-")[0]
	node := fc.NodeCfg.OssNodes[nodeId]
	return &node
}

//获取一个可以存储的节点（来实现扩容）；业务方每次上传时，先请求该接口动态获取块，然后存入磁盘容量富裕的块
func (fc *FileCenter) blockGetStorageNode() *cfgByWatchFile.Node {
	node_num := len(fc.NodeCfg.BlockNodes)
	nodes := []cfgByWatchFile.Node{}
	for _, v := range fc.NodeCfg.BlockNodes {
		nodes = append(nodes, v)
	}
	if node_num == 1 {
		return &nodes[0]
	}

	// todo 先用随机数做一个简单的，后面可以用hash啊、加权啊实现更严谨的
	//rand.Seed(time.Now().UnixNano())	//默认加锁
	r := rand.New(rand.NewSource(time.Now().UnixNano())) //不加锁
	return &nodes[r.Intn(node_num)]
}

/*执行脚本（可对文件/文件夹）
@eg："python /Users/xm/Desktop/11/dicomToDescribeJson.py -jp /Users/xm/Desktop/11/ct_1915259_20220524_.json -ojp /Users/xm/Desktop/ct_1915259_20220524.json -ot del -ft seg -plist \"/Users/xm/Desktop/ct_1915259_20220524/mask2seg_9a6a10388d8d414490aab73d4017b38b.dcm\",\"/Users/xm/Desktop/ct_1915259_20220524/mask2seg_9d6efd33e907414ba26f94a981675a3d.dcm\",\"/Users/xm/Desktop/ct_1915259_20220524/mask2seg_0676fb6c6a3a422095b71933c9112b25.dcm\",\"/Users/xm/Desktop/ct_1915259_20220524/mask2seg_93432e8cb60d4981b37107cbb641d258.dcm\",\"/Users/xm/Desktop/ct_1915259_20220524/mask2seg_8127198c65e74f8092c7fb3353437846.dcm\",\"/Users/xm/Desktop/ct_1915259_20220524/mask2seg_ae97426b736742bf9c7380e8a7d0d3b8.dcm\""
*/
func (fc *FileCenter) execScript(script string, filePath string, args map[string]string) error {
	// 构建sh
	var argsStr string
	doTime := tool.GetNowTime(tool.TimeSlimFmt)
	if script == "dicomToDescribeJson" {
		argsStr = fmt.Sprintf("-jp %s -p %s -ojp %s -ot %s -ft %s -plist %s -fid %s -doTime %s",
			args["jsonPath"],
			filePath,
			args["oldJsonPath"],
			args["operation"],
			args["fileType"],
			args["pathList"],
			args["fid"],
			doTime,
		)
	}
	sh := fmt.Sprintf("(cd %s && ./dicomToDescribeJson %s)", ScriptRootPath, argsStr)
	//sh := fmt.Sprintf("(cd %s && exec python3 dicomToDescribeJson.py %s)", ScriptRootPath, argsStr)
	cmd := exec.Command("/bin/sh", "-c", sh)
	err := cmd.Start()
	fmt.Printf("[fc.execScript-pid:%d,err:%+v],\nsh:%s\n", cmd.Process.Pid, err, sh)
	waitStartT := time.Now()
	for {
		if time.Now().Sub(waitStartT) > time.Second*8 {
			return errors.New("执行时间超时")
		} else {
			strJson, err := tool.Read2StrJson(fmt.Sprintf("%s/dicomToDescribeJson-pid/%s-%s.pid", ScriptRootPath, doTime, args["fid"]))
			if err == nil {
				switch strJson["state"] {
				case "running":
					fmt.Printf("还在运行中，等待中...\n")
				case "error":
					return errors.New(strJson["err"])
				case "success":
					return nil
				default:
					return errors.New("该此脚本执行的pid文件找不到state字段")
				}
			}
			time.Sleep(1 * time.Second)
		}
	}
	//return err
}

func dicomToDescribeJsonInit(fc *FileCenter, in *ExecScriptReq, filePath string, doTime string) (string, error) {
	// dicomToDescribeJson执行命令初始化
	if _, ok := in.Args["jsonPath"]; !ok {
		return "", errors.New("jsonPath不能为空")
	}
	oldJsonPath := ""
	if in.Args["operation"] != "" {
		//从头构建模式
		if _, ok := in.Args["oldJsonPath"]; !ok {
			return "", errors.New("当为从头构建模式，oldJsonPath不能为空")
		}
		filePath = ""
		oldJsonPath = fc.getNormalFilepath(in.Args["oldJsonPath"])
	}
	jsonPath := fc.getNormalFilepath(in.Args["jsonPath"])

	return fmt.Sprintf("./dicomToDescribeJson -jp %s -p %s -ojp %s -ot %s -ft %s -plist %s -fid %s -doTime %s",
		jsonPath,
		filePath,
		oldJsonPath,
		in.Args["operation"],
		in.Args["fileType"],
		in.Args["pathList"],
		in.Fid,
		doTime), nil
	//return fmt.Sprintf("python3 dicomToDescribeJson.py -jp %s -p %s -ojp %s -ot %s -ft %s -plist %s -fid %s",
	//	jsonPath,
	//	filePath,
	//	oldJsonPath,
	//	in.Args["operation"],
	//	in.Args["fileType"],
	//	in.Args["pathList"],
	//	in.Fid ), nil
}

func (fc *FileCenter) execScript_(in *ExecScriptReq) (err error) {
	var filePath string
	if in.FileType == globalkey.BlockSaveType {
		filePath = fc.getBlockFilepath(in.Fid)
	} else if in.FileType == globalkey.ObjSaveType {
		filePath = fc.getSaveFilepath(in.Fid)
	} else {
		return errors.New("不支持该文件存储类型")
	}
	//todo 判断该脚本是否为命令行执行
	if in.Script == "dicomToDescribeJson" {
		execSh := ""
		doTime := tool.GetNowTime(tool.TimeSlimFmt)
		if in.Script == "dicomToDescribeJson" {
			in.Args["fid"] = in.Fid
			if execSh, err = dicomToDescribeJsonInit(fc, in, filePath, doTime); err != nil {
				return err
			}
		}
		if execSh == "" {
			return errors.New("该脚本有问题")
		}
		sh := fmt.Sprintf("cd %s && %s", ScriptRootPath, execSh)
		cmd := exec.Command("/bin/sh", "-c", sh)
		err = cmd.Start()
		fmt.Printf("[fc.execScript-pid:%d,err:%+v],\nsh:%s\n", cmd.Process.Pid, err, sh)
		waitStartT := time.Now()
		for {
			if time.Now().Sub(waitStartT) > time.Second*10 {
				return errors.New("执行时间超时")
			} else {
				strJson, err := tool.Read2StrJson(fmt.Sprintf("%s/dicomToDescribeJson-pid/%s-%s.pid", ScriptRootPath, doTime, in.Fid))
				if err == nil {
					switch strJson["state"] {
					case "running":
						fmt.Print("还在运行中，等待中...\n")
					case "error":
						return errors.New(strJson["err"])
					case "success":
						return nil
					default:
						return errors.New("该此脚本执行的pid文件找不到state字段")
					}
				}
				fmt.Printf("运行中，strJson:%+v，err:%+v，\n", strJson, err)
				time.Sleep(1 * time.Second)
			}
		}
	} else {
		//可以直接调用
	}
	return nil
}

//通过路径获取obj存储的文件Fid
func (fc *FileCenter) GetObjFidByPath(path string) string {
	// todo 判断normal/分片
	//info := strings.Split(fid,"-")
	//path.Join(fc.NodeCfg.OssNodes[info[0]].DiskMountedOn, "sliceFile", info[1][:4], info[1][4:6], info[1][6:8], info[2])
	return ""
}

//通过路径获取block存储的文件Fid
func (fc *FileCenter) GetBlockFidByPath(path string) string {

	return path
}

func (fc *FileCenter) GetAiSourcePath(path string, nodeId string) string {
	return strings.Replace(path, fc.NodeCfg.BlockNodes[nodeId].DiskMountedOn+"/", "", 1)
}

// 加载元数据文件信息
func (fc *FileCenter) LoadMetadataByFid(fid string) (*ServerFileMetadata, error) {
	return LoadMetadata(fc.getMetadataFilepath(fid))
}

// 写元数据文件信息
func (fc *FileCenter) StoreMetadataByFid(fid string, metadata *ServerFileMetadata) error {
	return StoreMetadata(fc.getMetadataFilepath(fid), metadata)
}

// 更新元数据文件信息
func (fc *FileCenter) UpdateMetadataByFid(fid string, upData map[string]string) error {
	filePath := fc.getMetadataFilepath(fid)
	metadata, err := LoadMetadata(filePath)
	if err != nil {
		return err
	}
	metadata.Info = map[string]string{}
	for k, v := range upData {
		metadata.Info[k] = v
	}
	err = StoreMetadata(filePath, metadata)
	return err
}

func (fc *FileCenter) RmStartUploadSliceByFid(fid string) error {
	err := os.RemoveAll(fc.getSaveFilepath(fid))
	if err != nil {
		return err
	}
	err = os.Remove(fc.getMetadataFilepath(fid))
	return err
}

func (fc *FileCenter) sliceClearTimeOut(fid string) {
	fmt.Printf("文件%s设置上传时限：%+v", fid, sliceClearExpiryDate)
	time.Sleep(sliceClearExpiryDate)
	metadata, err := fc.LoadMetadataByFid(fid)
	if err != nil {
		fmt.Printf("文件%s超过设置的上传时限，进行处理时加载元数据失败", fid)
		return
	}
	if metadata.State == "uploading" {
		fmt.Printf("文件%s超过设置的上传时限，进行删除处理", fid)
		err = os.RemoveAll(fc.getSaveFilepath(fid))
		if err != nil {
			fmt.Printf("文件%s超过设置的上传时限，进行处理时删除真实数据失败", fid)
			return
		}
		err = os.Remove(fc.getMetadataFilepath(fid))
		if err != nil {
			fmt.Printf("文件%s超过设置的上传时限，进行处理时删除元数据失败", fid)
		}
		fmt.Printf("文件%s超过设置的上传时限，进行删除处理成功", fid)
	}
	fmt.Printf("文件%s上传时限结束", fid)
}
