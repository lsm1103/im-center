package compression

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"

	"im-center/common/tool"
)

// 打包成zip文件
func Compression(srcDir string, dstPath string) error {
	if tool.IsFile(dstPath) {
		return errors.New("该文件已经在压缩中")
	}
	// 预防：旧文件无法覆盖
	os.RemoveAll(dstPath)

	// 创建：zip文件
	zipfile, _ := os.Create(dstPath)
	defer zipfile.Close()

	// 打开：zip文件
	archive := zip.NewWriter(zipfile)
	defer archive.Close()

	// 遍历路径信息
	err := filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil{
			fmt.Printf("err:%s", err)
			return err
		}
		// 如果是源路径，提前进行下一个遍历
		if path == srcDir {
			return nil
		}

		// 获取：文件头信息
		header, _ := zip.FileInfoHeader(info)
		header.Name = strings.TrimPrefix(path, srcDir+`/`)

		// 判断：文件是不是文件夹
		if info.IsDir() {
			header.Name += `/`
		} else {
			// 设置：zip的文件压缩算法
			header.Method = zip.Deflate
		}

		// 创建：压缩包头部信息
		writer, err := archive.CreateHeader(header)
		if err != nil{
			fmt.Printf("err:%s", err)
		}
		if !info.IsDir() {
			file, err := os.Open(path)
			if err != nil{
				fmt.Printf("err:%s", err)
			}
			defer file.Close()
			io.Copy(writer, file)
		}
		return nil
	})
	if err != nil {
		os.RemoveAll(dstPath)
	}
	return err
}