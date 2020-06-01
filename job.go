package jobmgr

import (
	"context"
	"log"
	"time"
)

type Task interface {
	Run() TaskResult
}

type Job struct {
	State     State
	taskQueue chan Task
	pauseCh   chan chan bool
	resumeCh  chan bool
	auditers  []Auditer
	Sleep     time.Duration
}

func New() *Job {
	s := State{}
	return &Job{
		State: s,
		//TODO this 100 is somewhat arbitrary... should the user be able to set?
		taskQueue: make(chan Task, 100),
		pauseCh:   make(chan chan bool),
	}
}

func (j *Job) AddAuditer(a Auditer) {
	j.auditers = append(j.auditers, a)
}

func (j *Job) AddTask(t Task) {
	j.taskQueue <- t
}

func (j *Job) Pause() {
	confirmPaused := make(chan bool)
	j.pauseCh <- confirmPaused
	// We must wait for confirmPaused
	// So that this function can return
	// Synchronously (ie. when it's actually paused)
	<-confirmPaused
}

func (j *Job) Resume() {
	close(j.resumeCh)
}

func (j *Job) NoMoreTasks() {
	close(j.taskQueue)
}

func (j *Job) Start(ctx context.Context) error {
	j.State.Start()
	paused := false
	confirmPaused := make(chan bool)
	checkForPause := make(chan bool)
	go func() {
		for {
			select {
			case resp := <-j.pauseCh:
				j.resumeCh = make(chan bool)
				confirmPaused = resp
				paused = true
				j.State.SetPaused(true)
				checkForPause <- true
				<-j.resumeCh
				paused = false
				j.State.SetPaused(false)
			}
		}
	}()
	for {
		if paused {
			confirmPaused <- true
			<-j.resumeCh
		}
		select {
		case <-ctx.Done():
			log.Println("here")
			return nil
		case t, more := <-j.taskQueue:
			if !more {
				return nil
			}
			result := t.Run()
			j.State.Update(result)
			for _, a := range j.auditers {
				a.Write(result)
			}
			time.Sleep(j.Sleep)
		case <-checkForPause:
		}
	}
}
