package qbis

import (
	"time"

	"github.com/flipb/qbis-time/pkg/qbis/api"
)

//getWeekSpan returns the first and last datetime of the week containing the given date
// the datetimes returned are so called "qbis dates", datetimes with timestamp 00:00:00 and local timezone
func getWeekSpan(date time.Time) (start, end time.Time, err error) {
	start = date
	for start.Weekday() != time.Monday {
		start = start.AddDate(0, 0, -1)
	}

	start, err = api.GetDateForDateTime(start)
	if err != nil {
		return start, date, err
	}
	/*
		start = start.Add(-1 * time.Duration(start.Hour()) * time.Hour)
		start = start.Add(-1 * time.Duration(start.Minute()) * time.Minute)
		start = start.Add(-1 * time.Duration(start.Second()) * time.Second)
	*/
	end = date
	for end.Weekday() != time.Sunday {
		end = end.AddDate(0, 0, 1)
	}
	end, err = api.GetDateForDateTime(end)
	if err != nil {
		return start, end, err
	}
	/*
		end = end.Add(time.Duration(23-end.Hour()) * time.Hour)
		end = end.Add(time.Duration(59-end.Minute()) * time.Minute)
		end = end.Add(time.Duration(59-end.Second()) * time.Second)
	*/
	return start, end, nil
}
