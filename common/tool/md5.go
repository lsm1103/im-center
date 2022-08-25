package tool

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"os"
	"path"
	"strconv"
	"strings"
)

func md5Sum(uploadDir string, isDir bool) (string, error) {
	if !isDir{
		// 计算md5
		file, err := os.Open(uploadDir)
		if err != nil {
			fmt.Println("os Open error")
			return "", err
		}
		md5 := md5.New()
		_, err = io.Copy(md5, file)
		if err != nil {
			fmt.Println("io copy error")
			return "", err
		}
		return hex.EncodeToString(md5.Sum(nil)), nil
	} else {
		hash := md5.New()
		files, err := ioutil.ReadDir(uploadDir)
		for _, file := range files {
			fmt.Printf("file.Name:%s\n", file.Name())
			//对分片可能出现错误进行处理
			if strings.HasSuffix(file.Name(), ".dcm") {
				seq := strings.Split(file.Name(), ".dcm")[0]
				if _, err := strconv.Atoi(seq); err != nil {
					fmt.Printf("该分片id不符合规则:%s, err:%+v\n", seq, err)
					continue
				}
			} else if _, err = strconv.Atoi(file.Name()); err != nil {
				fmt.Printf("文件片有错:%v, name:%s\n", err, file.Name())
				continue
			}
			sliceFilePath := path.Join(uploadDir, file.Name() )
			sliceFile, err := os.Open(sliceFilePath)
			if err != nil {
				fmt.Printf("读取文件%s失败, err: %s\n", sliceFilePath, err)
				return "", err
			}
			defer sliceFile.Close()
			_, err = io.Copy(hash, sliceFile)
			if err != nil {
				fmt.Println("计算文件md5值失败")
			}
		}
		return hex.EncodeToString(hash.Sum(nil)), nil
	}
}

const SliceBytes = 1024*1024*1   // 分片大小

func md5SumSlice(uploadPath string)(string, error){
	fileStat, err := os.Stat(uploadPath)
	if err != nil {
		fmt.Printf("读取文件%s失败, err: %s\n", uploadPath, err)
		return "",nil
	}
	filesize := fileStat.Size()
	if filesize <= 0 {
		fmt.Printf("%s文件是空文件，不能上传\n", uploadPath)
		return "",nil
	}
	// 计算文件切片数量
	sliceNum := int(math.Ceil(float64(filesize) / float64(SliceBytes)))

	f, err := os.Open(uploadPath)
	if err != nil {
		return "", err
	}
	hash := md5.New()
	for i := 0; i < sliceNum; i++ {
		tmpData := make([]byte, SliceBytes)
		nr, err := f.Read(tmpData[:])
		if err != nil {
			fmt.Printf("read file error\n")
			return "",err
		}
		hash.Write(tmpData[:nr])
		//tmpData = tmpData[:nr]
	}

	// 计算文件md5
	return hex.EncodeToString(hash.Sum(nil)), nil
}

func mergeSlice(uploadPath string)  {
	dir, err := os.ReadDir(uploadPath)
	if err != nil {
		fmt.Printf("err:%s\n", err)
	}
	zf, err := os.OpenFile(fmt.Sprintf("/Users/xm/Desktop/%s.zip", path.Base(uploadPath)), os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Printf("err:%s\n", err)
	}
	defer zf.Close()
	for i := 0; i < len(dir); i++ {
		fmt.Printf("%d\n", i)
		sliceFile, err := os.Open(path.Join(uploadPath, strconv.Itoa(i)))
		if err != nil {
			fmt.Printf("err:%s\n", err)
		}
		defer sliceFile.Close()
		// 偏移量需要重新进行调整
		sliceFile.Seek(0, 0)
		_, err = io.Copy(zf, sliceFile)
		if err != nil {
			fmt.Printf("err:%s\n", err)
		}
	}
}