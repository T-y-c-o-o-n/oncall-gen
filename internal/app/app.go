package app

import (
	"oncall-gen/internal/model"
	"oncall-gen/internal/oncall"
	"oncall-gen/internal/parser"
	"time"
)

type App interface {
	CreateTeams(fileName string) error
}

type Impl struct {
	onCallClient oncall.Client
}

func NewAppImpl(onCallApiUrl string) App {
	return &Impl{
		onCallClient: oncall.NewClientImpl(onCallApiUrl),
	}
}

var dateLayout = "02/01/2006"
var oneDayDuration = time.Hour * 24

func (a Impl) CreateTeams(fileName string) error {
	teams, err := parser.ParseYaml(fileName)
	if err != nil {
		return err
	}

	token, cookie, err := a.onCallClient.Login("root", "root")

	if err != nil {
		return err
	}

	for _, team := range teams.Teams {
		err = a.onCallClient.CreateTeam(token, cookie, team)
		if err == oncall.ErrTeamExists {
			//err = a.onCallClient.UpdateTeam(token, cookie, team)
		} else if err != nil {
			return err
		}

		rosterName := team.Name + "-roster"
		err = a.onCallClient.CreateRoster(token, cookie, team.Name, rosterName)
		if err == oncall.ErrRosterExists {
		} else if err != nil {
			return err
		}

		for _, user := range team.Users {
			user.Contacts.Call = user.PhoneNumber
			user.Contacts.Email = user.Email

			err = a.onCallClient.CreateUser(token, cookie, user)
			if err == oncall.ErrUserExists {
				//err = a.onCallClient.UpdateUser(token, cookie, user)
			} else if err != nil {
				return err
			}

			err = a.onCallClient.UpdateUser(token, cookie, user)
			if err != nil {
				return err
			}

			// TODO maybe this is not necessary
			err = a.onCallClient.CreateTeamUser(token, cookie, team.Name, user.Name)
			if err == oncall.ErrUserExists {
			} else if err != nil {
				return err
			}

			err = a.onCallClient.CreateRosterUser(token, cookie, team.Name, rosterName, user.Name)
			if err == oncall.ErrUserExists {
			} else if err != nil {
				return err
			}

			for _, duty := range user.Duty {
				var startTime time.Time
				startTime, err = time.Parse(dateLayout, duty.Date)
				if err != nil {
					return nil
				}
				start := startTime.Unix()

				var existEvent bool
				existEvent, err = a.onCallClient.ExistEvent(token, cookie, team.Name, user.Name, duty.Role, start)
				if err != nil {
					return nil
				}
				if !existEvent {
					end := startTime.Add(oneDayDuration).Unix()
					event := &model.Event{
						Start: start,
						End:   end,
						User:  user.Name,
						Team:  team.Name,
						Role:  duty.Role,
					}
					err = a.onCallClient.CreateEvent(token, cookie, event)
					if err != nil {
						return err
					}
				}
			}
		}
	}

	return nil
}
