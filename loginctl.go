package loginctl

type ILoginctl interface {
	GetSessionInfo() (*SessionInfo, error)
	GetUsersList(filter bool) ([]User, error)
}

func New(trackedUsers []string) ILoginctl {
	loginctl := &Loginctl{
		trackedUsers: trackedUsers,
	}

	return loginctl
}

func NewFromRegularUsers() (ILoginctl, error) {
	loginctl := &Loginctl{}
	users, err := loginctl.GetUsersList(true)
	userNames := []string{}

	for _, u := range users {
		userNames = append(userNames, u.Name)
	}
	if err != nil {
		return nil, err
	}

	loginctl.trackedUsers = userNames

	return loginctl, nil
}
