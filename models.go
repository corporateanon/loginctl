package loginctl

type SessionInfo struct {
	Platform       string
	UserActivities map[string]bool
}

type Session struct {
	ID     string
	UID    uint64
	User   string
	Active bool
}

type User struct {
	ID   uint64
	Name string
}
