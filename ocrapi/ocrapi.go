package ocrapi

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"projects/nsfctool/util"
	"runtime"
)

func init() {
	// 确保目录存在
	dir := "./temp"
	_, err := os.Stat(dir)
	if err == nil {
		return
	}
	if os.IsNotExist(err) {
		os.Mkdir(dir, 0777)
	}

	// 确保是windows系统
	runningOS := runtime.GOOS
	if runningOS != "windows" {
		panic("程序需要运行在windows环境.")
	}
}

// tesseract路径
var Tess = filepath.Join(util.BaseDir(), "Tesseract-OCR/tesseract.exe")

// 临时文件
const IMG = "./temp/img.jpg"
const RESULT = "./temp/result"
const TXT = RESULT + ".txt"

// func clean(s string) string {
// 	s = strings.TrimSpace(s)
// 	s = strings.Replace(s, " ", "", -1)
// 	return s
// }

func Predict(imgBytes []byte) (string, error) {
	if util.PathExist(IMG) {
		os.Remove(IMG)
	}
	err := ioutil.WriteFile(IMG, imgBytes, 0777)
	if err != nil {
		return "", err
	}
	// defer os.Remove(IMG)

	if util.PathExist(TXT) {
		os.Remove(TXT)
	}
	err = RunTesseract(IMG, RESULT, "eng")
	if err != nil {
		return "", nil
	}
	// defer os.Remove(TXT)

	// 等待文本文件出现
	condition := func() bool {
		return util.PathExist(TXT)
	}
	util.Wait(condition, -1, -1)

	textBytes, err := ioutil.ReadFile(TXT)
	if err != nil {
		return "", nil
	}
	s := util.Clean(string(textBytes))

	return s, nil
}

func RunTesseract(imgPath, resPath, lang string) error {
	line := Tess + " " + imgPath + " " + RESULT + " -l " + lang
	// fmt.Println(line)
	cmd := exec.Command("cmd", "/C", line)
	err := cmd.Start()
	return err
}

// 实现了预测的接口
type TessPredictor struct{}

func (tess TessPredictor) Predict(byteData []byte) (string, error) {
	return Predict(byteData)
}
