package main

import (
	"bytes"
	"testing"
)

func TestEscaper(t *testing.T) {
	for _, s := range []struct{ in, out string }{
		{`"John Doe"<john@doe.com>`, `&quot;John Doe&quot;&lt;john@doe.com&gt;`},
	} {
		var b bytes.Buffer
		n, err := EscapeContent(&b).Write([]byte(s.in))
		if err != nil {
			t.Error(err)
		} else if n != len(s.in) {
			t.Fail()
		} else if (&b).String() != s.out {
			t.Error((&b).String())
		}
	}
}
