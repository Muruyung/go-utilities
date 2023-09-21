package converter

import (
	"time"
)

const (
	// DateLayoutFull layout full date
	DateLayoutFull = "2006-01-02 15:04:05"

	// DateLayoutFull2 layout full date 2
	DateLayoutFull2 = "02-01-2006 15:04:05"

	// DateLayoutSimple layout simple date
	DateLayoutSimple = "2006-01-02"

	// DateLayoutStrMonth date layout str month
	DateLayoutStrMonth = "02 Jan 2006"
)

// ConvertDatetime convertion date time to custom layout
func ConvertDatetime(parsedTime time.Time) (time.Time, error) {
	loc, _ := time.LoadLocation("Asia/Jakarta")
	result, err := time.ParseInLocation(DateLayoutFull, parsedTime.Format(DateLayoutFull), loc)
	if err == nil {
		return result, nil
	}

	return time.ParseInLocation(DateLayoutSimple, parsedTime.Format(DateLayoutSimple), loc)
}

// ConvertDateToString convertion date time to string
func ConvertDateToString(parsedTime time.Time) string {
	return parsedTime.Format(DateLayoutFull)
}

// ConvertDateToStringCustom convertion date time to string
func ConvertDateToStringCustom(parsedTime time.Time, layout string) string {
	return parsedTime.Format(layout)
}

// ConvertStringToDate convertion string to date
func ConvertStringToDate(dateString string) (time.Time, error) {
	loc, _ := time.LoadLocation("Asia/Jakarta")
	result, err := time.ParseInLocation(DateLayoutFull, dateString, loc)
	if err == nil {
		return result, nil
	}

	result, err = time.ParseInLocation(time.RFC3339, dateString, loc)
	if err == nil {
		return result, nil
	}

	return time.ParseInLocation(DateLayoutSimple, dateString, loc)
}
