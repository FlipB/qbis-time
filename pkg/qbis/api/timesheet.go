package api

// TimesheetData contains data about the week
type TimesheetData struct {
	ActionButtonOptionslist []struct {
		Action           string `json:"Action"`
		Active           bool   `json:"Active"`
		Bold             bool   `json:"Bold"`
		Disabled         bool   `json:"Disabled"`
		GetAsJSONString  string `json:"GetAsJsonString"`
		Icon             string `json:"Icon"`
		Key              string `json:"Key"`
		LineAfter        bool   `json:"LineAfter"`
		MobileRestricted bool   `json:"MobileRestricted"`
		Parameters       string `json:"Parameters"`
		Title            string `json:"Title"`
	} `json:"actionButtonOptionslist"`
	DaySettings            []DaySetting  `json:"daySettings"`
	ListOfProjectTime      []ProjectTime `json:"listOfProjectTime"`
	ListOfSalaryActivities []struct {
		Key   int    `json:"Key"`
		Value string `json:"Value"`
	} `json:"listOfSalaryActivities"`
	ListOfSalaryTime   []SalaryTime  `json:"listOfSalaryTime"`
	PublicHolidaysList []interface{} `json:"publicHolidaysList"`
	SummaryData        struct {
		BillableTime struct {
			ChargeableHours         int         `json:"ChargeableHours"`
			ChartChargeableHours    int         `json:"ChartChargeableHours"`
			ChartNonChargeableHours int         `json:"ChartNonChargeableHours"`
			ChartPercentage         int         `json:"ChartPercentage"`
			ChartTargetHours        int         `json:"ChartTargetHours"`
			Name                    interface{} `json:"Name"`
			PotentialHours          int         `json:"PotentialHours"`
			TargetHours             int         `json:"TargetHours"`
			TargetPercentage        interface{} `json:"TargetPercentage"`
		} `json:"BillableTime"`
		HasSchedule       bool    `json:"HasSchedule"`
		ScheduledHours    float64 `json:"ScheduledHours"`
		ScheduledTime     string  `json:"ScheduledTime"`
		ShowBillableChart bool    `json:"ShowBillableChart"`
		ShowProgressChart bool    `json:"ShowProgressChart"`
		ShowWorkedHours   bool    `json:"ShowWorkedHours"`
		WeekStatus        int     `json:"WeekStatus"`
		WorkedHours       float64 `json:"WorkedHours"`
		WorkedTime        string  `json:"WorkedTime"`
	} `json:"summaryData"`
	TimeSettings struct {
		AllowTimeManagers                     bool `json:"allowTimeManagers"`
		DepManagerID                          int  `json:"depManagerId"`
		IgnoreProjectActivityFactorValidation bool `json:"ignoreProjectActivityFactorValidation"`
		IsActive                              bool `json:"isActive"`
		IsShowSchedule                        bool `json:"isShowSchedule"`
		IsUsingHrSchedules                    bool `json:"isUsingHrSchedules"`
		ManagerID                             int  `json:"managerId"`
		ModulePermission                      int  `json:"modulePermission"`
		ProjectModuleAccess                   bool `json:"projectModuleAccess"`
		ProjectactivityAccess                 bool `json:"projectactivityAccess"`
		SalarytimeAccess                      bool `json:"salarytimeAccess"`
		ShowActionLocation                    bool `json:"showActionLocation"`
		WorkingtimeAccess                     bool `json:"workingtimeAccess"`
	} `json:"timeSettings"`
	WeekHistoryList []struct {
		ChangeDateString string `json:"ChangeDateString"`
		ChangedBy        string `json:"ChangedBy"`
		ChangedDate      string `json:"ChangedDate"`
		Message          string `json:"Message"`
		Status           int    `json:"Status"`
	} `json:"weekHistoryList"`
	WorkingTimeBreakList []WorkingTimeBreak `json:"workingTimeBreakList"`
	WorkingTimeDays      []WorkingTime      `json:"workingTimeDays"`
}

// TimesheetDataOld contains data about the week
type TimesheetDataOld struct {
	ActionButtonOptionslist []struct {
		Action           string `json:"Action"`
		Active           bool   `json:"Active"`
		Bold             bool   `json:"Bold"`
		Disabled         bool   `json:"Disabled"`
		GetAsJSONString  string `json:"GetAsJsonString"`
		Icon             string `json:"Icon"`
		LineAfter        bool   `json:"LineAfter"`
		MobileRestricted bool   `json:"MobileRestricted"`
		Parameters       string `json:"Parameters"`
		Title            string `json:"Title"`
	} `json:"actionButtonOptionslist"`
	DaySettings            []DaySetting  `json:"daySettings"`
	ListOfProjectTime      []ProjectTime `json:"listOfProjectTime"`
	ListOfSalaryActivities []struct {
		Key   int    `json:"Key"`
		Value string `json:"Value"`
	} `json:"listOfSalaryActivities"`
	ListOfSalaryTime   []SalaryTime  `json:"listOfSalaryTime"`
	PublicHolidaysList []interface{} `json:"publicHolidaysList"`
	SummaryData        struct {
		BillableTime struct {
			ChargeableHours         int         `json:"ChargeableHours"`
			ChartChargeableHours    int         `json:"ChartChargeableHours"`
			ChartNonChargeableHours int         `json:"ChartNonChargeableHours"`
			ChartPercentage         int         `json:"ChartPercentage"`
			ChartTargetHours        int         `json:"ChartTargetHours"`
			Name                    interface{} `json:"Name"`
			PotentialHours          int         `json:"PotentialHours"`
			TargetHours             int         `json:"TargetHours"`
			TargetPercentage        interface{} `json:"TargetPercentage"`
		} `json:"BillableTime"`
		HasSchedule       bool    `json:"HasSchedule"`
		ScheduledHours    float64 `json:"ScheduledHours"`
		ScheduledTime     string  `json:"ScheduledTime"`
		ShowBillableChart bool    `json:"ShowBillableChart"`
		ShowProgressChart bool    `json:"ShowProgressChart"`
		ShowWorkedHours   bool    `json:"ShowWorkedHours"`
		WeekStatus        int     `json:"WeekStatus"`
		WorkedHours       float64 `json:"WorkedHours"`
		WorkedTime        string  `json:"WorkedTime"`
	} `json:"summaryData"`
	TimeSettings struct {
		AllowTimeManagers     bool `json:"allowTimeManagers"`
		DepManagerID          int  `json:"depManagerId"`
		IsActive              bool `json:"isActive"`
		IsShowSchedule        bool `json:"isShowSchedule"`
		IsUsingHrSchedules    bool `json:"isUsingHrSchedules"`
		ManagerID             int  `json:"managerId"`
		ModulePermission      int  `json:"modulePermission"`
		ProjectModuleAccess   bool `json:"projectModuleAccess"`
		ProjectactivityAccess bool `json:"projectactivityAccess"`
		SalarytimeAccess      bool `json:"salarytimeAccess"`
		ShowActionLocation    bool `json:"showActionLocation"`
		WorkingtimeAccess     bool `json:"workingtimeAccess"`
	} `json:"timeSettings"`
	WeekHistoryList []struct {
		ChangeDateString string `json:"ChangeDateString"`
		ChangedBy        string `json:"ChangedBy"`
		ChangedDate      string `json:"ChangedDate"`
		Message          string `json:"Message"`
		Status           int    `json:"Status"`
	} `json:"weekHistoryList"`
	WorkingTimeDays []WorkingTime `json:"workingTimeDays"`
}

//ScheduleDay contains schedule information for a day
type ScheduleDay struct {
	Activities               []interface{} `json:"Activities"`
	Arrive                   string        `json:"Arrive"`
	DayDate                  string        `json:"DayDate"`
	HasScheduleArriveOrLeave bool          `json:"HasScheduleArriveOrLeave"`
	Leave                    string        `json:"Leave"`
	LunchFrom                string        `json:"LunchFrom"`
	LunchMinutes             int           `json:"LunchMinutes"`
	LunchTo                  string        `json:"LunchTo"`
	OverridePublicHolidays   bool          `json:"OverridePublicHolidays"`
	PrefillTime              bool          `json:"PrefillTime"`
	ScheduleEmployeeFrom     string        `json:"ScheduleEmployeeFrom"`
	ScheduleEmployeeTo       string        `json:"ScheduleEmployeeTo"`
	SkipDeviationValidation  bool          `json:"SkipDeviationValidation"`
	TotalMinutes             int           `json:"TotalMinutes"`
	Tracking                 int           `json:"Tracking"`
}

//DaySetting is included in TimesheetData for every day of the week
type DaySetting struct {
	DayComment                            interface{} `json:"DayComment"`
	DayDate                               string      `json:"DayDate"`
	DayNameString                         string      `json:"DayNameString"`
	DisabledTooltipProjectTime            string      `json:"DisabledTooltipProjectTime"`
	DisabledTooltipSalaryTime             string      `json:"DisabledTooltipSalaryTime"`
	DisabledTooltipWorkingTime            string      `json:"DisabledTooltipWorkingTime"`
	HasDayComment                         bool        `json:"HasDayComment"`
	HasLunchPolicy                        bool        `json:"HasLunchPolicy"`
	HasSchedule                           bool        `json:"HasSchedule"`
	HideDay                               bool        `json:"HideDay"`
	IsAllowedToRegisterWorkingTimeFromWeb bool        `json:"IsAllowedToRegisterWorkingTimeFromWeb"`
	IsDisabledProjectTime                 bool        `json:"IsDisabledProjectTime"`
	IsDisabledSalaryTime                  bool        `json:"IsDisabledSalaryTime"`
	IsDisabledWorkingTime                 bool        `json:"IsDisabledWorkingTime"`
	IsEmployeeInactive                    bool        `json:"IsEmployeeInactive"`
	IsHoliday                             bool        `json:"IsHoliday"`
	IsMonthClosedProjectTime              bool        `json:"IsMonthClosedProjectTime"`
	IsMonthClosedWorkingTime              bool        `json:"IsMonthClosedWorkingTime"`
	IsOutsideEmploymentPeriod             bool        `json:"IsOutsideEmploymentPeriod"`
	IsReadOnlyProjectTime                 bool        `json:"IsReadOnlyProjectTime"`
	IsReadOnlySalaryTime                  bool        `json:"IsReadOnlySalaryTime"`
	IsReadOnlyWorkingTime                 bool        `json:"IsReadOnlyWorkingTime"`
	IsSaved                               bool        `json:"IsSaved"`
	IsScheduleHourOnly                    bool        `json:"IsScheduleHourOnly"`
	IsToday                               bool        `json:"IsToday"`
	IsWeekSaved                           bool        `json:"IsWeekSaved"`
	IsWorkingDay                          bool        `json:"IsWorkingDay"`
	LunchMaximum                          int         `json:"LunchMaximum"`
	LunchMinimum                          int         `json:"LunchMinimum"`
	MonthName                             string      `json:"MonthName"`
	MySchedule                            ScheduleDay `json:"MySchedule"`
	WeekStatus                            int         `json:"WeekStatus"`
}
