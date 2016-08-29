package utils

import (
	"fmt"
	"os"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/go-ini/ini"
)

// Fatal prints the error's details if it is a libcontainer specific error type
// then exits the program with an exit status of 1.
func Fatal(err error) {
	// make sure the error is written to the logger
	logrus.Error(err)
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}

func IsExists(file string) (ret bool, err error) {
	// equivalent to Python's `if not os.path.exists(filename)`
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return false, err
	} else {
		return true, nil
	}
}

var logLevel string
var reqDoms []string
var cachePath string
var cacheTTL int

func LoadConfig(path string) error {

	if _, err := IsExists(path); err != nil {
		return fmt.Errorf("Unable to locate local configuration file. path: %s", path)
	}

	iniCfg, err := ini.Load(path)
	if err != nil || iniCfg == nil {
		return err
	}

	if sect, err := iniCfg.GetSection("MAIN"); sect == nil || err != nil {
		logLevel = "WARN"
	} else {
		if ok := sect.HasKey("log-level"); ok {
			logLevel = sect.Key("log-level").In("WARN", []string{"DEBUG", "INFO", "ERROR", "CRITICAL"})
		}

		if ok := sect.HasKey("req-doms"); ok {
			var doms []string
			reqdoms := strings.Split(sect.Key("req-doms").String(), ",")
			for _, dom := range reqdoms {
				dom = strings.TrimSpace(dom)
				doms = append(doms, dom)
			}

			if len(doms) <= 0 {
				return fmt.Errorf("Malformed req-doms in config file: %s", path)
			}
			reqDoms = doms
		}
	}

	if sect, err := iniCfg.GetSection("CACHE"); sect != nil && err == nil {
		if ok := sect.HasKey("cache-path"); ok {
			cachePath = sect.Key("cache-path").Validate(func(in string) string {
				if len(in) == 0 {
					return cachePath
				}
				return in
			})
		}
	}

	if sect, err := iniCfg.GetSection("DNS"); sect != nil && err == nil {
		if ok := sect.HasKey("cache-ttl"); ok {
			if v, err := sect.Key("cache-ttl").Int(); err == nil {
				cacheTTL = v
			}
		}
	}

	return nil
}
