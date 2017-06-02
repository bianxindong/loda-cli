package cmd

import (
	"fmt"
	"os"
	"regexp"

	"github.com/oiooj/cli"
)

var CmdMachine = cli.Command{
	Name:        "machine",
	Usage:       "搜索机器",
	Description: "正则搜索机器",
	Action:      runMachine,
}

func runMachine(c *cli.Context) {
	if len(c.Args()) > 0 {
		hostname := c.Args()[0]
		var serverList ServerList
		for _, server := range serverList.getServerList("loda", "machine") {
			if matched, err := regexp.MatchString(hostname, server.Hostname); matched {
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				fmt.Println(server.IP)
			}
		}

	} else {
		fmt.Println("Input anything plz")
	}
}
