package user

import (
	"discuz"
	"log"
	"testing"
)

func init() {
	discuz.Register("2", "你的key", "http://bbs.iyf.cc/uc_server", "uc") //替换为自己的信息
}
func TestRegister(t *testing.T) {
	_, err := Register("sdfsa", "123456", "sdfdsf%40qq.com")
	log.Println(err)
}
