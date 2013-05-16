//积分操作方法
package integral

import (
	"errors"
	"fmt"
	"github.com/iyf/godiscuz/discuz"
	"strconv"
)

//积分兑换请求
//uid:用户ID
//from:原积分
//to:目标积分
//toappid:目标应用ID
//amount:积分数额
func ExchangeRequest(uid, from, to, toappid, amount int, args ...string) (err error) {
	input := fmt.Sprintf("uid=%v&from=%v&to=%v&toappid=%v&amount=%v", uid, from, to, toappid, amount)
	data := discuz.ApiPost("credit", "request", input, args...)
	i, err := strconv.Atoi(data)
	if err == nil {
		if i == 0 {
			err = errors.New("请求失败")
		}
	} else {
		err = fmt.Errorf("请求错误：", err)
	}
	return
}
