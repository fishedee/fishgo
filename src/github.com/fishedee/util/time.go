package util

import "time"

// 取某天的某个小时的时间
func GetTimeInDay(day time.Time, hour int) time.Time {
	return getTime(day, hour, 0, 0, 0)
}

func GetTodayHour(hour int) time.Time {
	now := time.Now()
	return GetTimeInDay(now, hour)
}

func GetTodayZero() time.Time {
	return GetTodayHour(0)
}

func getTime(someDay time.Time, hour, minute, sec, msec int) time.Time {
	location := someDay.Location()
	year, month, day := someDay.Date()
	return time.Date(year, month, day, hour, minute, sec, msec, location)
}
