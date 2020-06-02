package jobmgr

type Result struct {
	values  []string
	headers []string
	success bool
}

func (r *Result) Success() bool {
	return r.success
}

func (r *Result) SetSuccess(success bool) {
	r.success = success
}

func (r *Result) AddValue(key, value string) {
	r.headers = append(r.headers, key)
	r.values = append(r.values, value)
}

func (r *Result) AuditHeaders() []string {
	return r.headers
}

func (r *Result) Audit() []string {
	return r.values
}
