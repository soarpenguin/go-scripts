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
	program_version := flag.String("V", "", "Module program version for deploy.")
	extra_vars := flag.String("e", "", "Extra vars for ansible-playbook.")
	section := flag.String("S", "", "Inventory section for distinguish hosts or tags.")
	retry_file := flag.String("r", "", "Retry file for ansible redo failed hosts.")
	inventory_file := flag.String("i", "", "Specify inventory host file.")
	operation_file := flag.String("f", "", "File name for module configure(yml format).")
	version := flag.Bool("v", false, "show version")

	flag.Parse()

	if *version {
		fmt.Printf("%s: %s\n", os.Args[0], VERSION)
		os.Exit(retOK)
	}

	var action string = flag.Arg(0)

	fmt.Printf("single_mode   : %s\n", *single_mode)
	fmt.Printf("concurrent   : %s\n", *concurrent)
	fmt.Printf("program_version   : %s\n", *program_version)
	fmt.Printf("extra_vars   : %s\n", *extra_vars)
	fmt.Printf("section   : %s\n", *section)
	fmt.Printf("retry_file   : %s\n", *retry_file)
	fmt.Printf("inventory_file   : %s\n", *inventory_file)
	fmt.Printf("operation_file   : %s\n", *operation_file)
	fmt.Printf("action   : %s\n", action)

}
