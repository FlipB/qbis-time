package qbis

import "github.com/flipb/qbis-time/pkg/qbis/api"

//ProjectActivityList contains all companies, projects and activities available to the employee
type ProjectActivityList struct {
	week *Week

	companies []api.ProjectCompany
}

//newProjectActivityList initialized a new ProjectActivityList for the given week
func newProjectActivityList(week *Week) (*ProjectActivityList, error) {
	companies, err := week.client.apiClient.GetProjects(week.client.employeeID, week.start, week.end)
	if err != nil {
		return nil, err
	}

	return &ProjectActivityList{
		week:      week,
		companies: companies,
	}, nil
}

//Companies returns the list of companies with projects
func (pal *ProjectActivityList) Companies() []ProjectCompany {
	var pcs = make([]ProjectCompany, 0)
	for i := range pal.companies {
		pcs = append(pcs, ProjectCompany{
			list:    pal,
			company: &pal.companies[i],
		})
	}

	return pcs

}

//Projects gets a list of all projects for all companies
func (pal *ProjectActivityList) Projects() ([]Project, error) {
	var projects = make([]Project, 0)
	for _, comp := range pal.Companies() {

		projs, err := comp.Projects()
		if err != nil {
			return nil, err
		}

		for _, proj := range projs {
			projects = append(projects, proj)
		}
	}
	return projects, nil
}

//Activities returns a list of all activities for all projects and all companies
func (pal ProjectActivityList) Activities() ([]ProjectActivity, error) {
	var activities = make([]ProjectActivity, 0)

	projects, err := pal.Projects()
	if err != nil {
		return nil, err
	}
	for _, x := range projects {
		acts, err := x.Activities()
		if err != nil {
			return nil, err
		}
		activities = append(activities, acts...)
	}

	return activities, nil
}

//ProjectCompany has projects, projects have activities.
//You can access all activities and projects for a company through ProjectCompany
type ProjectCompany struct {
	list    *ProjectActivityList
	company *api.ProjectCompany
}

//Name gets the name of the company
func (pc ProjectCompany) Name() string {
	return pc.company.CompanyName
}

//CompanyID gets the company ID in qbis for the company
func (pc ProjectCompany) CompanyID() int {
	return pc.company.CompanyID
}

//Projects gets a list of all projects for the company
func (pc *ProjectCompany) Projects() ([]Project, error) {
	var projects = make([]Project, 0)
	for i, p := range pc.company.Projects {

		// get activities
		list, err := pc.list.week.client.apiClient.GetProjectActivityList(pc.list.week.client.employeeID, p.ID, pc.list.week.start, pc.list.week.end)
		if err != nil {
			return nil, err
		}

		projects = append(projects, Project{
			company:    pc,
			project:    &pc.company.Projects[i],
			activities: list,
		})
	}
	return projects, nil
}

//Activities returns a list of all activities of all projects for the company
func (pc *ProjectCompany) Activities() ([]ProjectActivity, error) {
	activities := make([]ProjectActivity, 0)
	companyProjects, err := pc.Projects()
	if err != nil {
		return nil, err
	}
	for i := range companyProjects {
		projectActivities, err := companyProjects[i].Activities()
		if err != nil {
			return nil, err
		}
		activities = append(activities, projectActivities...)
	}
	return activities, nil
}

//Project has ProjectActivity(s)
type Project struct {
	company    *ProjectCompany
	project    *api.Project
	activities []api.ProjectActivityListItem
}

//Name gets the name of the project
func (p Project) Name() string {
	return p.project.Name
}

//ID gets the project ID of the project
// I think these are different from Activity IDs.
func (p Project) ID() int {
	return p.project.ID
}

//Code returns the project qbis code
//This seems to be a voluntary thing
func (p Project) Code() string {
	return p.project.Code
}

//Activities returns a list of all activities included in the project
func (p *Project) Activities() ([]ProjectActivity, error) {
	var activities = make([]ProjectActivity, 0)
	for _, a := range p.activities {
		activities = append(activities, ProjectActivity{
			//project: p,
			id:   a.ID,
			name: a.Name,
		})
	}
	return activities, nil
}

//ProjectActivity represents an activity in a project
//Currently you cant do much fun here, just get the ID and name
type ProjectActivity struct {
	week *Week // but we do have a week
	id   int
	name string
}

//ActivityID returns the ActivityID of the project activity
func (p ProjectActivity) ActivityID() int {
	return p.id
}

//Name returns the name of the project activity
func (p ProjectActivity) Name() string {
	return p.name
}

//InWeek returns true if the project activity is added to the week timesheet
func (p ProjectActivity) InWeek() bool {
	for _, x := range p.week.sheet.ListOfProjectTime {
		if x.ActivityID == p.id {
			return true
		}
	}
	return false
}
