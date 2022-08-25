package fileCenter

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
	"golang.org/x/xerrors"

	"im-center/common/compression"
	"im-center/common/dynamicConfig/cfgByWatchFile"
	"im-center/common/globalkey"
	"im-center/common/result"
	"im-center/common/uniqueid"
	"im-center/common/xerr"
	"im-center/model"
)


type (
	UploadReq struct {
		User_id   string        `form:"user_id,optional"` // 用户id
		File_name string        `form:"file_name"`        // 文件名
		File_size int64         `form:"file_size"`        // 大小
		File_type string        `form:"file_type"`        // 类型
		File_md5  string        `form:"file_md5"`         // md5
		File      io.ReadCloser `form:"file"`             // 文件
	}
	FileCenter struct {
		CanWork   bool
		NodeCfg   *cfgByWatchFile.Config
		FileModel model.FileModel
	}
)

//const sliceClearExpiryDate = 2*time.Minute
const sliceClearExpiryDate = 12 * time.Hour

func init() {
	getwd, err := os.Getwd()
	if err != nil {
		panic(fmt.Sprintf("获取当前执行地址失败：%+v", err))
	}
	ScriptRootPath = fmt.Sprintf("%s/dataFileScript", getwd)
	fmt.Println("当前执行地址:", getwd, ScriptRootPath)
}

func NewFileCenter(fileModel model.FileModel) FileCenter {
	return FileCenter{
		CanWork:   true,
		FileModel: fileModel,
	}
}

func (fc *FileCenter) GetMetadataFilepath(fid string) string {
	return fc.getMetadataFilepath(fid)
}
func (fc *FileCenter) GetSaveFilepath(fid string) string {
	return fc.getSaveFilepath(fid)
}
func (fc *FileCenter) GetNormalFilepath(fid string) string {
	return fc.getNormalFilepath(fid)
}
func (fc *FileCenter) ExecScript_(script string, filePath string, args map[string]string) error {
	return fc.execScript(script, filePath, args)
}
func (fc *FileCenter) CopyObj2BlockFileByFid(fid string, dstNode string, fun func(nodePath string, fPath string, fromNodePath string) string) (dstPath string, err error) {
	fPath := fc.getSaveFilepath(fid)
	node := fc.NodeCfg.BlockNodes[dstNode]
	dstPath = path.Join(node.DiskMountedOn, fPath)
	fromNode := fc.objGetNodeByFid(fid)
	if fun != nil {
		dstPath = fun(node.DiskMountedOn, fPath, fromNode.DiskMountedOn)
	}
	if IsDir(dstPath) {
		return dstPath, nil
	}
	return dstPath, CopyFile(fPath, dstPath)
}

/*确认文件是否存在
@method：get
@req：fid、fileType
@resp：http-200
*/
func (fc *FileCenter) CheckFileExist(w http.ResponseWriter, r *http.Request) {
	fid := r.FormValue("fid")
	if fid == "" {
		result.ParamErrorResult(r, w, xerrors.New("fid不能为空"))
		return
	}
	fileType := r.FormValue("fileType")
	if fileType == "" {
		result.ParamErrorResult(r, w, xerrors.New("fileType不能为空"))
		return
	}
	var err error
	if fileType == globalkey.ObjSaveType {
		if nPath := fc.getNormalFilepath(fid); !IsDir(nPath) {
			//normal
			if exist, e := fc.checkFileExistByMetadata(fid); !exist {
				//slice
				err = errors.Wrapf(xerr.NewErrMsg("文件不存在"), "err:%s", e)
			}
		}
	} else if fileType == globalkey.BlockSaveType {
		if nPath := fc.getBlockFilepath(fid); !IsFile(nPath) {
			err = xerr.NewErrMsg("文件不存在")
		}
	} else {
		err = xerr.NewErrMsg("不支持该模式")
	}
	result.HttpResult(r, w, err, nil)
}

/*上传单个文件（小于等于1m的小文件）
@method：post
@req：file
@resp：http-200
*/
func (fc *FileCenter) Upload(w http.ResponseWriter, r *http.Request) {
	log := logx.WithContext(r.Context())
	file, handler, err := r.FormFile("file")
	if err != nil {
		result.ParamErrorResult(r, w, errors.Wrap(xerr.NewErrMsg("读取file失败"), err.Error()))
		return
	}
	defer file.Close()
	data := []byte{}
	_, err = file.Read(data)
	if err != nil {
		result.HttpResult(r, w, errors.Wrap(xerr.NewErrMsg("读取文件失败"), err.Error()), nil)
		return
	}
	md5_num := hex.EncodeToString(md5.New().Sum(data))
	log.Errorf("md5_num:%s", md5_num)

	//todo 需要支持block存储
	node := fc.objGetStorageNode()
	fid := GetFid(node.Id, md5_num) + path.Ext(handler.Filename)
	if IsFile(fc.getNormalFilepath(fid)) {
		result.HttpResult(r, w, xerr.NewErrMsg("文件已存在"), nil)
		return
	}
	if err = writeFile(file, fc.getNormalFilepath(fid)); err != nil {
		result.HttpResult(r, w, errors.Wrap(xerr.NewErrMsg("写入文件失败"), err.Error()), nil)
		return
	}
	resp := map[string]string{"fid": fid}
	result.HttpResult(r, w, nil, &resp)
}

/* ---------------------------分片上传步骤---------------------------*/
/*1、创建上传文件夹
@method：post
@req：{
	Filesize   int64
	Filename   string
	SliceNum   int
	Md5sum     string
	ModifyTime string
	FileDataType string		// 文件数据类型, "dicom"
	info       map[string]string
}
@resp：{
	Fid         string      // 操作文件ID，随机生成的UUID
	Filesize    int64       // 文件大小（字节单位）
	Filename    string      // 文件名称
	SliceNum    int         // 切片数量
	Md5sum      string      // 文件md5值
	ModifyTime  string   	// 文件修改时间
	FileDataType string		// 文件数据类型
	Info       	map[string]string
}
@如果是上传dicom等特殊的，需要 FileDataType："dicom"
需要一个延迟任务来删除创建了又长时间不上传的文件夹， asynq, ✅
*/
func (fc *FileCenter) CreateUploadDir(w http.ResponseWriter, r *http.Request) {
	// 验证参数
	var cMetadata ClientFileMetadata
	if err := httpx.Parse(r, &cMetadata); err != nil {
		result.ParamErrorResult(r, w, err)
		return
	}
	// 处理类似dicom文件这类无法在前端获取md5的文件
	if cMetadata.Md5sum == "" {
		if cMetadata.FileDataType == "" {
			result.ParamErrorResult(r, w, xerr.NewErrMsg("当Md5sum为空时，FileDataType不能为空"))
			return
		}
		cMetadata.Md5sum = uniqueid.GenUid()
	}
	cMetadata.Fid = GetFid(fc.objGetStorageNode().Id, cMetadata.Md5sum)
	// 通过元数据文件，检查文件是否已存在
	metadataPath := fc.getMetadataFilepath(cMetadata.Fid)
	if IsFile(metadataPath) {
		//不同用户上传md5相同的文件的处理
		metadata, err := LoadMetadata(metadataPath)
		if err == nil {
			if cMetadata.UserId != metadata.UserId {
				_, err = fc.FileModel.Insert(nil, &model.File{
					Id:           metadata.Fid,
					UserId:       cMetadata.UserId,
					FileType:     globalkey.ObjSaveType,
					FileDataType: metadata.FileDataType,
					Name:         metadata.Filename,
					Size:         strconv.FormatInt(metadata.Filesize, 10),
					SliceNum:     strconv.Itoa(metadata.SliceNum),
					Md5Sum:       metadata.Md5sum,
				})
				if err != nil {
					result.HttpResult(r, w, xerr.NewErrMsg("添加file表失败, 请重试"), nil)
					return
				}
			}
		}
		result.HttpResult(r, w, xerr.NewErrCode(xerr.ALREADY_EXISTS), &cMetadata)
		return
	}
	// 创建文件夹
	err := os.MkdirAll(fc.getSaveFilepath(cMetadata.Fid), 0766)
	if err != nil {
		result.HttpResult(r, w, xerr.NewErrMsg("创建文件夹失败"), nil)
		return
	}
	//写元数据文件
	if err = StoreMetadata(metadataPath, &ServerFileMetadata{
		ClientFileMetadata: cMetadata,
		State:              "uploading",
	}); err != nil {
		result.HttpResult(r, w, xerr.NewErrMsg("写入元数据文件失败"), nil)
		return
	}
	go fc.sliceClearTimeOut(cMetadata.Fid)
	result.HttpResult(r, w, nil, &cMetadata)
}

func (fc *FileCenter) CreateUploadDir_(cMetadata *ClientFileMetadata) (*ClientFileMetadata, error) {
	// 处理类似dicom文件这类无法在前端获取md5的文件
	if cMetadata.Md5sum == "" {
		if cMetadata.FileDataType == "" {
			return nil, xerr.NewErrMsg("当Md5sum为空时，FileDataType不能为空")
		}
		cMetadata.Md5sum = uniqueid.GenUid()
	}
	cMetadata.Fid = GetFid(fc.objGetStorageNode().Id, cMetadata.Md5sum)
	// 通过元数据文件，检查文件是否已存在
	metadataPath := fc.getMetadataFilepath(cMetadata.Fid)
	if IsFile(metadataPath) {
		return cMetadata, xerr.NewErrCode(xerr.ALREADY_EXISTS)
	}
	// 创建文件夹
	err := os.MkdirAll(fc.getSaveFilepath(cMetadata.Fid), 0766)
	if err != nil {
		return nil, xerr.NewErrMsg("创建文件夹失败")
	}
	//写元数据文件
	if err = StoreMetadata(metadataPath, &ServerFileMetadata{
		ClientFileMetadata: *cMetadata,
		State:              "uploading",
	}); err != nil {
		return nil, xerr.NewErrMsg("写入元数据文件失败")
	}
	return cMetadata, nil
}

/*2、接收分片文件（分片上传）
@method：post
@req：{
	Fid   string
	Index int
	Data  []byte
}
@resp：http-200
*/
func (fc *FileCenter) ReceiveSliceFile(w http.ResponseWriter, r *http.Request) {
	// 验证参数
	fid := r.PostFormValue("fid")
	if fid == "" {
		result.ParamErrorResult(r, w, errors.New("字段fid为空"))
		return
	}
	index := r.PostFormValue("index")
	if index == "" {
		result.ParamErrorResult(r, w, errors.New("字段index为空"))
		return
	}
	metadata, err := LoadMetadata(fc.getMetadataFilepath(fid))
	if err != nil {
		result.HttpResult(r, w, xerr.NewErrMsg("加载元数据失败"), nil)
		return
	}
	//todo 可以考虑分片号不单单为数字，也可以为字符串；考虑其他特殊类型处理，可以改造成通过 FileDataType 来执行不同的函数
	var index_ int
	if metadata.FileDataType == "dicom" && strings.HasSuffix(index, ".dcm") {
		index_, err = strconv.Atoi(strings.Split(index, ".dcm")[0])
	} else {
		index_, err = strconv.Atoi(index)
	}
	if err != nil {
		result.HttpResult(r, w, xerr.NewErrMsgf("分片号%s转化为int失败:%+v", index, err), nil)
		return
	}
	if index_ < 0 || index_ > metadata.SliceNum {
		result.HttpResult(r, w, xerr.NewErrMsg("分片号不正确,不能为负或大于总分片数"), nil)
		return
	}
	sliceFilename := path.Join(fc.getSaveFilepath(fid), index)
	if IsFile(sliceFilename) {
		err = errors.Wrapf(xerr.NewErrCode(xerr.ALREADY_EXISTS), "%s分片文件已存在，直接丢弃", sliceFilename)
	} else {
		file, _, err := r.FormFile("data")
		if err != nil {
			err = errors.Errorf("获取data出错：%+v", err)
		} else if err = writeFile(file, sliceFilename); err != nil {
			err = errors.Wrapf(xerr.NewErrCode(xerr.USER_OPERATION_ERR), "存储分片失败：%+v, fid: %+v, index: %+v\n", err, fid, index)
		}
	}
	result.HttpResult(r, w, err, nil)
	return
}

/*3、获取上传中需要重新上传的分片序号列表
todo 支持普通分片上传的文件可以通过md5来进行获取上传进度
@method：get
@req：fid、filename
@resp：{
	Slices []int
}
*/
func (fc *FileCenter) GetUploadingStat(w http.ResponseWriter, r *http.Request) {
	idType := r.FormValue("idType")
	if idType == "" {
		result.ParamErrorResult(r, w, errors.New("idType不能为空"))
		return
	}
	id := r.FormValue("id")
	if id == "" {
		result.ParamErrorResult(r, w, errors.New("id不能为空"))
		return
	}
	var metadataPath string
	if idType == "fid" {
		metadataPath = fc.getMetadataFilepath(id)
	} else if idType == "md5" {
		//todo 多节点查询
		metadataPath = fmt.Sprintf("%s/slice/%s.slice", fc.objGetStorageNode().DiskMountedOn, id)
	}
	retrySeq := SliceSeq{Slices: []int{}}
	// 校验fid是否匹配
	metadata, err := LoadMetadata(metadataPath)
	if err != nil {
		err = errors.Wrapf(xerr.NewErrCode(xerr.DATA_NOT_FIND), "fid不正确，找不到元数据，%v", err)
	} else {
		tmpDir := fc.getSaveFilepath(metadata.Fid)
		if IsDir(tmpDir) {
			retrySeq.Slices = findRetrySeq(tmpDir, metadata)
		}
	}
	result.HttpResult(r, w, err, &retrySeq)
}

/*4、合并分片文件（不真正进行合并，只计算md5进行数据准确性校验）
@method：post
@req：{
	Fid        string
	Filesize   int64
	Filename   string
	SliceNum   int
	Md5sum     string
	ModifyTime time.Time
}
@resp：{
	ServerFileMetadata
}
*/
func (fc *FileCenter) MergeSliceFiles(w http.ResponseWriter, r *http.Request) {
	// 验证参数
	var cMetadata ClientFileMetadata
	if err := httpx.Parse(r, &cMetadata); err != nil {
		result.ParamErrorResult(r, w, err)
		return
	}
	if cMetadata.UserId == "" {
		result.HttpResult(r, w, xerr.NewErrMsg("user_id不能为空"), nil)
		return
	}
	//加载元数据
	metadataPath := fc.getMetadataFilepath(cMetadata.Fid)
	metadata, err := LoadMetadata(metadataPath)
	if err != nil {
		result.HttpResult(r, w, errors.Wrapf(xerr.NewErrCode(xerr.USER_OPERATION_ERR), "加载元数据失败: %+v", err), nil)
		return
	}
	//判断是否已经存在
	if metadata.State == "active" {
		_, err = fc.FileModel.Insert(nil, &model.File{
			Id:           metadata.Fid,
			UserId:       cMetadata.UserId,
			FileType:     globalkey.ObjSaveType,
			FileDataType: metadata.FileDataType,
			Name:         metadata.Filename,
			Size:         strconv.FormatInt(metadata.Filesize, 10),
			SliceNum:     strconv.Itoa(metadata.SliceNum),
			Md5Sum:       metadata.Md5sum,
		})
		if err != nil {
			logx.WithContext(r.Context()).Errorf("添加file表失败：%+v", err)
		}
		httpx.OkJson(w, &result.ResponseSuccessBean{
			Code: xerr.OK,
			Msg:  "文件已合并",
			Data: &metadata,
		})
		return
	}
	uploadDir := fc.getSaveFilepath(cMetadata.Fid)
	files, err := ioutil.ReadDir(uploadDir)
	if err != nil {
		result.HttpResult(r, w, errors.Wrapf(xerr.NewErrCode(xerr.USER_OPERATION_ERR), "获取文件列表出错err:%+v", err), nil)
		return
	}
	if len(files) != metadata.SliceNum {
		result.HttpResult(r, w, xerr.NewErrMsg("还有未上传的分片"), nil)
		return
	}

	IsReBuildFid := false
	// 计算md5
	hash := md5.New()
	for i := 0; i < metadata.SliceNum; i++ {
		sliceFilePath := path.Join(uploadDir, strconv.Itoa(i))
		//对分片可能出现错误进行处理
		if metadata.FileDataType == "dicom" {
			IsReBuildFid = true
			sliceFilePath = sliceFilePath + ".dcm"
			if ok := IsFile(sliceFilePath); !ok {
				fmt.Printf("该分片id不符合规则:%s, err:%+v\n", sliceFilePath, err)
				continue
			}
		} else if ok := IsFile(sliceFilePath); !ok {
			fmt.Printf("文件片有错:%v, name:%s\n", err, sliceFilePath)
			continue
		}

		sliceFile, err := os.Open(sliceFilePath)
		if err != nil {
			result.HttpResult(r, w, errors.Wrapf(xerr.NewErrMsg("读取文件失败"), "读取文件%s失败, err: %s\n", sliceFilePath, err), nil)
			return
		}
		defer sliceFile.Close()
		_, err = io.Copy(hash, sliceFile)
		if err != nil {
			result.HttpResult(r, w, errors.Wrap(xerr.NewErrMsg("计算文件md5值失败"), err.Error()), nil)
		}
	}

	hashBt := hash.Sum(nil)
	md5Sum := hex.EncodeToString(hashBt)
	//fmt.Printf("md5校验成功，原始md5：%s, 计算的md5：%s\n", cMetadata.Md5sum, md5Sum)
	// 更新元数据信息
	metadata.Filesize = int64(hash.Size())
	metadata.Md5sum = md5Sum
	metadata.State = "active"
	metadata.Info = map[string]string{}
	metadata.UserId = cMetadata.UserId
	if IsReBuildFid {
		metadata.Fid = GetFid(fc.objGetNodeByFid(cMetadata.Fid).Id, md5Sum)
		metadata.Info[globalkey.DescribeJsonFid] = metadata.Fid + ".json"
		//通过计算md5值生成新的fid，来进行去重
		if IsDir(fc.getSaveFilepath(metadata.Fid)) {
			if err = os.RemoveAll(uploadDir); err != nil {
				fmt.Printf("删除uploadDir:%+v", err)
			}
			if err = os.RemoveAll(metadataPath); err != nil {
				fmt.Printf("删除metadataPath:%+v", err)
			}

			//不同用户上传md5相同的文件的处理
			if cMetadata.UserId != metadata.UserId {
				_, err = fc.FileModel.Insert(nil, &model.File{
					Id:           metadata.Fid,
					UserId:       cMetadata.UserId,
					FileType:     globalkey.ObjSaveType,
					FileDataType: metadata.FileDataType,
					Name:         metadata.Filename,
					Size:         strconv.FormatInt(metadata.Filesize, 10),
					SliceNum:     strconv.Itoa(metadata.SliceNum),
					Md5Sum:       metadata.Md5sum,
				})
				if err != nil {
					result.HttpResult(r, w, xerr.NewErrMsg("添加file表失败, 请重试"), nil)
					return
				}
			}

			result.HttpResult(r, w, nil, &metadata)
			return
		}

		newUploadDir := fc.getSaveFilepath(metadata.Fid)
		if newMetadataPath := fc.getMetadataFilepath(metadata.Fid); metadataPath != newMetadataPath {
			if err = os.Rename(metadataPath, newMetadataPath); err != nil {
				fmt.Printf("dicom元数据文件按md5规则重命名失败, err：%s\n", err)
				httpx.OkJson(w, &result.ResponseSuccessBean{
					xerr.OK,
					"文件上传成功，但是dicom元数据文件按md5规则重命名失败，请手动重命名",
					&metadata})
				return
			}
		}
		if uploadDir != newUploadDir {
			if err = os.Rename(uploadDir, newUploadDir); err != nil {
				fmt.Printf("dicom文件按md5规则重命名失败, err：%s\n", err)
				httpx.OkJson(w, &result.ResponseSuccessBean{
					xerr.OK,
					"文件上传成功，但是dicom文件按md5规则重命名失败，请手动重命名",
					&metadata})
				return
			}
		}
		uploadDir = newUploadDir
	} else if metadata.FileDataType == "dicom-zip" {
		metadata.Info[globalkey.DescribeJsonFid] = metadata.Fid + ".json"
		//合并分片, 解压为想要的格式
		zipP := fc.getNormalFilepath(cMetadata.Fid + ".zip")
		err = mergeSlice(zipP, uploadDir, metadata.SliceNum)
		if err != nil {
			os.Remove(zipP)
			result.HttpResult(r, w, errors.Wrap(xerr.NewErrMsg("写入文件失败"), err.Error()), nil)
			return
		}

		file, err := os.Open(zipP)
		if err != nil {
			result.HttpResult(r, w, errors.Wrap(xerr.NewErrMsg("合并文件失败"), err.Error()), nil)
			return
		}
		zipMd5 := md5.New()
		_, err = io.Copy(zipMd5, file)
		if err != nil {
			result.HttpResult(r, w, errors.Wrap(xerr.NewErrMsg("计算文件md5值失败"), err.Error()), nil)
			return
		}
		metadata.Md5sum = hex.EncodeToString(zipMd5.Sum(nil))

		err = os.Rename(uploadDir, uploadDir+"_tmp")
		if err != nil {
			result.HttpResult(r, w, errors.Wrap(xerr.NewErrMsg("重命名失败"), err.Error()), nil)
			return
		}
		err = compression.DeCompression(zipP, uploadDir, false, true)
		if err != nil {
			result.HttpResult(r, w, errors.Wrap(xerr.NewErrMsg("解压失败"), err.Error()), nil)
			return
		}
		upFiles, err := os.ReadDir(uploadDir)
		if err != nil {
			result.HttpResult(r, w, errors.Wrap(xerr.NewErrMsg("读取文件数量失败"), err.Error()), nil)
			return
		}
		metadata.SliceNum = len(upFiles)
		//uploadDir = uploadDir+"_tmp"
		err = os.RemoveAll(uploadDir + "_tmp")
		if err != nil {
			result.HttpResult(r, w, errors.Wrap(xerr.NewErrMsg("删除原文件失败"), err.Error()), nil)
			return
		}
	} else {
		if md5Sum != cMetadata.Md5sum {
			// 删除保存文件夹
			if err = os.RemoveAll(uploadDir); err != nil {
				fmt.Printf("删除uploadDir:%+v", err)
			}
			if err = os.RemoveAll(metadataPath); err != nil {
				fmt.Printf("删除metadataPath:%+v", err)
			}
			result.HttpResult(r, w, xerr.NewErrMsgf("md5校验失败，请重新上传，原始md5：%s, 计算的md5：%s", cMetadata.Md5sum, md5Sum), nil)
			return
		}
	}
	err = StoreMetadata(metadataPath, metadata)
	if err != nil {
		result.HttpResult(r, w, errors.Wrapf(xerr.NewErrCode(xerr.USER_OPERATION_ERR), "更新元数据文件失败: %+v", err), nil)
		return
	}

	//执行脚本
	if jsonPath := fc.getNormalFilepath(metadata.Info[globalkey.DescribeJsonFid]); !IsFile(jsonPath) {
		args := map[string]string{}
		args["jsonPath"] = jsonPath
		args["fid"] = metadata.Fid
		err = fc.execScript("dicomToDescribeJson", uploadDir, args)
		if err != nil {
			fmt.Printf("执行失败: %+v", err)
			httpx.OkJson(w, &result.ResponseSuccessBean{
				Code: xerr.OK,
				Msg:  "文件上传成功，但是dicom构建json描述文件失败，请手动构建",
				Data: &metadata,
			})
			return
		}
	}

	logx.WithContext(r.Context()).Errorf("metadata:%+v", metadata)
	_, err = fc.FileModel.Insert(nil, &model.File{
		Id:           metadata.Fid,
		UserId:       cMetadata.UserId,
		FileType:     globalkey.ObjSaveType,
		FileDataType: metadata.FileDataType,
		Name:         metadata.Filename,
		Size:         strconv.FormatInt(metadata.Filesize, 10),
		SliceNum:     strconv.Itoa(metadata.SliceNum),
		Md5Sum:       metadata.Md5sum,
	})
	if err != nil {
		result.HttpResult(r, w, xerr.NewErrMsg("添加file表失败, 请重试"), nil)
		return
	}
	result.HttpResult(r, w, nil, &metadata)
}

//作为logic
func (fc *FileCenter) MergeSliceFiles_(cMetadata *ClientFileMetadata) (*ServerFileMetadata, error) {
	//加载元数据
	metadataPath := fc.getMetadataFilepath(cMetadata.Fid)
	metadata, err := LoadMetadata(metadataPath)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.USER_OPERATION_ERR), "加载元数据失败: %+v", err)
	}
	//去重
	if metadata.State == "active" {
		return metadata, nil
		//return metadata, xerr.NewErrCodeMsg(xerr.OK, "文件已合并")
	}
	uploadDir := fc.getSaveFilepath(cMetadata.Fid)
	files, err := ioutil.ReadDir(uploadDir)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.USER_OPERATION_ERR), "获取文件列表出错err:%+v", err)
	}
	IsReBuildFid := false
	// 计算md5
	hash := md5.New()
	for _, file := range files {
		//对分片可能出现错误进行处理, todo 可以把出问题的分片传回前端，进行重新上传
		if metadata.FileDataType == "dicom" && strings.HasSuffix(file.Name(), ".dcm") {
			IsReBuildFid = true
			seq := strings.Split(file.Name(), ".dcm")[0]
			if _, err := strconv.Atoi(seq); err != nil {
				fmt.Printf("该分片id不符合规则:%s, err:%+v\n", seq, err)
				continue
			}
		} else if _, err = strconv.Atoi(file.Name()); err != nil {
			fmt.Printf("文件片有错:%v, name:%s\n", err, file.Name())
			continue
		}
		sliceFilePath := path.Join(uploadDir, file.Name())
		sliceFile, err := os.Open(sliceFilePath)
		if err != nil {
			return nil, errors.Wrapf(xerr.NewErrMsg("读取文件失败"), "读取文件%s失败, err: %s\n", sliceFilePath, err)
		}
		defer sliceFile.Close()
		_, err = io.Copy(hash, sliceFile)
		if err != nil {
			return nil, errors.Wrap(xerr.NewErrMsg("计算文件md5值失败"), err.Error())
		}
	}
	hashBt := hash.Sum(nil)
	md5Sum := hex.EncodeToString(hashBt)
	//fmt.Printf("md5校验成功，原始md5：%s, 计算的md5：%s\n", cMetadata.Md5sum, md5Sum)
	// 更新元数据信息
	metadata.Filesize = int64(hash.Size())
	metadata.Md5sum = md5Sum
	metadata.State = "active"
	metadata.Info = map[string]string{}
	if IsReBuildFid {
		metadata.Fid = GetFid(fc.objGetNodeByFid(cMetadata.Fid).Id, md5Sum)
		metadata.Info[globalkey.DescribeJsonFid] = metadata.Fid + ".json"
		if IsDir(fc.getSaveFilepath(metadata.Fid)) {
			if err = os.RemoveAll(uploadDir); err != nil {
				fmt.Printf("删除uploadDir:%+v", err)
			}
			if err = os.RemoveAll(metadataPath); err != nil {
				fmt.Printf("删除metadataPath:%+v", err)
			}
			return metadata, nil
		}
	} else {
		if md5Sum != cMetadata.Md5sum {
			// 删除保存文件夹
			if err = os.RemoveAll(uploadDir); err != nil {
				fmt.Printf("删除uploadDir:%+v", err)
			}
			if err = os.RemoveAll(metadataPath); err != nil {
				fmt.Printf("删除metadataPath:%+v", err)
			}
			return nil, xerr.NewErrMsgf("md5校验失败，请重新上传，原始md5：%s, 计算的md5：%s", cMetadata.Md5sum, md5Sum)
		}
	}
	err = StoreMetadata(metadataPath, metadata)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.USER_OPERATION_ERR), "更新元数据文件失败: %+v", err)
	}

	if metadata.FileDataType == "dicom-zip" {
		//合并分片, 解压为想要的格式
		zipP := fc.getNormalFilepath(cMetadata.Fid + ".zip")
		err = ioutil.WriteFile(zipP, hashBt, 0666)
		if err != nil {
			return nil, errors.Wrap(xerr.NewErrMsg("写入文件失败"), err.Error())
		}
		err := compression.DeCompression(zipP, uploadDir+"_tmp", false, true)
		if err != nil {
			return nil, errors.Wrap(xerr.NewErrMsg("解压失败"), err.Error())
		}
		err = os.RemoveAll(uploadDir)
		if err != nil {
			return nil, errors.Wrap(xerr.NewErrMsg("删除原文件失败"), err.Error())
		}
		err = os.Rename(uploadDir+"_tmp", uploadDir)
		if err != nil {
			return nil, errors.Wrap(xerr.NewErrMsg("重命名失败"), err.Error())
		}
	} else if IsReBuildFid { // dicom文件夹
		newUploadDir := fc.getSaveFilepath(metadata.Fid)
		if newMetadataPath := fc.getMetadataFilepath(metadata.Fid); metadataPath != newMetadataPath {
			if err = os.Rename(metadataPath, newMetadataPath); err != nil {
				return metadata, xerr.NewErrCodeMsg(xerr.OK, "文件上传成功，但是dicom元数据文件按md5规则重命名失败，请手动重命名")
			}
		}
		if newUploadDir != uploadDir {
			if err = os.Rename(uploadDir, newUploadDir); err != nil {
				return metadata, xerr.NewErrCodeMsg(xerr.OK, "文件上传成功，但是dicom文件按md5规则重命名失败，请手动重命名")
			}
			uploadDir = newUploadDir
		}
	}
	//执行脚本
	if jsonPath := fc.getNormalFilepath(metadata.Info[globalkey.DescribeJsonFid]); !IsFile(jsonPath) {
		args := map[string]string{}
		args["jsonPath"] = jsonPath
		args["fid"] = metadata.Fid
		err = fc.execScript("dicomToDescribeJson", uploadDir, args)
		if err != nil {
			return metadata, errors.Wrapf(xerr.NewErrCodeMsg(xerr.OK, "文件上传成功，但是dicom构建json描述文件失败，请手动构建"), "执行失败: %+v", err)
		}
	}

	return metadata, nil
}

/*列出文件信息
@method：get
@req：nodeId
@resp：{
	Files    []{
		Filename    string  // 文件名
		Filesize    int64   // 文件大小
		Filetype    string  // 文件类型（目前有普通文件和切片文件两种）
	}
}
*/
func (fc *FileCenter) ListFiles(w http.ResponseWriter, r *http.Request) {
	nodeId := r.FormValue("nodeId")
	if nodeId == "" {
		result.ParamErrorResult(r, w, xerr.NewErrMsg("nodeId不能为空"))
		return
	}
	fileType := r.FormValue("fileType")
	if fileType == "" {
		fileType = globalkey.ObjSaveType
	}
	var node cfgByWatchFile.Node
	var ok bool
	if fileType == globalkey.ObjSaveType {
		node, ok = fc.NodeCfg.OssNodes[nodeId]
	} else if fileType == globalkey.BlockSaveType {
		node, ok = fc.NodeCfg.BlockNodes[nodeId]
	}
	if !ok {
		result.ParamErrorResult(r, w, xerr.NewErrMsg("nodeId错误"))
		return
	}
	route := r.FormValue("route")
	//todo 可能要处理一下危禁的参数
	files, err := ioutil.ReadDir(path.Join(node.DiskMountedOn, route))
	if err != nil {
		result.HttpResult(r, w, xerr.NewErrMsgf("获取%s文件列表失败", node.DiskMountedOn), nil)
		return
	}

	fileinfos := ListFileInfos{
		Files: []FileInfo{},
	}
	for _, file := range files {
		tmpFile := path.Join(node.DiskMountedOn, file.Name())
		finfo := FileInfo{
			Filename: file.Name(),
			Filesize: file.Size(),
			Filetype: "normal",
		}
		if file.IsDir() {
			finfo.Filetype = "dir"
		}
		if strings.HasSuffix(file.Name(), ".slice") {
			// 切片文件
			metadata, err := LoadMetadata(tmpFile)
			if err != nil || metadata.State != "active" {
				continue
			}
			finfo.Filename = metadata.Filename
			finfo.Filesize = metadata.Filesize
			finfo.Filetype = "slice"
		}
		fileinfos.Files = append(fileinfos.Files, finfo)
	}
	result.HttpResult(r, w, nil, &fileinfos)
}

/*下载单个文件（小于等于1m的小文件）
@method：get
@req：fid
@resp：文件流
*/
func (fc *FileCenter) Download(w http.ResponseWriter, r *http.Request) {
	//文件名
	fid := r.FormValue("fid")
	if fid == "" {
		result.ParamErrorResult(r, w, errors.New("fid为空"))
		return
	}
	//打开文件
	filePath := fc.getNormalFilepath(fid)
	file, err := os.Open(filePath)
	if err != nil {
		result.HttpResult(r, w, xerr.NewErrMsgf("读文件%s失败:%s", filePath, err), nil)
		return
	}
	defer file.Close() //结束后关闭文件

	//设置响应的header头
	w.Header().Add("Content-type", "application/octet-stream")
	w.Header().Add("content-disposition", "attachment; filename=\""+fid+"\"")
	//将文件写至responseBody
	_, err = io.Copy(w, file)
	if err != nil {
		http.Error(w, "文件下载失败", http.StatusInternalServerError)
		result.HttpResult(r, w, xerr.NewErrMsg("文件下载失败"), nil)
		return
	}
}

/* ---------------------------分片下载步骤---------------------------*/
/*1、获取文件信息
@method：get
@req：filename
@resp：{
	Filename string
	Filesize int64
	Filetype string
}
*/
func (fc *FileCenter) GetFileInfo(w http.ResponseWriter, r *http.Request) {
	fid := r.FormValue("fid")
	mPath := fc.getMetadataFilepath(fid)
	finfo := FileInfo{
		Filetype: "normal",
	}
	if !IsFile(mPath) {
		normalFilepath := fc.getNormalFilepath(fid)
		fstate, err := os.Stat(normalFilepath)
		if err != nil {
			fmt.Println("读取文件失败", normalFilepath)
			http.Error(w, "读取文件失败", http.StatusBadRequest)
			return
		}
		finfo.Filename = fstate.Name()
		finfo.Filesize = fstate.Size()
	} else {
		// 切片文件
		metadata, err := LoadMetadata(mPath)
		if err != nil {
			http.Error(w, "获取文件元数据信息失败", http.StatusInternalServerError)
			return
		}
		finfo.Filename = metadata.Filename
		finfo.Filesize = metadata.Filesize
		finfo.Filetype = "slice"
	}
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(finfo)
	if err != nil {
		fmt.Println("编码文件基本信息失败")
		http.Error(w, "服务异常", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

/*2、获取文件元数据信息
@method：get
@req：filename
@resp：{
	Fid        string
	Filesize   int64
	Filename   string
	SliceNum   int
	Md5sum     string
	ModifyTime time.Time
}
*/
func (fc *FileCenter) GetFileMetainfo(w http.ResponseWriter, r *http.Request) {
	fid := r.FormValue("fid")
	metaPath := fc.getMetadataFilepath(fid)
	if !IsFile(metaPath) {
		http.Error(w, "file not exist", http.StatusBadRequest)
		return
	}
	metadata, err := LoadMetadata(metaPath)
	if err != nil {
		http.Error(w, "文件损坏", http.StatusBadRequest)
		return
	}
	cMetadata := metadata.ClientFileMetadata
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(cMetadata)
	if err != nil {
		fmt.Println("编码文件基本信息失败")
		http.Error(w, "服务异常", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

/*3、通过分片下载文件
@method：get
@req：filename、sliceIndex
@resp：文件流
*/
func (fc *FileCenter) DownloadBySlice(w http.ResponseWriter, r *http.Request) {
	fid := r.FormValue("fid")
	if fid == "" {
		result.ParamErrorResult(r, w, errors.New("fid为空"))
		return
	}
	sliceIndex := r.FormValue("sliceIndex")
	if sliceIndex == "" {
		result.ParamErrorResult(r, w, errors.New("sliceIndex为空"))
		return
	}
	var filepath string
	fileType := r.FormValue("fileType")
	if fileType == "" {
		//todo 需要实现一个fid检查功能，需要实现一个对一般、分片文件的区分功能
		nodeId := strings.SplitN(fid, "-", 2)[0]
		if _, ok := fc.NodeCfg.OssNodes[nodeId]; ok {
			fileType = globalkey.ObjSaveType
		} else {
			fileType = globalkey.BlockSaveType
		}
	}
	if fileType == globalkey.ObjSaveType {
		filepath = fc.GetSaveFilepath(fid)
	} else if fileType == globalkey.BlockSaveType {
		filepath = fc.getBlockFilepath(fid)
	}
	sliceFilePath := path.Join(filepath, sliceIndex)
	if !IsFile(sliceFilePath) {
		result.HttpResult(r, w, xerr.NewErrMsg("文件切片不存在"), nil)
		return
	}

	file, err := os.Open(sliceFilePath)
	if err != nil {
		result.HttpResult(r, w, xerr.NewErrMsg("读取文件切片错误"), nil)
		return
	}
	//结束后关闭文件
	defer file.Close()

	//设置响应的header头
	w.Header().Add("Content-type", "application/octet-stream")
	w.Header().Add("content-disposition", fmt.Sprintf("attachment; filename=%s-%s", fid, sliceIndex))
	_, err = io.Copy(w, file)
	if err != nil {
		result.HttpResult(r, w, xerr.NewErrMsgf("下载文件分片%s失败, err:%+v", sliceFilePath, err), nil)
		return
	}
}

// 删除（可批量）

// 复制

// 移动（可批量），数据库逻辑实现

// 重命名（要检查是否被使用）， 更改元数据里面的文件名

// 压缩
func (fc *FileCenter) Compression(w http.ResponseWriter, r *http.Request) {
	var in CompressionReq
	if err := httpx.Parse(r, &in); err != nil {
		result.ParamErrorResult(r, w, err)
		return
	}
	fid := fmt.Sprintf("%s.zip", in.Fid)
	dstPath := fc.getNormalFilepath(fid)
	if IsFile(dstPath) {
		result.HttpResult(r, w, xerr.NewErrCodeMsg(xerr.OK, "已存在或正在压缩"), map[string]string{"fid": fid})
		return
	}
	var fPath string
	if in.FileType == globalkey.ObjSaveType {
		fPath = fc.getSaveFilepath(in.Fid)
		if fPath == "" {
			result.HttpResult(r, w, errors.Wrapf(xerr.NewErrMsg("fid不正确"), "fid:%s和fileType:%s不符合", in.Fid, in.FileType), nil)
			return
		}
	} else if in.FileType == globalkey.BlockSaveType {
		fPath = fc.getBlockFilepath(in.Fid)
	}
	err := compression.Compression(fPath, dstPath)
	if err != nil {
		result.HttpResult(r, w, errors.Wrap(xerr.NewErrCode(xerr.USER_OPERATION_ERR), err.Error()), nil)
		return
	}
	result.HttpResult(r, w, nil, map[string]string{"fid": fid})
}

// 解压

// 执行脚本（可对文件/文件夹）
func (fc *FileCenter) ExecScript(w http.ResponseWriter, r *http.Request) {
	var execScriptReq ExecScriptReq
	if err := httpx.Parse(r, &execScriptReq); err != nil {
		result.ParamErrorResult(r, w, err)
		return
	}
	err := fc.execScript_(&execScriptReq)
	if err != nil {
		result.HttpResult(r, w, xerr.NewErrMsgf("执行失败: %+v", err), nil)
		return
	}
	result.HttpResult(r, w, nil, &map[string]string{
		"fid":      execScriptReq.Fid,
		"jsonPath": execScriptReq.Args["jsonPath"],
	})
}

//var filePath string
//if execScriptReq.FileType == "block"{
//	//块存储还是有点问题
//	filePath = fc.getBlockFilepath(execScriptReq.Fid)
//}else if execScriptReq.FileType == "obj"{
//	filePath = fc.getSaveFilepath(execScriptReq.Fid)
//}else if execScriptReq.FileType == "update"{
//	fmt.Print("修改json的情况")
//}
//execScriptReq.Args["fid"] = execScriptReq.Fid
//if _,ok := execScriptReq.Args["jsonPath"]; !ok{
//	jsonFid := GetFid(fc.objGetStorageNode().Id, path.Base(filePath) )
//	execScriptReq.Args["jsonPath"] = fc.getNormalFilepath(jsonFid+".json")
//}
//err := fc.execScript(execScriptReq.Script, filePath, execScriptReq.Args)
//if err != nil {
//	result.HttpResult(r,w,xerr.NewErrMsgf("执行失败: %+v", err), nil)
//	return
//}
//result.HttpResult(r,w,nil,&map[string]string{
//	"fid":execScriptReq.Fid,
//	"jsonPath":execScriptReq.Args["jsonPath"],
//})
