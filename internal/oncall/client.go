package oncall

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/buger/jsonparser"
	"io"
	"net/http"
	"oncall-gen/internal/model"
	"strings"
)

type Client interface {
	Login(username string, password string) (csrfToken string, cookie *http.Cookie, err error)
	CreateTeam(csrfToken string, cookie *http.Cookie, team *model.Team) error
	UpdateTeam(csrfToken string, cookie *http.Cookie, team *model.Team) error
	CreateUser(csrfToken string, cookie *http.Cookie, user *model.User) error
	UpdateUser(csrfToken string, cookie *http.Cookie, user *model.User) error
	CreateTeamUser(csrfToken string, cookie *http.Cookie, teamName string, userName string) error
	CreateRoster(csrfToken string, cookie *http.Cookie, teamName string, rosterName string) error
	CreateRosterUser(csrfToken string, cookie *http.Cookie, teamName string, rosterName string, userName string) error
	ExistEvent(csrfToken string, cookie *http.Cookie, teamName string, userName string, role string, start int64) (bool, error)
	CreateEvent(csrfToken string, cookie *http.Cookie, event *model.Event) error
}

type ClientImpl struct {
	url        string
	httpClient *http.Client
}

func NewClientImpl(url string) *ClientImpl {
	return &ClientImpl{
		url:        url,
		httpClient: http.DefaultClient,
	}
}

var ErrTeamExists = errors.New("team already exist")

var ErrUserExists = errors.New("user already exist")
var ErrRosterExists = errors.New("roster already exist")
var ErrUserNotFound = errors.New("user not found")

func (c ClientImpl) Login(username string, password string) (csrfToken string, cookie *http.Cookie, err error) {
	data := fmt.Sprintf("username=%s&password=%s", username, password)
	response, err := c.httpClient.Post(
		c.url+"/login",
		"application/x-www-form-urlencoded",
		strings.NewReader(data),
	)
	if err != nil {
		return "", nil, err
	}

	if response.StatusCode == http.StatusOK {
		for _, ck := range response.Cookies() {
			if ck.Name == "oncall-auth" {
				cookie = ck
			}
		}
		var body []byte
		body, err = io.ReadAll(response.Body)
		if err != nil {
			return "", nil, err
		}
		var csrfTokenBytes []byte
		csrfTokenBytes, _, _, err = jsonparser.Get(body, "csrf_token")
		if err != nil {
			return "", nil, err
		}
		csrfToken = string(csrfTokenBytes)
	}
	return
}

func (c ClientImpl) CreateTeam(csrfToken string, cookie *http.Cookie, team *model.Team) error {
	data, err := json.Marshal(team)
	if err != nil {
		return err
	}

	request, err := http.NewRequest("POST", c.url+"/api/v0/teams", strings.NewReader(string(data)))
	if err != nil {
		return err
	}
	request.AddCookie(cookie)
	request.Header.Set("x-csrf-token", csrfToken)

	response, err := c.httpClient.Do(request)
	if err != nil {
		return err
	}

	if response.StatusCode == http.StatusCreated {
		return nil
	} else if response.StatusCode == http.StatusBadRequest {
		return fmt.Errorf("bad request: %v", response)
	} else if response.StatusCode == http.StatusUnprocessableEntity {
		return ErrTeamExists
	}
	return fmt.Errorf("unexpected code %d: %v", response.StatusCode, response)
}

func (c ClientImpl) UpdateTeam(csrfToken string, cookie *http.Cookie, team *model.Team) error {
	data, err := json.Marshal(team)
	if err != nil {
		return err
	}

	request, err := http.NewRequest("PUT", c.url+"/api/v0/teams/"+team.Name, strings.NewReader(string(data)))
	if err != nil {
		return err
	}
	request.AddCookie(cookie)
	request.Header.Set("x-csrf-token", csrfToken)

	fmt.Println(string(data))
	fmt.Println(request)

	response, err := c.httpClient.Do(request)
	if err != nil {
		return err
	}

	if response.StatusCode == http.StatusOK {
		return nil
	} else if response.StatusCode == http.StatusBadRequest {
		return fmt.Errorf("bad request: %v", response)
	} else if response.StatusCode == http.StatusUnprocessableEntity {
		return ErrTeamExists
	}
	return fmt.Errorf("unexpected code %d: %v", response.StatusCode, response)
}

func (c ClientImpl) CreateUser(csrfToken string, cookie *http.Cookie, user *model.User) error {
	data, err := json.Marshal(user)
	if err != nil {
		return err
	}

	request, err := http.NewRequest("POST", c.url+"/api/v0/users", strings.NewReader(string(data)))
	if err != nil {
		return err
	}
	request.AddCookie(cookie)
	request.Header.Set("x-csrf-token", csrfToken)

	response, err := c.httpClient.Do(request)
	if err != nil {
		return err
	}

	if response.StatusCode == http.StatusCreated {
		return nil
	} else if response.StatusCode == http.StatusBadRequest {
		return fmt.Errorf("bad request: %v", response)
	} else if response.StatusCode == http.StatusUnprocessableEntity {
		return ErrUserExists
	}
	return fmt.Errorf("unexpected code %d: %v", response.StatusCode, response)
}

func (c ClientImpl) UpdateUser(csrfToken string, cookie *http.Cookie, user *model.User) error {
	data, err := json.Marshal(user)
	if err != nil {
		return err
	}

	request, err := http.NewRequest("PUT", c.url+"/api/v0/users/"+user.Name, strings.NewReader(string(data)))
	if err != nil {
		return err
	}
	request.AddCookie(cookie)
	request.Header.Set("x-csrf-token", csrfToken)

	response, err := c.httpClient.Do(request)
	if err != nil {
		return err
	}

	if response.StatusCode == http.StatusNoContent {
		return nil
	} else if response.StatusCode == http.StatusNotFound {
		return ErrUserNotFound
	}
	return fmt.Errorf("unexpected code %d: %v", response.StatusCode, response)
}

func (c ClientImpl) CreateTeamUser(csrfToken string, cookie *http.Cookie, teamName string, userName string) error {
	data, err := json.Marshal(map[string]string{
		"team": teamName,
		"user": userName,
	})
	if err != nil {
		return err
	}

	request, err := http.NewRequest("POST", c.url+"/api/v0/teams/"+teamName+"/users", strings.NewReader(string(data)))
	if err != nil {
		return err
	}
	request.AddCookie(cookie)
	request.Header.Set("x-csrf-token", csrfToken)

	response, err := c.httpClient.Do(request)
	if err != nil {
		return err
	}

	if response.StatusCode == http.StatusCreated {
		return nil
	} else if response.StatusCode == http.StatusBadRequest {
		return fmt.Errorf("bad request: %v", response)
	} else if response.StatusCode == http.StatusUnprocessableEntity {
		return ErrUserExists
	}
	return fmt.Errorf("unexpected code %d: %v", response.StatusCode, response)
}

func (c ClientImpl) CreateRoster(csrfToken string, cookie *http.Cookie, teamName string, rosterName string) error {
	data, err := json.Marshal(map[string]string{
		"name": rosterName,
	})

	request, err := http.NewRequest("POST", c.url+"/api/v0/teams/"+teamName+"/rosters", strings.NewReader(string(data)))
	if err != nil {
		return err
	}
	request.AddCookie(cookie)
	request.Header.Set("x-csrf-token", csrfToken)

	response, err := c.httpClient.Do(request)
	if err != nil {
		return err
	}

	if response.StatusCode == http.StatusCreated {
		return nil
	} else if response.StatusCode == http.StatusUnprocessableEntity {
		return ErrRosterExists
	}
	return fmt.Errorf("unexpected code %d: %v", response.StatusCode, response)
}

func (c ClientImpl) CreateRosterUser(csrfToken string, cookie *http.Cookie, teamName string, rosterName string, userName string) error {
	data, err := json.Marshal(map[string]string{
		"name": userName,
	})
	if err != nil {
		return err
	}

	request, err := http.NewRequest("POST", c.url+"/api/v0/teams/"+teamName+"/rosters/"+rosterName+"/users", strings.NewReader(string(data)))
	if err != nil {
		return err
	}
	request.AddCookie(cookie)
	request.Header.Set("x-csrf-token", csrfToken)

	response, err := c.httpClient.Do(request)
	if err != nil {
		return err
	}

	if response.StatusCode == http.StatusCreated {
		return nil
	} else if response.StatusCode == http.StatusBadRequest {
		return fmt.Errorf("bad request: %v", response)
	} else if response.StatusCode == http.StatusUnprocessableEntity {
		return ErrUserExists
	}
	return fmt.Errorf("unexpected code %d: %v", response.StatusCode, response)
}

func (c ClientImpl) ExistEvent(csrfToken string, cookie *http.Cookie, teamName string, userName string, role string, start int64) (bool, error) {
	path := fmt.Sprintf("/api/v0/events?team=%s&user=%s&role=%s&start=%d", teamName, userName, role, start)
	request, err := http.NewRequest("GET", c.url+path, nil)
	request.AddCookie(cookie)
	request.Header.Set("x-csrf-token", csrfToken)

	response, err := c.httpClient.Do(request)
	if err != nil {
		return false, err
	}

	if response.StatusCode == http.StatusOK {
		var body []byte
		body, err = io.ReadAll(response.Body)
		if err != nil {
			return false, err
		}
		eventCount := 0
		_, err = jsonparser.ArrayEach(body, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
			eventCount += 1
		})
		if err != nil {
			return false, err
		}
		return eventCount > 0, nil
	}
	return false, fmt.Errorf("unexpected code %d: %v", response.StatusCode, response)
}

func (c ClientImpl) CreateEvent(csrfToken string, cookie *http.Cookie, event *model.Event) error {
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	request, err := http.NewRequest("POST", c.url+"/api/v0/events", strings.NewReader(string(data)))
	if err != nil {
		return err
	}
	request.AddCookie(cookie)
	request.Header.Set("x-csrf-token", csrfToken)

	response, err := c.httpClient.Do(request)
	if err != nil {
		return err
	}

	if response.StatusCode == http.StatusCreated {
		return nil
	} else if response.StatusCode == http.StatusBadRequest {
		return fmt.Errorf("bad request: %v", response)
	}
	return fmt.Errorf("unexpected code %d: %v", response.StatusCode, response)
}
