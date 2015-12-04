package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	// thirdparty lib
	"github.com/mgutz/ansi"
)

const (
	VERSION = "0.1.0"
)

const (
	retOK = iota
	retFail
)

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func isProcessOK(err error) {
	var msg string

	if err != nil {
		msg = ansi.Color("     [FAIL]", "red+b")
		fmt.Println(msg)
	} else {
		msg = ansi.Color("     [OK]", "green+b")
		fmt.Println(msg)
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

func dealFileWithWhiteList(filename string, cmd string, suffixs []string) {
	if cap(suffixs) > 0 {
		ext := filepath.Ext(filename)

		if !stringInSlice(ext, suffixs) {
			fmt.Printf("skip deal with file of: %s\n", filename)
			return
		}
	}

	cmd += filename
	execCmd(cmd, true)
}

func dealDirWithWhiteList(path string, cmd string, suffixs []string) {
	fmt.Printf("deal with dir of: %s\n", path)
	err := filepath.Walk(path, func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		if f.IsDir() {
			if strings.HasPrefix(f.Name(), ".") {
				return filepath.SkipDir
			} else {
				return nil
			}
		} else {
			if !strings.HasPrefix(f.Name(), ".") {
				dealFileWithWhiteList(path, cmd, suffixs)
			}
		}
		return nil
	})

	if err != nil {
		fmt.Printf("filepath.Walk() returned %v\n", err)
	}
}

func isExists(file string) (ret bool, err error) {
	// equivalent to Python's `if not os.path.exists(filename)`
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return false, err
	} else {
		return true, nil
	}
}

var usage = func() {
	fmt.Fprintf(os.Stdout, "Usage:\n  %s [options] file/dir\n\n", os.Args[0])
	fmt.Fprintf(os.Stdout, "Options:\n")
	flag.PrintDefaults()

	fmt.Fprintf(os.Stdout, "\nRequires:\n")
	fmt.Fprintf(os.Stdout, "  file/dir    file or dir to deal with.\n")

	os.Exit(retOK)
}

var (
	cmd         string
	suffixArray []string

	suffixs = flag.String("s", "", "white list of file suffixs for deal.")
	version = flag.Bool("v", false, "show version")
)

func init() {
	flag.Usage = usage
	flag.Parse()

	if *version {
		fmt.Printf("%s: %s\n", os.Args[0], VERSION)
		os.Exit(retOK)
	}

	*suffixs = strings.TrimSpace(*suffixs)

	if *suffixs != "" {
		suffixArray = strings.Split(*suffixs, ",")
	}

	switch runtime.GOOS {
	case "windows":
		fmt.Printf("[Error] not supported under windows.\n")
		os.Exit(retFail)
	case "darwin", "freebsd":
		cmd = "/usr/bin/sed -i \"\" \"s/[ ]*$//g\" "
		//fallthrough
	default:
		cmd = "sed -i \"s/[ \t]*$//g\" "
	}
}

func main() {

	if flag.NArg() == 0 {
		fmt.Printf("[Error] path or file must provide one.\n\n")
		flag.Usage()
		os.Exit(retFail)
	}

	for i := 0; i < flag.NArg(); i++ {
		path := strings.TrimSpace(flag.Arg(i))
		switch dir, err := os.Stat(path); {
		case err != nil:
			//fmt.Printf("[Error] error in Stat(): %s.\n", err)
			panic(err)
		case dir.IsDir():
			if path, err = filepath.Abs(path); err != nil {
				panic(err)
			}
			dealDirWithWhiteList(path, cmd, suffixArray)
		default:
			if path, err = filepath.Abs(path); err != nil {
				panic(err)
			}
			dealFileWithWhiteList(path, cmd, suffixArray)
		}
	}
}
