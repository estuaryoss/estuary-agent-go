package models

type CommandDetails struct {
	Out  string `json:"out"`
	Err  string `json:"err"`
	Pid  int    `json:"pid"`
	Code int    `json:"code"`
	Args string `json:"args"`
}

func NewCommandDetails() *CommandDetails {
	cd := &CommandDetails{}
	return cd
}

func (cd *CommandDetails) GetOut() string {
	return cd.Out
}

func (cd *CommandDetails) SetOut(out string) {
	cd.Out = out
}

func (cd *CommandDetails) GetErr() string {
	return cd.Err
}

func (cd *CommandDetails) SetErr(err string) {
	cd.Err = err
}

func (cd *CommandDetails) GetPid() int {
	return cd.Pid
}

func (cd *CommandDetails) SetPid(pid int) {
	cd.Pid = pid
}

func (cd *CommandDetails) GetCode() int {
	return cd.Code
}

func (cd *CommandDetails) SetCode(code int) {
	cd.Code = code
}

func (cd *CommandDetails) GetArgs() string {
	return cd.Args
}

func (cd *CommandDetails) SetArgs(args string) {
	cd.Args = args
}
