package jobmgr_test

import (
	"bytes"
	"io"
	"sync"
	"testing"
	"time"

	"github.com/lordtatty/jobmgr"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

func TestNew(t *testing.T) {
	assert := assert.New(t)
	m := jobmgr.New()
	assert.IsType(&jobmgr.Job{}, m)
	assert.IsType(jobmgr.State{}, m.State)
	assert.IsType(time.Second*0, m.Sleep)
}

type testTask struct {
	Result jobmgr.TaskResult
	WasRun bool
}

func (t *testTask) Run() jobmgr.TaskResult {
	t.WasRun = true
	return t.Result
}

type testTaskNoAudit struct {
	WasRun bool
}

func (t *testTaskNoAudit) Run() jobmgr.TaskResult {
	t.WasRun = true
	return nil
}

func TestStartHappyPath(t *testing.T) {
	assert := assert.New(t)
	m := jobmgr.New()
	ctx := context.Background()
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		m.Start(ctx)
		wg.Done()
	}()
	task := &testTaskNoAudit{}
	task2 := &testTaskNoAudit{}
	m.AddTask(task)
	m.AddTask(task2)
	m.NoMoreTasks()
	wg.Wait()
	assert.Equal(true, task.WasRun)
	assert.Equal(true, task2.WasRun)
}

func TestPauseResume(t *testing.T) {
	assert := assert.New(t)
	m := jobmgr.New()
	ctx := context.Background()
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		m.Start(ctx)
		wg.Done()
	}()
	task := &testTaskNoAudit{}
	task2 := &testTaskNoAudit{}
	m.AddTask(task)
	time.Sleep(1 * time.Second)
	assert.Equal(true, task.WasRun)
	m.Pause()
	m.AddTask(task2)
	assert.Equal(false, task2.WasRun)
	m.Resume()
	m.NoMoreTasks()
	wg.Wait()
	assert.Equal(true, task2.WasRun)
}

func TestStartHappyPathWithAudit(t *testing.T) {
	assert := assert.New(t)

	// Configure CSV auditer with result buf
	resultBuf := &bytes.Buffer{}
	auditer := &jobmgr.CsvAuditer{
		Writer: io.Writer(resultBuf),
	}

	// Set up tasks and task results
	result := &jobmgr.Result{}
	result.AddValue("KeyOne", "ValueOne")
	result.AddValue("KeyTwo", "ValueOne")
	result2 := &jobmgr.Result{}
	result2.AddValue("KeyOne", "ValueTwo")
	result2.AddValue("KeyTwo", "ValueTwo")
	task := &testTask{}
	task.Result = result
	task2 := &testTask{}
	task2.Result = result2

	// Run tasks
	m := jobmgr.New()
	m.AddAuditer(auditer)
	ctx := context.Background()
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		m.Start(ctx)
		wg.Done()
	}()
	m.AddTask(task)
	m.AddTask(task2)
	m.NoMoreTasks()
	wg.Wait()

	// Check output is correct
	expected := `KeyOne,KeyTwo
ValueOne,ValueOne
ValueTwo,ValueTwo
`
	assert.Equal(expected, resultBuf.String())
}

func TestCtxCancel(t *testing.T) {
	assert := assert.New(t)
	m := jobmgr.New()
	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		m.Start(ctx)
		wg.Done()
	}()
	cancel()
	timeout := waitTimeout(wg, time.Second*5)
	assert.Equal(false, timeout, "context did not cancel in a reasonable time")
}

// Thanks to https://stackoverflow.com/questions/32840687/timeout-for-waitgroup-wait
// waitTimeout waits for the waitgroup for the specified max timeout.
// Returns true if waiting timed out.
func waitTimeout(wg *sync.WaitGroup, timeout time.Duration) bool {
	c := make(chan struct{})
	go func() {
		defer close(c)
		wg.Wait()
	}()
	select {
	case <-c:
		return false // completed normally
	case <-time.After(timeout):
		return true // timed out
	}
}
