package model

type Duty struct {
	Date string `yaml:"date"`
	Role string `yaml:"role"`
}

type Contacts struct {
	Call  string `json:"call"`
	Email string `json:"email"`
}

type User struct {
	Name        string   `yaml:"name" json:"name"`
	FullName    string   `yaml:"full_name" json:"full_name"`
	PhoneNumber string   `yaml:"phone_number" json:"-"`
	Email       string   `yaml:"email" json:"-"`
	Contacts    Contacts `json:"contacts"`
	Duty        []*Duty  `yaml:"duty" json:"-"`
}

type Team struct {
	Name               string  `yaml:"name" json:"name"`
	SchedulingTimezone string  `yaml:"scheduling_timezone" json:"scheduling_timezone,omitempty"`
	Email              string  `yaml:"email" json:"email"`
	SlackChannel       string  `yaml:"slack_channel" json:"slack_channel"`
	Users              []*User `yaml:"users" json:"-"`
}

type Teams struct {
	Teams []*Team `yaml:"teams"`
}
