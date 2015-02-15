package sysstat

import (
	"fmt"
	"os"
)

func GetHost() string {
	host, err := os.Hostname()
	if err != nil {
		return fmt.Sprint(err)
	}
	return host
}
