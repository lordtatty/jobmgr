package cli

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/lordtatty/jobmgr"
)

type CliLoop struct {
	State      *jobmgr.State
	done       chan chan bool
	lineLength int
}

func NewCliLoop(j *jobmgr.Job) *CliLoop {
	return &CliLoop{
		State: &j.State,
	}
}

func (c *CliLoop) Start(ctx context.Context) {
	c.done = make(chan chan bool)
	nextUpdate := time.After(1 * time.Second)
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("\r")
			c.printProgress()
			return
		case <-nextUpdate:
			nextUpdate = time.After(1 * time.Second)
			fmt.Printf("\r")
			c.printProgress()
		}
	}
}

func (c *CliLoop) printProgress() {
	elapsed := passiveColour(c.State.TimeElapsed())
	expected := expectedColour(fmt.Sprintf("Expected: %d", c.State.Expected))
	progress := expectedColour(fmt.Sprintf("(%d%%)", c.State.Percentage()))
	processed := fmt.Sprintf("Processed: %d", c.State.Completed)
	successes := successColour(fmt.Sprintf("Successes: %d", c.State.Successes))
	failures := errorColour(fmt.Sprintf("Failures: %d", c.State.Failures))
	paused := ""
	if c.State.Paused == true {
		paused = "(PAUSED)"
	}
	m := fmt.Sprintf("%s - %s %s %s %s %s %s", elapsed, expected, progress, processed, successes, failures, paused)
	c.lineLength = len(m)
	fmt.Printf("%-"+strconv.FormatInt(int64(c.lineLength), 10)+"v", m)
}
