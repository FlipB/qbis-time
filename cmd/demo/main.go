package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/flipb/qbis-time/pkg/qbis"
)

//Config holds QBis configuration
type Config struct {
	URL      string
	Company  string
	User     string
	Password string
}

func main() {

	data, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Fatalf("unable to read config file: %v", err)
	}
	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		log.Fatalf("unable to parse config file: %v", err)
	}

	q, err := qbis.NewClient(config.Company, config.User, config.Password)
	if err != nil {
		log.Fatal(err)
	}

	// Get the next week's timesheet
	w, err := q.Week(time.Now().AddDate(0, 0, 7))
	if err != nil {
		log.Fatal(err)
	}

	d, err := w.Weekday(time.Monday)
	if err != nil {
		log.Fatal(err)
	}

	mondayNoon := d.Date.Add(time.Hour * 12)
	err = d.SetArrival(mondayNoon)
	if err != nil {
		log.Fatalf("error setting arrival time: %v\n", err)
	}

	s := w.SalaryTimeActivities()
	for _, s := range s {
		fmt.Printf("SalaryActivity ID %d %s - %v\n", s.ActivityID(), s.Name(), s.InWeek())
	}

	p, err := w.ProjectTimeActivities()
	if err != nil {
		log.Fatal(err)
	}
	projects, err := p.Projects()
	if err != nil {
		log.Fatal(err)
	}
	for _, project := range projects {
		fmt.Printf("Project %s (ID: %d) Code: %s\n", project.Name(), project.ID(), project.Code())
		activities, err := project.Activities()
		if err != nil {
			log.Fatal(err)
		}
		for _, activity := range activities {
			fmt.Printf("\t%s (ID: %d)\n", activity.Name(), activity.ActivityID())
		}
	}

	salt := d.LoggedSalaryTimeActivities()
	for _, s := range salt {
		fmt.Printf("SalaryTime Activity: %s, %dm\n", s.Name(), d.SalaryTimeMinutes(s.ActivityID()))
	}

	pts := d.LoggedProjectTimeActivities()
	for _, pt := range pts {
		fmt.Printf("ProjectTime Activity: %s, %dm : %s\n", pt.Name(), d.ProjectTimeMinutes(pt.ActivityID()), d.ProjectTimeInternalNotes(pt.ActivityID()))
	}

	err = w.Save()
	if err != nil {
		log.Fatalf("error saving week: %v\n", err)
	}

	pts = d.LoggedProjectTimeActivities()
	for _, pt := range pts {
		fmt.Printf("ProjectTime Activity: %s, %dm : %s\n", pt.Name(), d.ProjectTimeMinutes(pt.ActivityID()), d.ProjectTimeInternalNotes(pt.ActivityID()))
	}

	// This is the activity ID found by calling week.SalaryTimeActivities()
	sickActivityID := 10
	err = d.SetSalaryTime(sickActivityID, -90)
	if err != nil {
		log.Fatal(err)
	}

	err = w.Save()
	if err != nil {
		log.Fatal(err.Error())
	}
	salt = d.LoggedSalaryTimeActivities()
	for _, s := range salt {
		fmt.Printf("SalaryTime Activity: %s, %dm\n", s.Name(), d.SalaryTimeMinutes(s.ActivityID()))
	}

	// RESET WORKING TIME BELOW
	err = d.SetArrival(d.Date.Add(time.Hour * 8))
	if err != nil {
		println(err.Error())
	}
	err = d.SetDeparture(d.Date.Add(time.Hour * 17))
	if err != nil {
		println(err.Error())
	}

	err = w.Save()
	if err != nil {
		println(err.Error())
	}

	fmt.Printf("Time spent working = %d (want: 480)", d.LoggedMinutes())
}
