package tool

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)

/*
1、判断文件是否存在
if _, err := os.Stat(path); os.IsNotExist(err) {
	fmt.Println("read to fd fail", err)
}
2、写入文件
ioutil.WriteFile(path, data, 0666)
3、func OpenFile(name string, flag int, perm FileMode) (file *File, err error)
下面列举了一些常用的 flag 文件处理参数：
os.O_RDONLY int = 0		 			// 只读模式打开文件
os.O_WRONLY int = 1		 			// 只写模式打开文件
os.O_RDWR int = 2 					// 读写模式打开文件
os.O_APPEND int = 8		 			// 写操作时将数据附加到文件尾部
os.O_CREATE int = 512		 		// 如果不存在将创建一个新文件
os.O_EXCL int = 2048 				// 和O_CREATE配合使用，文件必须不存在
os.O_SYNC int = 128 				// 打开文件用于同步I/O
os.O_TRUNC int = 1024		 		// 如果可能，打开时清空文件
*/

// 判断所给路径是否为文件夹
func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

// 判断所给文件是否存在
func IsFile(path string) bool {
	s, err := os.Stat(path)
	if err != nil || s.IsDir() {
		//fmt.Println("文件不存在或为文件夹：", path)
		return false
	}
	return true
}

func Read2Str(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "read file fail", err
	}
	defer f.Close()

	fd, err := ioutil.ReadAll(f)
	if err != nil {
		return "read to fd fail", err
	}
	return string(fd), nil
}

func Read2Byte(path string) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return []byte("read file fail"), err
	}
	defer f.Close()
	fd, err := ioutil.ReadAll(f)
	if err != nil {
		return []byte("read file fail"), err
	}
	return fd, nil
}

func WriteFile(path string, data string, flag int) error {
	/*
		path string:写入地址,
		data string:写入数据,
		flag int:写入模式;
	*/
	err := os.MkdirAll(path[:strings.LastIndex(path, "/")], 0744)
	//fmt.Println("way:", path, "\rpath:", path, "\rdata:", data, "\rjson:", []byte(data))
	if err == nil {
		f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|flag, 0666)
		//及时关闭f句柄
		defer f.Close()
		if err == nil {
			_, err1 := f.WriteString(data)
			err = err1
			// fmt.Println("n", n)
		}
	}
	return err
}

func Read2Json(path string) (map[string]interface{}, error) {
	data, err := Read2Byte(path)
	m := make(map[string]interface{})
	if err == nil {
		err = json.Unmarshal([]byte(data), &m)
	}
	return m, err
}

func Read2StrJson(path string) (map[string]string, error) {
	data, err := Read2Byte(path)
	m := make(map[string]string)
	if err == nil {
		err = json.Unmarshal([]byte(data), &m)
	}
	return m, err
}

func WriteJson(path string, data string) error {
	/*
		path string:写入地址,
		data string:写入数据,
	*/
	e := os.MkdirAll(path[:strings.LastIndex(path, "/")], 0744)
	if e != nil {
		return e
	}
	//fmt.Println("way:", path, "\rdata:", data, "\rjson:", []byte(data))
	return ioutil.WriteFile(path, []byte(data), 0777)
}

func BuildTimePath() string {
	r := time.Now().In(Loc)
	return fmt.Sprintf("%s/%s/%s", strconv.Itoa(r.Year()), strconv.Itoa(int(r.Month())), strconv.Itoa(r.Day()))
}

func MkdirPath(path string) {
	e := os.MkdirAll(path[:strings.LastIndex(path, "/")], 0744)
	if e != nil {
		fmt.Printf("创建文件夹【%s】失败，e:%s", path, e)
	}
}

//func main() {
//	path := "atpFile/fhfghfgh1/2021/10/9/parameter.json"
//	data, _ := Read2Str(path)
//	//m := make(map[string]interface{})
//	//jsonerr := json.Unmarshal([]byte(data), &m)
//	fmt.Println("data:", data, "\n\r")
//	//fmt.Println("jsonerr:", jsonerr, "\n\r")
//	//fmt.Println("m:", m, "\n\r")
//	//e := WriteFile(path, data, 8)
//	//fmt.Println(`m["l1"]`, m["l1"], e, "\n\r")
//}
