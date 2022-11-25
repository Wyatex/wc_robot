package common

import (
	"math/rand"
	"time"
)

var pai = []string{
	"拍坏了，赔钱！",
	"再拍叫群主打屎你，哼(￢︿̫̿￢)！",
	"(。´・ω・)ん?",
	"正在定位您的真实地址...\n定位成功，轰炸机已起飞",
	"喂(#`O′)，拍人家干嘛！",
	"要被拍坏啦！！！",
	"我爪巴，球球您别拍了",
	"你再拍！",
	"烦诶！你拍什么呢",
}

func GetPaiYiPai() string {
	rand.Seed(time.Now().Unix())
	return pai[rand.Intn(len(pai)-1)]
}
