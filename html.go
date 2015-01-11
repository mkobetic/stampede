package main

import (
	"bytes"
	"io"
)

// HTML escaping for writing into element content only!
// See https://www.owasp.org/index.php/XSS_(Cross_Site_Scripting)_Prevention_Cheat_Sheet#RULE_.231_-_HTML_Escape_Before_Inserting_Untrusted_Data_into_HTML_Element_Content
type EscapeWriter struct {
	w io.Writer
}

func EscapeContent(w io.Writer) io.Writer {
	return &EscapeWriter{w}
}

func (w *EscapeWriter) Write(p []byte) (n int, err error) {
	for len(p) != 0 {
		i := bytes.IndexAny(p, `<>"'/`)
		if i == -1 {
			m, err := w.w.Write(p)
			return n + m, err
		}
		if m, err := w.w.Write(p[:i]); err != nil {
			return n, err
		} else {
			n += m
		}
		var esc []byte
		switch p[i] {
		case '<':
			esc = []byte("&lt;")
		case '>':
			esc = []byte("&gt;")
		case '"':
			esc = []byte("&quot;")
		case '\'':
			esc = []byte("&#x27;")
		case '/':
			esc = []byte("&#x2F;")
		}
		if _, err := w.w.Write(esc); err != nil {
			return n, err
		}
		n += 1
		p = p[i+1:]
	}
	return n, nil
}
