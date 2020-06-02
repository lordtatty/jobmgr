package jobmgr_test

import (
	"bytes"
	"errors"
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

func TestCsvWriterError(t *testing.T) {
	assert := assert.New(t)
	sut := jobmgr.CsvAuditer{
		Writer: &ErroringWriter{
			Message: "Successful Error Message",
		},
	}
	err := sut.Write(TaskAuditString{
		"one",
		"two",
	})
	assert.EqualError(err, "Successful Error Message")
}

type ErroringWriter struct {
	Message string
}

func (e *ErroringWriter) Write(p []byte) (n int, err error) {
	return 0, errors.New(e.Message)
}
