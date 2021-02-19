package models

type ProcessInfo struct {
	Pid  int    `json:"pid"`
	PPid int    `json:"parent"`
	Name string `json:"name"`
}

func NewProcessInfo() *ProcessInfo {
	pi := &ProcessInfo{}
	return pi
}

func (pi *ProcessInfo) GetPid() int {
	return pi.Pid
}

func (pi *ProcessInfo) SetPid(pid int) {
	pi.Pid = pid
}

func (pi *ProcessInfo) GetPPid() int {
	return pi.PPid
}

func (pi *ProcessInfo) SetPPid(ppid int) {
	pi.PPid = ppid
}

func (pi *ProcessInfo) GetName() string {
	return pi.Name
}

func (pi *ProcessInfo) SetName(name string) {
	pi.Name = name
}
