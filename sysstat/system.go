package sysstat

type System struct {
	Host string  `json:"host"`
	Disk Disk    `json:"disk"`
	Cpu  CpuInfo `json:"cpuinfo"`
	Load Load    `json:"load"`
	Ram  Ram     `json:"ram"`
	Time string  `json:"time"`
}

func GetSystem(s string) System {
	return System{GetHost(), GetDisk(s), GetCpuinfo(), GetLoad(), GetRam(), GetNow()}
}
