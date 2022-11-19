package common

import (
	"net/http"
	"wc_robot/common/utils"
)

type HitokotoSentence struct {
	Hitokoto string `json:"hitokoto"`
	FromWho  string `json:"from_who"`
	From     string `json:"from"`
}

// GetHitokotoResponse 查询一言
func GetHitokotoResponse() (*HitokotoSentence, error) {
	rsp, err := http.Get("https://v1.hitokoto.cn/")
	if err != nil {
		return nil, err
	}
	var cs HitokotoSentence
	if err := utils.ScanJson(rsp, &cs); err != nil {
		return nil, err
	}
	return &cs, nil
}
