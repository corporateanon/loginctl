// +build windows,amd64
package loginctl

import (
	"runtime"
	"strings"

	wapi "github.com/iamacarpet/go-win64api"
)

type Loginctl struct {
	trackedUsers []string
}

func (lctl *Loginctl) GetSessionInfo() (*SessionInfo, error) {
	loggedInUsers, err := wapi.ListLoggedInUsers()
	if err != nil {
		return nil, err
	}

	userActivity := map[string]bool{}

	for _, user := range lctl.trackedUsers {
		userActivity[user] = false
	}

	for _, loggedInUser := range loggedInUsers {
		loggedInUserName := strings.ToLower(loggedInUser.Username)

		if active, ok := userActivity[loggedInUserName]; !ok || active {
			continue
		}
		userActivity[loggedInUserName] = true
	}
	return &SessionInfo{
		Platform:       runtime.GOOS,
		UserActivities: userActivity,
	}, nil
}

func (lctl *Loginctl) GetUsersList(filter bool) ([]User, error) {
	winUsers, err := wapi.ListLocalUsers()
	if err != nil {
		return nil, err
	}
	users := []User{}

	for _, winUser := range winUsers {
		if filter && !winUser.IsEnabled {
			continue
		}
		user := User{
			ID:   0,
			Name: strings.ToLower(winUser.Username),
		}
		users = append(users, user)
	}

	return users, nil
}
