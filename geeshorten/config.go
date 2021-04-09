package main
import (
	"encoding/json"
	"fmt"
	"os"
	"geeshorten/db"
)

var MyConfig Config // 配置

type Config struct {
	Db            db.Db  //mysql配置
	ServerAddress string // 服务器监听地址
	BaseUrl       string // 短链接网址
}

func LoadConfig(path string) {
	fmt.Println("加载配置文件...")
	file, err := os.Open(path)
	if err != nil {
		fmt.Errorf("打开配置文件错误", err)
	}

	decoder := json.NewDecoder(file)

	err = decoder.Decode(&MyConfig)
	if err != nil {
		fmt.Errorf("配置文件json解码错误", err)
	}
}
