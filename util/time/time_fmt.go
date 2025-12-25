package time

import "time"

func TimestampToDateTime(timestamp int64) string {
	if timestamp == 0 {
		return ""
	}
	t := time.Unix(timestamp, 0)

	loc, _ := time.LoadLocation("Asia/Shanghai")
	tInLocal := t.In(loc)

	return tInLocal.Format("2006-01-02 15:04:05")
}

func TimestampToDate(timestamp int64) string {
	if timestamp == 0 {
		return ""
	}
	t := time.Unix(timestamp, 0)

	loc, _ := time.LoadLocation("Asia/Shanghai")
	tInLocal := t.In(loc)

	return tInLocal.Format("2006-01-02")
}
