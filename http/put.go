package http

import (
	"bytes"
	"io"
	"net/http"
)

func DoPut(url string, data []byte) ([]byte, error) {
	// 创建一个http客户端
	client := &http.Client{}

	// 创建一个 PUT 请求
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(data))
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
	return body, nil
}
