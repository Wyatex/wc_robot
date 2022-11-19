package common

import (
	"errors"
	"net/http"
	"strings"
	"wc_robot/common/utils"
)

type News struct {
	Zt int    `json:"zt"`
	Wb string `json:"wb"`
}

// 查询新闻
func GetNewsResponse() (string, error) {
	rsp, err := http.Get("http://bjb.yunwj.top/php/qq.php")
	if err != nil {
		return "", err
	}
	var news News
	if err := utils.ScanJson(rsp, &news); err != nil {
		return "", err
	}
	if news.Zt != 0 {
		return "", errors.New("接口请求失败")
	}
	return toMessage(news.Wb), nil
}

func toMessage(oldStr string) string {
	sb := strings.Builder{}
	sb.WriteString("吹雪BOT的每日简报，每天60秒带你读懂世界₍ᐢ..ᐢ₎♡")
	newStr := strings.ReplaceAll(oldStr, "【换行】", "\n\n")
	sb.WriteString(newStr)
	return sb.String()
}
