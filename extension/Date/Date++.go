package Date

				import (
					"DairoNPS/DebugTimer"
				)

import (
	"time"
)

// 日期格式化
func FormatByTimespan(timespan int64) string {
DebugTimer.Add110()
	t := time.Unix(timespan/1000, 0)
	return Format(t)
}

// 日期格式化
func FormatDateByTimespan(timespan int64) string {
DebugTimer.Add111()
	t := time.Unix(timespan/1000, 0)
	return FormatDate(t)
}

// 日期格式化
func Format(t time.Time) string {
DebugTimer.Add112()
	return t.Format("2006-01-02 15:04:05")
}

// 日期格式化
func FormatDate(t time.Time) string {
DebugTimer.Add113()
	return t.Format("2006-01-02")
}
