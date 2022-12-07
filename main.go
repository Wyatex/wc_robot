package main

import (
	"fmt"
	"log"
	"math/rand"
	"regexp"
	"strings"
	"time"
	"wc_robot/common/nba"

	"wc_robot/common"
	"wc_robot/common/alapi"
	"wc_robot/common/covid"
	"wc_robot/common/weather"
	"wc_robot/robot"
)

// 日志设置初始化
func init() {
	log.SetFlags(log.Llongfile | log.Ldate | log.Ltime)

	// 部署在 linux 上可直接通过 nohup ./wc_robot > robot.log & 运行并打印日志
	// 本机测试运行可取消下方注释，记录 log 便于观察

	// // 打印日志到本地 wc_robot.log
	// outputLogPath := "wc_robot.log"
	// f, err := os.Create(outputLogPath)
	// if err != nil {
	// 	log.Println("[WARN]创建日志文件失败, 日志仅输出在控制台")
	// }
	// w := io.MultiWriter(os.Stdout, f)
	// log.SetOutput(w)
}

var begin = time.Now()

func main() {
	defer func() {
		log.Printf("[INFO]本次机器人运行时间为: %s", time.Since(begin).String())
	}()
	config := common.GetConfig()

	r := robot.NewRobot()
	r.Chain.RegisterHandler("菜单|功能|会什么回复", onMenu)
	r.Chain.RegisterHandler("存活时间回复", onSurvivalTime)
	if config.WeatherMsgHandle.SwitchOn {
		r.Chain.RegisterHandler("天气回复", onWeather)
		r.Chain.RegisterHandler("空气质量回复", onAQI)
	}
	if config.ALAPI.SwitchOn {
		r.Chain.RegisterHandler("名言回复", onMingYan)
		r.Chain.RegisterHandler("情话回复", onQingHua)
		r.Chain.RegisterHandler("鸡汤回复", onSoul)
	}
	if config.CovidMsgHandle.SwitchOn {
		r.Chain.RegisterHandler("疫情回复", onCovid)
	}
	r.Chain.RegisterHandler("一言", onHitokoto)
	r.Chain.RegisterHandler("我老婆是谁", onWife)
	r.Chain.RegisterHandler("今日新闻", onNews)
	r.Chain.RegisterHandler("NBA赛程", onNBA)
	r.Chain.RegisterHandler("疯狂星期四", onKFC)
	r.Chain.RegisterHandler("拍一拍", onPai)

	if err := r.Login(); err != nil {
		log.Println(err)
	}
	robot.InitTasks(config)
	r.Block()
}

// 下面一些匹配：就strings.Contains()和正则匹配二者的性能来说，前者较优

// 判断是否匹配，匹配返回 true, 不匹配返回 false
func checkMatch(msg *robot.Message, keyword string) bool {
	//config := common.GetConfig()
	if msg.IsFromGroup() {
		//if !(strings.Contains(msg.Content, "@"+config.RobotName) && strings.Contains(msg.Content, keyword)) {
		if !strings.Contains(msg.Content, keyword) {
			return false
		}
	}
	if msg.IsFromMember() {
		if !strings.Contains(msg.Content, keyword) {
			return false
		}
	}
	return true
}

// 监听菜单｜功能｜会什么相关的文字进行回复
func onMenu(msg *robot.Message) error {
	config := common.GetConfig()
	if msg.IsFromGroup() {
		if !(strings.Contains(msg.Content, "@"+config.RobotName) &&
			(strings.Contains(msg.Content, "菜单") || strings.Contains(msg.Content, "功能") || strings.Contains(msg.Content, "会什么"))) {
			return nil
		}
	}
	if msg.IsFromMember() {
		if !(strings.Contains(msg.Content, "菜单") || strings.Contains(msg.Content, "功能") || strings.Contains(msg.Content, "会什么")) {
			return nil
		}
	}
	_, err := msg.ReplyText("你好呀👋\n" + `目前只支持"天气"、"空气质量(指标含义)"、"XX(城市、省份、国家)疫情"、"情话"、"鸡汤"、"名言"相关的问题哦`)
	return err
}

func onKFC(msg *robot.Message) error {
	if !checkMatch(msg, "疯狂星期四") {
		return nil
	}
	_, err := msg.ReplyText(common.GetKFC())
	return err
}

func onPai(msg *robot.Message) error {
	if !checkMatch(msg, "拍了拍我") {
		return nil
	}
	_, err := msg.ReplyText(common.GetPaiYiPai())
	return err
}

func onNBA(msg *robot.Message) error {
	if !checkMatch(msg, "NBA赛程") || msg.IsFromGroup() {
		return nil
	}
	str, err := nba.GetNBAResponse()
	if err != nil {
		return err
	}
	_, err = msg.ReplyText(str)
	return err
}

func onWife(msg *robot.Message) error {
	if !checkMatch(msg, "我老婆是谁") {
		return nil
	}
	// 拿出所有成员
	j := 0
	keys := make([]string, len(robot.Storage.MemberMap))
	for k := range robot.Storage.MemberMap {
		keys[j] = k
		j++
	}
	// 抽一个成员
	rand.Seed(time.Now().Unix())
	user := robot.Storage.MemberMap[keys[rand.Intn(len(keys))]]
	fmt.Println(user)
	sb := strings.Builder{}
	sb.WriteString("今天你的群老婆是：")
	sb.WriteString(user.NickName)
	_, err := msg.ReplyText(sb.String())
	return err
}

// 监听天气相关的文字进行回复
func onWeather(msg *robot.Message) error {
	config := common.GetConfig()
	if !checkMatch(msg, "天气") {
		return nil
	}

	w, err := weather.GetWeather(config.WeatherMsgHandle.CityCode)
	if err != nil {
		return err
	}
	_, err = msg.ReplyText(w.GetCurrentWeatherInfo())
	return err
}

// 监听空气质量(指标含义) 的文字进行回复
func onAQI(msg *robot.Message) error {
	config := common.GetConfig()
	if !checkMatch(msg, "空气质量") {
		return nil
	}

	if strings.Contains(msg.Content, "指标含义") {
		msg.ReplyText(weather.AQIIndicesDesc())
		return nil
	}
	w, err := weather.GetWeather(config.WeatherMsgHandle.CityCode)
	if err != nil {
		return err
	}
	_, err = msg.ReplyText(w.GetAQIInfo())
	return err
}

// 监听心灵鸡汤相关的文字进行回复
func onSoul(msg *robot.Message) error {
	if !checkMatch(msg, "鸡汤") {
		return nil
	}

	s, err := alapi.GetSoul()
	if err != nil {
		return err
	}
	_, err = msg.ReplyText(s)
	return err
}

// 监听情话相关的文字进行回复
func onQingHua(msg *robot.Message) error {
	if !checkMatch(msg, "情话") {
		return nil
	}

	content, err := alapi.GetQinghua()
	if err != nil {
		return err
	}
	_, err = msg.ReplyText(content)
	return err
}

// 监听名言相关的文字进行回复
func onMingYan(msg *robot.Message) error {
	if !checkMatch(msg, "名言") {
		return nil
	}

	content, err := alapi.GetMingYan()
	if err != nil {
		return err
	}
	_, err = msg.ReplyText(content)
	return err
}

func onNews(msg *robot.Message) error {
	if !checkMatch(msg, "今天新闻") {
		return nil
	}
	news, err := common.GetNewsResponse()
	if err != nil {
		msg.ReplyText("接口异常：" + err.Error())
		return err
	}
	_, err = msg.ReplyText(news)
	return err
}

func onHitokoto(msg *robot.Message) error {
	if !checkMatch(msg, "一言") {
		return nil
	}
	cr, err := common.GetHitokotoResponse()
	if err != nil {
		msg.ReplyText("接口异常：" + err.Error())

		return err
	}
	if cr.FromWho == "" {
		_, err = msg.ReplyText(cr.Hitokoto + " —— " + "「" + cr.From + "」")
		return err
	} else {
		_, err = msg.ReplyText(cr.Hitokoto + " —— " + cr.FromWho + "「" + cr.From + "」")
		return err
	}
}

var locationRE = regexp.MustCompile("([\u4e00-\u9fa5]{1,6})疫情")

// 监听疫情相关的文字进行回复
func onCovid(msg *robot.Message) error {
	if !checkMatch(msg, "疫情") {
		return nil
	}

	hits := locationRE.FindStringSubmatch(msg.Content)
	if len(hits) != 2 {
		return nil
	}
	cr, err := covid.GetCovidResponse(hits[1])
	if err != nil {
		//msg.ReplyText("非常抱歉，未检索到该地区疫情数据")
		return err
	}
	_, err = msg.ReplyText(covid.PrintCovidSituation(cr))
	return err
}

func onSurvivalTime(msg *robot.Message) error {
	if !checkMatch(msg, "存活时间") {
		return nil
	}

	now := time.Now()
	nowString := now.Format(common.TimeFormat)
	d := now.Sub(begin)
	second := int(d.Seconds()) % 60
	min := int(d.Minutes()) % 60
	hour := int(d.Hours())
	text := fmt.Sprintf("截止至 %s , 机器人已经存活了 %d 小时 %d 分 %d 秒",
		nowString, hour, min, second)
	_, err := msg.ReplyText(text)
	return err
}
