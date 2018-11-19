package cmd

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"github.com/urfave/cli"
)

// CmdMachine cmd
var CmdMachine = cli.Command{
	Name:        "machine",
	Usage:       "搜索机器",
	Description: "正则搜索机器",
	Action:      runMachine,
	Flags: []cli.Flag{
		flagFile,
		flagOutput,
	},
}

func runMachine(c *cli.Context) {
	if len(c.Args()) > 0 {
		if len(c.String("file")) != 0 {
			fmt.Println("转换文件查询 与 命令行查询 不能同时进行")
		} else {
			output := runOutput(c)
			hostname := c.Args()[0]
			var serverList ServerList
			strMap := make(map[string]string)
			for _, server := range serverList.getServerList("loda", "machine") {
				if server.Hostname == "" {
					continue
				}
				if matched, err := regexp.MatchString(hostname, server.Hostname); matched {
					if err != nil {
						fmt.Println(err)
						os.Exit(1)
					}
					strMap[server.IP] = fmt.Sprintf("%s %s", server.IP, server.Hostname)
				}

				if server.IP == "" {
					continue
				}
				if matched, err := regexp.MatchString(hostname, server.IP); matched {
					if err != nil {
						fmt.Println(err)
						os.Exit(1)
					}
					strMap[server.IP] = fmt.Sprintf("%s %s", server.IP, server.Hostname)
				}
			}
			if output == "" {
				for _, v := range strMap {
					fmt.Println(v)
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
		}
	} else {
		if len(c.String("file")) != 0 {
			runFile(c)
		} else {
			fmt.Println("请输入查询的关键词(ip/hostname),支持正则")
		}
	}
}

var flagFile = cli.StringFlag{
	Name:  "file, f",
	Usage: "加载写有ip/hostname的文件",
}

var flagOutput = cli.StringFlag{
	Name:  "output, o",
	Usage: "输出转换成的ip+hostname到 文件路径",
}

func runFile(c *cli.Context) {
	file := c.String("file")
	if len(file) == 0 {
		fmt.Println("请输入 转换文件 的正确路径!")
	} else {
		if PathExist(file) {
			f, err := ioutil.ReadFile(file)
			if err != nil {
				fmt.Println(file, "文件读取失败")
				return
			}
			fString := strings.Replace(string(f), " ", "", -1)
			fArray := strings.Split(fString, "\n")
			var serverList ServerList
			strMap := make(map[string]string)
			for _, server := range serverList.getServerList("loda", "machine") {
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
			output := runOutput(c)
			if output == "" {
				for _, v := range strMap {
					fmt.Println(v)
				}
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

func runOutput(c *cli.Context) string {
	output := c.String("output")
	if len(output) == 0 {
		return ""
	}
	fileParent, _ := GetParentDirectory(output)
	if !PathExist(fileParent) {
		fmt.Println(fileParent, "文件夹不存在，请输入正确的文件路径")
		os.Exit(1)
	}
	return output
}
