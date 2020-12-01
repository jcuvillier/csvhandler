package csvhandler

type Record struct {
	headers map[string]int
	values  []string
}

func (r Record) Get(key string) string {
	i, ok := r.headers[key]
	if !ok {
		return ""
	}
	if i < 0 || i > len(r.values) {
		return ""
	}
	return r.values[i]
}
