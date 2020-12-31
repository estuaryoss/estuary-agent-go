package models

type CommandStatus struct {
	Status         string          `json:"status"`
	CommandDetails *CommandDetails `json:"details"`
	Startedat      string          `json:"startedat"`
	Finishedat     string          `json:"finishedat"`
	Duration       float64         `json:"duration"`
}

func NewCommandStatus() *CommandStatus {
	cs := &CommandStatus{}
	return cs
}

func (cs *CommandStatus) GetStatus() string {
	return cs.Status
}

func (cs *CommandStatus) SetStatus(status string) {
	cs.Status = status
}

func (cs *CommandStatus) GetCommandDetails() *CommandDetails {
	return cs.CommandDetails
}

func (cs *CommandStatus) SetCommandDetails(cd *CommandDetails) {
	cs.CommandDetails = cd
}

func (cs *CommandStatus) GetStartedAt() string {
	return cs.Startedat
}

func (cs *CommandStatus) SetStartedAt(startedAt string) {
	cs.Startedat = startedAt
}

func (cs *CommandStatus) GetFinishedAt() string {
	return cs.Finishedat
}

func (cs *CommandStatus) SetFinishedAt(finishedAt string) {
	cs.Finishedat = finishedAt
}

func (cs *CommandStatus) GetDuration() float64 {
	return cs.Duration
}

func (cs *CommandStatus) SetDuration(duration float64) {
	cs.Duration = duration
}
