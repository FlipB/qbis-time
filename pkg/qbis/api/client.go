package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// Client implements a Qbis client
type Client struct {
	client *http.Client
	url    string
}

type authFormdata struct {
	Authenticate    string // Log+in
	Company         string // companynamehere
	CurrentLanguage string // lang_english
	Password        string // imnottellingyouthat
	RememberMe      bool
	Username        string // usernamegoeshere
}

type cookieData struct {
	SessionID string // lotsofrandomcharacters
}

// New constructs a new Client
func New(url string) *Client {
	client := new(Client)
	client.client = http.DefaultClient

	if !strings.HasSuffix(url, "/") {
		url = url + "/"
	}
	client.url = url

	// after url is set
	client.loadCookies()

	return client
}

// WithHTTPClient sets the http.Client to be used by the qbis client
func (c *Client) WithHTTPClient(client *http.Client) *Client {
	c.client = client
	return c
}

//loadCookies requires Client.url to be set
func (c Client) loadCookies() error {
	var err error
	if c.client.Jar == nil {
		c.client.Jar, err = cookiejar.New(&cookiejar.Options{})
		if err != nil {
			return err
		}
	}
	_, err = c.get("/Login/Login")
	if err != nil {
		return err
	}
	return nil
}

//GetUserEmployeeID returns the string representation of the logged in users employee id
func (c Client) GetUserEmployeeID() (string, error) {
	res, err := c.get("/Time/TimeOverview")
	if err != nil {
		return "", err
	}

	scripts, err := getEmbeddedScriptsInHTML(res.Body)
	if err != nil {
		return "", err
	}
	if len(scripts) == 0 {
		return "", fmt.Errorf("error loading scripts from TimeOverview page")
	}

	employeeID := ""
	for _, script := range scripts {
		str, err := getCurrentUserFromCurrentLoginBlock(script)
		if err != nil {
			//println(err.Error())
			continue
		}

		employeeID = str
	}
	if employeeID == "" {
		return "", fmt.Errorf("unable to find EmployeeID in scripts")
	}

	return employeeID, nil
}

func (c Client) printCookies() error {
	u, err := url.Parse(c.url)
	if err != nil {
		return err
	}
	cookies := c.client.Jar.Cookies(u)
	fmt.Printf("%+v\n", cookies)
	return nil
}

func (c Client) get(resource string) (*http.Response, error) {
	if strings.HasPrefix(resource, "/") {
		resource = strings.TrimLeft(resource, "/")
	}

	resp, err := c.client.Get(c.url + resource)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c Client) postData(resource string, data url.Values) (*http.Response, error) {
	if strings.HasPrefix(resource, "/") {
		resource = strings.TrimLeft(resource, "/")
	}

	resp, err := c.client.PostForm(c.url+resource, data)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c Client) postJSON(resource string, data io.Reader) (*http.Response, error) {
	if strings.HasPrefix(resource, "/") {
		resource = strings.TrimLeft(resource, "/")
	}

	resp, err := c.client.Post(c.url+resource, "application/json", data)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// GetTimesheet returns a timesheet for an employee spanning the two dates
func (c *Client) GetTimesheet(employee string, from time.Time, to time.Time) (*TimesheetData, error) {
	url := "/Time/Timesheet/GetTimeSheetData?employeeId=%s&fromDate=%s&toDate=%s&_=%s"
	fromString := from.UTC().Format("2006-01-02T15:04:05.000Z")
	toString := to.UTC().Format("2006-01-02T15:04:05.000Z")
	timestamp := strconv.Itoa(int(time.Now().UnixNano()))

	url = fmt.Sprintf(url, employee, fromString, toString, timestamp)
	response, err := c.get(url)
	if err != nil {
		return nil, err
	}
	timesheet := &TimesheetData{}

	err = json.NewDecoder(response.Body).Decode(timesheet)
	if err != nil {
		return nil, err
	}
	return timesheet, nil
}

// Login to Qbis
func (c *Client) Login(company string, user string, password string) error {
	urlvals := make(url.Values)
	urlvals.Set("Authenticate", "Log+in")
	urlvals.Set("Company", company)
	urlvals.Set("CurrentLanguage", "lang_english")
	urlvals.Set("Password", password)
	urlvals.Set("RememberMe", "true")
	urlvals.Set("Username", user)
	response, err := c.postData("/Login/Login/Authenticate", urlvals)
	if err != nil {
		return err
	}
	if response.StatusCode < 200 || response.StatusCode > 299 {
		return errors.New("unexpected status code")
	}

	// TODO implement check to see if login was successful.
	// currently no error is returned if you enter invalid credentials because you still get a 200 OK
	/*
			u, err := response.Location()
			if err != nil {
				//TODO handle
				println(err.Error())
			}
			println(u.Path)
		// https://stackoverflow.com/questions/23297520/how-can-i-make-the-go-http-client-not-follow-redirects-automatically
	*/

	return nil
}
