package main

import (
	"os"

	"github.com/marcsauter/cli"
)

const (
	Prompt   = ">>> "
	Histsize = 255
)

var (
	Commands map[string]func(*cli.CLI, []string) error
)

func init() {
	Commands = map[string]func(*cli.CLI, []string) error{
		"help":   help,
		"greet":  greet,
		"secret": secret,
		"exit":   exit,
	}
}

func main() {
	c := cli.NewCLI(Commands, Prompt, Histsize)
	c.Info("Type help ...\n")
	for {
		c.Prompt()
		c.Exec(c.ReadLine())
	}
}

func help(c *cli.CLI, command []string) error {
	helptext := `
help            this text
greet <name>    greets <name>
secret          enter a secret - no echo
hist            show history of commands (builtin)
run <index>     run the command with <index> in history (builtin)
exit            exit program

`
	c.Info(helptext)
	return nil
}

func greet(c *cli.CLI, command []string) error {
	if len(command) < 2 {
		c.Info("Usage: greet <name>\n")
		return nil
	}
	c.Info("Hello %s\n", command[1])
	return nil
}

func secret(c *cli.CLI, command []string) error {
	prompt := c.SetPrompt("Enter your Secret: ")
	c.Prompt()
	if err := c.EchoOff(); err != nil {
		c.Error("could not turn echo off\n")
		return nil
	}
	secret := c.ReadLine()
	c.NewLine()
	c.Info("Your Secret is: %s\n", secret)
	if err := c.EchoOn(); err != nil {
		c.Error("could not turn echo on again\n")
	}
	c.SetPrompt(prompt)
	return nil
}

func exit(c *cli.CLI, command []string) error {
	if c.YesNo("Are you sure? (y)es/(n)o", []string{"y"}) {
		os.Exit(0)
	}
	return nil
}
