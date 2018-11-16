package flag

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"github.com/oiooj/cli"
	"github.com/oiooj/loda-cli/cmd"
)

var FlagFile = cli.StringFlag{
	Name:  "file, f",
	Usage: "加载写有ip/hostname的文件",
}

func RunFile(c *cli.Context, output string) {
	file := c.String("file")
	if len(file) == 0 {
		fmt.Println("请输入 转换文件 的正确路径!")
	} else {
		if cmd.PathExist(file) {
			f, err := ioutil.ReadFile(file)
			if err != nil {
				fmt.Println(file, "文件读取失败")
				return
			}
			fString := strings.Replace(string(f), " ", "", -1)
			fArray := strings.Split(fString, "\n")
			var serverList cmd.ServerList
			strMap := make(map[string]string)
			for _, server := range serverList.GetServerList("loda", "machine") {
				if server.Hostname == "" {
					continue
				}
				if server.IP == "" {
					continue
				}
				for _, hostname := range fArray {
					if hostname == "" {
						continue
					}
					if matched, err := regexp.MatchString(hostname, server.Hostname); matched {
						if err != nil {
							fmt.Println(err)
							os.Exit(1)
						}
						strMap[server.IP] = fmt.Sprintf("%-15s %s", server.IP, server.Hostname)
					}
					if matched, err := regexp.MatchString(hostname, server.IP); matched {
						if err != nil {
							fmt.Println(err)
							os.Exit(1)
						}
						strMap[server.IP] = fmt.Sprintf("%-15s %s", server.IP, server.Hostname)
					}
				}
			}
			if output == "" {
				for _, v := range strMap {
					fmt.Println(v)
				}
				fmt.Println("")
				for _, hostname := range fArray {
					flag := false
					for _, v := range strMap {
						if matched, _ := regexp.MatchString(hostname, v); matched {
							flag = true
						}
					}
					if !flag {
						fmt.Println(hostname)
					}
				}
			} else {
				var body string
				os.Remove(output)
				fout, err := os.Create(output)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				w := bufio.NewWriter(fout)

				for _, v := range strMap {
					body = fmt.Sprintf("%s%s\n", body, v)
				}
				for _, hostname := range fArray {
					flag := false
					for _, v := range strMap {
						if matched, _ := regexp.MatchString(hostname, v); matched {
							flag = true
						}
					}
					if !flag {
						body = fmt.Sprintf("%s%s\n", body, hostname)
					}
				}
				_, err = w.WriteString(body)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				err = w.Flush()
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
			}
		} else {
			fmt.Println(file, "文件不存在，请输入正确的文件路径")
		}
	}
}
