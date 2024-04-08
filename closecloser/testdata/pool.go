package testdata

type manager struct {
	pool []resource
}

func (m *manager) Close() error {
	for i := range m.pool {
		m.pool[i].Close()
	}
	return nil
}

func (m *manager) Resource() *resource {
	return &m.pool[0]
}

type resource struct{}

func (r *resource) Close() error { return nil }

func (r *resource) Method() {}

func _() {
	var m manager // want ".+"
	r := m.Resource()

	r.Method()
}
