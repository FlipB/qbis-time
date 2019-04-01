package qbis

import (
	"fmt"
	"time"

	"github.com/flipb/qbis-time/pkg/qbis/api"
)

//Week represents a week
//weeks and days are the primary ways to interact with the Qbis api
type Week struct {
	start time.Time
	end   time.Time

	client Client

	sheet   *api.TimesheetData
	changed bool
}

func (w *Week) projectTimeDays() []api.ProjectTime {
	days := make([]api.ProjectTime, 0)

	for i := range w.sheet.ListOfProjectTime {
		activity := w.sheet.ListOfProjectTime[i]
		for day := range activity.Days {

			dayDate, err := api.DateStringToTime(activity.Days[day].DayDate)
			if err != nil {
				// TODO improve
				panic(err)
			}
			activity.Days[day].DayDate = api.TimeToISODateString(dayDate)
		}
		days = append(days, activity)
	}
	return days
}

func (w *Week) salaryTimeDays() []api.SalaryTimeBase {
	days := make([]api.SalaryTimeBase, 0)

	for i := range w.sheet.ListOfSalaryTime {
		activity := w.sheet.ListOfSalaryTime[i].SalaryTimeBase
		for day := range activity.Days {

			dayDate, err := api.DateStringToTime(activity.Days[day].DayDate)
			if err != nil {
				// TODO improve
				panic(err)
			}
			activity.Days[day].DayDate = api.TimeToISODateString(dayDate)
		}
		days = append(days, activity)
	}
	return days
}

func (w *Week) workingTimeDays() []api.WorkingTimeBase {
	var days = make([]api.WorkingTimeBase, 0)

	for i := range w.sheet.WorkingTimeDays {
		x := w.sheet.WorkingTimeDays[i].WorkingTimeBase
		dayDate, err := api.DateStringToTime(x.DayDate)
		if err != nil {
			// TODO improve
			panic(err)
		}
		x.DayDate = api.TimeToISODateString(dayDate)

		for _, b := range x.Breaks {
			breakDate, err := api.DateStringToTime(b.BreakDate)
			if err != nil {
				// TODO improve
				panic(err)
			}
			b.BreakDate = api.TimeToISODateString(breakDate)
		}

		days = append(days, x)
	}
	return days
}

//salaryTime gets a pointer to the activity with the given ID. returns error if not found
func (w *Week) salaryTime(activityID int) (*api.SalaryTime, error) {

	for i, x := range w.sheet.ListOfSalaryTime {
		if x.ActivityID == activityID {
			return &w.sheet.ListOfSalaryTime[i], nil
		}
	}

	// unable to find the activity sheet, it's a new activity for this week. Fetch and add it
	tempSalaryActivity, err := w.client.apiClient.GetSalaryActivity(w.client.employeeID, activityID, w.start, w.end)
	if err != nil {
		return nil, fmt.Errorf("unable to find salary time activity with ActivityID %d: %v", activityID, err)
	}

	w.sheet.ListOfSalaryTime = append(w.sheet.ListOfSalaryTime, *tempSalaryActivity)
	salaryTime := &w.sheet.ListOfSalaryTime[len(w.sheet.ListOfSalaryTime)-1]

	return salaryTime, nil
}

//projectTime gets a pointer to the activity with the given ID. returns error if not found
func (w *Week) projectTime(activityID int) (*api.ProjectTime, error) {

	for i, x := range w.sheet.ListOfProjectTime {
		if x.ActivityID == activityID {
			return &w.sheet.ListOfProjectTime[i], nil
		}
	}

	// project found with the given activity id was not found. Lets fetch and add it
	tempProjectTime, err := w.client.apiClient.GetProjectActivity(w.client.employeeID, activityID, w.start, w.end)
	if err != nil {
		return nil, fmt.Errorf("unable to find project activity with ActivityID %d", activityID)
	}

	// add the project activity to the weeks project activities
	// TODO check if this causes an issue if it's not used
	w.sheet.ListOfProjectTime = append(w.sheet.ListOfProjectTime, *tempProjectTime)
	projectTime := &w.sheet.ListOfProjectTime[len(w.sheet.ListOfProjectTime)-1]

	return projectTime, nil
}

//ErrorSaveWorkingTimeResponse is an error that embeds api.SaveWorkingTimeResponse
//allowing you to get the finer details of the error and warning messages etc
type ErrorSaveWorkingTimeResponse struct {
	*api.SaveWorkingTimeResponse
	message string
}

func (e ErrorSaveWorkingTimeResponse) Error() string {
	return e.message
}

func (w *Week) saveWorkingTime() (*api.SaveWorkingTimeResponse, error) {

	t := api.EmployeeWorkingTime{
		Days:       w.workingTimeDays(),
		EmployeeID: w.client.employeeID,
		FromDate:   api.TimeToISODateString(w.start),
		ToDate:     api.TimeToISODateString(w.end),
	}
	response, err := w.client.apiClient.SaveWorkingTime(t)
	if err != nil {
		return nil, fmt.Errorf("unable to save working time: %v", err)
	}
	if response.Saved != "" {
		// this means it failed to save, i think
		return nil, ErrorSaveWorkingTimeResponse{SaveWorkingTimeResponse: response, message: "failed to save"}
	}

	return response, nil
}

//ErrorSaveSalaryTimeResponse is an error that embeds api.SaveSalaryTimeResponse
//allowing you to get the finer details of the error and warning messages etc
type ErrorSaveSalaryTimeResponse struct {
	*api.SaveSalaryTimeResponse
	message string
}

func (e ErrorSaveSalaryTimeResponse) Error() string {
	return e.message
}

func (w *Week) saveSalaryTime() (*api.SaveSalaryTimeResponse, error) {

	t := api.EmployeeSalaryTime{
		EmployeeID:  w.client.employeeID,
		FromDate:    api.TimeToISODateString(w.start),
		ToDate:      api.TimeToISODateString(w.end),
		SalaryTime:  w.salaryTimeDays(),
		WorkingTime: w.workingTimeDays(),
	}
	response, err := w.client.apiClient.SaveSalaryTime(t)
	if err != nil {
		return nil, fmt.Errorf("unable to save salary time: %v", err)
	}
	if response.WasSaved != true {
		return nil, ErrorSaveSalaryTimeResponse{SaveSalaryTimeResponse: response, message: "error saving"}
	}
	return response, nil
}

//ErrorSaveProjectTimeResponse is an error that embeds api.SaveProjectTimeResponse
//allowing you to get the finer details of the error and warning messages etc
type ErrorSaveProjectTimeResponse struct {
	*api.SaveProjectTimeResponse
	message string
}

func (e ErrorSaveProjectTimeResponse) Error() string {
	return e.message
}

func (w *Week) saveProjectTime() (*api.SaveProjectTimeResponse, error) {

	t := api.EmployeeProjectTime{
		EmployeeID: w.client.employeeID,
		FromDate:   api.TimeToISODateString(w.start),
		ToDate:     api.TimeToISODateString(w.end),
		List:       w.projectTimeDays(),
	}
	response, err := w.client.apiClient.SaveProjectTime(t)
	if err != nil {
		return nil, fmt.Errorf("unable to save project time: %v", err)
	}
	if response.WasSaved != true {
		return nil, ErrorSaveProjectTimeResponse{SaveProjectTimeResponse: response, message: "error saving"}
	}
	return response, nil
}

//Save saves the week
func (w *Week) Save() error {
	if !w.changed {
		return fmt.Errorf("week has not changed (according to 'changed' flag)")
	}

	salRes, err := w.saveSalaryTime()
	if err != nil {
		salErr, ok := err.(ErrorSaveSalaryTimeResponse)
		if ok {
			return fmt.Errorf("error saving working time: %s (wasSaved: %b)", salErr.message, salErr.WasSaved)
		}
		return fmt.Errorf("error saving working time: %v", err)
	}

	workRes, err := w.saveWorkingTime()
	if err != nil {
		workErr, ok := err.(ErrorSaveWorkingTimeResponse)
		if ok {
			return fmt.Errorf("error saving working time: %s (Saved: %s)", workErr.message, workErr.Saved)
		}
		return fmt.Errorf("error saving working time: %v", err)
	}

	projRes, err := w.saveProjectTime()
	if err != nil {
		projErr, ok := err.(ErrorSaveProjectTimeResponse)
		if ok {
			return fmt.Errorf("error saving working time: %s (wasSaved: %b)", projErr.message, projErr.WasSaved)
		}
		return fmt.Errorf("error saving working time: %v", err)
	}

	// data was saved
	fmt.Println("Debug printing warning messages etc.")
	fmt.Printf("Working time response: %+v\n", workRes)
	fmt.Printf("Salary time response: %+v\n", salRes)
	fmt.Printf("Project time response: %+v\n", projRes)

	err = w.Update()
	if err != nil {
		return err
	}
	w.changed = false
	return nil
}

//Update refreshes timesheet data
func (w *Week) Update() error {
	sheet, err := w.client.apiClient.GetTimesheet(w.client.employeeID, w.start, w.end)
	if err != nil {
		return err
	}
	w.sheet = sheet

	return nil
}

// Closed returns true if week is closed in QBis
func (w Week) Closed() bool {
	return w.sheet.SummaryData.WeekStatus == 1
}

// Approved returns true if week is approved by a manager in QBis
func (w Week) Approved() bool {
	return w.sheet.SummaryData.WeekStatus == 2
}

//Open returns true if the week is open
func (w Week) Open() bool {
	return w.sheet.SummaryData.WeekStatus == 0
}

//ScheduledMinutes returns the number of minutes the employee was scheduled for this week
func (w Week) ScheduledMinutes() uint {
	if !w.sheet.SummaryData.HasSchedule {
		return 0 // i just assume this is the right call if you dont have a schedule set
	}

	return uint(w.sheet.SummaryData.ScheduledHours * 60)
}

//WorkedMinutes returns the number of minutes the employee has worked this week.
func (w Week) WorkedMinutes() uint {
	var workedHours float64
	workedHours = w.sheet.SummaryData.WorkedHours

	return uint(workedHours * 60)
}

//OvertimeMinutes either positive, 0 or negative.
func (w Week) OvertimeMinutes() int {
	return int(w.WorkedMinutes() - w.ScheduledMinutes())
}

// DAY

//Day gets the day of the timestamp
func (w *Week) Day(dayOf time.Time) (*Day, error) {
	year, month, day := dayOf.Date()
	var foundDaySettings *api.DaySetting
	foundDaySettings = nil
	foundDayIndex := -1
	for i, ds := range w.sheet.DaySettings {
		t, err := api.DateStringToTime(ds.DayDate)
		if err != nil {
			return nil, fmt.Errorf("unable to parse DayDate string '%s': %v", ds.DayDate, err)
		}
		iyear, imonth, iday := t.Date()
		if iyear == year && imonth == month && day == iday {
			foundDaySettings = &ds
			foundDayIndex = i
			break
		}
	}
	if foundDaySettings == nil {
		return nil, fmt.Errorf("day not found in week")
	}

	var pDay *Day
	pDay = &Day{}
	pDay.indexInWeek = foundDayIndex
	pDay.week = w

	qbisDate, err := api.GetDateForDateTime(dayOf)
	if err != nil {
		return nil, err
	}
	pDay.Date = qbisDate

	return pDay, nil
}

//Weekday returns the day in the week matching the desired weekday
func (w *Week) Weekday(weekDay time.Weekday) (*Day, error) {
	weekdayDatetime, err := getNextWeekday(w.start, w.end, weekDay)
	if err != nil {
		return nil, err
	}
	return w.Day(weekdayDatetime)
}

//getNextWeekday returns the next occurance of the specified dayInWeek within the time span
func getNextWeekday(startDate time.Time, endDate time.Time, dayInWeek time.Weekday) (time.Time, error) {
	dayDate := startDate
	// Remember weekday 0 is sunday, 1 is monday etc.
	for dayDate.Weekday() != dayInWeek {
		dayDate = dayDate.AddDate(0, 0, 1)
		if dayDate.After(endDate) {
			return time.Unix(0, 0), fmt.Errorf("unable to find weekday %d in week")
		}
	}

	return dayDate, nil
}

// SALARY TIME

//SalaryTimeActivities returns a list of all available SalaryTime activities
func (w *Week) SalaryTimeActivities() []SalaryActivity {
	return salaryTimeActivities(w)
}

// PROJECT TIME

//ProjectTimeActivities returns a list of all project activities available to the employee
func (w *Week) ProjectTimeActivities() (*ProjectActivityList, error) {
	return newProjectActivityList(w)
}
