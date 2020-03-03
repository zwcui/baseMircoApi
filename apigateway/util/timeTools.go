package util

import (
	"time"
	"strconv"
	"math"
)

func BeijingTime() time.Time {
	localTime := time.Now()
	location, _ := time.LoadLocation("Asia/Shanghai")
	beijingTime := localTime.In(location)
	return beijingTime
}

func UnixOfBeijingTime() int64 {
	localTime := time.Now()
	location, _ := time.LoadLocation("Asia/Shanghai")
	beijingTime := localTime.In(location)
	return beijingTime.Unix()
}

// 格式化时间
func TimeDurationFormat(timestamp int64) (timeFormat string) {
	if timestamp == 0 {
		timeFormat = "0分0秒"
		return timeFormat
	}

	var secondStr, minuteStr, hourStr string

	secondLeft := timestamp % 60
	totalMinute := (timestamp - secondLeft) / 60
	minuteLeft := totalMinute % 60
	totalHour := (totalMinute - minuteLeft) / 60
	if totalHour > 0 {
		hourStr = strconv.FormatInt(totalHour, 10) + "小时"
	}
	if minuteLeft > 0 {
		minuteStr = strconv.FormatInt(minuteLeft, 10) + "分"
	}
	if secondLeft > 0  {
		secondStr = strconv.FormatInt(secondLeft, 10) + "秒"
	}

	timeFormat = hourStr + minuteStr + secondStr
	return timeFormat
}

// 计算两者相差年份
// 改为四舍五入
func TimeDifferenceByYear(dateStr string) (yearNum int) {
	date, _ := time.Parse("2006-01-02", dateStr)

	now := BeijingTime()
	formatYear := date.Year()
	currentYear := now.Year()

	//formatYearDay := date.YearDay()
	//currentYearDay := now.YearDay()

	var timeDifference int

	if currentYear - formatYear <= 1 {
		timeDifference = 1
	} else {
		timeDifference = int(math.Floor(float64(now.Unix() - date.Unix()) / float64(365 * 24 * 60 * 60) + 0.5))
	}

	return timeDifference
}

//时间的加减
func TimeParseDuration(second int, addFlag bool) time.Time {
	now := BeijingTime()
	var secondStr = strconv.Itoa(second) + "s"
	if !addFlag {
		secondStr += "-"
	}
	timeFormat, _ := time.ParseDuration(secondStr)
	formatTime := now.Add(timeFormat)
	return formatTime
}


//获取传入的时间所在月份的第一天，即某月第一天的0点。如传入time.Now(), 返回当前月份的第一天0点时间。
func GetFirstDateOfMonth(d time.Time) int64 {
	d = d.AddDate(0, 0, -d.Day() + 1)
	return GetZeroTime(d).Unix()


}
//获取传入的时间所在月份的最后一天，即某月最后一天的0点。如传入time.Now(), 返回当前月份的最后一天0点时间。
func GetLastDateOfMonth(d time.Time) int64 {
	return FirstDateOfMonth(d).AddDate(0, 1, -1).Unix()

}

//获取某一天的0点时间
func GetZeroTime(d time.Time) time.Time {
	return time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, d.Location())
}

//获取传入的时间所在月份的第一天，即某月第一天的0点。如传入time.Now(), 返回当前月份的第一天0点时间。
func FirstDateOfMonth(d time.Time) time.Time {
	d = d.AddDate(0, 0, -d.Day() + 1)
	return GetZeroTime(d)
}
