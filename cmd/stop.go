package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"go.core/utl"
)

var (
	Stop = &cobra.Command{
		Use:   "stop",
		Short: "停止",
		Run: func(c *cobra.Command, args []string) {
			name := cmd.Use
			_, err := utl.Command(fmt.Sprintf("systemctl stop %s", name))
			if err != nil {
				fmt.Println(err)
				return
			}
		},
	}
)
