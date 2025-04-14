package ytdlp

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

// Init 空函数, 目的是触发包内的初始化
func Init() {}

// Extract 使用指定的 formatCode 解析 url
func Extract(url, formatCode string) (string, error) {
	if !execOk {
		return "", errors.New("yt-dlp 环境未初始化")
	}

	// 构造命令
	cmd := exec.Command(
		execPath,
		"-f", formatCode,
		url,
		"--get-url",
	)
	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf

	// 执行获取输出
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("执行命令失败: %v, err: %s", err, errBuf.String())
	}

	// 校验结果
	output := outBuf.String()
	scanner := bufio.NewScanner(strings.NewReader(output))
	if !scanner.Scan() {
		return "", fmt.Errorf("解析出非预期结果, 原始输出: %s", output)
	}

	line := scanner.Text()
	if !strings.HasPrefix(line, "http") {
		return "", fmt.Errorf("解析出非预期结果, 原始输出: %s", output)
	}

	return line, nil
}
