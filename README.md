A Simple golang Job Manager

!!Consider this in alpha at the moment, subject to breaking changes as it gets used in real-life projects.  But the intention is to stabilise quickly.

Create a simple task
```go
type LogTask struct {
    Message string
}

func (l *LogTask) Run() jobmgr.TaskResult {
	log.Println(l.Message)
    r := jobmgr.NewResult()
    r.AddValue("message", message)
    r.AddValye("via", "logger")
    r.SetSuccess(true)
    return r
}
```

A Basic happy path with the manager
```go
m := jobmgr.New()
ctx := context.Background()
wg := sync.WaitGroup{}
wg.Add(1)
go func() {
    m.Start(ctx)
    wg.Done()
}()
task := &LogTask{}
task2 := &LogTask{}
m.AddTask(task)
m.AddTask(task2)
// Tell the manager that it should not expect any more tasks
// It will close itself down once it has completed
// All those on it's queue
m.NoMoreTasks() 
wg.Wait() // Will continue once the job manager is done
```

Closing the manager via context

```go
m := jobmgr.New()
ctx, cancel := context.WithCancel(context.Background())
wg := sync.WaitGroup{}
wg.Add(1)
go func() {
    m.Start(ctx)
    wg.Done()
}()
cancel()
wg.Wait() // Will continue once the job manager has exited
```