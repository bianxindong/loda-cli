package main

import (
	"os"

	"github.com/oiooj/loda-cli/cmd"
	"github.com/oiooj/loda-cli/setting"
	"github.com/urfave/cli"
)

func main() {

	app := cli.NewApp()
	app.EnableBashCompletion = true
	app.Name = setting.AppName
	app.Usage = setting.Usage
	app.Version = setting.Version
	app.Author = setting.Author
	app.Email = setting.Email

	app.Commands = []cli.Command{
		cmd.CmdTree,
		cmd.CmdMachine,
		cmd.CmdNamed,
	}
	app.Flags = append(app.Flags, []cli.Flag{}...)
	app.Run(os.Args)
}
