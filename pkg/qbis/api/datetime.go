package api

import (
	"fmt"
	"regexp"
	"strconv"
	"time"
)

//GetDateForDateTime returns the timestamp of the "date"-datetimes used in qbis. Assumes qbis respects local users TZ.
//QBis represents dates as datetimes with the time set to 00:00:00. They seem to use the local timezone for this.
//On the wire the datetime is serialized in UTC - this can be confusing because the date of the datetime can change.
func GetDateForDateTime(date time.Time) (time.Time, error) {
	/*
		utc, err := time.LoadLocation("UTC")
		if err != nil {
			return date, fmt.Errorf("unable to get UTC location: %v", err)
		}
	*/
	midnightLocalTime := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.Local)

	// We can keep the internal representation in localtime (it's less confusing) and just make sure we serialize the time in UTC before sending to qbis
	//return midnightLocalTime.In(utc), nil
	return midnightLocalTime, nil
}

//DateStringToTime converts a qbis datetime string ( /Date(1520809200000)/ ) to a time
func DateStringToTime(date string) (time.Time, error) {
	//Example: "/Date(1520809200000)/"
	re, err := regexp.Compile(`^\/Date\((\d+?)000\)\/$`)

	matches := re.FindStringSubmatch(date)
	if len(matches) != 2 {
		return time.Unix(0, 0), fmt.Errorf("error parsing date string")
	}

	posix, err := strconv.ParseInt(matches[1], 10, 64)
	if err != nil {
		return time.Unix(0, 0), err
	}

	return time.Unix(posix, 0), nil
}

//TimeToDateString returns the a string representing the time in qbis format
func TimeToDateString(t time.Time) string {
	return fmt.Sprintf("/Date(%d000)/", t.Unix())
}

//TimeToISODateString formats timestamp like javascript
func TimeToISODateString(t time.Time) string {
	return t.UTC().Format("2006-01-02T15:04:05.000Z")
}
