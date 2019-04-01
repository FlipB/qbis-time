package qbis

import (
	"fmt"
	"time"

	"github.com/flipb/qbis-time/pkg/qbis/api"
)

//Day represents a day in a the week. Each week has... wait for it... 7 days! //wizeass
//Interacting with a Day is how you change values in the Qbis sheet.
//Once you have changed what you want to do you can use Week.Save to persist the changes.
//Days interact with and modify the live week struct - if you make changes and then read some values
//you might not get valid data back. Best way to interact is:
//1. Read data
//2. Make changes
//3. Save week
type Day struct {
	week *Week

	indexInWeek int
	Date        time.Time
}

func (d *Day) workingTime() *api.WorkingTime {
	return &d.week.sheet.WorkingTimeDays[d.indexInWeek]
}

func (d *Day) daySetting() *api.DaySetting {
	return &d.week.sheet.DaySettings[d.indexInWeek]
}

func (d *Day) containsTime(t time.Time) bool {
	dy, dm, dd := d.Date.UTC().Date()
	ty, tm, td := t.UTC().Date()
	if ty == dy && tm == dm && td == dd {
		return true
	}
	return false
}

//SetArrival sets time of arrival for the employee
func (d *Day) SetArrival(time time.Time) error {
	if !d.containsTime(time) {
		fmt.Errorf("date of day does not match date of supplied time: got %v, expected %v", time, d.Date)
	}
	d.workingTime().Arrive = time.Hour()*60 + time.Minute()
	d.workingTime().IsModified = true
	d.week.changed = true

	return nil
}

//SetDeparture sets time of departure for the employee
func (d *Day) SetDeparture(time time.Time) error {
	if !d.containsTime(time) {
		fmt.Errorf("date of day does not match date of supplied time: got %v, expected %v", time, d.Date)
	}
	d.workingTime().Leave = time.Hour()*60 + time.Minute()
	d.workingTime().IsModified = true
	d.week.changed = true
	return nil
}

//SetBreakMinutes sets the number of minutes the employee has been on lunch break
func (d *Day) SetBreakMinutes(minutes uint) {
	d.workingTime().Lunch = int(minutes)
	d.workingTime().IsModified = true
	d.week.changed = true
}

//ScheduledMinutes returns the number of minutes the employer thinks the employee is supposed to work
func (d *Day) ScheduledMinutes() uint {
	if !d.daySetting().HasSchedule {
		return 0
	}
	return uint(d.daySetting().MySchedule.TotalMinutes)
}

//LoggedMinutes represents the number of minutes that the employee says s/he has worked
func (d *Day) LoggedMinutes() uint {
	if !d.workingTime().IsModified {
		return uint(d.workingTime().Total)
	}

	// week has been modified, Total is not guaranteed to be up to date
	// we have to calculate it

	// check if it makes sense.
	if (d.workingTime().Leave - d.workingTime().Arrive - d.workingTime().Lunch) < 0 {
		panic(fmt.Errorf("Logged minutes calculation less than 0"))
	}

	return uint(d.workingTime().Leave - d.workingTime().Arrive - d.workingTime().Lunch)
}

//Holiday returns true if the day is a holiday
func (d *Day) Holiday() bool {
	return d.daySetting().IsHoliday
}

// SALARY TIME

//LoggedSalaryTimeActivities returns the "used" salarytime activities on that day
func (d *Day) LoggedSalaryTimeActivities() []SalaryActivity {
	var dayActivities = make([]SalaryActivity, 0)
	for _, x := range d.week.sheet.ListOfSalaryTime {
		if x.Days[d.indexInWeek].DayMinutes != 0 || x.Days[d.indexInWeek].DayDays != 0 || x.IsDefault {
			dayActivities = append(dayActivities, SalaryActivity{
				week:      d.week,
				key:       x.ActivityID,
				value:     x.ActivityName,
				isDefault: x.IsDefault,
			})
		}
	}
	return dayActivities
}

//SalaryTimeMinutes returns the number of minutes registered on the given salary activity that day
func (d *Day) SalaryTimeMinutes(activityID int) int {
	salaryTime, err := d.week.salaryTime(activityID)
	if err != nil {
		return 0
	}
	return salaryTime.Days[d.indexInWeek].DayMinutes
}

//SetSalaryTime sets the number of minutes spent on the activity that day
func (d *Day) SetSalaryTime(activityID int, minutes int) error {
	salaryTime, err := d.week.salaryTime(activityID)
	if err != nil {
		return fmt.Errorf("error getting activity with id %d : %v", activityID, err)
	}

	// attempt to do some validation
	if minutes < 0 && !salaryTime.AllowNegative {
		return fmt.Errorf("salaryTime activity %s (%d) does not allow negative minutes: %d", salaryTime.ActivityName, activityID, minutes)
	}
	if minutes > 0 && !salaryTime.AllowPositive {
		return fmt.Errorf("salaryTime activity %s (%d) does not allow positive minutes: %d", salaryTime.ActivityName, activityID, minutes)
	}

	salaryTime.Days[d.indexInWeek].DayMinutes = minutes
	d.week.changed = true
	return nil
}

// PROJECT TIME

//LoggedProjectTimeActivities returns the "used" project time activities on that day
func (d *Day) LoggedProjectTimeActivities() []ProjectActivity {
	var dayActivities = make([]ProjectActivity, 0)
	for _, x := range d.week.sheet.ListOfProjectTime {
		if x.Days[d.indexInWeek].DayMinutes != 0 {
			dayActivities = append(dayActivities, ProjectActivity{
				week: d.week,
				id:   x.ActivityID,
				name: x.ActivityName,
			})
		}
	}
	return dayActivities
}

//ProjectTimeMinutes returns the number of minutes registered on the given project activity that day
func (d *Day) ProjectTimeMinutes(activityID int) int {
	projectTime, err := d.week.projectTime(activityID)
	if err != nil {
		return 0
	}
	return projectTime.Days[d.indexInWeek].DayMinutes
}

//ProjectTimeInternalNotes returns internal notes on the project activity the given day
func (d *Day) ProjectTimeInternalNotes(activityID int) string {
	projectTime, err := d.week.projectTime(activityID)
	if err != nil {
		return ""
	}
	return projectTime.Days[d.indexInWeek].InternalNotes
}

//ProjectTimeExternalNotes returns external notes on the project activity the given day
func (d *Day) ProjectTimeExternalNotes(activityID int) string {
	projectTime, err := d.week.projectTime(activityID)
	if err != nil {
		return ""
	}
	return projectTime.Days[d.indexInWeek].ExternalNotes
}

//SetProjectTime sets the number of minutes spent on the activity that day
func (d *Day) SetProjectTime(activityID int, minutes int) error {
	projectTime, err := d.week.projectTime(activityID)
	if err != nil {
		return fmt.Errorf("error getting activity with id %d : %v", activityID, err)
	}

	projectTime.Days[d.indexInWeek].DayMinutes = minutes
	d.week.changed = true
	return nil
}

//SetProjectTimeInternalNote sets the internal note on the activity that day
func (d *Day) SetProjectTimeInternalNote(activityID int, note string) error {
	projectTime, err := d.week.projectTime(activityID)
	if err != nil {
		return fmt.Errorf("error getting activity with id %d : %v", activityID, err)
	}

	projectTime.Days[d.indexInWeek].InternalNotes = note
	d.week.changed = true
	return nil
}

//SetProjectTimeExternalNote sets the external note on the activity that day
func (d *Day) SetProjectTimeExternalNote(activityID int, note string) error {
	projectTime, err := d.week.projectTime(activityID)
	if err != nil {
		return fmt.Errorf("error getting activity with id %d : %v", activityID, err)
	}

	projectTime.Days[d.indexInWeek].ExternalNotes = note
	d.week.changed = true
	return nil
}
