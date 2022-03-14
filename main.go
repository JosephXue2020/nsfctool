package main

import (
	"fmt"
	"projects/nsfctool/ocrapi"
	"projects/nsfctool/office"
	"projects/nsfctool/util"
	"projects/nsfctool/webapi"
	"strconv"
	"strings"
)

func main() {

	// 获取任务
	tempPath := "./template.xlsx"
	tasks := getTask(tempPath)[1:]

	// 执行任务
	resPath := "./result.xlsx"
	var predictor ocrapi.TessPredictor
	webapi.RunTask(tasks, 3, 200, predictor, resPath)

	// 任务停住
	fmt.Print("按任意键退出: ")
	fmt.Scanf("%s")
}

func getTask(path string) [][]string {
	items, err := office.ReadExcel(path, "Sheet1")
	if err != nil {
		panic("excel表格不合要求.")
	}

	var result [][]string
	for _, item := range items {
		years := item[2]
		years = util.Clean(years)
		if strings.Contains(years, "-") {
			segs := strings.Split(years, "-")
			if len(segs) > 1 {
				start, err := strconv.Atoi(segs[0])
				if err != nil {
					panic("excel表格中年份列存在问题.")
				}
				end, err := strconv.Atoi(segs[1])
				if err != nil {
					panic("excel表格中年份列存在问题.")
				}
				if start >= end {
					panic("excel表格中年份列存在问题.")
				}

				for i := start; i <= end; i++ {
					temp := []string{item[0], item[1], strconv.Itoa(i)}
					result = append(result, temp)
				}
			} else if len(segs) == 1 {
				temp := []string{item[0], item[1], segs[0]}
				result = append(result, temp)
			} else {
				panic("excel表格中年份列存在问题.")
			}
		} else {
			result = append(result, item)
		}
	}

	return result
}

func test() {
	// 获取任务
	tempPath := "./template.xlsx"
	tasks := getTask(tempPath)
	fmt.Println("查询任务总数：", len(tasks))

	// 网络交互
	var predictor ocrapi.TessPredictor
	capt := webapi.NewCaptcha(3, 300, predictor)
	err := capt.GetImg()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(capt)
	capt.CheckIn()
	// capt.CollectTraingData("./example", 2)
}
