// +build linux

package loginctl

import (
	"bufio"
	"io"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
)

var reSessionLine = regexp.MustCompile(`^\s*(\w+)\s*(\d+)\s*([\w\-]+)\s*([\w]+)\s*$`)
var reActive = regexp.MustCompile(`Active=yes`)
var reIdle = regexp.MustCompile(`IdleHint=yes`)
var reUser = regexp.MustCompile(`^([^:]+):([^:]+):(\d+):`)

type Loginctl struct {
	trackedUsers []string
}

func (this *Loginctl) GetSessionInfo() (*SessionInfo, error) {
	sessions, err := this.getSessionsList()
	if err != nil {
		return nil, err
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
	}, nil
}

func (this *Loginctl) getSessionsList() ([]Session, error) {
	cmd := exec.Command("loginctl")

	sessions := []Session{}

	scanLinesCmd(
		cmd,
		func(line string) error {
			parts := reSessionLine.FindStringSubmatch(line)
			if parts != nil && len(parts) > 0 {

				uid, err := strconv.ParseUint(parts[2], 10, 64)
				if err != nil {
					return err //TODO: add line number and a line to the error
				}

				sess := Session{
					ID:   parts[1],
					UID:  uid,
					User: parts[3],
				}

				active, err := this.getSessionActive(sess.ID)
				if err == nil {
					sess.Active = active
				}
				sessions = append(sessions, sess)
			}
			return nil
		},
	)

	return sessions, nil
}

func (this *Loginctl) GetUsersList(filter bool) ([]User, error) {
	file, err := os.Open("/etc/passwd")
	if err != nil {
		return nil, err
	}
	users := []User{}
	if err := scanLines(file, func(line string) error {
		parts := reUser.FindStringSubmatch(line)
		if parts != nil && len(parts) == 4 {
			id, err := strconv.ParseUint(parts[3], 10, 64)
			if err != nil {
				return err //TODO: add line number and a line to the error
			}

			name := parts[1]

			if !filter || id >= 1000 && name != "nobody" {
				users = append(users, User{
					ID:   id,
					Name: parts[1],
				})
			}
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return users, nil
}

func (this *Loginctl) getSessionActive(sessionID string) (bool, error) {
	cmd := exec.Command("loginctl", "show-session", sessionID)
	out, err := cmd.Output()
	if err != nil {
		return false, err
	}
	return reActive.Match(out) && !reIdle.Match(out), nil
}

func scanLines(stdout io.ReadCloser, cb func(line string) error) error {
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		err := cb(scanner.Text())
		if err != nil {
			return err
		}
	}
	return nil
}

func scanLinesCmd(
	cmd *exec.Cmd,
	cb func(line string) error,
) error {
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	defer stdout.Close()

	if err := cmd.Start(); err != nil {
		return err
	}

	if err := scanLines(stdout, cb); err != nil {
		cmd.Wait()
		return err
	}

	return cmd.Wait()
}
