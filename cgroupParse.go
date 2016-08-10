package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func findCgroupMountpointAndRoot(pid int, subsystem string) (string, string, error) {
	f, err := os.Open(fmt.Sprintf("/proc/%d/mountinfo", pid))
	if err != nil {
		return "", "", err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		txt := scanner.Text()
		fields := strings.Split(txt, " ")
		for _, opt := range strings.Split(fields[len(fields)-1], ",") {
			if opt == subsystem {
				return fields[4], fields[3], nil
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return "", "", err
	}

	return "", "", fmt.Errorf("cgroup path for %s not found", subsystem)
}

func parseCgroupFile(path string) (map[string]string, error) {
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

func main() {
	mountpoint, hostRoot, err := findCgroupMountpointAndRoot(os.Getpid(), "memory")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("%s:%s\n", mountpoint, hostRoot)
	}

	cgroups, err := parseCgroupFile(fmt.Sprintf("/proc/%d/cgroup", os.Getpid()))
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("-----------------------\n")
		for k, v := range cgroups {
			fmt.Printf("%s:%s\n", k, v)
		}
	}
}
