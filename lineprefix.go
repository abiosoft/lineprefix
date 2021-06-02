package lineprefix

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/fatih/color"
)

// New creates a new lineprefix writer with options.
func New(options ...Option) io.WriteCloser {
	var l lineWriter
	for _, option := range options {
		option.apply(&l)
	}

	// default to stdout if not set
	if l.out == nil {
		l.out = os.Stdout
	}
	return &l
}

// Writer sets the writer to use. The default writer is os.Stdout.
func Writer(w io.Writer) Option {
	return optionFunc(func(l *lineWriter) {
		l.out = w
	})
}

// Prefix sets the prefix to use.
// Can be called multiple times to set multiple prefixes.
func Prefix(s string) Option {
	return optionFunc(func(l *lineWriter) {
		l.prefixes = append(l.prefixes, func() string { return s })
	})
}

// PrefixFunc is like prefix but with the ability to make it dynamic.
// Can be called multiple times to set multiple prefixes.
func PrefixFunc(f func() string) Option {
	return optionFunc(func(l *lineWriter) {
		l.prefixes = append(l.prefixes, f)
	})
}

// Color sets the colour for the line outputs excluding the prefix.
func Color(c *color.Color) Option {
	return optionFunc(func(l *lineWriter) {
		l.color = c
	})
}

// RenderEscaped (if true) enables the rendering of escaped whitespace characters.
// e.g. \\t appears as tab instead of \t, \\n appears as newline instead of \n etc.
func RenderEscaped(b bool) Option {
	return optionFunc(func(l *lineWriter) {
		l.renderEscaped = true
	})
}

type optionFunc func(*lineWriter)

func (o optionFunc) apply(l *lineWriter) { o(l) }

// Option is the configuration option for a new instance of lineprefix writer.
type Option interface {
	apply(*lineWriter)
}

var _ io.Writer = (*lineWriter)(nil)
var _ io.Closer = (*lineWriter)(nil)

type prefix func() string
type prefixes []prefix

func (p prefixes) Bytes() []byte {
	var b bytes.Buffer
	for _, prefix := range p {
		fmt.Fprint(&b, prefix())
	}
	// add a traling space char
	fmt.Fprint(&b, " ")
	return b.Bytes()
}

// lineWriter is a simple writer that only writes to an underlying writer when
// a newline is encountered.
type lineWriter struct {
	out      io.Writer
	prefixes prefixes
	color    *color.Color

	renderEscaped bool

	buf bytes.Buffer
	sync.Mutex

	open   bool
	closed bool
}

func (l *lineWriter) Write(b []byte) (int, error) {
	l.Lock()
	defer l.Unlock()
	if l.closed {
		return 0, io.ErrClosedPipe
	}

	if !l.open {
		l.open = true

		// enable color for the first line
		if l.color != nil {
			l.color.SetWriter(&l.buf)
		}
	}

	for i := 0; i < len(b); i++ {

		if l.renderEscaped {
			// special case: replace escaped chars with their real value
			// newline, tab, quote, backslack
			if b[i] == '\\' {
				// peek if available
				if i+1 < len(b) {
					i++
					switch b[i] {
					case 'n':
						b[i] = '\n'
					case 't':
						b[i] = '\t'

					// do nothing for these, escape char already skipped
					case '\\':
					case '"':
					case '\'':

					// otherwise, don't skip escape char
					default:
						i--
					}
				}
			}
		}

		eol := b[i] == '\n' // end of line

		if eol && l.color != nil {
			// reset color
			l.color.UnsetWriter(&l.buf)
		}

		// cache the char
		l.buf.WriteByte(b[i])

		// write to underlying writer if newline is encountered
		if eol {
			// write prefixes
			_, err := l.out.Write(l.prefixes.Bytes())
			if err != nil {
				return 0, err
			}

			// write line
			n, err := l.buf.WriteTo(l.out)
			if err != nil {
				return int(n), err
			}

			// truncate the buffer
			l.buf.Truncate(0)

			// enable color for the next line
			if l.color != nil {
				l.color.SetWriter(&l.buf)
			}
		}
	}

	// if it gets here, all bytes are successfully written
	return len(b), nil
}

func (l *lineWriter) Close() error {
	l.Lock()
	defer l.Unlock()

	l.closed = true
	_, err := l.buf.WriteTo(l.out)
	return err
}
