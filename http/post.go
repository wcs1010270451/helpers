package http

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
)

func DoPost(url string, data []byte) ([]byte, error) {
	// 创建一个http客户端
	client := &http.Client{}

	// 创建一个POST请求
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")

	// 发送请求并接收响应
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return body, errors.New(fmt.Sprintf("状态码为 %d", resp.StatusCode))
	}
	return body, nil
}
