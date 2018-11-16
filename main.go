package main

import (
	"fmt"
	"os"

	"github.com/oiooj/cli"
	"github.com/oiooj/loda-cli/cmd"
	"github.com/oiooj/loda-cli/flag"
	"github.com/oiooj/loda-cli/setting"
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
	app.Flags = []cli.Flag{
		flag.FlagFile,
		flag.FlagOutput,
	}
	app.Action = func(c *cli.Context) {
		if len(c.Args()) != 0 {
			fmt.Println("loda-cli -h查看使用用法")
		} else {
			output := flag.RunOutput(c)
			flag.RunFile(c, output)
		}
	}
	app.Run(os.Args)
}
