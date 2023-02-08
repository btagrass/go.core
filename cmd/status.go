package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"go.core/utl"
)

var (
	Status = &cobra.Command{
		Use:   "status",
		Short: "状态",
		Run: func(c *cobra.Command, args []string) {
			name := cmd.Use
			_, err := utl.Command(fmt.Sprintf("systemctl status %s", name))
			if err != nil {
				fmt.Println(err)
				return
			}
		},
	}
)
