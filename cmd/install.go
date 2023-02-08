package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"go.core/app"
	"go.core/utl"
)

var (
	Install = &cobra.Command{
		Use:   "install",
		Short: "安装",
		Run: func(c *cobra.Command, args []string) {
			// 添加配置
			name := cmd.Use
			err := os.WriteFile(
				fmt.Sprintf("/etc/systemd/system/%s.service", name),
				[]byte(fmt.Sprintf(`
[Unit]
Description=%s
After=network.target

[Service]
Type=simple
WorkingDirectory=%s
ExecStart=%s run
Restart=always
RestartSec=30s

[Install]
WantedBy=multi-user.target
`,
					strings.ToUpper(name),
					app.Dir,
					filepath.Join(app.Dir, name),
				)),
				os.ModePerm,
			)
			if err != nil {
				fmt.Println(err)
				return
			}
			// 启动服务
			_, err = utl.Command(
				fmt.Sprintf("systemctl enable %s", name),
				"systemctl daemon-reload",
				fmt.Sprintf("systemctl restart %s", name),
			)
			if err != nil {
				fmt.Println(err)
				return
			}
		},
	}
)
