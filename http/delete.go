package http

import (
	"bytes"
	"io"
	"net/http"
)

func DoDelete(url string, data []byte) ([]byte, error) {
	// 创建一个http客户端
	client := &http.Client{}

	var buffer *bytes.Buffer
	if data != nil {
		buffer = bytes.NewBuffer(data)
	} else {
		buffer = nil
	}

	// 创建一个 DELETE 请求
	req, err := http.NewRequest("DELETE", url, buffer)
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
