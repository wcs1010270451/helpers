package time

import (
	"fmt"
	"time"
)

// TimeNowInTimezone 获取当前时间，支持时区
func TimeNowInTimezone() time.Time {
	chinaTimezone, _ := time.LoadLocation("Asia/Shanghai")
	nowTime := time.Now().In(chinaTimezone)
	return nowTime
}

// StringTimeAddDay 时间加n天
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

// MicrosecondsStr 将 time.Duration 类型（nano seconds 为单位）输出为小数点后 3 位的 ms （microsecond 毫秒，千分之一秒）
func MicrosecondsStr(elapsed time.Duration) string {
	return fmt.Sprintf("%.3fms", float64(elapsed.Nanoseconds())/1e6)
}
