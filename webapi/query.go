package webapi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"projects/nsfctool/office"
	"projects/nsfctool/util"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Person struct {
	Name string
	Org  string
	Year string
}

type query struct {
	url             string
	capt            *captcha
	headers         map[string]string
	queryStr        string
	encodedQueryStr string
}

func NewQuery(capt *captcha) *query {
	capt.CheckIn()
	headers := util.CopyMap(capt.headers)
	headers["Referer"] = queryUrl // 环境变量
	headers["csrfToken"] = ""

	qry := &query{
		url:     queryUrl,
		capt:    capt,
		headers: headers,
	}
	return qry
}

func (qry *query) UpdateHeaders() {
	qry.headers["Cookie"] = qry.capt.headers["Cookie"]
}

func (qry *query) LoadTask(p Person) {
	s := fmt.Sprintf("psnName:%s,orgName:%s,checkcode:%s,year:%s", p.Name, p.Org, qry.capt.code, p.Year)
	qry.queryStr = s
	qry.encodedQueryStr = url.QueryEscape(s)
}

func (qry *query) FrameReq() error {
	data := map[string][]string{
		"resultDate": {qry.queryStr},
		"checkcode":  {qry.capt.code},
	}
	resp, err := postForm(qry.url, qry.headers, data, qry.capt.timeoutSec)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("检索页面请求返回错误码：%d", resp.StatusCode)
	}
	defer resp.Body.Close()

	newCookie := resp.Header.Get("Set-Cookie")
	if newCookie != "" {
		qry.headers["Cookie"] = newCookie
	}
	return nil
}

func (qry *query) TableReq() (string, error) {
	searchStr := "resultDate^:" + qry.encodedQueryStr + "[tear]sort_name1^:psnName[tear]sort_name2^:title[tear]sort_order^:desc"
	data := map[string][]string{
		"_search":      {"false"},
		"nd":           {strconv.Itoa(int(time.Now().Unix()))},
		"rows":         {"10"},
		"page":         {"1"},
		"sidx":         {""},
		"sord":         {"desc"},
		"searchString": {searchStr},
	}
	url := qry.url + "?flag=grid&checkcode="
	resp, err := postForm(url, qry.headers, data, qry.capt.timeoutSec)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	byteData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(byteData), nil
}

func (qry *query) CorrectResp(s string) bool {
	sign := "<rows><page>1</page><total>"
	return strings.Contains(s, sign)
}

func ParseTable(s string) [][]string {
	var res [][]string
	pat := regexp.MustCompile(`<row id[\s\S]*?</row>`)
	segs := pat.FindAllString(s, -1)
	if len(segs) != 0 {
		patCell := regexp.MustCompile("<cell>(.*?)</cell>")
		for _, seg := range segs {
			items := patCell.FindAllStringSubmatch(seg, -1)
			if len(items) != 0 {
				var temp []string
				for _, v := range items {
					temp = append(temp, v[1])
				}
				res = append(res, temp)
			}
		}
	}

	return res
}

func RunTask(tasks [][]string, timeoutSec int, sleepMilliSec int, predictor Predictor, resPath string) [][]string {
	fmt.Println("需要检索次数：", len(tasks))
	fmt.Println("检索任务开始...")
	// 准备工作
	capt := NewCaptcha(timeoutSec, sleepMilliSec, predictor)
	qry := NewQuery(capt)

	var result [][]string
	colName := []string{
		// "序号",
		"姓名",
		"单位",
		"年份",
		"表格字符串",
		"解析结果",
		"项目名称",
		"项目负责人",
		"依托单位",
		"直接费用",
		"批准年份",
	}
	for i, task := range tasks {
		fmt.Printf("进度: %d, 姓名: %s, 单位:%s, 年份: %s\n", i, task[0], task[1], task[2])
		p := Person{task[0], task[1], task[2]}
		for {
			qry.LoadTask(p)

			qry.FrameReq()
			qry.capt.RandSleep()

			tableXML, _ := qry.TableReq()
			if !qry.CorrectResp(tableXML) {
				qry.capt.CheckIn()
				qry.UpdateHeaders()
				continue
			}

			tableSli := ParseTable(tableXML)
			tableBytes, _ := json.Marshal(tableSli)
			tableStr := string(tableBytes)
			temp := []string{task[0], task[1], task[2], tableXML, tableStr}
			if tableSli != nil {
				for _, inner := range tableSli {
					innerTemp := append(temp, inner...)
					result = append(result, innerTemp)
				}
			} else {
				result = append(result, append(temp, []string{"", "", "", "", ""}...))
			}

			break
		}
		// qry.capt.RandSleep()
		// 每1个循环保存结果
		office.WriteExcel(resPath, result, colName)
	}

	return result
}
