package compression

import (
	"archive/zip"
	"bytes"
	"fmt"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

func DeCompression(srcPath string, dstPath string, isOriginalName bool, isInRoot bool) error {
	//isOriginalName: 是否保持原文件的文件名，在dstPath文件夹里面
	//打开并读取压缩文件中的内容
	fr, err := zip.OpenReader(srcPath)
	if err != nil {
		return err
	}
	defer fr.Close()
	srcName := path.Base(srcPath)
	dirName := strings.SplitN(srcName, ".", 2 )[0]
	if dstPath == ""{
		dstPath = path.Join(path.Dir(srcPath), dirName)
	}
	err = os.MkdirAll(dstPath, 0766)
	if err != nil {
		return err
	}
	//fmt.Printf("srcName:%s, dstPath:%s, path.Base(dstPath):%s\n", srcName, dstPath, path.Base(dstPath))
	//r.reader.file 是一个集合，里面包括了压缩包里面的所有文件
	var decodeName string
	count := 0
	for _, file := range fr.Reader.File {
		//处理中文编码
		if file.Flags == 0{
			//如果标致位是0  则是默认的本地编码   默认为gbk
			i:= bytes.NewReader([]byte(file.Name))
			decoder := transform.NewReader(i, simplifiedchinese.GB18030.NewDecoder())
			content,_:= ioutil.ReadAll(decoder)
			decodeName = string(content)
		}else{
			//如果标志为是 1 << 11也就是 2048  则是utf-8编码
			decodeName = file.Name
		}

		fPath := path.Join(dstPath, decodeName)
		if !isOriginalName{
			fPath = path.Join(dstPath, strings.Replace(decodeName, dirName,"", 1))
		}
		//fmt.Println("unzip: ", fPath)
		//判断文件该目录文件是否为文件夹
		if file.FileInfo().IsDir() {
			err := os.MkdirAll(fPath, 0766)
			if err != nil {
				fmt.Println("MkdirAll:",fPath, err)
			}
			continue
		}
		//把所有的文件都解压到根目录
		if isInRoot{
			fPath = path.Join(dstPath, dirName, path.Base(decodeName) )
			if !isOriginalName{
				fPath = path.Join(dstPath, path.Base(decodeName))
			}
		}
		//为文件时，打开文件
		r, err := file.Open()
		if err != nil {
			fmt.Println("Open:",fPath, err)
			continue
		}
		if path.Ext(fPath) == ".dcm" || path.Ext(fPath) == ".Dcm"{
			fPath = fmt.Sprintf("%s/%d.dcm", path.Dir(fPath), count)
			count++
		}
		//在对应的目录中创建相同的文件
		NewFile, err := os.Create(fPath)
		if err != nil {
			fmt.Println("Create:", fPath, err)
			continue
		}
		//将内容复制
		io.Copy(NewFile, r)
		//关闭文件
		NewFile.Close()
		r.Close()
	}

	fdir, err := os.ReadDir(dstPath)
	if err != nil {
		return err
	}
	for _,f := range fdir{
		if f.IsDir(){
			os.RemoveAll(path.Join(dstPath, f.Name() ) )
		}
	}
	return nil
}