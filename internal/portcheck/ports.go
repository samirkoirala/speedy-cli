package portcheck

import (
	"fmt"
	"os/exec"
	"runtime"
	"strconv"
	"strings"

	"speedy-cli/internal/common"
)

type PortUsage struct {
	Port    int    `json:"port"`
	InUse   bool   `json:"inUse"`
	PID     int    `json:"pid,omitempty"`
	Process string `json:"process,omitempty"`
}

func Check(port int) (common.Result, []PortUsage) {
	ports := []int{3000, 5432, 6379, 8080}
	if port > 0 {
		ports = []int{port}
	}

	usages := make([]PortUsage, 0, len(ports))
	busy := 0
	for _, p := range ports {
		u := inspectPort(p)
		if u.InUse {
			busy++
		}
		usages = append(usages, u)
	}

	res := common.Result{Status: common.StatusSuccess, Message: "port scan complete"}
	if busy > 0 {
		res.Status = common.StatusWarning
		res.Suggestion = "Stop conflicting process before starting your service"
	}
	return res, usages
}

func inspectPort(port int) PortUsage {
	if runtime.GOOS == "windows" {
		return PortUsage{Port: port, InUse: false}
	}
	cmd := exec.Command("lsof", "-nP", "-iTCP:"+strconv.Itoa(port), "-sTCP:LISTEN")
	out, err := cmd.Output()
	if err != nil || len(out) == 0 {
		return PortUsage{Port: port, InUse: false}
	}
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	if len(lines) < 2 {
		return PortUsage{Port: port, InUse: false}
	}
	fields := strings.Fields(lines[1])
	if len(fields) < 2 {
		return PortUsage{Port: port, InUse: true}
	}
	pid, _ := strconv.Atoi(fields[1])
	proc := fields[0]
	if proc == "" {
		proc = fmt.Sprintf("pid-%d", pid)
	}
	return PortUsage{Port: port, InUse: true, PID: pid, Process: proc}
}
