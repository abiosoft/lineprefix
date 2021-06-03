// Package lineprefix provides a `io.Writer` wrapper with line prefix and color customizations.
//
// Static Prefix
//  prefix := lineprefix.Prefix("app |")
//  writer := lineprefix.New(prefix)
//
//  fmt.Fprintln(writer, "hello world")
//
// Dynamic Prefix
//  now := func() string { return time.Now().UTC().Format("2006-01-02 15:04:05") }
//  prefix := lineprefix.PrefixFunc(now)
//  writer := lineprefix.New(prefix)
//
//  for i := 0; i<3; i++ {
//      fmt.Fprintln(writer, "hello world")
//      time.Sleep(time.Second)
//  }
//
// Colors can be added with github.com/faith/color
//
// A prefix with blue text on white background
//  blue := color.New(color.FgBlue, color.BgWhite).SprintFunc()
//  prefix := lineprefix.Prefix(blue("app"))
//  writer := lineprefix.New(prefix)
//
//  fmt.Fprintln(writer, "this outputs blue on white prefix text")
//
// A blue colored output
//  blue := color.New(color.FgBlue)
//  option := lineprefix.Color(blue)
//  writer := lineprefix.New(option)
//
//  fmt.Fprintln(writer, "this outputs blue color text")
package lineprefix
