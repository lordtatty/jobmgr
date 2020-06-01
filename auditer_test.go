package jobmgr_test

import (
	"bytes"
	"testing"

	"github.com/lordtatty/jobmgr"
	"github.com/stretchr/testify/assert"
)

type TaskAuditString []string

func (t TaskAuditString) Audit() []string {
	return t
}

func (t TaskAuditString) AuditHeaders() []string {
	h := []string{}
	for _, v := range t {
		h = append(h, "header-"+v)
	}
	return h
}

func TestCsvWriter(t *testing.T) {
	assert := assert.New(t)
	var buf bytes.Buffer
	sut := jobmgr.CsvAuditer{
		Writer: &buf,
	}
	sut.Write(TaskAuditString{
		"one",
		"two",
	})
	sut.Write(TaskAuditString{
		"three",
		"four",
	})
	assert.Equal("header-one,header-two\none,two\nthree,four\n", buf.String())
}
