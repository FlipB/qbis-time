package qbis

//SalaryActivity represents a Salary Activity (Sick leave, vacation etc.)
type SalaryActivity struct {
	week  *Week
	key   int    // is a number (activity id)
	value string // is the name of the activity ("Övertid x1.5")

	isDefault bool // this is the default salary activity
	//allowPositive bool // this salary activity allows positive values
	//allowNegative bool // this salary activity allows negative values

	//typeCode int	// type of activity. Defines how time is calculated
}

//salaryTimeActivities returns a list of all available SalaryTime activities
func salaryTimeActivities(week *Week) []SalaryActivity {
	var stActivities = make([]SalaryActivity, 0)

	// get default first (special case)
	foundDefault := false
	for _, x := range week.sheet.ListOfSalaryTime {
		if x.IsDefault && !foundDefault {
			foundDefault = true
			stActivities = append(stActivities, SalaryActivity{
				week:      week,
				key:       x.ActivityID,
				value:     x.ActivityName,
				isDefault: true,
				//allowNegative: x.AllowNegative,
				//allowPositive: x.AllowPositive,
			})
			continue
		}
		if x.IsDefault {
			panic("more than one default SalaryActivity possible??")
		}
	}

	for _, y := range week.sheet.ListOfSalaryActivities {
		stActivities = append(stActivities, SalaryActivity{
			week:  week,
			key:   y.Key,
			value: y.Value,
		})
	}
	return stActivities
}

//ActivityID returns the ActivityID of the salary activity
func (s SalaryActivity) ActivityID() int {
	return s.key
}

//Name returns the name of the salary activity
func (s SalaryActivity) Name() string {
	// s.week.sheet.ListOfSalaryActivities does not include "Komp-tid", probably because it's the default
	return s.value
}

//InWeek returns true if the salary activity is added to the week timesheet
func (s SalaryActivity) InWeek() bool {
	for _, x := range s.week.sheet.ListOfSalaryTime {
		if x.ActivityID == s.key {
			return true
		}
	}
	return false
}
