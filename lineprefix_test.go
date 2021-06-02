package lineprefix

import (
	"bytes"
	"fmt"
	"io"
	"testing"
)

func TestPrefix(t *testing.T) {
	tests := []struct {
		name     string
		opts     []Option
		run      func(w io.Writer)
		expected string
	}{
		{
			name: "prefix",
			opts: []Option{Prefix("hello")},
			run: func(w io.Writer) {
				fmt.Fprintln(w, "world")
			},
			expected: "hello world\n",
		},
		{
			name: "no prefix",
			run: func(w io.Writer) {
				fmt.Fprintln(w, "hello world")
			},
			expected: "hello world\n",
		},
		{
			name: "render escaped false",
			run: func(w io.Writer) {
				fmt.Fprintln(w, `\tsome\nand more in \"quote\" and \'single\'`)
			},
			expected: `\tsome\nand more in \"quote\" and \'single\'` + "\n",
		},
		{
			name: "render escaped true",
			opts: []Option{RenderEscaped(true)},
			run: func(w io.Writer) {
				fmt.Fprintln(w, `\tsome\nand more in \"quote\" and \'single\'`)
			},
			expected: "\tsome\nand more in \"quote\" and 'single'\n",
		},
		{
			name: "render escaped true multiple writes",
			opts: []Option{RenderEscaped(true)},
			run: func(w io.Writer) {
				fmt.Fprint(w, `\tsome\nand more in \"quote\" and \'single\`)
				fmt.Fprint(w, `' and another new\`)
				fmt.Fprintln(w, `n line`)
			},
			expected: "\tsome\nand more in \"quote\" and 'single' and another new\n line\n",
		},
		{
			name: "close and flush",
			opts: []Option{Prefix("this")},
			run: func(w io.Writer) {
				fmt.Fprintln(w, "is it")
				fmt.Fprint(w, "should appear also")
				w.(io.Closer).Close()
			},
			expected: "this is it\nshould appear also",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var b bytes.Buffer
			writer := New(append(tt.opts, Writer(&b))...)
			tt.run(writer)

			if tt.expected != b.String() {
				t.Errorf("got: %v, want: %v", b.String(), tt.expected)
			}
		})
	}
}
