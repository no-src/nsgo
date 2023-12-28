package browser

import (
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/no-src/nsgo/osutil"
)

// commands returns a list of possible commands to use to open an url.
func commands() [][]string {
	var cmds [][]string
	if exe := os.Getenv("BROWSER"); exe != "" {
		cmds = append(cmds, []string{exe})
	}
	switch runtime.GOOS {
	case "darwin":
		cmds = append(cmds, []string{"/usr/bin/open"})
	case "windows":
		cmds = append(cmds, []string{"cmd", "/c", "start"})
	default:
		if os.Getenv("DISPLAY") != "" {
			// xdg-open is only for use in a desktop environment.
			cmds = append(cmds, []string{"xdg-open"})
		}
	}
	cmds = append(cmds,
		[]string{"chrome"},
		[]string{"google-chrome"},
		[]string{"chromium"},
		[]string{"msedge"},
		[]string{"firefox"},
	)
	return cmds
}

// OpenBrowser tries to open url in a browser and reports whether it succeeded.
func OpenBrowser(url string) bool {
	for _, args := range commands() {
		if osutil.IsWindows() {
			url = strings.ReplaceAll(url, "&", "^&")
		} else {
			url = strings.ReplaceAll(url, "&", `\&`)
		}
		cmd := exec.Command(args[0], append(args[1:], url)...)
		if cmd.Start() == nil && appearsSuccessful(cmd, 3*time.Second) {
			return true
		}
	}
	return false
}

// appearsSuccessful reports whether the command appears to have run successfully.
// If the command runs longer than the timeout, it's deemed successful.
// If the command runs within the timeout, it's deemed successful if it exited cleanly.
func appearsSuccessful(cmd *exec.Cmd, timeout time.Duration) bool {
	errc := make(chan error, 1)
	go func() {
		errc <- cmd.Wait()
	}()

	select {
	case <-time.After(timeout):
		return true
	case err := <-errc:
		return err == nil
	}
}
