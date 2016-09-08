package lib

import (
	"bufio"
	"errors"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const cgroupNamePrefix = "name="

func Systemexec(s string) (string, error) {
	cmd := exec.Command("/bin/sh", "-c", s)

	log.Println(s)

	buff, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(buff), err
}

func getControllerPath(subsystem string, cgroups map[string]string) (string, error) {

	if p, ok := cgroups[subsystem]; ok {
		return p, nil
	}

	if p, ok := cgroups[cgroupNamePrefix+subsystem]; ok {
		return p, nil
	}

	return "", errors.New("could not find the subsystem")
}

func GetContainerIDFromCgroup(cgroup map[string]string) (string, error) {
	value, err := getControllerPath("cpu", cgroup)
	if err != nil {
		return "", nil
	}
	return value, nil

}

func ParseCgroupFile(path string) (map[string]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	s := bufio.NewScanner(f)
	cgroups := make(map[string]string)

	for s.Scan() {
		if err := s.Err(); err != nil {
			return nil, err
		}

		text := s.Text()
		parts := strings.Split(text, ":")

		for _, subs := range strings.Split(parts[1], ",") {
			cgroups[subs] = parts[2]
		}
	}
	return cgroups, nil
}

//using sys operation to get the pid according to the command lsof
//should have the root privilage
func GetInfofromPortbylsof(portip int) (string, error) {
	//cmd := "lsof -i :" + strconv.Itoa(portip) + " > null"
	cmd := "lsof -i :" + strconv.Itoa(portip) + "|grep LISTEN |awk '{print $2}'"
	returnString, err := Systemexec(cmd)
	if err != nil {
		return "", err
	}
	return returnString, nil
}

//if the process of pid is in cgroup control, return true,containerid,error else return false,"",err
func IfinCgroupControl(pid int) (bool, string, error) {
	containerID, err := getContainerIDFromPid(pid)
	if err != nil {
		return false, "", err
	}
	return true, containerID, nil
}

func getContainerIDFromPid(pid int) (string, error) {
	port := strconv.Itoa(pid)
	path := "/proc/" + port + "/cgroup"
	cgroupMap, err := ParseCgroupFile(path)
	if err != nil {
		return "", err
	}

	rawcontainerID, err := GetContainerIDFromCgroup(cgroupMap)
	if err != nil {
		return "", err
	}
	//log.Println(rawcontainerID)
	strList := strings.Split(rawcontainerID, "/")
	if len(strList) < 3 {
		return "", errors.New("not in cgroup control")
	}
	return strList[2], nil
}
