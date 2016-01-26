package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"

	flag "github.com/docker/docker/pkg/mflag"
)

// var and struct define section
const (
	retOk = iota
	retFailed
	retInvaidArgs
)

var (
	help                bool
	tplfile, appName    string
	consulIp, virtualIp string
)

type VIPPort struct {
	VirtualIp string
	Port      string
}

type Application struct {
	AppName  string
	ConsulIp string
	Servers  []VIPPort
}

// function define section.
func isExists(file string) (ret bool, err error) {
	// equivalent to Python's `if not os.path.exists(filename)`
	if _, err := os.Stat(file); err != nil {
		return false, err
	} else {
		return true, nil
	}
}

func checkIPPortValid(ipPort string) (ret bool, err error) {
	ipPort = strings.TrimSpace(ipPort)

	if ipPort == "" {
		return false, fmt.Errorf("ip and port string is nil")
	}

	ipPortRegx := "^(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?).){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(:[0-9]+)?$"
	if m, _ := regexp.MatchString(ipPortRegx, ipPort); !m {
		return false, fmt.Errorf("invalid ip and port string")
	}
	return true, nil
}

func checkIPValid(ipv4 string) (ret bool, err error) {
	ipv4 = strings.TrimSpace(ipv4)

	if ipv4 == "" {
		return false, fmt.Errorf("ipv4 string is nil")
	}

	ipRegx := "^(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?).){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$"
	if m, _ := regexp.MatchString(ipRegx, ipv4); !m {
		return false, fmt.Errorf("invalid ipv4 string")
	}
	return true, nil
}

func checkIPS(virtualIp string) (vips []VIPPort, err error) {
	virtualIp = strings.TrimSpace(virtualIp)

	if virtualIp == "" {
		return nil, fmt.Errorf("virtualIp string is nil")
	}

	ipports := strings.Split(virtualIp, ",")

	for _, ipport := range ipports {
		ipport = strings.TrimSpace(ipport)
		if ipport == "" {
			continue
		}

		_, err = checkIPPortValid(ipport)
		if err == nil {
			if strings.Count(ipport, ":") > 0 {
				result := strings.SplitN(ipport, ":", 2)
				vips = append(vips, VIPPort{VirtualIp: result[0], Port: result[1]})
			} else {
				vips = append(vips, VIPPort{VirtualIp: ipport, Port: "80"})
			}
		} else {
			return nil, err
		}
	}

	return vips, nil
}

func checkAppName(appName string) (appname string, err error) {
	appname = strings.TrimSpace(appName)

	if appName == "" {
		return "", fmt.Errorf("appName string is nil")
	}

	nameRegx := `^[\w][\w.-]{0,127}$`
	if m, _ := regexp.MatchString(nameRegx, appname); !m {
		return "", fmt.Errorf("invalid application name")
	}

	return appname, nil
}

func init() {
	flag.StringVar(&tplfile, []string{"t", "-template"}, "server.conf.template", "Template file name for produce config file.")
	flag.StringVar(&appName, []string{"-app"}, "appname", "App name for upstream/logfile name/conf name/consul key.")
	flag.StringVar(&consulIp, []string{"c", "-consul"}, "10.10.10.10:8500", "Consul server 'url:port' for get upstream info.")
	flag.StringVar(&virtualIp, []string{"v", "-virtualip"}, "0.0.0.0,0.0.0.1", "Virtual IP list for this app.")
	flag.BoolVar(&help, []string{"h", "-help"}, false, "Display the help")
	flag.Parse()
}

func Usage() {
	fmt.Printf("Usage: %s [options]\n\n%s", os.Args[0], "Options:")
	flag.PrintDefaults()
	os.Exit(retOk)
}

// main route
func main() {
	if help {
		Usage()
	}

	// check valid for appname.
	appName, err := checkAppName(appName)
	if err != nil {
		log.Fatalf("App name check failed: %s", err)
	} else if appName == "appname" {
		log.Fatalf("Please provide app name, not use the default value.")
	}

	// check valid for consul ip.
	_, err = checkIPPortValid(consulIp)
	if err != nil {
		log.Fatalf("Consul server IP check failed: %s", err)
	}

	// check valid for virtual list.
	vips, err := checkIPS(virtualIp)
	if err != nil {
		log.Fatalf("Virtual IP check failed: %s", err)
	} else if vips == nil || vips[0].VirtualIp == "0.0.0.0" {
		log.Fatalf("Please provide one virtual ip at least.")
	}

	// check valid of template file.
	if tplfile, err = filepath.Abs(tplfile); err != nil {
		log.Fatalf("Get Abs path of %s failed: %s\n", tplfile, err)
	} else if _, err := isExists(tplfile); err != nil {
		log.Fatalf("%s\n", err)
	}

	app := Application{
		AppName:  appName,
		ConsulIp: consulIp,
		Servers:  vips,
	}

	t, err := template.ParseFiles(tplfile)
	if err != nil {
		log.Fatalf("Template parse failed:%s\n", err)
	}

	err = t.Execute(os.Stdout, app)
	if err != nil {
		log.Fatalf("Template execute failed:%s\n", err)
	}
}
