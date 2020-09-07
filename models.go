package loginctl

type SessionInfo = struct {
	Platform       string
	UserActivities map[string]bool
}

type Session = struct {
	ID     string
	UID    string
	User   string
	Active bool
}
