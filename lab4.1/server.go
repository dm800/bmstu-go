package main

import (
	"fmt"
	"github.com/gliderlabs/ssh"
	"io"
	"log"
	"os/exec"
	"strings"
)

func main() {
	ssh.Handle(func(s ssh.Session) {

		fmt.Println(s.RawCommand())
		command := strings.Index(s.RawCommand(), " ")
		cmd := exec.Command(s.RawCommand()[:command], s.RawCommand()[command:])
		stdout, err := cmd.Output()
		if err != nil {
			io.WriteString(s, fmt.Sprintf("ERROR OCCURED %s", err.Error()))
		} else {
			io.WriteString(s, fmt.Sprintf(string(stdout)))
		}
	})

	log.Println("starting ssh server on port 9000...")
	log.Fatal(ssh.ListenAndServe(":9000", nil))
}
