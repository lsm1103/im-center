package cfgByWatchFile

import (
	"encoding/json"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"golang.org/x/xerrors"
	"io/ioutil"
	"log"
	"os"
)

type (
	Node struct {
		Id            string `json:"id"`
		DiskSize      string `json:"diskSize"`
		DiskUsed      string `json:"diskUsed"`
		DiskAvail     string `json:"diskAvail"`
		DiskUse       string `json:"diskUse%"`
		DiskMountedOn string `json:"diskMountedOn"`
	}
	Config struct {
		OssNodes map[string]Node `json:"ossNodes"`
		BlockNodes map[string]Node `json:"blockNodes"`
	}
)

func ReadDynamicCfg(path string) (*Config, error) {
	f, err := os.Open(path)
	defer f.Close()
	if err != nil { return nil, xerrors.New("read file fail") }
	fd, err := ioutil.ReadAll(f)
	if err != nil { return nil, xerrors.New("read file fail") }
	m := Config{}
	err = json.Unmarshal(fd, &m)
	if err != nil { return nil, xerrors.New("jsonUnmarshal file fail") }
	return &m,nil
}

func Server(watchFilePath string, fun func(data *Config) ) {
	fmt.Printf("开启监听文件：%s", watchFilePath)
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	//done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				log.Println("event:", event)
				if event.Op&fsnotify.Write == fsnotify.Write {

					dcfg, err := ReadDynamicCfg(watchFilePath)
					if err != nil {
						log.Println("err:", err)
					}
					log.Printf("modified file:%s, dcfg:%+v", event.Name, dcfg)
					//ctx.FileCenter.NodeCfg = dcfg
					fun(dcfg)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add(watchFilePath)
	//err = watcher.Add("/Users/xm/Desktop/work_project/im-center/etc/fileCenter.json")
	if err != nil {
		log.Fatal(err)
	}
	//<-done
}
