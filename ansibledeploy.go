package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	//"path/filepath"
	//"runtime"
	//"strings"
)

const (
	VERSION = "0.1.0"
)

const (
	ANSIBLE_CMD = "/usr/bin/ansible-playbook"
)

const (
	retOK = iota
	retFaied
	retInvaidArgs
)

func execCmd(cmd string, shell bool) (out []byte, err error) {
	fmt.Printf("run command: %s", cmd)
	if shell {
		out, err = exec.Command("bash", "-c", cmd).Output()
	} else {
		out, err = exec.Command(cmd).Output()
	}
	return out, err
}

func main() {

	single_mode := flag.Bool("s", false, "Single mode in deploy one host for observation.")
	concurrent := flag.Int("c", 1, "Process nummber for run the command at same time.")
	version := flag.Bool("v", false, "show version")

	flag.Parse()

	if *version {
		fmt.Printf("%s: %s\n", os.Args[0], VERSION)
		os.Exit(retOK)
	}

	var action string = flag.Arg(0)

	fmt.Printf("single_mode   : %s\n", *single_mode)
	fmt.Printf("concurrent   : %s\n", *concurrent)
	fmt.Printf("action   : %s\n", action)
}
