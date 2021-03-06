package sysstat

import (
	"io/ioutil"
)

type Load struct {
	Avg1 float64 `json:"avg1"`
	Avg2 float64 `json:"avg2"`
	Avg3 float64 `json:"avg3"`
}

func GetLoad() Load {
	b, err := ioutil.ReadFile("/proc/loadavg")
	if err != nil {
		return Load{}
	}
	return Load{toFloat(b[0:4]), toFloat(b[5:9]), toFloat(b[10:14])}
}
