package nba

import (
	"net/http"
	"strings"
	"wc_robot/common/utils"
)

// NBA赛程数据响应结构体
type NBAResponse struct {
	Reason string `json:"reason"`
	Result struct {
		Title    string `json:"title"`
		Duration string `json:"duration"`
		Matchs   []struct {
			Date string `json:"date"`
			Week string `json:"week"`
			List []struct {
				TimeStart  string `json:"time_start"`
				StatusText string `json:"status_text"`
				Team1      string `json:"team1"`
				Team1Score string `json:"team1_score"`
				Team2      string `json:"team2"`
				Team2Score string `json:"team2_score"`
			} `json:"list"`
		} `json:"matchs"`
	} `json:"result"`
}

func GetNBAResponse() (string, error) {
	rsp, err := http.Get("http://apis.juhe.cn/fapig/nba/query?key=2a02a066183289ac6fdb31c783353447")
	if err != nil {
		return "", err
	}
	var cs NBAResponse
	if err := utils.ScanJson(rsp, &cs); err != nil {
		return "", err
	}
	return toString(&cs), nil
}

func toString(res *NBAResponse) string {
	sb := strings.Builder{}
	sb.WriteString(res.Reason)
	sb.WriteString("当前赛事：")
	sb.WriteString(res.Result.Title + "-" + res.Result.Duration + "\n\n")
	for _, match := range res.Result.Matchs {
		sb.WriteString(match.Date + match.Week + ":\n")
		for _, s := range match.List {
			sb.WriteString(s.TimeStart + " " + s.StatusText + " ")
			sb.WriteString(s.Team1 + "：" + s.Team2 + " ")
			sb.WriteString(s.Team1Score + "：" + s.Team2Score + "\n")
		}
	}
	return sb.String()
}
