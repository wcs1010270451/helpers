package time

import (
	"time"
)

// 时间加n天
func StringTimeAddDay(timeStr string, days int) (string, error) {
	format := "2006-01-02"
	// 解析字符串为time.Time类型
	parsedTime, err := time.Parse(format, timeStr)
	if err != nil {
		return "", err
	}

	// 加一天
	parsedTime = parsedTime.AddDate(0, 0, days)

	// 格式化回字符串
	newStrTime := parsedTime.Format(format)
	return newStrTime, nil
}
