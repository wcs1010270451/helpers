package http

import (
	"errors"
	"fmt"
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
	if resp.StatusCode != http.StatusOK {
		return body, errors.New(fmt.Sprintf("状态码为：%d,错误信息：%s", resp.StatusCode, string(body)))
	}
	return body, nil
}
