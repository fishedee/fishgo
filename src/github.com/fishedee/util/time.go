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

// 获取上个月 1号的时间
func GetLastMonth() time.Time {
	return GetMonthZero(time.Now()).AddDate(0, -1, 0)
}

// 获取某个月的 1号的时间
func GetMonthZero(someDay time.Time) time.Time {
	location := someDay.Location()
	year, month, _ := someDay.Date()
	return time.Date(year, month, 1, 0, 0, 0, 0, location)
}

func GetTodayZero() time.Time {
	return GetTodayHour(0)
}

func getTime(someDay time.Time, hour, minute, sec, msec int) time.Time {
	location := someDay.Location()
	year, month, day := someDay.Date()
	return time.Date(year, month, day, hour, minute, sec, msec, location)
}

//判断时间，Local或UTC时区是否为1年1月1日0时0分0秒
func IsTimeZero(t time.Time) (bool, error) {
	if t.Location() == time.UTC {
		return t.IsZero(), nil
	}

	loc, err := time.LoadLocation("Local")
	if err != nil {
		return false, err
	}
	locZeroTime, err := time.ParseInLocation("2006-01-02 15:04:05", "0001-01-01 00:00:00", loc)
	if err != nil {
		return false, err
	}
	return t.Equal(locZeroTime), err
}
