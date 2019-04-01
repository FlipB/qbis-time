package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"time"
)

//SalaryTime contained in TimesheetData
type SalaryTime struct {
	SalaryTimeBase
	MyScheduleDays []ScheduleDay `json:"MyScheduleDays"`
}

//SalaryTimeBase is used as input to save SalaryTime
//SalaryTime activities seem to be activities like any other except they are special
//SalaryTime defines the classification of your work, eg. if you were on leave, sick, or overtime.
//The activities have different types.
//I am guessing that type 3 cancels out scheduled time? Eg. scheduled 8 hours + -8 hours komp = 0 hours worked.
//Komptid is set as default and it will automatically adjust when your working time is adjusted
//A Negative Type 0 will reduce your scheduled working time by that amount.
//A Positive of Type 0 requires you to have increased your working time by that amount.
type SalaryTimeBase struct {
	ActivityBase
	AllowNegative                 bool   `json:"AllowNegative"`
	AllowPositive                 bool   `json:"AllowPositive"`
	AutoFill                      int    `json:"AutoFill"`
	CalculationUnit               int    `json:"CalculationUnit"`
	DisplayFormat                 int    `json:"DisplayFormat"`
	HasParentGroupActivity        bool   `json:"HasParentGroupActivity"`
	IsDefault                     bool   `json:"IsDefault"`
	IsDeletable                   bool   `json:"IsDeletable"`
	IsHoursWorkedIntervalActivity bool   `json:"IsHoursWorkedIntervalActivity"`
	IsIntervalActivity            bool   `json:"IsIntervalActivity"`
	Locked                        bool   `json:"Locked"`
	LowerLimit                    int    `json:"LowerLimit"`
	ParentGroupActivityID         int    `json:"ParentGroupActivityId"`
	PresentationUnit              string `json:"PresentationUnit"`
	ShowErrorMessage              bool   `json:"ShowErrorMessage"`
	ShowWarningMessage            bool   `json:"ShowWarningMessage"`
	SpecifyClockTimes             bool   `json:"SpecifyClockTimes"`
	Type                          int    `json:"Type"`
	UpperLimit                    int    `json:"UpperLimit"`
}

//EmployeeSalaryTime struct represents the payload used when saving the matrix containing overtime, sick-time etc.
type EmployeeSalaryTime struct {
	EmployeeID        string             `json:"employeeId"`
	FromDate          string             `json:"fromDate"`
	ToDate            string             `json:"toDate"`
	SalaryTime        []SalaryTimeBase   `json:"salaryTime"`
	WorkingTime       []WorkingTimeBase  `json:"workingTime"`
	WorkingTimeBreaks []WorkingTimeBreak `json:"workingTimeBreaks"`
}

//SaveSalaryTimeResponse represents the object returned on attempted save of SalaryTime
type SaveSalaryTimeResponse struct {
	FromToValidationResult  string        `json:"fromToValidationResult"`
	InvalidFullDaysError    string        `json:"invalidFullDaysError"`
	InvalidPartialDaysError string        `json:"invalidPartialDaysError"`
	LimitErrorResults       []interface{} `json:"limitErrorResults"`
	LimitWarningResults     []interface{} `json:"limitWarningResults"`
	ResetWarning            string        `json:"resetWarning"`
	WasSaved                bool          `json:"wasSaved"`
}

//SalaryActivityOverview ...
type SalaryActivityOverview struct {
	ActivityOverviewBase
	ShowInMyTime bool `json:"ShowInMyTime"`
}

//TotalBudgetMinutes returns the total budget for all employees on the activity (in minutes)
func (a SalaryActivityOverview) TotalBudgetMinutes() (uint, error) {
	return a.getPropertyMinutes("activityHours")
}

//BudgetMinutes returns the budget for the employee on the activity (in minutes)
func (a SalaryActivityOverview) BudgetMinutes() (uint, error) {
	return a.getPropertyMinutes("allocatedHours")
}

//SpentMinutes returns the number of minutes registred by the employee on the activity
func (a SalaryActivityOverview) SpentMinutes() (uint, error) {
	return a.getPropertyMinutes("registeredHours")
}

//GetSalaryActivityOverview ..
func (c *Client) GetSalaryActivityOverview(employee string, activityID int, from time.Time, to time.Time) (*SalaryActivityOverview, error) {
	url := "/Time/TimesheetSalaryTime/GetActivityOverview?activityId=%d&employeeId=%s&fromDate=%s&toDate=%s&_=%s"

	fromString := from.UTC().Format("2006-01-02T15:04:05.000Z")
	toString := to.UTC().Format("2006-01-02T15:04:05.000Z")
	timestamp := strconv.Itoa(int(time.Now().UnixNano()))

	url = fmt.Sprintf(url, activityID, employee, fromString, toString, timestamp)
	response, err := c.get(url)
	if err != nil {
		return nil, err
	}
	var overview SalaryActivityOverview
	err = json.NewDecoder(response.Body).Decode(&overview)
	if err != nil {
		return nil, err
	}

	return &overview, nil
}

//SaveSalaryTime saves the matrix with hours of overtime, sickleave, timebank etc.
func (c *Client) SaveSalaryTime(time EmployeeSalaryTime) (*SaveSalaryTimeResponse, error) {

	var b bytes.Buffer
	err := json.NewEncoder(&b).Encode(time)
	if err != nil {
		return nil, err
	}
	//fmt.Printf("payload: \n%s\n", b.String())
	response, err := c.postJSON("/Time/TimesheetSalaryTime/SaveSalaryTime", &b)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}
	err = response.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("error closing response body: %v", err)
	}

	//fmt.Printf("body:\n\n%s\n\n\n", body)

	var saveResponse SaveSalaryTimeResponse
	err = json.Unmarshal(body, &saveResponse)
	//err = json.NewDecoder(response.Body).Decode(&saveResponse)
	if err != nil {
		return nil, err
	}
	return &saveResponse, nil
}

//GetSalaryActivity returns the full details of the SalaryActivity for the employee
func (c *Client) GetSalaryActivity(employee string, activityID int, from time.Time, to time.Time) (*SalaryTime, error) {
	url := "/Time/TimesheetSalaryTime/GetActivityInformation?activityId=%d&employeeId=%s&fromDate=%s&toDate=%s&_=%s"

	fromString := from.UTC().Format("2006-01-02T15:04:05.000Z")
	toString := to.UTC().Format("2006-01-02T15:04:05.000Z")
	timestamp := strconv.Itoa(int(time.Now().UnixNano()))

	url = fmt.Sprintf(url, activityID, employee, fromString, toString, timestamp)
	response, err := c.get(url)
	if err != nil {
		return nil, err
	}

	var details SalaryTime
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}
	err = response.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("error closing response body: %v", err)
	}

	err = json.Unmarshal(body, &details)
	//err = json.NewDecoder(response.Body).Decode(&details)
	if err != nil {
		return nil, err
	}

	return &details, nil
}
