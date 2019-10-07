package main

import (
	"log"
	"os/exec"
	"time"
)

func main() {
	cmd := exec.Command("ping", "www.baidu.com")

	cmd.Stdout = log.Writer()
	cmd.Stderr = log.Writer()
	go func() {
		check(cmd)
	}()
	cmd.Run()

	select {}
}

func check(cmd *exec.Cmd) {
	for {
		log.Println(cmd.ProcessState)
		time.Sleep(time.Duration(1) * time.Second)
	}

}
