package jobmgr_test

import (
	"testing"
	"time"

	"github.com/lordtatty/jobmgr"
	"github.com/stretchr/testify/assert"
)

func TestStateTimeElapsed(t *testing.T) {
	assert := assert.New(t)
	s := jobmgr.State{}
	s.Start()
	time.Sleep(time.Second * 1)
	assert.Equal("00m:01s", s.TimeElapsed())
}

func TestStateUpdate(t *testing.T) {
	assert := assert.New(t)
	s := jobmgr.State{}
	s.Start()
	r := &jobmgr.Result{}
	r.SetSuccess(true)
	s.Update(r)
	assert.Equal(uint32(1), s.Completed)
	assert.Equal(uint32(1), s.Successes)
	assert.Equal(uint32(0), s.Failures)
	r = &jobmgr.Result{}
	r.SetSuccess(false)
	s.Update(r)
	assert.Equal(uint32(2), s.Completed)
	assert.Equal(uint32(1), s.Successes)
	assert.Equal(uint32(1), s.Failures)
	r = &jobmgr.Result{}
	r.SetSuccess(true)
	s.Update(r)
	assert.Equal(uint32(3), s.Completed)
	assert.Equal(uint32(2), s.Successes)
	assert.Equal(uint32(1), s.Failures)
}

func TestPercentaage(t *testing.T) {
	assert := assert.New(t)
	s := jobmgr.State{}
	s.Completed = 10
	s.Expected = 50
	assert.Equal(20, s.Percentage())
}

func TestSetPaused(t *testing.T) {
	assert := assert.New(t)
	s := jobmgr.State{}
	s.SetPaused(true)
	assert.Equal(true, s.Paused)
	s.SetPaused(false)
	assert.Equal(false, s.Paused)
}
