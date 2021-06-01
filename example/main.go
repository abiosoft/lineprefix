package main

import (
	"fmt"
	"io"
	"os/exec"
	"sync"
	"time"

	"github.com/abiosoft/lineprefix"
	"github.com/fatih/color"
)

func stage(name string) {
	fmt.Println()
	fmt.Println("----" + name + "----")
}

func now() string { return time.Now().UTC().Format("2006-01-02 15:04:05") }

func main() {
	// static prefix
	{
		stage("static prefix")
		writer := lineprefix.New(lineprefix.Prefix("prefix"))
		fmt.Fprintln(writer, "hello world with static prefix")
		fmt.Fprintln(writer, "another line with static prefix")
	}

	// dynamic prefix
	{
		stage("dynamic prefix")
		writer := lineprefix.New(lineprefix.PrefixFunc(func() string {
			return "prefix " + now()
		}))

		for i := 0; i < 5; i++ {
			fmt.Fprintln(writer, "hello world with dynamic prefix")
			time.Sleep(time.Second * 1)
		}
	}

	// colors
	{
		stage("colors")
		colorWriter := func(c *color.Color, prefix string) io.Writer {
			prefix = c.Sprintf(" %s ", prefix)
			return lineprefix.New(lineprefix.Prefix(prefix))
		}
		blue := colorWriter(color.New(color.BgBlue, color.FgWhite), "blue")
		fmt.Fprintln(blue, "hello world with colored prefix")

		red := colorWriter(color.New(color.BgRed, color.FgWhite), "red ")
		fmt.Fprintln(red, "hello world with colored prefix")
	}

	// outputs for multiple commands
	{
		ping := func(wg *sync.WaitGroup, writer io.Writer, domain string) {
			cmd := exec.Command("ping", "-c", "5", domain)
			cmd.Stdout = writer

			cmd.Run()
			wg.Done()
		}
		var wg sync.WaitGroup

		stage("multiple outputs and commands")
		{
			wg.Add(2)
			gprefix := lineprefix.PrefixFunc(func() string { return " GOOG " + now() })
			gwriter := lineprefix.New(gprefix)
			go ping(&wg, gwriter, "google.com")

			aprefix := lineprefix.PrefixFunc(func() string { return " AAPL " + now() })
			awriter := lineprefix.New(aprefix)
			go ping(&wg, awriter, "apple.com")

			wg.Wait()
		}

		stage("multiple output with partial prefix color")
		{
			wg.Add(2)
			blue := color.New(color.BgBlue, color.FgWhite).SprintFunc()
			gwriter := lineprefix.New(
				lineprefix.Prefix(blue(" GOOG ")+" "),
				lineprefix.PrefixFunc(func() string { return now() }),
			)
			go ping(&wg, gwriter, "google.com")

			red := color.New(color.BgRed, color.FgWhite).SprintFunc()
			awriter := lineprefix.New(
				lineprefix.Prefix(red(" AAPL ")+" "),
				lineprefix.PrefixFunc(func() string { return now() }),
			)
			go ping(&wg, awriter, "apple.com")

			wg.Wait()
		}

		stage("multiple output with mixed prefix color and line color")
		{
			wg.Add(2)
			blue := color.New(color.BgBlue, color.FgWhite).SprintFunc()
			bluefg := color.New(color.FgBlue)
			gwriter := lineprefix.New(
				lineprefix.Prefix(blue(" GOOG ")+" "),
				lineprefix.PrefixFunc(func() string { return bluefg.Sprint(now()) }),
				lineprefix.Color(bluefg),
			)
			go ping(&wg, gwriter, "google.com")

			red := color.New(color.BgRed, color.FgWhite).SprintFunc()
			redfg := color.New(color.FgRed)
			awriter := lineprefix.New(
				lineprefix.Prefix(red(" AAPL ")+" "),
				lineprefix.PrefixFunc(func() string { return redfg.Sprint(now()) }),
				lineprefix.Color(redfg),
			)
			go ping(&wg, awriter, "apple.com")

			wg.Wait()
		}
	}
}
