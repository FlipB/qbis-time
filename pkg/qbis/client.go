//Package qbis provides a higher lever QBis client implementation
//as well as types for interacting with the API
package qbis

import (
	"fmt"
	"time"

	"github.com/flipb/qbis-time/pkg/qbis/api"
)

//Client is a highlevel qbis client
type Client struct {
	apiClient  api.Client
	employeeID string
}

//NewClient creates a new qbis client
func NewClient(qbisCompany string, qbisUser string, qbisPassword string) (*Client, error) {
	qbisClient := api.New("https://login.qbis.se")

	err := qbisClient.Login(qbisCompany, qbisUser, qbisPassword)
	if err != nil {
		return nil, fmt.Errorf("login to qbis failed: %v", err)
	}

	return NewClientFromAPIClient(*qbisClient)
}

//NewClientFromAPIClient creates a new qbis client
func NewClientFromAPIClient(client api.Client) (*Client, error) {

	employee, err := client.GetUserEmployeeID()
	if err != nil {
		return nil, err
	}

	return &Client{client, employee}, nil
}

//Week returns the week containing a specific point in time
func (q Client) Week(time time.Time) (*Week, error) {
	w := Week{}
	s, e, err := getWeekSpan(time)
	if err != nil {
		return nil, err
	}
	w.start = s
	w.end = e
	w.client = q

	err = w.Update()
	if err != nil {
		return nil, err
	}

	return &w, nil
}

//WeekNow returns the current week
func (q Client) WeekNow() (*Week, error) {
	return q.Week(time.Now())
}
