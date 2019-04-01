package api

import (
	"bytes"
	"encoding/json"
)

//WorkingTime ...
type WorkingTime struct {
	WorkingTimeBase
	ScheduleDay ScheduleDay `json:"ScheduleDay"`
}

//WorkingTimeBase the base working time struct is embedded in another struct in TimesheetData. Also used as input to save data
type WorkingTimeBase struct {
	Arrive int `json:"Arrive"`
	Breaks []struct {
		BreakDate        string `json:"BreakDate"`
		BreakFromMinutes int    `json:"BreakFromMinutes"`
		BreakID          int    `json:"BreakId"`
		BreakToMinutes   int    `json:"BreakToMinutes"`
		EmployeeID       int    `json:"EmployeeId"`
		Source           int    `json:"Source"`
	} `json:"Breaks"`
	DayDate                    string `json:"DayDate"`
	DayName                    string `json:"DayName"`
	HasSchedule                bool   `json:"HasSchedule"`
	ID                         int    `json:"ID"`
	IsModified                 bool   `json:"IsModified"`
	IsMonthClosed              bool   `json:"IsMonthClosed"`
	IsOutsideJoinAndLeaveDates bool   `json:"IsOutsideJoinAndLeaveDates"`
	IsPublicHoliday            bool   `json:"IsPublicHoliday"`
	IsSaved                    bool   `json:"IsSaved"`
	IsScheduleHourOnly         bool   `json:"IsScheduleHourOnly"`
	IsToday                    bool   `json:"IsToday"`
	Leave                      int    `json:"Leave"`
	Locked                     bool   `json:"Locked"`
	Lunch                      int    `json:"Lunch"`
	NextDay                    string `json:"NextDay"`
	Overmidnight               bool   `json:"Overmidnight"`
	OverridePublicHolidays     bool   `json:"OverridePublicHolidays"`
	PrefillSpecifiedBreak      bool   `json:"PrefillSpecifiedBreak"`
	ScheduledHours             int    `json:"ScheduledHours"`
	Total                      int    `json:"Total"`
}

// WorkingTimeBreak is read in timesheet and used in Save salary time
type WorkingTimeBreak struct {
	Days []struct {
		BreakDate        string `json:"BreakDate"`
		BreakFromMinutes int    `json:"BreakFromMinutes"`
		BreakID          int    `json:"BreakId"`
		BreakToMinutes   int    `json:"BreakToMinutes"`
		EmployeeID       int    `json:"EmployeeId"`
		Source           int    `json:"Source"`
	} `json:"Days"`
	IsNewRow bool `json:"IsNewRow"`
}

//EmployeeWorkingTime represents time spent by employee working. This is the matrix containing arrival, departure and lunch time
type EmployeeWorkingTime struct {
	Days       []WorkingTimeBase `json:"days"`
	EmployeeID string            `json:"employeeId"`
	FromDate   string            `json:"fromDate"`
	ToDate     string            `json:"toDate"`
}

//SaveWorkingTimeResponse ...
type SaveWorkingTimeResponse struct {
	Saved string `json:"saved"`
}

//SaveWorkingTime saves the matrix with hours spent working
func (c *Client) SaveWorkingTime(time EmployeeWorkingTime) (*SaveWorkingTimeResponse, error) {
	// https://login.qbis.se/Time/TimesheetWorkingTime/SaveWorkingTime
	var b bytes.Buffer
	err := json.NewEncoder(&b).Encode(time)
	if err != nil {
		return nil, err
	}
	response, err := c.postJSON("/Time/TimesheetWorkingTime/SaveWorkingTime", &b)
	if err != nil {
		return nil, err
	}
	var saveResponse SaveWorkingTimeResponse
	err = json.NewDecoder(response.Body).Decode(&saveResponse)
	if err != nil {
		return nil, err
	}
	return &saveResponse, nil
}
