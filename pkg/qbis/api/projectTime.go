package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

//ProjectActivityListItem is a category of work that belongs to a project
type ProjectActivityListItem struct {
	Factor     string `json:"Factor"`
	FactorLock bool   `json:"FactorLock"`
	ID         int    `json:"ID"`
	Name       string `json:"Name"`
	ReadOnly   bool   `json:"ReadOnly"`
}

//SaveProjectTimeResponse represents the response object when saving ProjectTime
type SaveProjectTimeResponse struct {
	DuplicatedActivitiesError string        `json:"duplicatedActivitiesError"`
	LimitErrorResults         []interface{} `json:"limitErrorResults"`
	LimitWarningResults       []interface{} `json:"limitWarningResults"`
	WasSaved                  bool          `json:"wasSaved"`
}

//EmployeeProjectTime represents time spent by employee on project activities
type EmployeeProjectTime struct {
	EmployeeID string        `json:"employeeId"`
	FromDate   string        `json:"fromDate"`
	List       []ProjectTime `json:"list"`
	ToDate     string        `json:"toDate"`
}

// ProjectTime ...
type ProjectTime struct {
	ProjectTimeBase
	Days []struct {
		DayDate                   string      `json:"DayDate"`
		DayID                     int         `json:"DayId"`
		DayMinutes                int         `json:"DayMinutes"`
		DayTime                   string      `json:"DayTime"`
		ExternalNotes             string      `json:"ExternalNotes"`
		InternalNotes             string      `json:"InternalNotes"`
		IsInvoiced                bool        `json:"IsInvoiced"`
		IsOutsideActivityDateSpan bool        `json:"IsOutsideActivityDateSpan"`
		IsReadOnly                bool        `json:"IsReadOnly"`
		IsStaffLedgerRegistration bool        `json:"IsStaffLedgerRegistration"`
		ProjectTimeApprovedBy     interface{} `json:"ProjectTimeApprovedBy"`
		ProjectTimeApprovedByID   int         `json:"ProjectTimeApprovedById"`
		ReadOnlyColor             string      `json:"ReadOnlyColor"`
	} `json:"Days"`
}

//ProjectTimeBase ...
type ProjectTimeBase struct {
	ActivityBase
	ActivityComplete            bool   `json:"ActivityComplete"`
	ActivityDateSpanString      string `json:"ActivityDateSpanString"`
	ActivitySalary              bool   `json:"ActivitySalary"`
	Autofill                    bool   `json:"Autofill"`
	CustomerFullName            string `json:"CustomerFullName"`
	CustomerName                string `json:"CustomerName"`
	CustomerProjectActivityName string `json:"CustomerProjectActivityName"`
	EndDate                     string `json:"EndDate"`
	FixedPrice                  bool   `json:"FixedPrice"`
	FromServiceRequest          bool   `json:"FromServiceRequest"`
	IsDeleteable                bool   `json:"IsDeleteable"`
	IsFixedPriceFactor          bool   `json:"IsFixedPriceFactor"`
	IsProjectTimeApproved       bool   `json:"IsProjectTimeApproved"`
	IsReadOnlyFactor            bool   `json:"IsReadOnlyFactor"`
	IsVisibleFactor             bool   `json:"IsVisibleFactor"`
	LockFactor                  bool   `json:"LockFactor"`
	PhaseName                   string `json:"PhaseName"`
	ProjectFullName             string `json:"ProjectFullName"`
	ProjectName                 string `json:"ProjectName"`
	ProjectTimeApprovedBy       string `json:"ProjectTimeApprovedBy"`
	ReadOnlyColor               string `json:"ReadOnlyColor"`
	ReadOnlyFactorColor         string `json:"ReadOnlyFactorColor"`
	StartDate                   string `json:"StartDate"`
	Week                        string `json:"Week"`
	YearWeek                    int    `json:"YearWeek"`
}

//ProjectCompany contains a company and it's projetcs
type ProjectCompany struct {
	CompanyID   int       `json:"CompanyID"`
	CompanyName string    `json:"CompanyName"`
	Projects    []Project `json:"Projects"`
}

//Project represents a qbis project that activities are tied to
type Project struct {
	Code string `json:"Code"`
	ID   int    `json:"ID"`
	Name string `json:"Name"`
}

//ProjectActivityOverview represents the data shown in the overlay when you are entering time in the ui (see comment for details!)
type ProjectActivityOverview struct {
	ActivityOverviewBase
	CompanyProjectDisplayName string `json:"CompanyProjectDisplayName"`
	Details                   string `json:"Details"`
}

//TotalBudgetMinutes returns the total budget for all employees on the activity (in minutes)
func (a ProjectActivityOverview) TotalBudgetMinutes() (uint, error) {
	return a.getPropertyMinutes("activityHours")
}

//BudgetMinutes returns the budget for the employee on the activity (in minutes)
func (a ProjectActivityOverview) BudgetMinutes() (uint, error) {
	return a.getPropertyMinutes("allocatedHours")
}

//SpentMinutes returns the number of minutes registred by the employee on the activity
func (a ProjectActivityOverview) SpentMinutes() (uint, error) {
	return a.getPropertyMinutes("registeredHours")
}

//GetProjects ....
func (c *Client) GetProjects(employee string, from time.Time, to time.Time) ([]ProjectCompany, error) {
	url := "/Time/TimesheetProjectTime/GetCustomerProjectDropDown?employeeId=%s&fromDate=%s&toDate=%s&selectedID=%s&_=%s"

	fromString := from.UTC().Format("2006-01-02T15:04:05.000Z")
	toString := to.UTC().Format("2006-01-02T15:04:05.000Z")
	timestamp := strconv.Itoa(int(time.Now().UnixNano()))

	selected := "0"

	url = fmt.Sprintf(url, employee, fromString, toString, selected, timestamp)
	response, err := c.get(url)
	if err != nil {
		return nil, err
	}
	projectCompanies := make([]ProjectCompany, 0)

	err = json.NewDecoder(response.Body).Decode(&projectCompanies)
	if err != nil {
		return nil, err
	}

	return projectCompanies, nil
}

//GetProjectActivityList returns a slice of activities for the given project id
func (c *Client) GetProjectActivityList(employee string, projectID int, from time.Time, to time.Time) ([]ProjectActivityListItem, error) {
	url := "/Time/TimesheetProjectTime/GetActivityDropDown?employeeId=%s&fromDate=%s&toDate=%s&selectedID=%s&projectID=%d&_=%s"

	fromString := from.UTC().Format("2006-01-02T15:04:05.000Z")
	toString := to.UTC().Format("2006-01-02T15:04:05.000Z")
	timestamp := strconv.Itoa(int(time.Now().UnixNano()))

	selected := "0"

	url = fmt.Sprintf(url, employee, fromString, toString, selected, projectID, timestamp)
	response, err := c.get(url)
	if err != nil {
		return nil, err
	}

	activies := make([]ProjectActivityListItem, 0)
	err = json.NewDecoder(response.Body).Decode(&activies)
	if err != nil {
		return nil, err
	}

	return activies, nil
}

//GetProjectActivity ...
func (c *Client) GetProjectActivity(employee string, activityID int, from time.Time, to time.Time) (*ProjectTime, error) {
	url := "/Time/TimesheetProjectTime/GetActivityInformation?activityId=%d&employeeId=%s&fromDate=%s&toDate=%s&_=%s"

	fromString := from.UTC().Format("2006-01-02T15:04:05.000Z")
	toString := to.UTC().Format("2006-01-02T15:04:05.000Z")
	timestamp := strconv.Itoa(int(time.Now().UnixNano()))

	url = fmt.Sprintf(url, activityID, employee, fromString, toString, timestamp)
	response, err := c.get(url)
	if err != nil {
		return nil, err
	}

	var details ProjectTime
	err = json.NewDecoder(response.Body).Decode(&details)
	if err != nil {
		return nil, err
	}

	return &details, nil
}

//GetProjectActivityOverview ..
func (c *Client) GetProjectActivityOverview(employee string, activityID int, from time.Time, to time.Time) (*ProjectActivityOverview, error) {
	url := "/Time/TimesheetProjectTime/GetActivityOverview?activityId=%d&employeeId=%s&fromDate=%s&toDate=%s&_=%s"

	fromString := from.UTC().Format("2006-01-02T15:04:05.000Z")
	toString := to.UTC().Format("2006-01-02T15:04:05.000Z")
	timestamp := strconv.Itoa(int(time.Now().UnixNano()))

	url = fmt.Sprintf(url, activityID, employee, fromString, toString, timestamp)
	response, err := c.get(url)
	if err != nil {
		return nil, err
	}
	var overview ProjectActivityOverview
	err = json.NewDecoder(response.Body).Decode(&overview)
	if err != nil {
		return nil, err
	}

	return &overview, nil
}

//SaveProjectTime saves the matrix with hours spent on project activities
func (c *Client) SaveProjectTime(time EmployeeProjectTime) (*SaveProjectTimeResponse, error) {
	var b bytes.Buffer
	err := json.NewEncoder(&b).Encode(time)
	if err != nil {
		return nil, err
	}
	response, err := c.postJSON("/Time/TimesheetProjectTime/SaveProjectTime", &b)
	if err != nil {
		return nil, err
	}
	var saveResponse SaveProjectTimeResponse
	err = json.NewDecoder(response.Body).Decode(&saveResponse)
	if err != nil {
		return nil, err
	}
	return &saveResponse, nil
}
