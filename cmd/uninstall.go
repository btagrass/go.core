package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"go.core/utl"
)

var (
	Uninstall = &cobra.Command{
		Use:   "uninstall",
		Short: "卸载",
		Run: func(c *cobra.Command, args []string) {
			// 停止服务
			name := cmd.Use
			_, err := utl.Command(fmt.Sprintf("systemctl stop %s", name), fmt.Sprintf("systemctl disable %s", name))
			if err != nil {
				fmt.Println(err)
				return
			}
			// 删除配置
			err = utl.Remove(fmt.Sprintf("/etc/systemd/system/%s.service", name))
			if err != nil {
				fmt.Println(err)
				return
			}
		},
	}
)
