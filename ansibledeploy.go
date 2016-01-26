package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"

	// thirdparty lib
	"github.com/go-ini/ini"
	"github.com/soarpenguin/log4go"
	//"gopkg.in/yaml.v2"
)

const (
	VERSION = "0.1.0"
)

var ANSIBLE_CMD string = "/usr/bin/ansible-playbook"
var gLogger log4go.Logger

const (
	retOk = iota
	retFailed
	retInvaidArgs
)

//de-init for all
func deinitLogger() {
	if nil != gLogger {
		gLogger.Close()
		gLogger = nil
	}
}

//init for logger
func initLogger(log4file bool) {
	var filenameOnly string
	var logFilename string

	if gLogger != nil {
		gLogger.Close()
		gLogger = nil
	}

	filenameOnly = GetCurFilename()
	logFilename = filenameOnly + ".log"

	gLogger = make(log4go.Logger)
	//for console
	gLogger.AddFilter("stdout", log4go.INFO, log4go.NewConsoleLogWriter())
	//for log file
	if log4file {
		if _, err := os.Stat(logFilename); err == nil {
			os.Remove(logFilename)
		}
		gLogger.AddFilter("logfile", log4go.FINEST, log4go.NewFileLogWriter(logFilename, false))
	}

	return
}

// GetCurFilename
// Get current file name, without suffix
func GetCurFilename() string {
	var filenameWithSuffix string
	var fileSuffix string
	var filenameOnly string

	_, fulleFilename, _, _ := runtime.Caller(0)
	filenameWithSuffix = path.Base(fulleFilename)
	fileSuffix = path.Ext(filenameWithSuffix)
	filenameOnly = strings.TrimSuffix(filenameWithSuffix, fileSuffix)

	return filenameOnly
}

func execCmd(cmd string, shell bool) (ret int, out []byte, err error) {
	fmt.Printf("run command: %s\n", cmd)
	if shell {
		//out, err = exec.Command("bash", "-c", cmd).Output()
		runcmd := exec.Command("bash", "-c", cmd)
		runcmd.Stdout = os.Stdout
		runcmd.Stderr = os.Stderr
		err := runcmd.Start()
		if err != nil {
			log.Fatal("%v", err)
		}

		if err = runcmd.Wait(); err != nil {
			if exiterr, ok := err.(*exec.ExitError); ok {
				// The program has exited with an exit code != 0

				if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
					ret = int(status)
				}
			} else {
				log.Fatal("cmd.Wait: %v", err)
			}
		}
	} else {
		out, err = exec.Command(cmd).Output()
	}
	return ret, out, err
}

func checkExistFiles(files ...string) bool {
	loginfo := "checkExistFiles"

	for _, file := range files {
		file = strings.TrimSpace(file)
		if _, err := isExists(file); err != nil {
			fmt.Printf("[ERROR] %s() - %s.\n", loginfo, err)
			return false
		}
	}
	return true
}

func isExists(file string) (ret bool, err error) {
	// equivalent to Python's `if not os.path.exists(filename)`
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return false, err
	} else {
		return true, nil
	}
}

func opYamlSyntaxCheck(action, inventory_file, operation_file, ext_vars, section string) {
	loginfo := "opYamlSyntaxCheck"
	var check_cmd string

	if _, err := isExists(operation_file); err != nil {
		fmt.Printf("[ERROR] %s() - %s.\n", loginfo, err)
		panic(err)
	}

	// ansible-playbook -i ~/hosts ~/test.yml --extra-vars "forks=1 hosts=update " --syntax-check
	if ext_vars != "" {
		check_cmd = fmt.Sprintf("%s -i %s %s --extra-vars \"forks=1 hosts=%s %s\" -t %s --syntax-check ",
			ANSIBLE_CMD, inventory_file, operation_file, action, ext_vars, action)
	} else {
		check_cmd = fmt.Sprintf("%s -i %s %s --extra-vars \"forks=1 hosts=%s \" -t %s --syntax-check ",
			ANSIBLE_CMD, inventory_file, operation_file, action, action)
	}

	if ret, _, err := execCmd(check_cmd, true); err != nil {
		fmt.Printf("[ERROR] %s() - %s.\n", loginfo, err)
		panic(err)
	} else if ret != 0 {
		fmt.Printf("[ERROR] syntax-check failed of %s.\n", operation_file)
		os.Exit(retFailed)
	}
}

func displayYamlFile(file string) {
	loginfo := "displayYamlFile"

	if _, err := isExists(file); err != nil {
		fmt.Printf("[ERROR] %s() - %s.\n", loginfo, err)
		panic(err)
	}

	contents, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}

	println(string(contents))
}

func getSingleHostname(file string, section string) string {
	var ini_cfg *ini.File
	loginfo := "getSingleHostname"

	hostname := ""

	ini_cfg, err := ini.Load(file)
	if err != nil || ini_cfg == nil {
		fmt.Printf("[ERROR] %s() - %s.\n", loginfo, err)
		return hostname
	}

	items, err := ini_cfg.GetSection(section)
	if err != nil {
		fmt.Printf("[ERROR] %s() from %s %s.\n", loginfo, file, err)
		os.Exit(retFailed)
	}

	keystrs := items.KeyStrings()
	if len(keystrs) <= 0 {
		fmt.Printf("[ERROR] %s() get single hostname from %s failed.", loginfo, file)
		os.Exit(retFailed)
	}

	for _, v := range keystrs {
		v = strings.TrimSpace(v)

		if v != "" {
			tokens := strings.Split(v, " ")

			if len(tokens) == 0 {
				continue
			}

			host := tokens[0]

			if strings.HasPrefix(host, "(") && strings.HasSuffix(host, ")") {
				srp := strings.NewReplacer("(", "", ")", "")
				host = srp.Replace(host)
			} else if strings.HasPrefix(host, "[") && strings.HasSuffix(host, "]") {
				srp := strings.NewReplacer("[", "", "]", "")
				host = srp.Replace(host)
			}

			tokens = strings.Split(host, " ")
			if len(tokens) == 0 {
				continue
			}
			hostname = tokens[0]

			// Three cases to check:
			// 0. A hostname that contains a range pesudo-code and a port
			// 1. A hostname that contains just a port
			if strings.Count(hostname, ":") > 1 {
				// Possible an IPv6 address, or maybe a host line with multiple ranges
				// IPv6 with Port  XXX:XXX::XXX.port
				// FQDN            foo.example.com
				if strings.Count(hostname, ".") == 1 {
					hostname = hostname[0:strings.LastIndex(hostname, ".")]
				}
			} else if (strings.Count(hostname, "[") > 0 && strings.Count(hostname, "]") > 0 &&
				(strings.LastIndex(hostname, "]") < strings.LastIndex(hostname, ":"))) ||
				((strings.Count(hostname, "]") <= 0) && (strings.Count(hostname, ":") > 0)) {
				hostname = hostname[0:strings.LastIndex(hostname, ":")]
			}
		}

		if hostname != "" {
			break
		}
	}

	return hostname
}

func doUpdateAction(action string, inventory_file string, operation_file string,
	version string, concurrent int) {

	loginfo := "doUpdateAction"
	var cmd, ext_vars string

	if action != "update" || inventory_file == "" || operation_file == "" {
		fmt.Printf("[ERROR] Error parameters in %s\n", loginfo)
		goto proc_error
	}

	if version != "" {
		version = strings.TrimSpace(version)
		ext_vars = fmt.Sprintf(" version=%s ", version)
	} else {
		ext_vars = ""
	}

	if ext_vars != "" {
		cmd = fmt.Sprintf("%s -i %s %s --extra-vars \"forks=%d hosts=%s %s\" -t %s -f %d ",
			ANSIBLE_CMD, inventory_file, operation_file, concurrent, action, ext_vars, action, concurrent)
	} else {
		cmd = fmt.Sprintf("%s -i %s %s --extra-vars \"forks=%d hosts=%s \" -t %s -f %d ",
			ANSIBLE_CMD, inventory_file, operation_file, concurrent, action, action, concurrent)
	}

	if _, _, err := execCmd(cmd, true); err != nil {
		fmt.Printf("[ERROR] %s() with error: %s\n", loginfo, err)
		goto proc_error
	}

	return

proc_error:
	os.Exit(retFailed)
}

func doDeployAction(action string, inventory_file string, operation_file string,
	singlemode bool, concurrent int, retry_file string, ext_vars string, section string) {

	loginfo := "doDeployAction"
	var cmd string

	if action != "deploy" || inventory_file == "" || operation_file == "" {
		fmt.Printf("[ERROR] Error parameters in %s", loginfo)
		goto proc_error
	}

	section = strings.TrimSpace(section)
	if section != "" {
		action = section
	}

	if ext_vars != "" {
		cmd = fmt.Sprintf("%s -i %s %s --extra-vars \"forks=%d hosts=%s %s\" -t %s -f %d ",
			ANSIBLE_CMD, inventory_file, operation_file, concurrent, action, ext_vars, action, concurrent)
	} else {
		cmd = fmt.Sprintf("%s -i %s %s --extra-vars \"forks=%d hosts=%s \" -t %s -f %d ",
			ANSIBLE_CMD, inventory_file, operation_file, concurrent, action, action, concurrent)
	}

	if singlemode {
		hostname := getSingleHostname(inventory_file, action)

		if hostname != "" {
			cmd = fmt.Sprintf("%s -l %s ", cmd, strings.TrimSpace(hostname))
		} else {
			fmt.Printf("[ERROR] %s() get single hostname from %s failed in single mode.", loginfo, inventory_file)
			goto proc_error
		}
	} else if retry_file != "" {
		fmt.Printf("%s\n", retry_file)
		if _, err := isExists(retry_file); err != nil {
			fmt.Printf("[ERROR] %s() please check the exists of retry file: %s.\n", loginfo, retry_file)
			goto proc_error
		} else {
			cmd = fmt.Sprintf("%s --limit @%s ", cmd, retry_file)
		}
	}

	if _, _, err := execCmd(cmd, true); err != nil {
		fmt.Printf("[ERROR] %s() with error: %s\n", loginfo, err)
		goto proc_error
	}

	return

proc_error:
	os.Exit(retFailed)
}

var Usage = func() {
	fmt.Fprintf(os.Stdout, "Usage:\n  %s [options] action\n\n", os.Args[0])
	fmt.Fprintf(os.Stdout, "Options:\n")
	flag.PrintDefaults()

	fmt.Fprintf(os.Stdout, "\n  action    action to do required:(check,update,deploy,rollback).\n")
	os.Exit(retOk)
}

var (
	concurrent                                       *int
	log4file, single_mode, version                   *bool
	program_version, extra_vars, section, opt_action *string
	retry_file, inventory_file, operation_file       *string
	action                                           string
)

func init() {
	flag.Usage = Usage

	log4file = flag.Bool("log4file", false, "Specify log in file for output.")
	single_mode = flag.Bool("s", false, "Single mode in deploy one host for observation.")
	concurrent = flag.Int("c", 1, "Process nummber for run the command at same time.")
	program_version = flag.String("V", "", "Module program version for deploy.")
	extra_vars = flag.String("e", "", "Extra vars for ansible-playbook.")
	section = flag.String("S", "", "Inventory section for distinguish hosts or tags.")
	opt_action = flag.String("action", "", "Action to do required:(check,update,deploy,rollback).")
	retry_file = flag.String("r", "", "Retry file for ansible redo failed hosts.")
	inventory_file = flag.String("i", "", "Specify inventory host file.")
	operation_file = flag.String("f", "", "File name for module configure(yml format).")
	version = flag.Bool("v", false, "Show version")

	flag.Parse()

	if *opt_action == "" {
		action = flag.Arg(0)
	} else {
		action = *opt_action
	}

	initLogger(*log4file)
	defer deinitLogger()
}

func main() {
	var ini_cfg *ini.File
	var err error

	if _, err = isExists(ANSIBLE_CMD); err != nil {
		ANSIBLE_CMD = "ansible-playbook"
	}

	if *version {
		fmt.Printf("%s: %s\n", os.Args[0], VERSION)
		os.Exit(retOk)
	}

	if *operation_file == "" || *inventory_file == "" {
		fmt.Printf("[ERROR] Not supported action: %s\n", action)
		//gLogger.Error("operation and inventory file must provide.\n")
		flag.Usage()
		os.Exit(retFailed)
	} else {
		ret := checkExistFiles(*operation_file, *inventory_file)
		if !ret {
			fmt.Printf("[ERROR] check exists of operation and inventory file.\n")
			//gLogger.Error("check exists of operation and inventory file.\n")
			os.Exit(retInvaidArgs)
		}
	}

	if action == "" {
		fmt.Printf("[ERROR] action(check,update,deploy,rollback) must provide one.\n\n")
		//gLogger.Error("action(check,update,deploy,rollback) must provide one.\n")
		flag.Usage()
		os.Exit(retInvaidArgs)
	}

	if *retry_file != "" {
		if *retry_file, err = filepath.Abs(*retry_file); err != nil {
			panic(fmt.Errorf("get Abs path of %s failed: %s\n", *retry_file, err))
			//gLogger.Error("get Abs path of %s failed: %s\n", *retry_file, err)
			os.Exit(retInvaidArgs)
		}
	}

	if *inventory_file, err = filepath.Abs(*inventory_file); err != nil {
		panic(fmt.Errorf("get Abs path of %s failed: %s\n", *inventory_file, err))
	} else {
		ini_cfg, err = ini.Load(*inventory_file)
		if err != nil {
			panic(fmt.Errorf("ini load conf failed: %s\n", err))
		}
	}

	if *operation_file, err = filepath.Abs(*operation_file); err != nil {
		panic(fmt.Errorf("get Abs path of %s failed: %s\n", *operation_file, err))
	}

	fmt.Printf("[%s] action on [%s]\n", action, *operation_file)
	//gLogger.Info("[%s] action on [%s]\n", action, *operation_file)
	switch action {
	case "check":
		fmt.Printf("-------------Now doing in action: %s\n", action)
		//gLogger.Info("-------------Now doing in action: %s\n", action)
		fmt.Printf("inventory: %s\n", *inventory_file)
		ini_cfg.WriteTo(os.Stdout)
		fmt.Println("")
		opYamlSyntaxCheck(action, *inventory_file, *operation_file, *extra_vars, "all")
		displayYamlFile(*operation_file)
	case "update":
		fmt.Printf("-------------Now doing in action: %s\n", action)
		//gLogger.Info("-------------Now doing in action: %s\n", action)
		doUpdateAction(action, *inventory_file, *operation_file, *program_version, *concurrent)
	case "deploy":
		fmt.Printf("-------------Inventory file is: %s\n", *inventory_file)
		//gLogger.Info("-------------Inventory file is: %s\n", *inventory_file)
		fmt.Printf("inventory: %s\n", *inventory_file)
		ini_cfg.WriteTo(os.Stdout)
		fmt.Printf("-------------Now doing in action:[%s], single mode:[%t]\n", action, *single_mode)
		doDeployAction(action, *inventory_file, *operation_file, *single_mode, *concurrent, *retry_file, *extra_vars, *section)
	case "rollback":
		fmt.Printf("-------------Now doing in action: %s\n", action)
		//gLogger.Info("-------------Now doing in action: %s\n", action)
		fmt.Println("rollback action do nothing now.")
	default:
		fmt.Printf("Not supported action: %s\n", action)
		//gLogger.Info("Not supported action: %s\n", action)
		os.Exit(retFailed)
	}
}
