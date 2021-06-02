package lineprefix

import (
	"io"
	"sync"

	"github.com/fatih/color"
)

// color.Output is a single instance, guard for concurrent calls.
var colorOutputMutex sync.Mutex

// colorWrapper is a wrapper around fatih's color.Color to enable setting
// and unsetting colors for a generic io.Writer.
type colorWrapper struct {
	*color.Color
}

// SetWriter unsets the color on the writer.
func (c *colorWrapper) SetWriter(w io.Writer) {
	colorOutputMutex.Lock()
	defer colorOutputMutex.Unlock()

	tmp := color.Output
	color.Output = w
	c.Set()
	color.Output = tmp
}

// UnsetWriter sets the color on the writer.
func (c *colorWrapper) UnsetWriter(w io.Writer) {
	colorOutputMutex.Lock()
	defer colorOutputMutex.Unlock()

	tmp := color.Output
	color.Output = w
	color.Unset()
	color.Output = tmp
}
