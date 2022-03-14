package util

import (
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/gofrs/uuid"
)

// 获取项目根路径
func BaseDir() string {
	execP, _ := os.Executable()
	dir := path.Dir(execP)
	baseDir, _ := filepath.Abs(dir)
	return baseDir
}

// 判断路径存在
func PathExist(p string) bool {
	_, err := os.Stat(p)
	if err != nil {
		return false
	}
	return true
}

// 等待条件出现
func Wait(condition func() bool, sleepMilliSec int, maxnum int) {
	if sleepMilliSec < 0 {
		sleepMilliSec = 500
	}
	if maxnum < 0 {
		maxnum = 1000
	}
	for i := 0; i < maxnum; i++ {
		if condition() {
			break
		}
		time.Sleep(time.Duration(sleepMilliSec) * 1000000)
	}
}

// 清空空格
func Clean(s string) string {
	s = strings.TrimSpace(s)
	s = strings.Replace(s, " ", "", -1)
	return s
}

// 复制map
func CopyMap(m map[string]string) map[string]string {
	n := make(map[string]string, 0)
	for k, v := range m {
		n[k] = v
	}
	return n
}

// 获取uuid string
func UUIDString() string {
	ui, _ := uuid.NewV4()
	return ui.String()
}
