package jobmgr

type Result struct {
	data    map[string]string
	success bool
}

func NewResult() *Result {
	return &Result{
		data: make(map[string]string),
	}
}

func (r *Result) Success() bool {
	return r.success
}

func (r *Result) SetSuccess(success bool) {
	r.success = success
}

func (r *Result) AddValue(key, value string) {
	r.data[key] = value
}

func (r *Result) AuditHeaders() []string {
	headers := make([]string, len(r.data))
	i := 0
	for k := range r.data {
		headers[i] = k
		i++
	}
	return headers
}

func (r *Result) Audit() []string {
	values := make([]string, len(r.data))
	i := 0
	for _, v := range r.data {
		values[i] = v
		i++
	}
	return values
}
