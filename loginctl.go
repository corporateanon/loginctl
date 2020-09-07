package loginctl

import (
	"bufio"
	"log"
	"os/exec"
	"regexp"
	"runtime"
)

var reSessionLine = regexp.MustCompile(`^\s*(\w+)\s*(\d+)\s*([\w\-]+)\s*([\w]+)\s*$`)
var reActive = regexp.MustCompile(`Active=yes`)
var reIdle = regexp.MustCompile(`IdleHint=yes`)

type Loginctl struct {
	trackedUsers []string
}

func New(trackedUsers []string) (*Loginctl, error) {
	loginctl := &Loginctl{
		trackedUsers: trackedUsers,
	}

	return loginctl, nil
}

func (this *Loginctl) GetSessionInfo() *SessionInfo {
	sessions, err := this.GetSessionsList()
	if err != nil {
		log.Println("error getting sessions list")
	}

	userActivity := map[string]bool{}

	for _, user := range this.trackedUsers {
		userActivity[user] = false
	}

	for _, sess := range sessions {
		if active, ok := userActivity[sess.User]; !ok || active {
			continue
		}
		if sess.Active {
			userActivity[sess.User] = true
		}
	}

	return &SessionInfo{
		Platform:       runtime.GOOS,
		UserActivities: userActivity,
	}
}

func (this *Loginctl) GetSessionsList() ([]Session, error) {
	cmd := exec.Command("loginctl")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	scanner := bufio.NewScanner(stdout)
	sessions := []Session{}
	for scanner.Scan() {
		line := scanner.Text()
		parts := reSessionLine.FindStringSubmatch(line)
		if parts != nil && len(parts) > 0 {
			sess := Session{
				ID:   parts[1],
				UID:  parts[2],
				User: parts[3],
			}

			active, err := this.getSessionActive(sess.ID)
			if err == nil {
				sess.Active = active
			}
			log.Println(sess)
			sessions = append(sessions, sess)
		}
	}
	return sessions, nil
}

func (this *Loginctl) getSessionActive(sessionID string) (bool, error) {
	cmd := exec.Command("loginctl", "show-session", sessionID)
	out, err := cmd.Output()
	if err != nil {
		return false, err
	}
	return reActive.Match(out) && !reIdle.Match(out), nil
}
