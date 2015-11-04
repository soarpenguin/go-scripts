package main

import (
    "fmt"
    "flag"
    "runtime"
    "os/exec"
    "os"
)

const (
    VERSION  = "0.1.0"
)

const (
    RetOK   = iota
    RetFail
)

func execCmd(cmd string, shell bool) []byte {
    fmt.Println("command is ", cmd)
    if shell {
        out, err := exec.Command("bash", "-c", cmd).Output()
        if err != nil {
            panic(err)
        }
        return out
    } else {
        out, err := exec.Command(cmd).Output()
        if err != nil {
            panic(err)
        }
        return out
    }
}

func dealWithFile(filename string, cmd string) {
    cmd += filename
    execCmd(cmd, true)
}

func dealWithDir(dir string) {


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
        os.Exit(RetOK)
    }

    switch runtime.GOOS {
        case "windows":
            fmt.Printf("Not supported under windows.\n")
            os.Exit(1)
        case "darwin", "freebsd":
            cmd = "/usr/bin/sed -i \"\" \"s/[ ]*$//g\" "
        default:
            cmd = "/usr/bin/sed -i \"s/[ \t]*$//g\" "
    }

    if *op_path == "" && *op_file == "" {
        fmt.Printf("path or file must provide one.\n")
        os.Exit(RetFail)
    } else if *op_file != "" {
        if _, err := isExists(*op_file); err == nil  {
            dealWithFile(*op_file, cmd)
        }
    } else if *op_path != "" {

    }
}
