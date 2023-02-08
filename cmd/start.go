package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"go.core/utl"
)

var (
	Start = &cobra.Command{
		Use:   "start",
		Short: "启动",
		Run: func(c *cobra.Command, args []string) {
			name := cmd.Use
			_, err := utl.Command(fmt.Sprintf("systemctl restart %s", name))
			if err != nil {
				fmt.Println(err)
				return
			}
		},
	}
)
