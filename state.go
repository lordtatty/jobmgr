package jobmgr

import (
	"fmt"
	"math"
	"time"
)

type TaskAudit interface {
	Audit() []string
	AuditHeaders() []string
}

type TaskResult interface {
	TaskAudit
	Success() bool
}

type State struct {
	StartTime time.Time
	Paused    bool
	Expected  uint32
	Completed uint32
	Successes uint32
	Failures  uint32
}

func (j *State) Start() {
	j.StartTime = time.Now()
}

func (j *State) SetPaused(p bool) {
	j.Paused = p
}

func (j *State) Percentage() int {
	return int((float64(j.Completed) / float64(j.Expected)) * 100)
}

func (j *State) TimeElapsed() string {
	elapsed := time.Since(j.StartTime)
	seconds := int64(math.Floor(elapsed.Seconds()))
	sRemainder := seconds % 60
	minutes := int64(math.Floor(elapsed.Minutes()))
	return fmt.Sprintf("%02dm:%02ds", minutes, sRemainder)
}

func (j *State) Update(r TaskResult) {
	j.Completed++
	if r != nil && r.Success() == true {
		j.Successes++
		return
	}
	j.Failures++
}
