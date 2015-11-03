package main

import (
    "fmt"
    "flag"
    "runtime"
    "strings"
    "os/exec"
    "os"
    "sync"
)

func exe_cmd(cmd string, wg *sync.WaitGroup) {
    fmt.Println("command is ",cmd)
    parts := strings.Fields(cmd)
    head := parts[0]
    parts = parts[1:len(parts)]
    
    out, err := exec.Command(head,parts...).Output()
    if err != nil {
        fmt.Printf("%s", err)
    }
    fmt.Printf("%s", out)
    wg.Done() // Need to signal to waitgroup that this goroutine is done
}

func deal_with_file(filename string) {
    wg := new(sync.WaitGroup)
    wg.Add(1)
  
    
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

    op_path := flag.String("d", "./test.log", "Dir to recursive remove spaces at the end of the line.")
    op_file := flag.String("f", "", "File name for remove spaces at the end of line.")
    nowait_flag := flag.Bool("W", false, "Do not wait until operation completes")

    flag.Parse()

    var cmd string = flag.Arg(0);

    fmt.Printf("action   : %s\n", cmd)
    fmt.Printf("data path: %s\n", *op_path)
    fmt.Printf("log file : %s\n", *op_file)
    fmt.Printf("nowait   : %v\n", *nowait_flag)

    fmt.Printf("-------------------------------------------------------\n")

    fmt.Printf("there are %d non-flag input param\n", flag.NArg())
    for i, param := range flag.Args(){
        fmt.Printf("#%d    :%s\n", i, param)
    }

    switch runtime.GOOS {
        case "windows":
            fmt.Printf("Not supported under windows.\n")
            os.Exit(1)
        case "darwin", "freebsd":
            cmd = "sed -i \"\" \"s/[ ]*$//g\"" 
        default:
            cmd = "sed -i \"s/[ \t]*$//g\"" 
    }
    
    if _, err := isExists(*op_path); err != nil {
        fmt.Printf("%s\n", err)
    }
}
