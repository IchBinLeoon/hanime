package utils

import (
	"fmt"
	"os"
	"os/signal"
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
		fmt.Println("\nCancelled")
		os.Exit(0)
	}()
}
