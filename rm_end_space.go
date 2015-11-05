package main

import (
    "fmt"
    "flag"
    "runtime"
    "os/exec"
    "os"
    "strings"
    "path/filepath"
)

const (
    VERSION  = "0.1.0"
)

const (
    retOK   = iota
    retFail
)

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

func dealWithFile(filename string, cmd string) {
    cmd += filename
    execCmd(cmd, true)
}

func dealWithDir(path string, cmd string) {
    fmt.Printf("deal with dir of: %s\n", path)
    err := filepath.Walk(path, func(path string, f os.FileInfo, err error) error {
        if ( f == nil ) {return err}
        if f.IsDir() {
            if strings.HasPrefix(f.Name(), ".") {
                return filepath.SkipDir
            } else {
                return nil
            }
        } else {
            if ! strings.HasPrefix(f.Name(), ".") {
                dealWithFile(path, cmd)
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

func main(){

    var cmd string

    op_path := flag.String("d", "", "dir to recursive remove spaces at the end of the line.")
    op_file := flag.String("f", "", "file name for remove spaces at the end of line.")
    version := flag.Bool("v", false, "show version")

    flag.Parse()

    if *version {
        fmt.Printf("%s: %s\n", os.Args[0], VERSION)
        os.Exit(retOK)
    }

    switch runtime.GOOS {
        case "windows":
            fmt.Printf("Not supported under windows.\n")
            os.Exit(retFail)
        case "darwin", "freebsd":
            cmd = "/usr/bin/sed -i \"\" \"s/[ ]*$//g\" "
        default:
            cmd = "sed -i \"s/[ \t]*$//g\" "
    }

    if *op_path == "" && *op_file == "" {
        fmt.Printf("path or file must provide one.\n")
        os.Exit(retFail)
    } else if *op_file != "" {
        if _, err := isExists(*op_file); err == nil  {
            if *op_file, err = filepath.Abs(*op_file); err != nil {
                panic(err)
            }
            dealWithFile(*op_file, cmd)
        }
    } else if *op_path != "" {
        if _, err := isExists(*op_path); err == nil  {
            if *op_path, err = filepath.Abs(*op_path); err != nil {
                panic(err)
            }
            dealWithDir(*op_path, cmd)
        }
    }
}
