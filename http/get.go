package http

import (
	"io"
	"net/http"
)

func DoGet(url string) ([]byte, error) {
	// 创建一个 GET 请求
	resp, err := http.Get(url)

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
