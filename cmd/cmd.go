package cmd

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	cmd = &cobra.Command{} // 命令
)

// 执行
func Execute(name, description string, cmds ...*cobra.Command) {
	cmd.Use = name
	cmd.Short = description
	for _, c := range cmds {
		cmd.AddCommand(c)
	}
	err := cmd.Execute()
	if err != nil {
		logrus.Fatal(err)
	}
}
