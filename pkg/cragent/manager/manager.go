package manager

import (
	"regexp"
)

var (
	DefaultCgroupDir = "/sys/fs/cgroup"
	DefaultDockerReg = regexp.MustCompile(`[a-zA-Z0-9-_+.]+:[a-fA-F0-9]+`)
)

type Manager struct {
}
