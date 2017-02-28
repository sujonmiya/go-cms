package models

type Metadata map[string]string

func (m Metadata) Get(key string) string {
	if v, ok := m[key]; !ok {
		return ""
	} else {
		return v
	}
}

func (m Metadata) Put(k, v string) {
	m[k] = v
}
