package api

import (
	"fmt"
	"strconv"
	"strings"
)

//ActivityBase contains fields common between SalaryTimeBase, SalaryActivity, ProjectActivity
type ActivityBase struct {
	ActivityActive bool   `json:"ActivityActive"`
	ActivityID     int    `json:"ActivityId"`
	ActivityName   string `json:"ActivityName"`
	Days           []struct {
		DayID            int    `json:"DayId"`
		DayDate          string `json:"DayDate"`
		DayMinutes       int    `json:"DayMinutes"`
		DayDays          int    `json:"DayDays"`
		Notes            string `json:"Notes"`
		Delete           bool   `json:"Delete"`
		Locked           bool   `json:"Locked"`
		IsReadOnly       bool   `json:"IsReadOnly"`
		IsPrefilled      bool   `json:"IsPrefilled"`
		AllowEmptyFromTo bool   `json:"AllowEmptyFromTo"`
		LunchOffset      int    `json:"LunchOffset"`
		DayFromMinutes   int    `json:"DayFromMinutes"`
		DayToMinutes     int    `json:"DayToMinutes"`
	} `json:"Days"`
	EmployeeID int      `json:"EmployeeId"`
	Factor     float64  `json:"Factor"`
	IsNewRow   bool     `json:"IsNewRow"`
	IsReadOnly bool     `json:"IsReadOnly"`
	Tooltip    []string `json:"Tooltip"`
}

/*ActivityOverviewBase see full comment for details
// Example of Properties
[
	{
		// activityHours seems to be the total activity budget (allocations for multiple people)
		TextIdentifier:	"activityHours",
		Value	"0",
		UnitIdentifier	"unitHours"
	},
	{
		// Allocated hours seems to be the current employees allocated budget hours
		TextIdentifier	"allocatedHours",
		Value	"0",
		UnitIdentifier	"unitHours"
	},
	{
		// Allocated hours seems to be Spent hours. Not sure if these are only for the employee or everyone
		TextIdentifier	"registeredHours",
		Value	"3,25",
		UnitIdentifier	"unitHours"
	},
	{
		TextIdentifier:	"startDate",
		Value:	"31<sup>st</sup> Oct 2017",
		UnitIdentifier:	"unitEmpty"
	},
	{
		TextIdentifier:	"endDate",
		Value:	"2<sup>nd</sup> Dec 2022",
		UnitIdentifier:	"unitEmpty"
	}
]
*/
//ActivityOverviewBase contains fields common between SalaryActivityOverview and ProjectActivityOverview
type ActivityOverviewBase struct {
	DisplayName string `json:"DisplayName"`
	Properties  []struct {
		TextIdentifier string `json:"TextIdentifier"`
		Value          string `json:"Value"`
		UnitIdentifier string `json:"UnitIdentifier"`
	} `json:"Properties"`
}

//getPropertyMinutes parses and returns the number of minutes of the given property in the activity overview
func (a ActivityOverviewBase) getPropertyMinutes(propIdentifier string) (uint, error) {
	for _, prop := range a.Properties {
		if prop.TextIdentifier != propIdentifier {
			continue
		}

		if prop.UnitIdentifier != "unitHours" {
			return 0, fmt.Errorf("parsing error: unexpected UnitIdentifier")
		}
		// we replace commas with periods because the decimal separator depends on the user settings
		decimal := strings.Replace(prop.Value, ",", ".", -1)
		float, err := strconv.ParseFloat(decimal, 64)
		if err != nil {
			return 0, fmt.Errorf("unable to parse float: %v", err)
		}
		return uint(float * 60), nil
	}
	return 0, fmt.Errorf("unable to find property '%s'", propIdentifier)
}
