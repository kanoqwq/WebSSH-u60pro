package utils

import (
	"math/rand"
	"time"
	"unicode/utf8"
	"path/filepath"
	"os"
)

// RandString 生成指定长度随机字符串
func RandString(length int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyz"
	data := []byte(str)
	var result []byte
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		result = append(result, data[r.Intn(len(data))])
	}
	return string(result)
}

func TruncateString(s string, length int) string {
	if utf8.RuneCountInString(s) <= length {
		return s
	}

	runes := []rune(s)
	if length < 1 {
		return ""
	}
	return string(runes[:length])
}

func GetExecDir() string {
	// 获取当前工作目录
	dir, err := os.Getwd()
	if err == nil && dir != "" {
		absDir, err2 := filepath.Abs(dir)
		if err2 == nil {
			return absDir
		}
		return dir
	}

	// os.Getwd 失败时 fallback
	tmpDir := os.TempDir()
	if tmpDir != "" {
		return tmpDir
	}

	// 最后兜底
	return "."
}