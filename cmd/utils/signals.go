package utils

import (
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

func CatchInterrupt(tmpPath *string) {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		err := CleanUp(*tmpPath)
		if err != nil {
			fmt.Println(err)
		}
		if runtime.GOOS == "windows" {
			fmt.Println("\n[â] Cancelled!")
		} else {
			fmt.Println("\n[\033[0;31mâ\033[0;m] Cancelled!")
		}
		os.Exit(0)
	}()
}
