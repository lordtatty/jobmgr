package jobmgr_test

import (
	"testing"

	"github.com/lordtatty/jobmgr"
	"github.com/stretchr/testify/assert"
)

func TestNewResult(t *testing.T) {
	assert := assert.New(t)
	r := jobmgr.NewResult()
	assert.IsType(jobmgr.Result{}, r)
}

func TestResultSuccess(t *testing.T) {
	assert := assert.New(t)
	r := jobmgr.NewResult()
	assert.False(r.Success())
	r.SetSuccess(true)
	assert.True(r.Success())
	r.SetSuccess(false)
	assert.False(r.Success())
}

func TestResultAuditHeaders(t *testing.T) {
	assert := assert.New(t)
	r := jobmgr.NewResult()
	r.AddValue("KeyOne", "ValueOne")
	r.AddValue("KeyTwo", "ValueTwo")
	expected := []string{"KeyOne", "KeyTwo"}
	assert.Equal(expected, r.AuditHeaders())
}

func TestResultAudit(t *testing.T) {
	assert := assert.New(t)
	r := jobmgr.NewResult()
	r.AddValue("KeyOne", "ValueOne")
	r.AddValue("KeyTwo", "ValueTwo")
	expected := []string{"ValueOne", "ValueTwo"}
	assert.Equal(expected, r.Audit())
}
