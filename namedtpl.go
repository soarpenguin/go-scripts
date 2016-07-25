package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
	"time"

	flag "github.com/docker/docker/pkg/mflag"
)

// var and struct define section
const (
	retOk = iota
	retFailed
	retInvaidArgs
)

var (
	help             bool
	output           string
	tplfile, appName string
)

type DNSRecord struct {
	DnsName string
	IPs     []string
}

type DNSRecords struct {
	AppName  string
	Services []DNSRecord
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

// generate rand string
const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

func RandStringBytesMaskImprSrc(n int) string {
	var src = rand.NewSource(time.Now().UnixNano())

	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}

func genConfigFile(tplfile string, app DNSRecords, output, appname string) {
	output = strings.TrimSpace(output)

	if output == "" || output == "appname.conf" {
		output = appname + ".conf"
	}

	if _, err := isExists(output); err == nil {
		log.Printf("[WARN] %s is existed, will redirect to tmp file!!!!", output)
		output = fmt.Sprintf("%s.%s", output, RandStringBytesMaskImprSrc(6))
	}

	t, err := template.ParseFiles(tplfile)
	if err != nil {
		log.Fatalf("Template parse failed:%s\n", err)
	}

	ofd, err := os.OpenFile(output, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Printf("[WARN] Open %s failed: %s\n", output, err)
		ofd = os.Stdout
	} else {
		defer ofd.Close()
	}

	//err = t.Execute(os.Stdout, app)
	err = t.Execute(ofd, app)
	if err != nil {
		log.Fatalf("Template execute failed:%s\n", err)
	}
	log.Printf("[NOTICE] please check content in %s", output)
}

func init() {
	flag.StringVar(&tplfile, []string{"t", "-template"}, "", "Template file name for produce config file.")
	flag.StringVar(&appName, []string{"-app"}, "appname", "App name for upstream/logfile name/conf name/consul key.")
	flag.StringVar(&output, []string{"o", "-output"}, "appname.conf", "Config file name for save config file.")
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

	// check valid of template file.
	if tplfile, err := filepath.Abs(tplfile); err != nil {
		log.Fatalf("Get Abs path of %s failed: %s\n", tplfile, err)
	} else if _, err := isExists(tplfile); err != nil {
		log.Fatalf("%s\n", err)
	}

	vips := []DNSRecord{
		{"test1", []string{"10.10.11.11", "10.101.10.11"}},
		{"test2", []string{"10.10.11.11", "10.101.10.11"}},
	}

	// new a application struct.
	app := DNSRecords{
		AppName:  appName,
		Services: vips,
	}
	// generate final config file.
	genConfigFile(tplfile, app, output, appName)
}
