package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	go signalListen()
	time.Sleep(time.Hour)
}

func signalListen() {
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		fmt.Println("get signal:", s)
		os.Exit(1)
	}
}
