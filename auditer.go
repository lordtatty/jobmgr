package jobmgr

import (
	"encoding/csv"
	"io"
	"strings"
	"sync"
)

type Auditer interface {
	Write(t TaskAudit) error
}

type CsvAuditer struct {
	Writer    io.Writer
	w         *csv.Writer
	setupOnce sync.Once
}

func (c *CsvAuditer) Write(t TaskAudit) error {
	var err error
	c.setupOnce.Do(func() {
		csv.NewWriter(c.Writer)
		c.w = csv.NewWriter(c.Writer)
		err = c.writeln(t.AuditHeaders())
	})
	if err != nil {
		return err
	}
	return c.writeln(t.Audit())
}

func (c *CsvAuditer) writeln(l []string) error {
	values := make([]string, 0, len(l))
	for _, val := range l {
		safeVal := strings.Replace(val, "\n", "", -1)
		values = append(values, safeVal)
	}
	err := c.w.Write(values)
	if err != nil {
		return err
	}
	c.w.Flush()
	return c.w.Error()
}
