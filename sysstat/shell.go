package sysstat

import (
	"log"
	"os/exec"
)

func ExecShell(s string) string {
	out, err := exec.Command("/bin/sh", "-c", s).Output()
	if err != nil {
		log.Println(err)
	}

	return string(out)
}
