package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"regexp"
)

const (
	VERSION = "0.1.0"
)

const (
	retOK = iota
	retFail
)

func isIP(ip string) (b bool) {
	ippattern := "^((2[0-4]\\d|25[0-5]|[01]?\\d\\d?)\\.){3}(2[0-4]\\d|25[0-5]|[01]?\\d\\d?)$"
	if m, _ := regexp.MatchString(ippattern, ip); !m {
		return false
	}
	return true
}

func isProcessOK(err error) {
	if err != nil {
		fmt.Println("     [FAIL]")
	} else {
		fmt.Println("     [OK]")
	}
}

func execCmd(cmd string, shell bool) (out []byte, err error) {
	fmt.Printf("run command: %s", cmd)
	if shell {
		out, err = exec.Command("bash", "-c", cmd).Output()
		isProcessOK(err)
	} else {
		out, err = exec.Command(cmd).Output()
		isProcessOK(err)
	}
	return out, err
}

func main() {

	version := flag.Bool("v", false, "show version")

	flag.Parse()

	if *version {
		fmt.Printf("%s: %s\n", os.Args[0], VERSION)
		os.Exit(retOK)
	}

	var nameOrIP string = flag.Arg(0)

	if nameOrIP == "" {
		fmt.Printf("Please provide a hostname or ip.")
		os.Exit(retFail)
	}

	if isIP(nameOrIP) {
		fmt.Printf("Match ip %s\n", nameOrIP)
	} else {
		fmt.Printf("Not match ip %s\n", nameOrIP)
	}
}
