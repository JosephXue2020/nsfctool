package webapi

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"projects/nsfctool/util"
	"strings"
	"time"
)

// ocr接口
type Predictor interface {
	Predict([]byte) (string, error)
}

type captcha struct {
	url           string
	validUrl      string
	headers       map[string]string // 及时更新
	timeoutSec    int
	sleepMilliSec int
	predictor     Predictor
	data          []byte
	code          string
}

var captchaUrl = "https://isisn.nsfc.gov.cn/egrantindex/validatecode.jpg"
var validUrl = "https://isisn.nsfc.gov.cn/egrantindex/funcindex/validate-checkcode"

func NewCaptcha(timeoutSec int, sleepMilliSec int, predictor Predictor) *captcha {
	capt := new(captcha)
	capt.url = captchaUrl
	capt.validUrl = validUrl
	capt.headers = headers
	capt.timeoutSec = timeoutSec
	capt.sleepMilliSec = sleepMilliSec
	capt.predictor = predictor

	return capt
}

// 获取当前时间
func NowTime() string {
	now := time.Now()
	s := now.Format("Mon Jan 02 2006 15:04:05 GMT 0800 (中国标准时间)")
	return s
}

// 获取一个验证码
func (capt *captcha) GetImg() error {
	param := map[string]string{"date": NowTime()}
	resp, err := get(capt.url, capt.headers, param, capt.timeoutSec)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	byteData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	capt.data = byteData
	if capt.predictor != nil {
		capt.code, _ = capt.predictor.Predict(byteData)
	}

	return nil
}

func (capt *captcha) Valid() bool {
	if len(capt.data) == 0 {
		return false
	}
	if capt.code == "" {
		return false
	}

	headers := util.CopyMap(capt.headers)
	headers["x-requested-with"] = "XMLHttpRequest"
	data := map[string][]string{
		"checkCode": {capt.code},
	}
	resp, err := postForm(capt.validUrl, headers, data, 3)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	byteData, _ := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false
	}
	msg := string(byteData)
	if msg == "error" || msg == "null" {
		return false
	}

	newCookie := resp.Header.Get("Set-Cookie")
	if newCookie != "" && strings.Contains(newCookie, "sessionidindex") {
		capt.headers["Cookie"] = newCookie
	}

	return true
}

// 返回睡眠时间1倍和2倍之间的随机数
func (capt *captcha) RandInt() int {
	return rand.Intn(capt.sleepMilliSec) + capt.sleepMilliSec
}

func (capt *captcha) RandSleep() {
	milliSec := capt.RandInt()
	time.Sleep(time.Duration(milliSec) * 10e6)
}

func (capt *captcha) CheckIn() {
	fmt.Println("验证码验证中...")
	for {
		err := capt.GetImg()
		if err != nil {
			capt.RandSleep()
			continue
		}
		valided := capt.Valid()
		if valided {
			break
		}
	}
	fmt.Println("验证码验证成功")
}

func (capt *captcha) CollectTraingData(dir string, num int) {
	if !util.PathExist(dir) {
		os.Mkdir(dir, 0777)
	}

	for i := 0; i < num; i++ {
		capt.CheckIn()
		fname := util.UUIDString() + "_" + capt.code + ".jpg"
		path := filepath.Join(dir, fname)
		ioutil.WriteFile(path, capt.data, 0777)
		fmt.Println("获取到训练数据：", i+1)
	}
}
