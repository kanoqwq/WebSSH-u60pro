package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

// CallUBUS 执行 ubus call 并返回 JSON 结果
func GetDataFromUbus(service, method string, params map[string]interface{}) (map[string]interface{}, error) {
	var cmdArgs []string
	if params != nil {
		// 转为 JSON 参数
		paramsJSON, err := json.Marshal(params)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal params: %w", err)
		}
		cmdArgs = []string{"call", service, method, string(paramsJSON)}
	} else {
		cmdArgs = []string{"call", service, method}
	}

	// 设置超时，避免 ubus 阻塞
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "ubus", cmdArgs...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("ubus call failed: %w, output: %s", err, string(output))
	}

	// 解析返回 JSON
	var result map[string]interface{}
	if err := json.Unmarshal(output, &result); err != nil {
		// 如果不是 JSON（比如空字符串），也返回原始文本
		return map[string]interface{}{"raw": strings.TrimSpace(string(output))}, nil
	}

	return result, nil
}