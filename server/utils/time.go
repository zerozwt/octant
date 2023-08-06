package utils

import (
	"fmt"
	"regexp"
	"time"
)

var timeReqChecker *regexp.Regexp = regexp.MustCompile(`^[0-9]+$`)

func IsValidTimeString(ts string) bool {
	return len(ts) == 14 && timeReqChecker.MatchString(ts)
}

var cstLoc *time.Location = time.FixedZone("CST", 8*3600)
var errInvalidTimeString error = fmt.Errorf("invalid time string")
var maxDayInMonths []int = []int{31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}

func TimeStringToUTC(ts string) (int64, error) {
	year := strToInt(ts, 0, 4)
	if year == 0 {
		return 0, errInvalidTimeString
	}
	month := strToInt(ts, 4, 6)
	if month < 1 || month > 12 {
		return 0, errInvalidTimeString
	}
	day := strToInt(ts, 6, 8)
	if day < 1 || !(day <= maxDayInMonths[month-1] || (month == 2 && day == 29 && isLunarYear(year))) {
		return 0, errInvalidTimeString
	}
	hour := strToInt(ts, 8, 10)
	if hour > 23 {
		return 0, errInvalidTimeString
	}
	minute := strToInt(ts, 10, 12)
	if minute > 59 {
		return 0, errInvalidTimeString
	}
	second := strToInt(ts, 12, 14)
	if second > 59 {
		return 0, errInvalidTimeString
	}
	return time.Date(year, time.Month(month), day, hour, minute, second, 0, cstLoc).Unix(), nil
}

func strToInt(value string, begin, end int) int {
	switch end - begin {
	case 4:
		return int(value[begin]-'0')*1000 + int(value[begin+1]-'0')*100 + int(value[begin+2]-'0')*10 + int(value[begin+3]-'0')
	case 2:
		return int(value[begin]-'0')*10 + int(value[begin+1]-'0')
	}
	return 0
}

func isLunarYear(year int) bool {
	return (year%400 == 0) || (year%4 == 0 && year%100 != 0)
}

func TimeToCSTString(timestamp int64) string {
	ts := time.Unix(timestamp, 0).In(cstLoc)
	return fmt.Sprintf("%04d-%02d-%02d %02d:%02d:%02d CST",
		ts.Year(), ts.Month(), ts.Day(),
		ts.Hour(), ts.Minute(), ts.Second())
}
