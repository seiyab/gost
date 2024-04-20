package testdata

import "io"

func _(s []io.WriteCloser, m map[string]io.WriteCloser) {
	sw := s[0]
	sw.Write(nil)

	for _, w := range s {
		w.Write(nil)
	}

	mw := m["key"]
	mw.Write(nil)

	for _, w := range m {
		w.Write(nil)
	}
}
