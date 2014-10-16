package cli

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// CLI represents the command line interpreter
type CLI struct {
	prompt      string
	consoleMode uint32
	hist        []string
	histSize    uint8
	commands    map[string]func(*CLI, []string) error
}

// Prompt writes the given prompt, or if no prompt is given, the default prompt to stderr
func (c CLI) Prompt() {
	fmt.Fprintf(os.Stderr, "%s", c.prompt)
}

// SetPrompt set the new default prompt
func (c *CLI) SetPrompt(prompt string) string {
	previous := c.prompt
	c.prompt = prompt
	return previous
}

func (c CLI) NewLine() {
	fmt.Fprintln(os.Stdout)
}

// Info writes a message to stdout
func (c CLI) Info(format string, args ...interface{}) {
	fmt.Fprintf(os.Stdout, format, args...)
}

// Error writes a message to stderr
func (c CLI) Error(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format, args...)
}

// Readline reads a line from stdin
func (c *CLI) ReadLine() string {
	bio := bufio.NewReader(os.Stdin)
	line, _, _ := bio.ReadLine()
	return string(line)
}

// InputString writes a label to stdout and reads the entered string from stdin
// It returns the entered string
func (c *CLI) InputString(label string, values []string) string {
	for {
		c.Info("%s: ", label)
		value := c.ReadLine()
		if len(value) > 0 {
			if len(values) == 0 {
				return value
			} else {
				for _, v := range values {
					if value == v {
						return value
					}
				}
			}
		}
	}
}

// YesNo writes a question to stdout and compares the answer agains
// values. If values is empty, YesNo evaluates always to true.
// It returns true on positiv match, otherwise false
func (c *CLI) YesNo(label string, values []string) bool {
	for {
		c.Info("%s: ", label)
		value := c.ReadLine()
		if len(value) > 0 {
			if len(values) == 0 {
				return true
			} else {
				for _, v := range values {
					if value == v {
						return true
					}
				}
				return false
			}
		}
	}
}

// Choice is a simple menu
// It returns the entered key
func (c *CLI) Choice(items map[string]string) string {
	label := ""
	for short, long := range items {
		label = fmt.Sprintf("%s%s  - %s\n", label, short, long)
	}
	for {
		c.Info("%s>>> ", label)
		value := c.ReadLine()
		if len(value) > 0 {
			if _, ok := items[value]; ok {
				return value
			}
		}
	}
}

// Exec executes a command
func (c *CLI) Exec(line string) {
	cmd := strings.Split(line, " ")
	// catch builtins
	switch cmd[0] {
	case "hist":
		c.history()
	case "run":
		c.run(cmd)
	default:
		c.command(cmd)
	}
}

// command runs the command, on success the command will be added to the history
func (c *CLI) command(cmd []string) {
	if len(cmd[0]) > 0 {
		if f, ok := c.commands[cmd[0]]; ok {
			if err := f(c, cmd); err != nil {
				c.Error("ERROR: %s", err.Error())
			} else {
				// add to history if execution was successful
				c.historyAdd(strings.Join(cmd, " "))
			}
		} else {
			c.Error("Unknown command: %s\n", strings.Join(cmd, " "))
		}
	}
}

// historyAdd adds a command to the history
func (c *CLI) historyAdd(cmd string) {
	if len(c.hist) < cap(c.hist) {
		c.hist = c.hist[0 : len(c.hist)+1]
	} else {
		copy(c.hist, c.hist[1:])
	}
	c.hist[len(c.hist)-1] = cmd
}

// history writes the history to stdout
func (c *CLI) history() {
	for i, v := range c.hist {
		c.Info("%5d %s\n", i, v)
	}
}

// run executes a selected command from the history
func (c *CLI) run(cmd []string) {
	if len(cmd) == 2 {
		i, err := strconv.ParseUint(cmd[1], 10, 8)
		if err == nil && i < uint64(len(c.hist)) {
			c.Exec(c.hist[i])
			return
		}
	}
	c.Error("Usage: run <history item number>\n")
}

// NewCLI creates a new command line interpreter
func NewCLI(commands map[string]func(*CLI, []string) error, prompt string, histSize uint8) *CLI {
	c := new(CLI)
	c.commands = commands
	c.prompt = prompt
	c.histSize = histSize
	if c.histSize < 1 {
		c.histSize = 1
	}
	c.hist = make([]string, 0, c.histSize)
	return c
}
