package models

type CommandDescription struct {
	Finished   bool                      `json:"finished"`
	Started    bool                      `json:"started"`
	Startedat  string                    `json:"startedat"`
	Finishedat string                    `json:"finishedat"`
	Duration   float64                   `json:"duration"`
	Pid        int                       `json:"pid"`
	Id         string                    `json:"id"`
	Commands   map[string]*CommandStatus `json:"commands"`
	Processes  []*ProcessInfo            `json:"processes"`
}

func NewCommandDescription() *CommandDescription {
	cd := &CommandDescription{}
	return cd
}

func (cd *CommandDescription) IsFinished() bool {
	return cd.Finished
}

func (cd *CommandDescription) SetFinished(isFinished bool) {
	cd.Finished = isFinished
}

func (cd *CommandDescription) IsStarted() bool {
	return cd.Started
}

func (cd *CommandDescription) SetStarted(isStarted bool) {
	cd.Started = isStarted
}

func (cd *CommandDescription) GetStartedAt() string {
	return cd.Startedat
}

func (cd *CommandDescription) SetStartedAt(startedAt string) {
	cd.Startedat = startedAt
}

func (cd *CommandDescription) GetFinishedAt() string {
	return cd.Finishedat
}

func (cd *CommandDescription) SetFinishedAt(finishedAt string) {
	cd.Finishedat = finishedAt
}

func (cd *CommandDescription) GetDuration() float64 {
	return cd.Duration
}

func (cd *CommandDescription) SetDuration(duration float64) {
	cd.Duration = duration
}

func (cd *CommandDescription) GetPid() int {
	return cd.Pid
}

func (cd *CommandDescription) SetPid(pid int) {
	cd.Pid = pid
}

func (cd *CommandDescription) GetId() string {
	return cd.Id
}

func (cd *CommandDescription) SetId(id string) {
	cd.Id = id
}

func (cd *CommandDescription) GetCommands() map[string]*CommandStatus {
	return cd.Commands
}

func (cd *CommandDescription) SetCommands(commands map[string]*CommandStatus) {
	cd.Commands = commands
}

func (cd *CommandDescription) GetProcesses() []*ProcessInfo {
	return cd.Processes
}

func (cd *CommandDescription) SetProcesses(processes []*ProcessInfo) {
	cd.Processes = processes
}
