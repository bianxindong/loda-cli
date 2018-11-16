package flag

import (
	"fmt"
	"os"

	"github.com/oiooj/cli"
	"github.com/oiooj/loda-cli/cmd"
)

var FlagOutput = cli.StringFlag{
	Name:  "output, o",
	Usage: "输出转换成的ip+hostname到 文件路径",
}

func RunOutput(c *cli.Context) string {
	output := c.String("output")
	if len(output) == 0 {
		return ""
	}
	fileParent := cmd.GetParentDirectory(output)
	if !cmd.PathExist(fileParent) {
		fmt.Println(fileParent, "文件夹不存在，请输入正确的文件路径")
		os.Exit(0)
	}
	return output
}
